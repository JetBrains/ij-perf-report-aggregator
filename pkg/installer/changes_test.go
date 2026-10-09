package installer

import (
	"slices"
	"testing"
)

func TestChangesRoundTrip(t *testing.T) {
	t.Parallel()
	commits := []string{"22596363b3de40b06f981fb85d82312e8c0ed511", "0123456789abcdef0123456789abcdef01234567"}
	encoded := make([]string, 0, len(commits))
	for _, c := range commits {
		e, ok := EncodeCommit(c)
		if !ok {
			t.Fatalf("EncodeCommit(%q) rejected a hex commit", c)
		}
		if len(e) >= len(c) {
			t.Errorf("EncodeCommit(%q) = %q is not shorter than hex", c, e)
		}
		encoded = append(encoded, e)
	}
	if decoded := DecodeChanges(encoded); !slices.Equal(decoded, commits) {
		t.Errorf("DecodeChanges(EncodeCommit) = %v, want %v", decoded, commits)
	}

	if _, ok := EncodeCommit("13 04 2022 12:14"); ok {
		t.Error("EncodeCommit accepted a private build change")
	}
	if decoded := DecodeChanges([]string{"not base64!", ""}); !slices.Equal(decoded, []string{"not base64!", ""}) {
		t.Errorf("DecodeChanges kept %v", decoded)
	}
}
