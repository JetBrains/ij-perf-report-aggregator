package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newFlakyTeamCity returns a collector for a TeamCity stub that fails the first failures requests with 500.
func newFlakyTeamCity(t *testing.T, failures int32) (*Collector, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if requests.Add(1) <= failures {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`{"items":[]}`))
	}))
	collector := &Collector{serverUrl: srv.URL, httpClient: srv.Client(), logger: slog.New(slog.DiscardHandler)}
	return collector, &requests
}

func TestDownloadStartUpReportWithRetriesSucceedsFirstTime(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		collector, requests := newFlakyTeamCity(t, 0)
		start := time.Now()

		data, err := collector.downloadStartUpReportWithRetries(t.Context(), Build{Id: 1}, collector.serverUrl+"/report.json")
		require.NoError(t, err)
		assert.JSONEq(t, `{"items":[]}`, string(data))
		assert.Equal(t, int32(1), requests.Load())
		assert.Zero(t, time.Since(start))
	})
}

func TestDownloadStartUpReportWithRetriesRecoversFromServerErrors(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		collector, requests := newFlakyTeamCity(t, 3)

		data, err := collector.downloadStartUpReportWithRetries(t.Context(), Build{Id: 1}, collector.serverUrl+"/report.json")
		require.NoError(t, err)
		assert.JSONEq(t, `{"items":[]}`, string(data))
		assert.Equal(t, int32(4), requests.Load())
	})
}

func TestDownloadStartUpReportWithRetriesGivesUpWithinMaxElapsedTime(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		collector, requests := newFlakyTeamCity(t, 1000)
		start := time.Now()

		_, err := collector.downloadStartUpReportWithRetries(t.Context(), Build{Id: 1}, collector.serverUrl+"/report.json")
		require.ErrorContains(t, err, "maximum retries reached")
		assert.Greater(t, requests.Load(), int32(1))
		// 15s max elapsed time plus at most one 5s max interval
		assert.LessOrEqual(t, time.Since(start), 20*time.Second)
	})
}
