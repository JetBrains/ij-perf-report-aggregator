package installer

import (
	"encoding/json"
	"slices"
	"testing"
)

func TestCheckedOutCommitsSkipsSettingsRepository(t *testing.T) {
	t.Parallel()
	// trimmed response for build 1095312065: the intellij-teamcity-config commits (root 93735) come first
	response := `{
		"revisions": {"revision": [{"vcs-root-instance": {"id": "106699"}}]},
		"changes": {"change": [
			{"version": "36828fb7ead9fd03872a48150527954794a43f8d", "vcsRootInstance": {"id": "93735"}},
			{"version": "071fff0cd530fa2a6a41eec85a500ab33c134e63", "vcsRootInstance": {"id": "93735"}},
			{"version": "9aa3614bdf772d12a2d888cf9b02b43106912a30", "vcsRootInstance": {"id": "106699"}},
			{"version": "abea2dec08000000000000000000000000000000", "vcsRootInstance": {"id": "106699"}}
		]}
	}`
	var build BuildChanges
	if err := json.Unmarshal([]byte(response), &build); err != nil {
		t.Fatal(err)
	}
	want := []string{"9aa3614bdf772d12a2d888cf9b02b43106912a30", "abea2dec08000000000000000000000000000000"}
	if commits := build.CheckedOutCommits(); !slices.Equal(commits, want) {
		t.Errorf("CheckedOutCommits() = %v, want %v", commits, want)
	}
}

func TestCheckedOutCommitsWithoutRevisionsKeepsAllChanges(t *testing.T) {
	t.Parallel()
	var build BuildChanges
	if err := json.Unmarshal([]byte(`{"changes": {"change": [{"version": "a", "vcsRootInstance": {"id": "1"}}, {"version": "b", "vcsRootInstance": {"id": "2"}}]}}`), &build); err != nil {
		t.Fatal(err)
	}
	if commits := build.CheckedOutCommits(); !slices.Equal(commits, []string{"a", "b"}) {
		t.Errorf("CheckedOutCommits() = %v, want [a b]", commits)
	}
}
