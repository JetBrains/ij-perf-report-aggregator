package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/bytebufferpool"
)

func TestCompressData(t *testing.T) {
	t.Parallel()
	rcm := NewResponseCacheManager(nil)

	testData := []byte("sample data to compress")

	compressedData, err := rcm.compressData(testData)
	require.NoError(t, err)
	assert.NotEmptyf(t, compressedData, "Expected compressed data, got empty")
}

// poolScribblingWriter stands in for concurrent requests: between the headers and the body it takes
// buffers from both pools and overwrites them, as generateCacheKey does with the request path. It
// takes several from each, because the scratch buffers released during compression come out first.
type poolScribblingWriter struct {
	*httptest.ResponseRecorder
}

func (w poolScribblingWriter) WriteHeader(code int) {
	w.ResponseRecorder.WriteHeader(code)
	for _, get := range []func() *bytebufferpool.ByteBuffer{bytebufferpool.Get, byteBufferPool.Get} {
		for range 4 {
			_, _ = get().WriteString("/api/q/other-request")
		}
	}
}

func TestResponseBufferIsNotReleasedBeforeWrite(t *testing.T) {
	t.Parallel()
	rcm := NewResponseCacheManager(NewPrometheusMetrics())
	const body = `[[["response body of this request"]]]`
	handler := rcm.CreateHandler(func(_ *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
		buffer := byteBufferPool.Get()
		_, _ = buffer.WriteString(body)
		return buffer, true, nil
	})

	for range 10 {
		recorder := httptest.NewRecorder()
		// no-cache forces the handler path, where the pooled buffer is written out
		request := httptest.NewRequest(http.MethodGet, "/api/q/this-request", http.NoBody)
		request.Header.Set("Cache-Control", "no-cache")
		handler.ServeHTTP(poolScribblingWriter{recorder}, request)
		require.Equal(t, body, recorder.Body.String())
	}
}
