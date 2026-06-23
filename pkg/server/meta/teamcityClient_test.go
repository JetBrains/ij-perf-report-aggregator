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
	prevRev      string   // revision the previous dot's build was built on ("" => build has no revision)
	prevRevID    int64    // change id resolved for prevRev (0 => not found in the configuration)
	sinceNewest  []string // commits strictly after prevRev, newest-first (as the changes endpoint returns)
	changeIDSeen *bool    // set true if the change-id (version) lookup was hit
}

// newGapTestClient returns a TeamCityClient wired to an httptest server that
// answers the four endpoints getChangesGap calls, according to stub.
func newGapTestClient(t *testing.T, stub gapStub) *TeamCityClient {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		locator := r.URL.Query().Get("locator")
		switch {
		case strings.HasPrefix(r.URL.Path, "/app/rest/builds/"):
			// getBuildRevision asks for the revisions field; getBuildType does not.
			if strings.Contains(r.URL.RawQuery, "revisions") {
				rev := []map[string]string{}
				if stub.prevRev != "" {
					rev = append(rev, map[string]string{"version": stub.prevRev})
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"revisions": map[string]any{"revision": rev}})
			} else {
				_ = json.NewEncoder(w).Encode(map[string]string{"buildTypeId": "BT"})
			}
		case strings.HasPrefix(r.URL.Path, "/app/rest/changes"):
			if strings.Contains(locator, "version:") {
				// getChangeID: resolve prevRev -> change id.
				if stub.changeIDSeen != nil {
					*stub.changeIDSeen = true
				}
				change := []map[string]int64{}
				if stub.prevRevID != 0 {
					change = append(change, map[string]int64{"id": stub.prevRevID})
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"change": change})
			} else {
				// getChangesSince: commits after prevRev, newest-first.
				change := make([]map[string]string, 0, len(stub.sinceNewest))
				for _, v := range stub.sinceNewest {
					change = append(change, map[string]string{"version": v})
				}
				_ = json.NewEncoder(w).Encode(map[string]any{"count": len(change), "change": change})
			}
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
		name         string
		currentFirst string
		stub         gapStub
		wantKnown    bool
		wantHasGap   bool
		wantGapCount int
	}{
		{
			name:         "gap: range start is newer than previous dot, leaving commits uncovered",
			currentFirst: "c3",
			// commits after prevRev, newest-first; c2 and c1 are older than the range start c3.
			stub:         gapStub{prevRev: "p", prevRevID: 100, sinceNewest: []string{"c5", "c4", "c3", "c2", "c1"}},
			wantKnown:    true,
			wantHasGap:   true,
			wantGapCount: 2,
		},
		{
			name:         "no gap: range start is the immediate successor of the previous dot",
			currentFirst: "c1",
			stub:         gapStub{prevRev: "p", prevRevID: 100, sinceNewest: []string{"c3", "c2", "c1"}},
			wantKnown:    true,
			wantHasGap:   false,
			wantGapCount: 0,
		},
		{
			name:         "no gap: range reaches back to or before the previous dot (start not among newer commits)",
			currentFirst: "older-than-prev",
			stub:         gapStub{prevRev: "p", prevRevID: 100, sinceNewest: []string{"c2", "c1"}},
			wantKnown:    true,
			wantHasGap:   false,
			wantGapCount: 0,
		},
		{
			name:         "unknown: no current range start",
			currentFirst: "",
			stub:         gapStub{prevRev: "p", prevRevID: 100, sinceNewest: []string{"c1"}},
			wantKnown:    false,
		},
		{
			name:         "unknown: previous dot's build has no revision",
			currentFirst: "c1",
			stub:         gapStub{prevRev: "", prevRevID: 100, sinceNewest: []string{"c1"}},
			wantKnown:    false,
		},
		{
			name:         "unknown: previous revision can't be located in the configuration",
			currentFirst: "c1",
			stub:         gapStub{prevRev: "p", prevRevID: 0, sinceNewest: []string{"c1"}},
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
		})
	}
}

// Empty currentFirstCommit must short-circuit before any TeamCity call.
func TestGetChangesGapEmptyFirstCommitSkipsRequests(t *testing.T) {
	t.Parallel()
	changeIDSeen := false
	client := newGapTestClient(t, gapStub{prevRev: "p", prevRevID: 100, sinceNewest: []string{"c1"}, changeIDSeen: &changeIDSeen})
	gap, err := client.getChangesGap(context.Background(), "CUR", "PREV", "")
	require.NoError(t, err)
	assert.False(t, gap.Known)
	assert.False(t, changeIDSeen, "should not query TeamCity when there is no range start")
}
