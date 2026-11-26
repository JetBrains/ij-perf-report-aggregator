package degradation_detector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"sync"
	"testing"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/meta"
	"github.com/stretchr/testify/assert"
)

// missingDataBackend is the backend that stores a missing metric once and answers 409 to the next request for it.
type missingDataBackend struct {
	mu       sync.Mutex
	known    map[meta.MissingDataInsertParams]bool
	requests int
}

func newMissingDataBackend(t *testing.T, known ...meta.MissingDataInsertParams) (*httptest.Server, *missingDataBackend) {
	t.Helper()
	backend := &missingDataBackend{known: map[meta.MissingDataInsertParams]bool{}}
	for _, params := range known {
		backend.known[params] = true
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params meta.MissingDataInsertParams
		if !assert.Equal(t, "/api/meta/missingData", r.URL.Path) || !assert.NoError(t, json.NewDecoder(r.Body).Decode(&params)) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		backend.mu.Lock()
		defer backend.mu.Unlock()
		backend.requests++
		if backend.known[params] {
			w.WriteHeader(http.StatusConflict)
			return
		}
		backend.known[params] = true
		_, _ = w.Write([]byte("1"))
	}))
	t.Cleanup(server.Close)
	return server, backend
}

func missingDataFor(test string, channel string) MissingData {
	settings := PerformanceSettings{Project: test, Metric: "attempt.mean.ms", Channel: channel}
	return MissingData{LastBuild: "1", LastTimestamp: 1_700_000_000_000, TCBuildType: "ijplatform_master_UnitPerfTests", Settings: settings}
}

// postMissingData returns the channels that get the missing metric of each test, in a stable order.
func postMissingData(server *httptest.Server, data ...MissingData) map[string][]string {
	input := make(chan MissingData, len(data))
	for _, datum := range data {
		input <- datum
	}
	close(input)
	channels := map[string][]string{}
	for datum := range PostMissingData(server.Client(), server.URL, input) {
		test := datum.Settings.GetProject()
		channels[test] = append(channels[test], datum.Settings.SlackChannel())
	}
	for _, list := range channels {
		slices.Sort(list)
	}
	return channels
}

func TestMissingDataReachesEachChannelOfTheSeries(t *testing.T) {
	t.Parallel()
	server, backend := newMissingDataBackend(t)

	channels := postMissingData(server,
		missingDataFor("com.intellij.java.Test", "idea-java-alerts"),
		missingDataFor("com.intellij.java.Test", "ij-perf-unit-tests-alerts"),
		missingDataFor("org.other.Test", "ij-perf-unit-tests-alerts"),
	)

	assert.Equal(t, map[string][]string{
		"com.intellij.java.Test": {"idea-java-alerts", "ij-perf-unit-tests-alerts"},
		"org.other.Test":         {"ij-perf-unit-tests-alerts"},
	}, channels)
	assert.Equal(t, 2, backend.requests, "one request per series")
}

func TestMissingDataReachesAChannelOnce(t *testing.T) {
	t.Parallel()
	server, _ := newMissingDataBackend(t)

	channels := postMissingData(server,
		missingDataFor("com.intellij.java.Test", "idea-java-alerts"),
		missingDataFor("com.intellij.java.Test", "idea-java-alerts"),
	)

	assert.Equal(t, map[string][]string{"com.intellij.java.Test": {"idea-java-alerts"}}, channels)
}

func TestMissingDataThatTheBackendKnowsIsNotReportedAgain(t *testing.T) {
	t.Parallel()
	known := missingDataFor("com.intellij.java.Test", "")
	server, backend := newMissingDataBackend(t, meta.MissingDataInsertParams{
		BuildType: known.TCBuildType, Project: "com.intellij.java.Test", Metric: "attempt.mean.ms", MissingSince: known.LastTimestamp,
	})

	channels := postMissingData(server,
		missingDataFor("com.intellij.java.Test", "idea-java-alerts"),
		missingDataFor("com.intellij.java.Test", "ij-perf-unit-tests-alerts"),
	)

	assert.Empty(t, channels)
	assert.Equal(t, 1, backend.requests)
}
