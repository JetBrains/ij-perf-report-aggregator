package meta

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tagNames(tags []Tag) []string {
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return names
}

func TestBuildTagsAlwaysIncludesCreatedByIjPerf(t *testing.T) {
	t.Parallel()
	tags := buildTags(IDEA, nil)
	assert.Contains(t, tagNames(tags), "created-by-ij-perf")
	assert.NotContains(t, tagNames(tags), "analysed-by-ij-perf")
}

func TestBuildTagsAppendsExtraTags(t *testing.T) {
	t.Parallel()
	tags := buildTags(IDEA, []Tag{analysedByIjPerfTag})
	names := tagNames(tags)
	assert.Contains(t, names, "created-by-ij-perf")
	assert.Contains(t, names, "analysed-by-ij-perf")
}

func TestBuildTagsExtraTagsForUnmappedProject(t *testing.T) {
	t.Parallel()
	// A project that hits no case in the switch still gets created-by-ij-perf plus extras.
	tags := buildTags("00-0", []Tag{analysedByIjPerfTag})
	names := tagNames(tags)
	assert.Contains(t, names, "created-by-ij-perf")
	assert.Contains(t, names, "analysed-by-ij-perf")
}

func TestAnalysedByIjPerfTagHasDistinctId(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "68-534888", analysedByIjPerfTag.ID)
	assert.NotEqual(t, "68-523929", analysedByIjPerfTag.ID, "must differ from created-by-ij-perf tag id")
}

func TestAddTagPostsToIssueTagsEndpoint(t *testing.T) {
	t.Parallel()
	var gotPath, gotMethod string
	var gotTag Tag
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotTag)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"1-1","name":"analysed-by-ij-perf"}`))
	}))
	defer srv.Close()

	client := NewYoutrackClient(srv.URL, "token")
	err := client.AddTag(context.Background(), "IJPL-1234", analysedByIjPerfTag)
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, gotMethod)
	assert.Equal(t, "/api/issues/IJPL-1234/tags", gotPath)
	assert.Equal(t, "68-534888", gotTag.ID)
}

func TestAddTagReturnsErrorOnFailure(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := NewYoutrackClient(srv.URL, "token")
	err := client.AddTag(context.Background(), "NOPE-1", analysedByIjPerfTag)
	require.Error(t, err)
}
