package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/model"

	"github.com/go-chi/chi/v5"
	"github.com/valyala/bytebufferpool"
)

func TestPrometheusMetricsEndpointExposesHTTPMetrics(t *testing.T) {
	t.Parallel()

	metrics := NewPrometheusMetrics()

	router := chi.NewRouter()
	router.Use(metrics.Middleware)
	router.Handle("/metrics", metrics.Handler())
	router.Get("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	testResponse := httptest.NewRecorder()
	router.ServeHTTP(testResponse, httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	if testResponse.Code != http.StatusNoContent {
		t.Fatalf("unexpected test route status: got %d, want %d", testResponse.Code, http.StatusNoContent)
	}

	metricFamilies := scrapeMetrics(t, router)

	requestCount := getCounterValue(t, metricFamilies["ij_perf_http_requests_total"], map[string]string{
		"method":      http.MethodGet,
		"route":       "/test",
		"status_code": "204",
	})
	if requestCount != 1 {
		t.Fatalf("unexpected request counter value: got %v, want 1", requestCount)
	}

	if _, ok := metricFamilies["go_goroutines"]; !ok {
		t.Fatal("go_goroutines metric was not exported")
	}

	if _, ok := metricFamilies["ij_perf_http_response_size_bytes"]; ok {
		t.Fatal("ij_perf_http_response_size_bytes metric should not be exported")
	}
}

func TestPrometheusMetricsEndpointExposesCacheMetrics(t *testing.T) {
	t.Parallel()

	metrics := NewPrometheusMetrics()

	cacheManager := NewResponseCacheManager(metrics)

	router := chi.NewRouter()
	router.Use(metrics.Middleware)
	router.Handle("/metrics", metrics.Handler())
	router.Handle("/cached", cacheManager.CreateHandler(func(_ *http.Request) (*bytebufferpool.ByteBuffer, bool, error) {
		return &bytebufferpool.ByteBuffer{B: []byte("payload")}, false, nil
	}))

	firstResponse := httptest.NewRecorder()
	router.ServeHTTP(firstResponse, httptest.NewRequest(http.MethodGet, "/cached", http.NoBody))
	if firstResponse.Code != http.StatusOK {
		t.Fatalf("unexpected first cached response status: got %d, want %d", firstResponse.Code, http.StatusOK)
	}

	secondResponse := httptest.NewRecorder()
	router.ServeHTTP(secondResponse, httptest.NewRequest(http.MethodGet, "/cached", http.NoBody))
	if secondResponse.Code != http.StatusOK {
		t.Fatalf("unexpected second cached response status: got %d, want %d", secondResponse.Code, http.StatusOK)
	}

	bypassRequest := httptest.NewRequest(http.MethodGet, "/cached", http.NoBody)
	bypassRequest.Header.Set("Cache-Control", "no-cache")
	bypassResponse := httptest.NewRecorder()
	router.ServeHTTP(bypassResponse, bypassRequest)
	if bypassResponse.Code != http.StatusOK {
		t.Fatalf("unexpected bypass cached response status: got %d, want %d", bypassResponse.Code, http.StatusOK)
	}

	cacheManager.Clear()

	metricFamilies := scrapeMetrics(t, router)

	missCount := getCounterValue(t, metricFamilies["ij_perf_response_cache_requests_total"], map[string]string{"result": "miss"})
	if missCount != 1 {
		t.Fatalf("unexpected cache miss counter value: got %v, want 1", missCount)
	}

	hitCount := getCounterValue(t, metricFamilies["ij_perf_response_cache_requests_total"], map[string]string{"result": "hit"})
	if hitCount != 1 {
		t.Fatalf("unexpected cache hit counter value: got %v, want 1", hitCount)
	}

	bypassCount := getCounterValue(t, metricFamilies["ij_perf_response_cache_requests_total"], map[string]string{"result": "bypass"})
	if bypassCount != 1 {
		t.Fatalf("unexpected cache bypass counter value: got %v, want 1", bypassCount)
	}

	clearCount := getCounterValue(t, metricFamilies["ij_perf_response_cache_clears_total"], map[string]string{})
	if clearCount != 1 {
		t.Fatalf("unexpected cache clear counter value: got %v, want 1", clearCount)
	}
}

func TestPrometheusMetricsEndpointCollapsesWildcardRoutes(t *testing.T) {
	t.Parallel()

	metrics := NewPrometheusMetrics()

	router := chi.NewRouter()
	router.Use(metrics.Middleware)
	router.Handle("/metrics", metrics.Handler())
	router.Handle("/api/q/*", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	requestPath := "/api/q/" + strings.Repeat("abc123", 40)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, requestPath, http.NoBody))

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected wildcard route status: got %d, want %d", response.Code, http.StatusOK)
	}

	metricFamilies := scrapeMetrics(t, router)

	requestCount := getCounterValue(t, metricFamilies["ij_perf_http_requests_total"], map[string]string{
		"method":      http.MethodGet,
		"route":       "/api/q",
		"status_code": "200",
	})
	if requestCount != 1 {
		t.Fatalf("unexpected wildcard route counter value: got %v, want 1", requestCount)
	}

	metricsText := metricsBody(t, router)
	if strings.Contains(metricsText, requestPath) {
		t.Fatalf("metrics output leaked raw request path %q", requestPath)
	}
}

