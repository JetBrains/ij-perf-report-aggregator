package degradation_detector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractDataFromRequest(t *testing.T) {
	t.Parallel()
	data, err := extractDataFromRequest([]byte(`[[[1000,2000],[10,20],["b1","b2"],["bt1","bt2"]]]`))
	require.NoError(t, err)
	assert.Equal(t, []int64{1000, 2000}, data.timestamps)
	assert.Equal(t, []int{10, 20}, data.values)
	assert.Equal(t, []string{"b1", "b2"}, data.builds)
	assert.Equal(t, []string{"bt1", "bt2"}, data.buildTypes)
}

func TestExtractDataFromRequestWithMissingColumn(t *testing.T) {
	t.Parallel()
	_, err := extractDataFromRequest([]byte(`[[[1000,2000],[10,20],["b1","b2"]]]`))
	require.Error(t, err)
}

func TestFetchMetricsFromClickhouseSkipsFailedQueries(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoded, err := util.DecodeQuery(strings.TrimPrefix(r.URL.Path, "/api/q/"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var queries []dataQuery.Query
		err = json.Unmarshal(decoded, &queries)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, filter := range queries[0].Filters {
			if filter.Field == "measures.name" && filter.Value == "broken" {
				http.Error(w, "boom", http.StatusInternalServerError)
				return
			}
		}
		_, _ = w.Write([]byte(`[[[1000,2000],[10,20],["b1","b2"],["bt1","bt2"]]]`))
	}))
	t.Cleanup(server.Close)

	settings := []Settings{
		PerformanceSettings{Metric: "broken"},
		PerformanceSettings{Metric: "ok"},
	}
	results := collectWithTimeout(t, FetchMetricsFromClickhouse(settings, server.Client(), server.URL))

	require.Len(t, results, 1)
	assert.Equal(t, settings[1], results[0].Settings)
	assert.Equal(t, []int{10, 20}, results[0].values)
}

func collectWithTimeout[T any](t *testing.T, ch <-chan T) []T {
	t.Helper()
	var result []T
	timeout := time.After(10 * time.Second)
	for {
		select {
		case item, ok := <-ch:
			if !ok {
				return result
			}
			result = append(result, item)
		case <-timeout:
			t.Fatal("channel was not closed")
			return nil
		}
	}
}
