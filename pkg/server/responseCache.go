package server

import (
  "compress/gzip"
  "fmt"
  "github.com/VictoriaMetrics/fastcache"
  "github.com/cespare/xxhash"
  "github.com/develar/errors"
  "github.com/valyala/bytebufferpool"
  "github.com/valyala/quicktemplate"
  "go.uber.org/zap"
  "io"
  "net/http"
  "net/url"
  "report-aggregator/pkg/util"
  "strconv"
  "strings"
  "time"
)

var byteBufferPool bytebufferpool.Pool

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

  lastModified time.Time
}

type CachingHandler struct {
  handler func(request *http.Request) ([]byte, error)
  manager *ResponseCacheManager
}

func (t *CachingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  t.manager.handle(w, r, t.handler)
}

func NewResponseCacheManager(logger *zap.Logger) *ResponseCacheManager {
  return &ResponseCacheManager{
    cache:        fastcache.New(512 * 1000 * 1000 /* 512 MB */),
    lastModified: time.Now(),
    logger:       logger,
  }
}

func (t *ResponseCacheManager) CreateHandler(handler func(request *http.Request) ([]byte, error)) http.Handler {
  return &CachingHandler{
    handler: handler,
    manager: t,
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
  var eTag string
  if result == nil {
    var err error
    result, err = handler(request)
    if err != nil {
      t.handleError(err, w)
      return
    }

    t.cache.SetBig(cacheKey, result)
    eTag = computeEtag(result)
  } else {
    prevEtag := request.Header.Get("If-None-Match")
    eTag = computeEtag(result)
    if prevEtag == eTag {
      w.WriteHeader(304)
      return
    }
  }

  w.Header().Set("ETag", eTag)
  //w.Header().Set("Last-Modified", t.lastModified.Format(http.TimeFormat))

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
  case *HttpError:
    w.WriteHeader(exception.Code)
    writehttpError(w, exception)

  default:
    t.logger.Error("cannot handle http request", zap.Error(err))
    http.Error(w, err.Error(), 503)
  }
}

func computeEtag(result []byte) string {
  // add length to reduce chance of collision
  return strconv.FormatUint(xxhash.Sum64(result), 36) + "-" + strconv.FormatInt(int64(len(result)), 36)
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

func CopyBuffer(buffer *bytebufferpool.ByteBuffer) []byte {
  result := make([]byte, len(buffer.B))
  copy(result, buffer.B)
  return result
}