func TestPrometheusMetricsDebouncesUserRequests(t *testing.T) {
	t.Parallel()

	metrics := NewPrometheusMetrics()
	fakeNow := time.Now()
	metrics.now = func() time.Time { return fakeNow }

	router := chi.NewRouter()
	router.Use(metrics.Middleware)
	router.Handle("/metrics", metrics.Handler())
	router.Get("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	send := func() {
		req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
		req.Header.Set("X-Auth-Request-Email", "Alice@example.com")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	for range 5 {
		send()
	}

	userLabels := map[string]string{"user": "alice"}
	if got := getCounterValue(t, scrapeMetrics(t, router)["ij_perf_http_user_requests_total"], userLabels); got != 1 {
		t.Fatalf("debounced user counter: got %v, want 1", got)
	}

	fakeNow = fakeNow.Add(2 * metrics.userSeenWindow)
	send()

	if got := getCounterValue(t, scrapeMetrics(t, router)["ij_perf_http_user_requests_total"], userLabels); got != 2 {
		t.Fatalf("user counter after window elapsed: got %v, want 2", got)
	}
}

func TestPrometheusMetricsSkipsAnonymousUserRequests(t *testing.T) {
	t.Parallel()

	metrics := NewPrometheusMetrics()

	router := chi.NewRouter()
	router.Use(metrics.Middleware)
	router.Handle("/metrics", metrics.Handler())
	router.Get("/test", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/test", http.NoBody))

	metricsText := metricsBody(t, router)
	if strings.Contains(metricsText, "ij_perf_http_user_requests_total") {
		t.Fatalf("anonymous request produced ij_perf_http_user_requests_total series:\n%s", metricsText)
	}
}

func TestPrometheusMetricsLabelsUnmatchedRoutes(t *testing.T) {
	t.Parallel()

	metrics := NewPrometheusMetrics()

	router := chi.NewRouter()
	router.Use(metrics.Middleware)
	router.Handle("/metrics", metrics.Handler())

	scannerPath := "/scanner/probe/" + strings.Repeat("xyz", 40)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, scannerPath, http.NoBody))

	if response.Code != http.StatusNotFound {
		t.Fatalf("unexpected unmatched route status: got %d, want %d", response.Code, http.StatusNotFound)
	}

	metricFamilies := scrapeMetrics(t, router)

	requestCount := getCounterValue(t, metricFamilies["ij_perf_http_requests_total"], map[string]string{
		"method":      http.MethodPost,
		"route":       "unmatched",
		"status_code": "404",
	})
	if requestCount != 1 {
		t.Fatalf("unexpected unmatched route counter value: got %v, want 1", requestCount)
	}

	metricsText := metricsBody(t, router)
	if strings.Contains(metricsText, scannerPath) {
		t.Fatalf("metrics output leaked raw scanner path %q", scannerPath)
	}
}

func scrapeMetrics(t *testing.T, handler http.Handler) map[string]*dto.MetricFamily {
	t.Helper()

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected metrics endpoint status: got %d, want %d", response.Code, http.StatusOK)
	}

	parser := expfmt.NewTextParser(model.UTF8Validation)
	metricFamilies, err := parser.TextToMetricFamilies(strings.NewReader(response.Body.String()))
	if err != nil {
		t.Fatalf("failed to parse metrics response: %v", err)
	}

	return metricFamilies
}

func metricsBody(t *testing.T, handler http.Handler) string {
	t.Helper()

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/metrics", http.NoBody))

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected metrics endpoint status: got %d, want %d", response.Code, http.StatusOK)
	}

	return response.Body.String()
}

func getCounterValue(t *testing.T, family *dto.MetricFamily, expectedLabels map[string]string) float64 {
	t.Helper()

	if family == nil {
		t.Fatal("metric family was not exported")
	}

	for _, metric := range family.GetMetric() {
		if hasExpectedLabels(metric, expectedLabels) {
			if metric.GetCounter() == nil {
				t.Fatalf("metric %s is not a counter", family.GetName())
			}
			return metric.GetCounter().GetValue()
		}
	}

	t.Fatalf("metric %s with labels %v was not exported", family.GetName(), expectedLabels)
	return 0
}

func hasExpectedLabels(metric *dto.Metric, expectedLabels map[string]string) bool {
	for key, value := range expectedLabels {
		found := false
		for _, label := range metric.GetLabel() {
			if label.GetName() == key && label.GetValue() == value {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
