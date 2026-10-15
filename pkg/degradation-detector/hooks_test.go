package degradation_detector

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func hookTestSettings(project string) PerformanceSettings {
	return PerformanceSettings{
		Db:          "perfintDev",
		Table:       "air",
		Project:     project,
		Machine:     "intellij-linux-performance-aws-%",
		Metric:      "open.XL.coverRemovedMs",
		Branch:      "master",
		Channel:     "ij-air-perf",
		ProductLink: "air",
	}
}

func hookTestDegradation() Degradation {
	return Degradation{
		Build:         "1090164265",
		timestamp:     time.Date(2026, 10, 14, 3, 0, 0, 0, time.UTC).UnixMilli(),
		medianValues:  CenterValues{previousValue: 3000, newValue: 4500},
		IsDegradation: true,
	}
}

func hookTestHook(urlEnv string) Hook {
	return Hook{Name: "test", Channel: "ij-air-perf", URLEnv: urlEnv, TokenEnv: urlEnv + "_TOKEN", AuthScheme: "ApiKey"}
}

type hookStub struct {
	mu       sync.Mutex
	requests []*http.Request
	bodies   [][]byte
}

// newHookStub answers with failStatus for the first failedCalls calls, then with 200.
func newHookStub(t *testing.T, failedCalls int, failStatus int) (*httptest.Server, *hookStub) {
	t.Helper()
	stub := &hookStub{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		stub.mu.Lock()
		stub.requests = append(stub.requests, r)
		stub.bodies = append(stub.bodies, body)
		calls := len(stub.requests)
		stub.mu.Unlock()
		if calls <= failedCalls {
			w.WriteHeader(failStatus)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server, stub
}

func (s *hookStub) count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.requests)
}

func TestHookMatchesTheChannelOfTheTeam(t *testing.T) {
	t.Parallel()
	hook := hookTestHook("X")
	assert.True(t, hook.matches(hookTestSettings("air/long-session/XL")))
	assert.True(t, hook.matches(hookTestSettings("other/project")), "the channel alone selects")
	withHash := hookTestSettings("air/long-session/XL")
	withHash.Channel = "#ij-air-perf"
	assert.True(t, hook.matches(withHash), "a channel with #")
	other := hookTestSettings("air/long-session/XL")
	other.Channel = "ij-idea-perf"
	assert.False(t, hook.matches(other), "another team")
}

func TestHookFiltersNarrowTheChannel(t *testing.T) {
	t.Parallel()
	hook := hookTestHook("X")
	hook.Db, hook.Table, hook.ProjectPrefix, hook.Branch, hook.MetricPrefix = "perfintDev", "air", "air/", "master", "open."
	assert.True(t, hook.matches(hookTestSettings("air/long-session/XL")))
	assert.True(t, hook.matches(hookTestSettings("air/long-session/M,air/long-session/XL")), "a merged degradation")
	for name, change := range map[string]func(*PerformanceSettings){
		"another database": func(s *PerformanceSettings) { s.Db = "perfint" },
		"another table":    func(s *PerformanceSettings) { s.Table = "idea" },
		"another project":  func(s *PerformanceSettings) { s.Project = "idea/indexing" },
		"another branch":   func(s *PerformanceSettings) { s.Branch = "263" },
		"another metric":   func(s *PerformanceSettings) { s.Metric = "heap.XL.heapLive.MiB" },
	} {
		settings := hookTestSettings("air/long-session/XL")
		change(&settings)
		assert.False(t, hook.matches(settings), name)
	}
}

func TestHookAuthorization(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "ApiKey t", Hook{AuthScheme: "ApiKey"}.authorization("t"))
	assert.Equal(t, "Bearer t", Hook{}.authorization("t"))
}

func TestNewHookEvent(t *testing.T) {
	t.Parallel()
	event := newHookEvent(hookTestSettings("air/long-session/XL"), hookTestDegradation())

	assert.Equal(t, "air/long-session/XL/open.XL.coverRemovedMs:1090164265", event.EventID)
	assert.Equal(t, "degradation", event.Kind)
	assert.Equal(t, "ij-perf", event.Source)
	assert.Equal(t, "air", event.Table)
	assert.Equal(t, "default", event.Mode)
	assert.Equal(t, "2026-10-14T03:00:00Z", event.Date)
	assert.InDelta(t, 3000, event.PreviousValue, 0)
	assert.InDelta(t, 4500, event.NewValue, 0)
	assert.InDelta(t, 50, event.ChangePercent, 0.001)
	assert.Contains(t, event.ChartURL, "https://ij-perf.labs.jb.gg/air/testsDev?")
	assert.Contains(t, event.ChartURL, "measure=open.XL.coverRemovedMs")
	assert.Contains(t, event.ChartURL, "point=1090164265")
	assert.NotContains(t, event.ChartURL, "|")
	assert.Contains(t, event.ReportURL, "https://ij-perf.labs.jb.gg/degradations/report?")
}

func TestSlackLinkURL(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "https://example.com/a?b=c", slackLinkURL("<https://example.com/a?b=c|See charts>"))
	assert.Equal(t, "https://example.com", slackLinkURL("<https://example.com>"))
	assert.Equal(t, "plain text", slackLinkURL("plain text"))
}

