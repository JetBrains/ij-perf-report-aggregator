package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/time/rate"
)

const (
	// Sized to comfortably hold a full client batch: the UI batches up to 10 events, each with
	// per-field caps (message 1000, stack 4000, url/pageUrl 1000 each, userAgent 300) that are
	// truncated by character count, so multibyte UTF-8 can inflate the byte size several-fold.
	// 256 KiB keeps a worst-case batch from being rejected with 400 while still bounding memory.
	clientErrorMaxBodyBytes              = 256 * 1024
	clientErrorMaxBatchSize              = 20
	clientErrorMaxEventCount             = 20
	clientErrorMaxReportedCountPerBatch  = 100
	clientErrorMaxMessageLength          = 1000
	clientErrorMaxStackLength            = 4000
	clientErrorMaxURLLength              = 1000
	clientErrorMaxUserAgentLength        = 300
	clientErrorMaxVersionLength          = 64
	clientErrorMaxTimestampLength        = 64
	clientErrorQueueSize                 = 512
	clientErrorRateLimitBurst            = 100
	clientErrorRateLimitPerMinute        = 100
	clientErrorMaxRateLimitTrackedKeys   = 4096
	clientErrorRateLimitTrackedKeyMaxAge = 10 * time.Minute
)

type ClientErrorReporter struct {
	metrics     *PrometheusMetrics
	queue       chan []clientErrorEvent
	rateLimiter *clientErrorRateLimiter
	wg          sync.WaitGroup
	closeOnce   sync.Once
	// mu guards queue sends against Close: closing a channel while a handler is mid-send panics.
	// Senders hold RLock; Close takes the write lock so it can't close until in-flight sends finish.
	mu     sync.RWMutex
	closed bool
}

type clientErrorRequest struct {
	Errors []clientErrorEvent `json:"errors"`
}

type clientErrorEvent struct {
	Source    string `json:"source"`
	Version   string `json:"version,omitempty"`
	Message   string `json:"message,omitempty"`
	Stack     string `json:"stack,omitempty"`
	URL       string `json:"url,omitempty"`
	PageURL   string `json:"pageUrl,omitempty"`
	UserAgent string `json:"userAgent,omitempty"`
	Line      int    `json:"line,omitempty"`
	Column    int    `json:"column,omitempty"`
	Count     int    `json:"count,omitempty"`
	FirstSeen string `json:"firstSeen,omitempty"`
	LastSeen  string `json:"lastSeen,omitempty"`
}

func NewClientErrorReporter(metrics *PrometheusMetrics) *ClientErrorReporter {
	reporter := &ClientErrorReporter{
		metrics:     metrics,
		queue:       make(chan []clientErrorEvent, clientErrorQueueSize),
		rateLimiter: newClientErrorRateLimiter(clientErrorRateLimitBurst, float64(clientErrorRateLimitPerMinute)/float64(time.Minute/time.Second), time.Now),
	}
	reporter.wg.Add(1)
	go reporter.run()
	return reporter
}

func (r *ClientErrorReporter) Close() {
	r.closeOnce.Do(func() {
		r.mu.Lock()
		r.closed = true
		close(r.queue)
		r.mu.Unlock()
		r.wg.Wait()
	})
}

func (r *ClientErrorReporter) Handler() http.Handler {
	return http.HandlerFunc(r.handle)
}

