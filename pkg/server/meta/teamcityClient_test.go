package meta

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// gapStub describes the canned TeamCity responses for one getChangesGap scenario.
type gapStub struct {
	prevRev        string           // revision the previous dot's build was built on ("" => build has no revision)
	branch         string           // current build's branch
	changeIDs      map[string]int64 // revision (version) -> change id, as getChangeID resolves them
	betweenChanges []changeRef      // changes attributed to the builds between the two dots
	requested      *bool            // set true when any TeamCity endpoint is hit
}

// versionFromLocator extracts the value of the version: dimension from a changes locator.
func versionFromLocator(locator string) string {
	_, after, found := strings.Cut(locator, "version:")
	if !found {
		return ""
	}
	version, _, _ := strings.Cut(after, ",")
	return version
}

// newGapTestClient returns a TeamCityClient wired to an httptest server that answers the
// endpoints getChangesGap calls, according to stub.
func newGapTestClient(t *testing.T, stub gapStub) *TeamCityClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if stub.requested != nil {
			*stub.requested = true
		}
		w.Header().Set("Content-Type", "application/json")
		locator := r.URL.Query().Get("locator")
		switch {
		case r.URL.Path == "/app/rest/builds":
			// getChangesBetweenBuilds: one build carrying every change between the two dots.
			changes := make([]map[string]any, 0, len(stub.betweenChanges))
			for _, c := range stub.betweenChanges {
				changes = append(changes, map[string]any{"id": c.Id, "version": c.Version})
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"build": []map[string]any{{"changes": map[string]any{"change": changes}}},
			})
		case strings.HasPrefix(r.URL.Path, "/app/rest/builds/"):
			// getBuildRevision asks for the revisions field; getBuildInfo does not.
			if strings.Contains(r.URL.RawQuery, "revisions") {
				rev := []map[string]string{}
				if stub.prevRev != "" {
					rev = append(rev, map[string]string{"version": stub.prevRev})
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"revisions": map[string]any{"revision": rev}})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]any{"buildTypeId": "BT", "branchName": stub.branch})
			}
		case strings.HasPrefix(r.URL.Path, "/app/rest/changes"):
			// getChangeID: resolve a revision (version:) to its change id.
			change := []map[string]int64{}
			if id, ok := stub.changeIDs[versionFromLocator(locator)]; ok && id != 0 {
				change = append(change, map[string]int64{"id": id})
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"change": change})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return &TeamCityClient{teamCityURL: srv.URL, authToken: "test-token"}
}

func TestGetChangesGap(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name              string
		currentFirst      string
		stub              gapStub
		wantKnown         bool
		wantHasGap        bool
		wantGapCount      int
		wantFirstAfterDot string
	}{
		{
			name:         "gap: commits landed between the two dots (out-of-range and duplicate ids ignored)",
			currentFirst: "cur",
			stub: gapStub{
				prevRev:   "p",
				branch:    "master",
				changeIDs: map[string]int64{"p": 100, "cur": 200},
				// 150,160,170 fall in (100,200); 90 precedes the previous dot, 250 is the current
				// build's own commit, and the repeated 160 is a duplicate — all excluded. The oldest
				// in range (150) is the first commit after the previous dot.
				betweenChanges: []changeRef{
					{Id: 170, Version: "v170"},
					{Id: 160, Version: "v160"},
					{Id: 150, Version: "v150"},
					{Id: 90, Version: "v90"},
					{Id: 250, Version: "v250"},
					{Id: 160, Version: "v160"},
				},
			},
			wantKnown:         true,
			wantHasGap:        true,
			wantGapCount:      3,
			wantFirstAfterDot: "v150",
		},
		{
			name:         "no gap: previous dot is the immediate predecessor (no builds between)",
			currentFirst: "cur",
			stub:         gapStub{prevRev: "p", branch: "master", changeIDs: map[string]int64{"p": 100, "cur": 110}, betweenChanges: nil},
			wantKnown:    true,
			wantHasGap:   false,
			wantGapCount: 0,
		},
		{
			name:         "no gap: current range starts at or before the previous dot (out-of-order builds)",
			currentFirst: "cur",
			stub:         gapStub{prevRev: "p", branch: "master", changeIDs: map[string]int64{"p": 200, "cur": 100}, betweenChanges: []changeRef{{Id: 150, Version: "v150"}}},
			wantKnown:    true,
			wantHasGap:   false,
			wantGapCount: 0,
		},
		{
			name:         "unknown: no current range start",
			currentFirst: "",
			stub:         gapStub{prevRev: "p", changeIDs: map[string]int64{"p": 100}},
			wantKnown:    false,
		},
		{
			name:         "unknown: previous dot's build has no revision",
			currentFirst: "cur",
			stub:         gapStub{prevRev: "", changeIDs: map[string]int64{"cur": 200}},
			wantKnown:    false,
		},
		{
			name:         "unknown: previous revision can't be located in the configuration",
			currentFirst: "cur",
			stub:         gapStub{prevRev: "p", changeIDs: map[string]int64{"cur": 200}},
			wantKnown:    false,
		},
		{
			name:         "unknown: current first commit can't be located in the configuration",
			currentFirst: "cur",
			stub:         gapStub{prevRev: "p", changeIDs: map[string]int64{"p": 100}},
			wantKnown:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := newGapTestClient(t, tt.stub)
			gap, err := client.getChangesGap(context.Background(), "CUR", "PREV", tt.currentFirst)
			require.NoError(t, err)
			require.NotNil(t, gap)
			assert.Equal(t, tt.wantKnown, gap.Known, "Known")
			assert.Equal(t, tt.wantHasGap, gap.HasGap, "HasGap")
			assert.Equal(t, tt.wantGapCount, gap.GapCommitCount, "GapCommitCount")
			assert.Equal(t, tt.wantFirstAfterDot, gap.FirstCommitAfterPreviousDot, "FirstCommitAfterPreviousDot")
		})
	}
}

// Empty currentFirstCommit must short-circuit before any TeamCity call.
func TestGetChangesGapEmptyFirstCommitSkipsRequests(t *testing.T) {
	t.Parallel()
	requested := false
	client := newGapTestClient(t, gapStub{prevRev: "p", changeIDs: map[string]int64{"p": 100}, requested: &requested})
	gap, err := client.getChangesGap(context.Background(), "CUR", "PREV", "")
	require.NoError(t, err)
	assert.False(t, gap.Known)
	assert.False(t, requested, "should not query TeamCity when there is no range start")
}
