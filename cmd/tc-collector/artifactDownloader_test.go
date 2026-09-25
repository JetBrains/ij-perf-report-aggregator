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

const reportBody = `{"items":[]}`

// newFlakyTeamCity returns a collector for a TeamCity stub that answers the first failures requests with failureStatus.
func newFlakyTeamCity(t *testing.T, failures int32, failureStatus int) (*Collector, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if requests.Add(1) <= failures {
			w.WriteHeader(failureStatus)
			return
		}
		_, _ = w.Write([]byte(reportBody))
	}))
	collector := &Collector{serverUrl: srv.URL, httpClient: srv.Client(), logger: slog.New(slog.DiscardHandler)}
	return collector, &requests
}

func TestDownloadStartUpReportWithRetriesSucceedsFirstTime(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		collector, requests := newFlakyTeamCity(t, 0, 0)
		start := time.Now()

		data, err := collector.downloadStartUpReportWithRetries(t.Context(), Build{Id: 1}, collector.serverUrl+"/report.json")
		require.NoError(t, err)
		assert.JSONEq(t, reportBody, string(data))
		assert.Equal(t, int32(1), requests.Load())
		assert.Zero(t, time.Since(start))
	})
}

func TestDownloadStartUpReportWithRetriesRecoversFromTransientErrors(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusInternalServerError, http.StatusBadGateway, http.StatusRequestTimeout, http.StatusTooManyRequests} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				collector, requests := newFlakyTeamCity(t, 3, status)

				data, err := collector.downloadStartUpReportWithRetries(t.Context(), Build{Id: 1}, collector.serverUrl+"/report.json")
				require.NoError(t, err)
				assert.JSONEq(t, reportBody, string(data))
				assert.Equal(t, int32(4), requests.Load())
			})
		})
	}
}

func TestDownloadStartUpReportWithRetriesGivesUpWithinMaxElapsedTime(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		collector, requests := newFlakyTeamCity(t, 1000, http.StatusInternalServerError)
		start := time.Now()

		_, err := collector.downloadStartUpReportWithRetries(t.Context(), Build{Id: 1}, collector.serverUrl+"/report.json")
		require.ErrorContains(t, err, "500 Internal Server Error")
		assert.Greater(t, requests.Load(), int32(1))
		// 15s max elapsed time plus at most one 5s max interval
		assert.LessOrEqual(t, time.Since(start), 20*time.Second)
	})
}

func TestDownloadStartUpReportWithRetriesSkipsClientErrorsWithoutRetry(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			t.Parallel()
			synctest.Test(t, func(t *testing.T) {
				collector, requests := newFlakyTeamCity(t, 1000, status)
				start := time.Now()

				_, err := collector.downloadStartUpReportWithRetries(t.Context(), Build{Id: 1, Status: "SUCCESS"}, collector.serverUrl+"/report.json")
				require.ErrorContains(t, err, http.StatusText(status))
				assert.Equal(t, int32(1), requests.Load())
				assert.Zero(t, time.Since(start))
			})
		})
	}
}

func TestMissingReportOfFailedBuildIsSkipped(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		collector, requests := newFlakyTeamCity(t, 1000, http.StatusNotFound)
		build := Build{Id: 1, Status: "FAILURE", Artifacts: Artifacts{File: []Artifact{{Url: "/app/rest/builds/id:1/artifacts/metadata/startup-stats.json"}}}}
		start := time.Now()

		reports, err := collector.downloadReports(t.Context(), build)
		require.NoError(t, err)
		assert.Empty(t, reports)
		assert.Equal(t, int32(1), requests.Load())
		assert.Zero(t, time.Since(start))
	})
}