func (r *ClientErrorReporter) handle(w http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	var payload clientErrorRequest
	decoder := json.NewDecoder(http.MaxBytesReader(w, request.Body, clientErrorMaxBodyBytes))
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		http.Error(w, "Invalid request body: multiple JSON values", http.StatusBadRequest)
		return
	}

	events, reportedCount := normalizeClientErrorBatch(payload)
	if len(events) == 0 {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	if !r.rateLimiter.Allow(clientErrorRateLimitKey(request), reportedCount) {
		http.Error(w, "client error report rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.closed {
		http.Error(w, "client error reporter is shutting down", http.StatusServiceUnavailable)
		return
	}
	select {
	case r.queue <- events:
		w.WriteHeader(http.StatusAccepted)
	default:
		http.Error(w, "client error report queue is full", http.StatusServiceUnavailable)
	}
}

func (r *ClientErrorReporter) run() {
	defer r.wg.Done()
	for events := range r.queue {
		r.process(events)
	}
}

func (r *ClientErrorReporter) process(events []clientErrorEvent) {
	for _, event := range events {
		r.metrics.ObserveClientError(event.Source, event.Version, event.Count)
		slog.Warn(
			"client UI error",
			"source", event.Source,
			"version", event.Version,
			"count", event.Count,
			"message", event.Message,
			"stack", event.Stack,
			"url", event.URL,
			"page_url", event.PageURL,
			"line", event.Line,
			"column", event.Column,
			"user_agent", event.UserAgent,
			"first_seen", event.FirstSeen,
			"last_seen", event.LastSeen,
		)
	}
}

func normalizeClientErrorBatch(payload clientErrorRequest) ([]clientErrorEvent, int) {
	limit := min(len(payload.Errors), clientErrorMaxBatchSize)
	events := make([]clientErrorEvent, 0, limit)
	reportedCount := 0
	for i := range limit {
		if reportedCount >= clientErrorMaxReportedCountPerBatch {
			break
		}

		event, ok := normalizeClientErrorEvent(payload.Errors[i])
		if !ok {
			continue
		}
		if reportedCount+event.Count > clientErrorMaxReportedCountPerBatch {
			event.Count = clientErrorMaxReportedCountPerBatch - reportedCount
		}
		if event.Count <= 0 {
			continue
		}

		reportedCount += event.Count
		events = append(events, event)
	}
	return events, reportedCount
}

func normalizeClientErrorEvent(event clientErrorEvent) (clientErrorEvent, bool) {
	event.Source = clientErrorSourceLabel(event.Source)
	event.Version = clientErrorVersionLabel(event.Version)
	event.Message = truncateClientErrorString(strings.TrimSpace(event.Message), clientErrorMaxMessageLength)
	event.Stack = truncateClientErrorString(strings.TrimSpace(event.Stack), clientErrorMaxStackLength)
	event.URL = truncateClientErrorString(strings.TrimSpace(event.URL), clientErrorMaxURLLength)
	event.PageURL = truncateClientErrorString(strings.TrimSpace(event.PageURL), clientErrorMaxURLLength)
	event.UserAgent = truncateClientErrorString(strings.TrimSpace(event.UserAgent), clientErrorMaxUserAgentLength)
	event.FirstSeen = truncateClientErrorString(strings.TrimSpace(event.FirstSeen), clientErrorMaxTimestampLength)
	event.LastSeen = truncateClientErrorString(strings.TrimSpace(event.LastSeen), clientErrorMaxTimestampLength)

	if event.Line < 0 {
		event.Line = 0
	}
	if event.Column < 0 {
		event.Column = 0
	}
	if event.Count < 1 {
		event.Count = 1
	} else if event.Count > clientErrorMaxEventCount {
		event.Count = clientErrorMaxEventCount
	}

	return event, event.Message != "" || event.Stack != "" || event.URL != "" || event.PageURL != ""
}

func truncateClientErrorString(value string, maxLength int) string {
	if len(value) <= maxLength {
		return value
	}
	// Back off to the nearest rune boundary at or before the byte limit so we never emit
	// invalid UTF-8 (a multibyte rune split in the middle) into the structured logs.
	end := maxLength
	for end > 0 && !utf8.RuneStart(value[end]) {
		end--
	}
	return value[:end]
}

func clientErrorSourceLabel(source string) string {
	switch source {
	case "window_error", "unhandled_rejection", "console_error", "resource_error", "vue_error":
		return source
	default:
		return "unknown"
	}
}

// clientErrorVersionLabel sanitizes the client-supplied build version before it is used as a
// Prometheus label. The value is untrusted (it crosses the wire from the browser), so it must be
// bounded the same way userLabel bounds the user label: a restricted charset and a length cap keep
// a buggy or hostile client from injecting newlines, label-breaking characters, or arbitrarily long
// strings. Anything outside that shape collapses into "unknown". This bounds the *shape* of each
// label; ObserveClientError additionally bounds the *number* of distinct version series.
func clientErrorVersionLabel(version string) string {
	version = strings.TrimSpace(version)
	if version == "" || len(version) > clientErrorMaxVersionLength || !isValidVersionLabel(version) {
		return "unknown"
	}
	return version
}

func isValidVersionLabel(s string) bool {
	for i := range len(s) {
		c := s[i]
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '.' || c == '-' || c == '_' || c == '+') {
			return false
		}
	}
	return true
}

