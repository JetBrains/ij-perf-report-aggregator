package meta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newFlakyYoutrack returns a client for a YouTrack stub that fails the first failures requests with 500.
func newFlakyYoutrack(t *testing.T, failures int32) (*YoutrackClient, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if requests.Add(1) <= failures {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`{"id":"1-1","idReadable":"X-1"}`))
	}))
	return &YoutrackClient{youTrackUrl: srv.URL, httpClient: srv.Client()}, &requests
}

func TestWaitIssueIsCreatedRetriesUntilIssueExists(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		client, requests := newFlakyYoutrack(t, 2)
		start := time.Now()

		require.NoError(t, client.waitIssueIsCreated(t.Context(), "X-1"))
		assert.Equal(t, int32(3), requests.Load())
		assert.Equal(t, 6*time.Second, time.Since(start))
	})
}

func TestWaitIssueIsCreatedGivesUpAfterFiveAttempts(t *testing.T) {
	t.Parallel()
	synctest.Test(t, func(t *testing.T) {
		client, requests := newFlakyYoutrack(t, 100)

		require.ErrorContains(t, client.waitIssueIsCreated(context.Background(), "X-1"), "after 5 retries")
		assert.Equal(t, int32(5), requests.Load())
	})
}
