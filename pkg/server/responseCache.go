package server

import (
	"context"
	"errors"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/http-error"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/valyala/bytebufferpool"
	"github.com/zeebo/xxh3"
	"log/slog"
	"net/http"
	"strconv"
)

var byteBufferPool bytebufferpool.Pool

type CachingHandler struct {
	handler func(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error)
	manager *ResponseCacheManager
}

func (ch *CachingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ch.manager.handle(w, r, ch.handler)
}

type ResponseCacheManager struct {
	cache *fastcache.Cache
}

func NewResponseCacheManager() (*ResponseCacheManager, error) {
	cacheSize := 1000 * 1000 * 1000
	cache := fastcache.New(cacheSize)
	return &ResponseCacheManager{
		cache: cache,
	}, nil
}

func (rcm *ResponseCacheManager) CreateHandler(handler func(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error)) http.Handler {
	return &CachingHandler{
		handler: handler,
		manager: rcm,
	}
}

func (rcm *ResponseCacheManager) handle(w http.ResponseWriter, request *http.Request, handler func(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error)) {
	w.Header().Set("Vary", "Accept-Encoding")

	cacheKey := generateCacheKey(request)
	value := rcm.cache.Get(nil, cacheKey)
	var result []byte
	if value != nil && request.Header.Get("Cache-Control") != "no-cache" {
		var err error
		result, err = decompressData(value)
		if err != nil {
			slog.Error("cannot decompress cached result", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		prevEtag := request.Header.Get("If-None-Match")
		eTag := computeEtag(result)
		if prevEtag == eTag {
			w.Header().Set("ETag", eTag)
			w.WriteHeader(http.StatusNotModified)
			return
		}
	} else {
		buffer, releaseBuffer, err := handler(request)
		if err != nil {
			rcm.handleError(err, w)
			return
		}
		result, err = rcm.compressData(buffer.B)
		if err != nil {
			slog.Error("cannot compress result", "error", err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		rcm.cache.Set(cacheKey, result)
		result = buffer.B
		if releaseBuffer {
			bytebufferpool.Put(buffer)
		}
	}

	w.Header().Set("ETag", computeEtag(result))
	w.Header().Set("Content-Length", strconv.Itoa(len(result)))
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(result)
	if err != nil {
		slog.Error("cannot write cached result", "error", err)
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
}

func generateCacheKey(request *http.Request) []byte {
	buffer := bytebufferpool.Get()
	defer bytebufferpool.Put(buffer)
	// if json requested, it means that handler can return data in several formats
	if request.Header.Get("Accept") == "application/json" {
		_, _ = buffer.WriteString("j:")
	}
	u := request.URL
	_, _ = buffer.WriteString(u.Path)
	// do not complicate, use RawQuery as is without sorting
	if u.RawQuery != "" {
		_, _ = buffer.WriteString(u.RawQuery)
	}
	return CopyBuffer(buffer)
}

func (rcm *ResponseCacheManager) handleError(err error, w http.ResponseWriter) {
	var httpError *http_error.HttpError
	if errors.As(err, &httpError) {
		w.WriteHeader(httpError.Code)
		writehttpError(w, httpError)
	} else {
		if errors.Is(err, context.Canceled) {
			http.Error(w, err.Error(), 499)
		} else {
			slog.Error("cannot handle http request", "error", err)
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
		}
	}
}

// we don't add here salt like "when server was started" to reflect changes in a new server logic,
// because if no data in cache (server restarted), in any case data will be recomputed
func computeEtag(result []byte) string {
	hash := xxh3.Hash128(result)
	return strconv.FormatUint(hash.Hi, 36) + "-" + strconv.FormatUint(hash.Lo, 36)
}

func (rcm *ResponseCacheManager) Clear() {
	rcm.cache.Reset()
}

func CopyBuffer(buffer *bytebufferpool.ByteBuffer) []byte {
	result := make([]byte, len(buffer.B))
	copy(result, buffer.B)
	return result
}
