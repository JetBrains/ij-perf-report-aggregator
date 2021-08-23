package server

import (
  "compress/gzip"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/http-error"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/dgraph-io/ristretto"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "github.com/zeebo/xxh3"
  "go.uber.org/zap"
  "io"
  "net/http"
  "strconv"
  "strings"
)

var byteBufferPool bytebufferpool.Pool

type ResponseCacheManager struct {
  cache  *ristretto.Cache
  logger *zap.Logger
}

type CachingHandler struct {
  handler func(request *http.Request) ([]byte, error)
  manager *ResponseCacheManager
}

func (t *CachingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.manager.handle(w, r, t.handler)
}

func NewResponseCacheManager(logger *zap.Logger) (*ResponseCacheManager, error) {
  // 256 MB
  cacheSize := 256 * 1000 * 1000
  cache, err := ristretto.NewCache(&ristretto.Config{
    NumCounters: int64((cacheSize / 50 /* assume that each response ~ 50 KB */) * 10) /* number of keys to track frequency of */,
    MaxCost:     int64(cacheSize),
    BufferItems: 64 /* number of keys per Get buffer */,
  })
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &ResponseCacheManager{
    cache:  cache,
    logger: logger,
  }, nil
}

func (t *ResponseCacheManager) CreateHandler(handler func(request *http.Request) ([]byte, error)) http.Handler {
  return &CachingHandler{
    handler: handler,
    manager: t,
  }
}

var jsonMarker = []byte("json:")

func generateKey(request *http.Request) []byte {
  buffer := quicktemplate.AcquireByteBuffer()
  defer quicktemplate.ReleaseByteBuffer(buffer)

  // if json requested, it means that handler can return data in several formats
  if request.Header.Get("Accept") == "application/json" {
    _, _ = buffer.Write(jsonMarker)
  }

  u := request.URL
  _, _ = io.WriteString(buffer, u.Path)
  // do not complicate, use RawQuery as is without sorting
  if len(u.RawQuery) > 0 {
    _, _ = io.WriteString(buffer, u.RawQuery)
  }

  result := make([]byte, len(buffer.B))
  copy(result, buffer.B)
  return result
}

func (t *ResponseCacheManager) handle(w http.ResponseWriter, request *http.Request, handler func(request *http.Request) ([]byte, error)) {
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Cache-Control", "public, must-revalidate, max-age=15")
	w.Header().Set("Vary", "Accept-Encoding")

  cacheKey := generateKey(request)
  value, found := t.cache.Get(cacheKey)
  var eTag string
  var result []byte
  if found {
    result = value.([]byte)
    prevEtag := request.Header.Get("If-None-Match")
    eTag = computeEtag(result)
    if prevEtag == eTag {
      w.Header().Set("ETag", eTag)
      w.WriteHeader(http.StatusNotModified)
      return
    }
  } else {
    var err error
    result, err = handler(request)
    if err != nil {
      t.handleError(err, w)
      return
    }

    t.cache.Set(cacheKey, result, int64(len(result)))
    eTag = computeEtag(result)
  }

  w.Header().Set("ETag", eTag)

  if len(result) > 8192 && strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
    w.Header().Set("Content-Encoding", "gzip")

    compressedData, err := t.gzipData(result)
    if err != nil {
      t.logger.Error("cannot compress result", zap.Error(err))
      http.Error(w, err.Error(), 503)
      return
    }

    result = compressedData
  }

  _, err := w.Write(result)
  if err != nil {
    t.logger.Error("cannot write cached result", zap.Error(err))
    http.Error(w, err.Error(), 503)
  }
}

func (t *ResponseCacheManager) handleError(err error, w http.ResponseWriter) {
  switch exception := errors.Cause(err).(type) {
  case *http_error.HttpError:
    w.WriteHeader(exception.Code)
    writehttpError(w, exception)

  default:
    //fmt.Printf("%+v", err)
    t.logger.Error("cannot handle http request", zap.Error(err))
    http.Error(w, err.Error(), 503)
  }
}

// we don't add here salt like "when server was started" to reflect changes in a new server logic,
// because if no data in cache (server restarted), in any case data will be recomputed
func computeEtag(result []byte) string {
  // add length to reduce chance of collision
  return strconv.FormatUint(xxh3.Hash(result), 36) + "-" + strconv.FormatInt(int64(len(result)), 36)
}

func (t *ResponseCacheManager) gzipData(value []byte) ([]byte, error) {
  buffer := bytebufferpool.Get()
  defer bytebufferpool.Put(buffer)
  gzipWriter, err := gzip.NewWriterLevel(buffer, 5)
  if err != nil {
    return nil, err
  }

  _, err = gzipWriter.Write(value)
  if err != nil {
    util.Close(gzipWriter, t.logger)
    return nil, err
  }

  util.Close(gzipWriter, t.logger)

  return CopyBuffer(buffer), nil
}

func (t *ResponseCacheManager) Clear() {
  t.cache.Clear()
}

func CopyBuffer(buffer *bytebufferpool.ByteBuffer) []byte {
  result := make([]byte, len(buffer.B))
  copy(result, buffer.B)
  return result
}
