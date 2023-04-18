package server

import (
  "context"
  "errors"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/http-error"
  e "github.com/develar/errors"
  "github.com/dgraph-io/ristretto"
  "github.com/valyala/bytebufferpool"
  "github.com/zeebo/xxh3"
  "go.uber.org/zap"
  "io"
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
  cache  *ristretto.Cache
  logger *zap.Logger
}

func NewResponseCacheManager(logger *zap.Logger) (*ResponseCacheManager, error) {
  cacheSize := 1000 * 1000 * 1000 // 1 GiB
  cache, err := ristretto.NewCache(&ristretto.Config{
    NumCounters: int64((cacheSize / 50 /* assume that each response ~ 50 KB */) * 10) /* number of keys to track frequency of */,
    MaxCost:     int64(cacheSize),
    BufferItems: 64 /* number of keys per Get buffer */,
  })
  if err != nil {
    return nil, e.WithStack(err)
  }
  return &ResponseCacheManager{
    cache:  cache,
    logger: logger,
  }, nil
}

func (rcm *ResponseCacheManager) CreateHandler(handler func(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error)) http.Handler {
  return &CachingHandler{
    handler: handler,
    manager: rcm,
  }
}

func (rcm *ResponseCacheManager) handle(w http.ResponseWriter, request *http.Request, handler func(request *http.Request) (*bytebufferpool.ByteBuffer, bool, error), ) {
  w.Header().Set("Vary", "Accept-Encoding")

  cacheKey := generateCacheKey(request)
  value, found := rcm.cache.Get(cacheKey)
  var result []byte
  if found {
    cacheData, ok := value.([]byte)
    if !ok {
      return
    }
    var err error
    result, err = decompressData(cacheData)
    if err != nil {
      rcm.handleError(err, w)
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
      rcm.logger.Error("cannot compress result", zap.Error(err))
      http.Error(w, err.Error(), http.StatusServiceUnavailable)
      return
    }
    rcm.cache.Set(cacheKey, result, int64(len(result)))
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
    rcm.logger.Error("cannot write cached result", zap.Error(err))
    http.Error(w, err.Error(), http.StatusServiceUnavailable)
  }
}

func generateCacheKey(request *http.Request) []byte {
  buffer := bytebufferpool.Get()
  defer bytebufferpool.Put(buffer)
  // if json requested, it means that handler can return data in several formats
  if request.Header.Get("Accept") == "application/json" {
    _, _ = buffer.Write([]byte("j:"))
  }
  u := request.URL
  _, _ = io.WriteString(buffer, u.Path)
  // do not complicate, use RawQuery as is without sorting
  if len(u.RawQuery) > 0 {
    _, _ = io.WriteString(buffer, u.RawQuery)
  }
  return CopyBuffer(buffer)
}

func (rcm *ResponseCacheManager) handleError(err error, w http.ResponseWriter) {
  cause := e.Cause(err)
  var httpError *http_error.HttpError
  if errors.As(cause, &httpError) {
    w.WriteHeader(httpError.Code)
    writehttpError(w, httpError)
  } else {
    if errors.Is(cause, context.Canceled) {
      http.Error(w, err.Error(), 499)
    } else {
      rcm.logger.Error("cannot handle http request", zap.Error(err))
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
  rcm.cache.Clear()
}

func CopyBuffer(buffer *bytebufferpool.ByteBuffer) []byte {
  result := make([]byte, len(buffer.B))
  copy(result, buffer.B)
  return result
}