func TestSendHookEventPostsTheEventWithTheAuthorization(t *testing.T) {
	t.Parallel()
	server, stub := newHookStub(t, 0, 0)
	event := newHookEvent(hookTestSettings("air/long-session/XL"), hookTestDegradation())

	require.NoError(t, sendHookEvent(context.Background(), server.Client(), server.URL, "ApiKey secret", event, time.Millisecond))

	require.Equal(t, 1, stub.count())
	assert.Equal(t, http.MethodPost, stub.requests[0].Method)
	assert.Equal(t, "ApiKey secret", stub.requests[0].Header.Get("Authorization"))
	assert.Equal(t, "application/json", stub.requests[0].Header.Get("Content-Type"))
	var sent HookEvent
	require.NoError(t, json.Unmarshal(stub.bodies[0], &sent))
	assert.Equal(t, event, sent)
}

func TestSendHookEventRetriesAServerErrorAndRateLimit(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusBadGateway, http.StatusTooManyRequests} {
		server, stub := newHookStub(t, 2, status)
		require.NoError(t, sendHookEvent(context.Background(), server.Client(), server.URL, "ApiKey s", HookEvent{EventID: "e"}, time.Millisecond))
		assert.Equal(t, 3, stub.count(), "status %d", status)
	}
	server, stub := newHookStub(t, 100, http.StatusBadGateway)
	require.Error(t, sendHookEvent(context.Background(), server.Client(), server.URL, "ApiKey s", HookEvent{EventID: "e"}, time.Millisecond))
	assert.Equal(t, maxHookAttempts, stub.count())
}

func TestSendHookEventDoesNotRetryAClientError(t *testing.T) {
	t.Parallel()
	server, stub := newHookStub(t, 100, http.StatusUnauthorized)

	require.Error(t, sendHookEvent(context.Background(), server.Client(), server.URL, "ApiKey wrong", HookEvent{EventID: "e"}, time.Millisecond))

	assert.Equal(t, 1, stub.count())
}

func TestSendDegradationsToHooksSendsOnlyMatchingAndConfigured(t *testing.T) {
	server, stub := newHookStub(t, 0, 0)
	t.Setenv("TEST_HOOK_URL", server.URL)
	t.Setenv("TEST_HOOK_URL_TOKEN", "secret")
	hooks := []Hook{
		hookTestHook("TEST_HOOK_URL"),
		// Its variables are not set, so it sends nothing.
		hookTestHook("TEST_HOOK_UNSET"),
	}

	degradations := make(chan DegradationWithSettings, 3)
	degradations <- DegradationWithSettings{Details: hookTestDegradation(), Settings: hookTestSettings("air/long-session/XL")}
	otherTeam := hookTestSettings("idea/indexing")
	otherTeam.Channel = "ij-idea-perf"
	degradations <- DegradationWithSettings{Details: hookTestDegradation(), Settings: otherTeam}
	degradations <- DegradationWithSettings{Details: hookTestDegradation(), Settings: FleetStartupSettings{}}
	close(degradations)

	SendDegradationsToHooks(degradations, server.Client(), hooks)

	require.Equal(t, 1, stub.count())
	assert.Equal(t, "ApiKey secret", stub.requests[0].Header.Get("Authorization"))
}

func TestSendDegradationsToHooksWithoutHooksDrainsTheChannel(t *testing.T) {
	t.Parallel()
	degradations := make(chan DegradationWithSettings, 1)
	degradations <- DegradationWithSettings{Details: hookTestDegradation(), Settings: hookTestSettings("air/long-session/XL")}
	close(degradations)

	SendDegradationsToHooks(degradations, http.DefaultClient, nil)

	_, open := <-degradations
	assert.False(t, open)
}

// A nightly series with a step of 50 % goes through detection and merging, as in the analyzer, and reaches the hook.
// The same series without the step sends nothing.
func TestADegradationOfANightlySeriesReachesTheHook(t *testing.T) {
	server, stub := newHookStub(t, 0, 0)
	t.Setenv("TEST_E2E_HOOK_URL", server.URL)
	t.Setenv("TEST_E2E_HOOK_URL_TOKEN", "secret")
	hooks := []Hook{hookTestHook("TEST_E2E_HOOK_URL")}

	series := func(step float64) QueryResultWithSettings {
		nights := 20
		result := queryResult{}
		for i := range nights {
			value := 3000.0
			if i >= nights-6 {
				value *= 1 + step
			}
			// A small, fixed noise, so the test does not depend on a random seed.
			result.values = append(result.values, int(value)+[]int{-40, 25, -10, 35, -20, 5}[i%6])
			result.builds = append(result.builds, strconv.Itoa(1000+i))
			result.timestamps = append(result.timestamps, time.Date(2026, 10, 1+i, 3, 0, 0, 0, time.UTC).UnixMilli())
		}
		settings := hookTestSettings("air/long-session/XL")
		settings.AnalysisSettings = AnalysisSettings{ReportType: DegradationEvent, MedianDifferenceThreshold: 15, EffectSizeThreshold: 2}
		return QueryResultWithSettings{queryResult: result, Settings: settings}
	}

	for _, c := range []struct {
		step     float64
		expected int
	}{{0, 0}, {0.5, 1}} {
		data := make(chan QueryResultWithSettings, 1)
		data <- series(c.step)
		close(data)
		before := stub.count()

		SendDegradationsToHooks(MergeDegradations(InferDegradations(data)), server.Client(), hooks)

		require.Equal(t, c.expected, stub.count()-before, "step %v", c.step)
	}
	var event HookEvent
	require.NoError(t, json.Unmarshal(stub.bodies[0], &event))
	assert.Equal(t, "degradation", event.Kind)
	assert.Equal(t, "1014", event.Build, "the first build after the step")
	assert.InDelta(t, 50, event.ChangePercent, 3)
}
