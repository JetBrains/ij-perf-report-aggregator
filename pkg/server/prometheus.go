package server

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusMetrics struct {
	handler http.Handler

	inFlightRequests    prometheus.Gauge
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
	httpUserRequests    *prometheus.CounterVec

	responseCacheRequestsTotal *prometheus.CounterVec
	responseCacheClearsTotal   prometheus.Counter

	userSeenWindow time.Duration
	userSeenMu     sync.Mutex
	userLastSeen   map[string]time.Time
}

const maxRouteLabelLength = 64

var chiPathParamPattern = regexp.MustCompile(`\{[^}/]+\}`)

func NewPrometheusMetrics() *PrometheusMetrics {
	registry := prometheus.NewRegistry()

	metrics := &PrometheusMetrics{
		inFlightRequests: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "ij_perf",
			Subsystem: "http",
			Name:      "in_flight_requests",
			Help:      "Current number of in-flight HTTP requests.",
		}),
		httpRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "ij_perf",
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "Total number of HTTP requests handled by the backend.",
		}, []string{"method", "route", "status_code"}),
		httpRequestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "ij_perf",
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Latency of HTTP requests handled by the backend.",
			Buckets:   []float64{0.01, 0.025, 0.05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60, 120, 300},
		}, []string{"method", "route", "status_code"}),
		httpUserRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "ij_perf",
			Subsystem: "http",
			Name:      "user_requests_total",
			Help:      "Active-user activity counter, debounced to roughly one increment per user per minute so it isn't dominated by chatty page loads. Only authenticated requests (with X-Auth-Request-Email) are counted. Label is the local-part of the email. Use count(count by (user) (rate(...[5m]))) for unique active users.",
		}, []string{"user"}),
		responseCacheRequestsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "ij_perf",
			Subsystem: "response_cache",
			Name:      "requests_total",
			Help:      "Total number of cache lookups grouped by result.",
		}, []string{"result"}),
		responseCacheClearsTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "ij_perf",
			Subsystem: "response_cache",
			Name:      "clears_total",
			Help:      "Total number of response-cache clear operations.",
		}),
		userSeenWindow: time.Minute,
		userLastSeen:   make(map[string]time.Time),
	}

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		metrics.inFlightRequests,
		metrics.httpRequestsTotal,
		metrics.httpRequestDuration,
		metrics.httpUserRequests,
		metrics.responseCacheRequestsTotal,
		metrics.responseCacheClearsTotal,
	)

	metrics.handler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	return metrics
}

func (m *PrometheusMetrics) Handler() http.Handler {
	return m.handler
}

func (m *PrometheusMetrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.inFlightRequests.Inc()
		defer m.inFlightRequests.Dec()

		start := time.Now()
		ww := chimiddleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		statusCode := ww.Status()
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		labels := prometheus.Labels{
			"method":      r.Method,
			"route":       routeLabel(sanitizeRoutePattern(routePattern(r))),
			"status_code": strconv.Itoa(statusCode),
		}

		m.httpRequestsTotal.With(labels).Inc()
		m.httpRequestDuration.With(labels).Observe(time.Since(start).Seconds())

		if user := userLabel(r.Header.Get("X-Auth-Request-Email")); user != "" && m.shouldRecordUser(user, time.Now()) {
			m.httpUserRequests.WithLabelValues(user).Inc()
		}
	})
}

func (m *PrometheusMetrics) ObserveCacheLookup(result string) {
	if m == nil {
		return
	}
	m.responseCacheRequestsTotal.WithLabelValues(result).Inc()
}

func (m *PrometheusMetrics) ObserveCacheClear() {
	if m == nil {
		return
	}
	m.responseCacheClearsTotal.Inc()
}

func (m *PrometheusMetrics) shouldRecordUser(user string, now time.Time) bool {
	m.userSeenMu.Lock()
	defer m.userSeenMu.Unlock()
	if last, ok := m.userLastSeen[user]; ok && now.Sub(last) < m.userSeenWindow {
		return false
	}
	m.userLastSeen[user] = now
	return true
}

func userLabel(email string) string {
	email = strings.ToLower(email)
	if at := strings.IndexByte(email, '@'); at > 0 {
		return email[:at]
	}
	return email
}

func routePattern(r *http.Request) string {
	routeContext := chi.RouteContext(r.Context())
	if routeContext == nil {
		return "unmatched"
	}

	pattern := routeContext.RoutePattern()
	if pattern == "" {
		return "unmatched"
	}
	return pattern
}

// routeLabel returns the Prometheus `route` label for a request.
// Anything chi did not route (404s, scanner traffic, early-middleware rejects)
// collapses into a single "unmatched" bucket — without this, security scans
// generate one series per unique path and break Prometheus scraping.
func routeLabel(pattern string) string {
	if pattern != "" && !isHighCardinalityRouteLabel(pattern) {
		return pattern
	}
	return "unmatched"
}

func sanitizeRoutePattern(pattern string) string {
	if pattern == "" || pattern == "unmatched" {
		return ""
	}

	pattern = chiPathParamPattern.ReplaceAllString(pattern, ":param")
	pattern = strings.TrimSuffix(pattern, "/*")
	pattern = strings.TrimSuffix(pattern, "*")

	if pattern == "" {
		return "/"
	}
	return pattern
}

func isHighCardinalityRouteLabel(label string) bool {
	return len(label) > maxRouteLabelLength ||
		strings.Contains(label, "://") ||
		strings.ContainsAny(label, "?%")
}
