package degradation_detector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDateLink(t *testing.T) {
	t.Parallel()
	mockNow := time.Date(2024, 12, 5, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name        string
		degradation Degradation
		expected    string
	}{
		{
			name: "One week range",
			degradation: Degradation{
				timestamp: mockNow.UnixMilli(),
			},
			expected: "timeRange=custom&customRange=2024-11-5:2024-12-5",
		},
		{
			name: "Different timestamp",
			degradation: Degradation{
				timestamp: mockNow.AddDate(0, 0, -2).UnixMilli(), // 2 days before mockNow
			},
			expected: "timeRange=custom&customRange=2024-11-3:2024-12-5",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := getCustomRange(tc.degradation.GetRangeStartTime(), mockNow)
			if result != tc.expected {
				t.Errorf("getCustomDateLinkBetweenDates() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// newSlackStub answers chat.postMessage with a 429 carrying retryAfter for the first rateLimitedCalls calls,
// then with finalBody.
func newSlackStub(t *testing.T, rateLimitedCalls int32, retryAfter string, finalBody string) (*slack.Client, *atomic.Int32) {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) <= rateLimitedCalls {
			w.Header().Set("Retry-After", retryAfter)
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(finalBody))
	}))
	t.Cleanup(server.Close)
	return slack.New("token", slack.OptionAPIURL(server.URL+"/")), &calls
}

const slackOk = `{"ok":true,"channel":"C1","ts":"1.0"}`

func TestSendSlackMessageRetriesWhenRateLimited(t *testing.T) {
	t.Parallel()
	api, calls := newSlackStub(t, 2, "0", slackOk)

	err := sendSlackMessage(t.Context(), api, SlackMessage{Channel: "C1", Text: "text"})

	require.NoError(t, err)
	assert.Equal(t, int32(3), calls.Load())
}

func TestSendSlackMessageGivesUpAfterMaxAttempts(t *testing.T) {
	t.Parallel()
	api, calls := newSlackStub(t, maxSlackAttempts+1, "0", slackOk)

	err := sendSlackMessage(t.Context(), api, SlackMessage{Channel: "C1", Text: "text"})

	var rateLimited *slack.RateLimitedError
	require.ErrorAs(t, err, &rateLimited)
	assert.Equal(t, int32(maxSlackAttempts), calls.Load())
}

func TestSendSlackMessageDoesNotRetryOtherErrors(t *testing.T) {
	t.Parallel()
	api, calls := newSlackStub(t, 0, "", `{"ok":false,"error":"channel_not_found"}`)

	err := sendSlackMessage(t.Context(), api, SlackMessage{Channel: "C1", Text: "text"})

	require.ErrorContains(t, err, "channel_not_found")
	assert.Equal(t, int32(1), calls.Load())
}

func TestSendSlackMessageDoesNotWaitPastDeadline(t *testing.T) {
	t.Parallel()
	api, calls := newSlackStub(t, 1, "30", slackOk)
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	start := time.Now()
	err := sendSlackMessage(ctx, api, SlackMessage{Channel: "C1", Text: "text"})

	var rateLimited *slack.RateLimitedError
	require.ErrorAs(t, err, &rateLimited)
	assert.Equal(t, int32(1), calls.Load())
	assert.Less(t, time.Since(start), time.Second)
}
