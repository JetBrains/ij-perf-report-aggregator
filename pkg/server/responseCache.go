package server

import (
  "compress/gzip"
  "fmt"
  "github.com/VictoriaMetrics/fastcache"
  "github.com/valyala/quicktemplate"
  "go.uber.org/zap"
  "io"
  "net/http"
  "net/url"
  "report-aggregator/pkg/util"
  "strings"
)

type HttpError struct {
  Code    int
  Message string
}

func (t *HttpError) Error() string {
  return fmt.Sprintf("code=%d, message: %s", t.Code, t.Message)
}

func NewHttpError(code int, message string) error {
  return &HttpError{
    Code:    code,
    Message: message,
  }
}

type ResponseCacheManager struct {
  cache  *fastcache.Cache
  logger *zap.Logger
}

type CachingHandler struct {
  handler   func(request *http.Request) ([]byte, error)
  manager   *ResponseCacheManager
}

func (t *CachingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.manager.handle(w, r, t.handler)
}

func NewResponseCacheManager() *ResponseCacheManager {
  return &ResponseCacheManager{
    cache: fastcache.New(500 * 1000 * 1000 /* 500 MB */),
  }
}

func (t *ResponseCacheManager) CreateHandler(handler func(request *http.Request) ([]byte, error)) http.Handler {
  return &CachingHandler{
    handler:   handler,
    manager:   t,
  }
}

func generateKey(url *url.URL) []byte {
  buffer := quicktemplate.AcquireByteBuffer()
  defer quicktemplate.ReleaseByteBuffer(buffer)

  _, _ = io.WriteString(buffer, url.Path)
  // do not complicate, use RawQuery as is without sorting
  if len(url.RawQuery) > 0 {
    _, _ = io.WriteString(buffer, url.RawQuery)
  }

  result := make([]byte, len(buffer.B))
  copy(result, buffer.B)
  return result
}

func (t *ResponseCacheManager) get(url *url.URL) ([]byte, []byte) {
  key := generateKey(url)
  var value []byte
  value = t.cache.GetBig(value, key)
  if len(value) == 0 {
    return key, nil
  }
  return nil, value
}

func (t *ResponseCacheManager) handle(w http.ResponseWriter, request *http.Request, handler func(request *http.Request) ([]byte, error)) {
  w.Header().Set("Content-Type", "application/json")

  cacheKey, result := t.get(request.URL)
  if result == nil {
    var err error
    result, err = handler(request)
    if err != nil {
      t.logger.Error("cannot handle http request", zap.Error(err))
      http.Error(w, err.Error(), 503)
      return
    }

    t.cache.SetBig(cacheKey, result)
  }

  ac := request.Header.Get("Accept-Encoding")
  if len(result) > 8192 && strings.Contains(ac, "gzip") {
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
  return
}

func (t *ResponseCacheManager) gzipData(value []byte) ([]byte, error) {
  buffer := quicktemplate.AcquireByteBuffer()
  defer quicktemplate.ReleaseByteBuffer(buffer)
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

  result := make([]byte, len(buffer.B))
  copy(result, buffer.B)
  return result, nil
}