func clientErrorRateLimitKey(request *http.Request) string {
	if user := userLabel(request.Header.Get("X-Auth-Request-Email")); user != "" {
		return "user:" + user
	}

	// Deliberately key on RemoteAddr (the real TCP peer) and never on client-supplied
	// X-Forwarded-For / X-Real-IP headers: those are trivially spoofable, so trusting them would
	// let an unauthenticated caller mint unlimited distinct rate-limit keys and bypass the per-client
	// limit entirely. The trade-off is that, behind a reverse proxy, all unauthenticated traffic
	// collapses onto the proxy's address and shares one bucket — which fails closed (more
	// restrictive), the safe direction for a best-effort, low-value error sink. Real traffic is
	// authenticated and keyed per user above.
	if host, _, err := net.SplitHostPort(request.RemoteAddr); err == nil && host != "" {
		return "ip:" + truncateClientErrorString(host, 64)
	}
	return "ip:" + truncateClientErrorString(request.RemoteAddr, 64)
}

type clientErrorRateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*clientErrorRateLimitEntry
	limit    rate.Limit
	burst    int
	now      func() time.Time
}

type clientErrorRateLimitEntry struct {
	limiter *rate.Limiter
	last    time.Time
}

func newClientErrorRateLimiter(capacity int, refillPerSecond float64, now func() time.Time) *clientErrorRateLimiter {
	if capacity < 1 {
		capacity = 1
	}
	return &clientErrorRateLimiter{
		limiters: make(map[string]*clientErrorRateLimitEntry),
		limit:    rate.Limit(refillPerSecond),
		burst:    capacity,
		now:      now,
	}
}

func (l *clientErrorRateLimiter) Allow(key string, cost int) bool {
	if cost < 1 {
		cost = 1
	}
	if cost > l.burst {
		cost = l.burst
	}
	now := l.now()

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := l.limiters[key]
	if entry == nil {
		// Only enforce the cap when adding a new key — existing keys don't grow the map. Since the
		// map grows by at most one key per call, making room for one new key needs a single scan.
		if len(l.limiters) >= clientErrorMaxRateLimitTrackedKeys {
			l.makeRoomForNewKey(now)
		}
		entry = &clientErrorRateLimitEntry{
			limiter: rate.NewLimiter(l.limit, l.burst),
			last:    now,
		}
		l.limiters[key] = entry
	}
	entry.last = now
	return entry.limiter.AllowN(now, cost)
}

// makeRoomForNewKey frees a slot in a single pass: it drops every stale entry and, in the same
// scan, remembers the oldest survivor. If nothing was stale (a flood of fresh, distinct keys —
// e.g. spoofed sources), it evicts that oldest survivor as a hard cap so the map can't grow
// unbounded. One scan, not two, even under a sustained flood.
func (l *clientErrorRateLimiter) makeRoomForNewKey(now time.Time) {
	cutoff := now.Add(-clientErrorRateLimitTrackedKeyMaxAge)
	var oldestKey string
	var oldest time.Time
	found := false
	removedStale := false
	for key, entry := range l.limiters {
		if entry.last.Before(cutoff) {
			delete(l.limiters, key)
			removedStale = true
			continue
		}
		if !found || entry.last.Before(oldest) {
			oldestKey = key
			oldest = entry.last
			found = true
		}
	}
	if !removedStale && found {
		delete(l.limiters, oldestKey)
	}
}
