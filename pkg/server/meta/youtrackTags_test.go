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
	tags, _ := buildTags(IDEA, "", nil)
	assert.Contains(t, tagNames(tags), "created-by-ij-perf")
	assert.NotContains(t, tagNames(tags), "analysed-by-ij-perf")
}

func TestBuildTagsAppendsExtraTags(t *testing.T) {
	t.Parallel()
	tags, _ := buildTags(IDEA, "", []Tag{analysedByIjPerfTag})
	names := tagNames(tags)
	assert.Contains(t, names, "created-by-ij-perf")
	assert.Contains(t, names, "analysed-by-ij-perf")
}

func TestBuildTagsExtraTagsForUnmappedProject(t *testing.T) {
	t.Parallel()
	// A project that hits no case in the switch still gets created-by-ij-perf plus extras.
	tags, _ := buildTags("00-0", "", []Tag{analysedByIjPerfTag})
	names := tagNames(tags)
	assert.Contains(t, names, "created-by-ij-perf")
	assert.Contains(t, names, "analysed-by-ij-perf")
}

func TestBuildTagsWithoutProductKeepsProjectDefaults(t *testing.T) {
	t.Parallel()
	tags, productTag := buildTags(IJPL, "", nil)
	assert.Equal(t, []string{"Regression", "blocking-release", "created-by-ij-perf"}, tagNames(tags))
	assert.Nil(t, productTag)

	tags, productTag = buildTags(KT, "", nil)
	assert.Equal(t, []string{"kotlin-regression", "blocking-release-idea", "created-by-ij-perf"}, tagNames(tags))
	assert.Nil(t, productTag)
}

func TestBuildTagsProductTagReplacesGenericBlockingRelease(t *testing.T) {
	t.Parallel()
	// Filed into the shared IJPL tracker from an IDEA dashboard: the product tag wins (AT-5039).
	// It is returned separately (applied after creation) and takes the payload slot of the default.
	tags, productTag := buildTags(IJPL, "idea", nil)
	assert.Equal(t, []string{"Regression", "created-by-ij-perf"}, tagNames(tags))
	require.NotNil(t, productTag)
	assert.Equal(t, "blocking-release-idea", productTag.Name)
}

func TestBuildTagsUnmappedProductFallsBackToProjectDefault(t *testing.T) {
	t.Parallel()
	// The product has no blocking tag in productBlockingTags yet: project default applies.
	tags, productTag := buildTags(IJPL, "rust", nil)
	assert.Equal(t, []string{"Regression", "blocking-release", "created-by-ij-perf"}, tagNames(tags))
	assert.Nil(t, productTag)
}

func TestBuildTagsNoBlockingTagForProjectsWithoutDefault(t *testing.T) {
	t.Parallel()
	// JBR-like target from a dashboard whose product has no tag: no blocking tag at all,
	// because such projects never got the generic one either.
	tags, productTag := buildTags("22-202", "jbr", nil)
	assert.Equal(t, []string{"created-by-ij-perf"}, tagNames(tags))
	assert.Nil(t, productTag)
}

func TestApplyProductBlockingTagSuccessLeavesNoExceptions(t *testing.T) {
	t.Parallel()
	var taggedWith []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tag Tag
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &tag)
		taggedWith = append(taggedWith, tag.ID)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewYoutrackClient(srv.URL, "token")
	productTag := productBlockingTags["pycharm"]
	var exceptions []string
	applyProductBlockingTag(context.Background(), client, "PY-1", "22-36", &productTag, &exceptions)

	assert.Equal(t, []string{productTag.ID}, taggedWith)
	assert.Empty(t, exceptions)
}

func TestApplyProductBlockingTagFallsBackToDefaultWithNotice(t *testing.T) {
	t.Parallel()
	productTag := productBlockingTags["pycharm"]
	var taggedWith []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tag Tag
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &tag)
		// The service account may not use the product tag; the default works.
		if tag.ID == productTag.ID {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		taggedWith = append(taggedWith, tag.ID)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client := NewYoutrackClient(srv.URL, "token")
	var exceptions []string
	applyProductBlockingTag(context.Background(), client, "IJPL-1", IJPL, &productTag, &exceptions)

	assert.Equal(t, []string{blockingReleaseTag.ID}, taggedWith, "the default blocking tag must be applied instead")
	require.Len(t, exceptions, 1)
	assert.Contains(t, exceptions[0], productTag.Name)
	assert.Contains(t, exceptions[0], "service account")
}

func TestApplyProductBlockingTagNoDefaultForProject(t *testing.T) {
	t.Parallel()
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := NewYoutrackClient(srv.URL, "token")
	productTag := productBlockingTags["pycharm"]
	var exceptions []string
	// A target project without a default blocking tag: notice only, no second AddTag attempt.
	applyProductBlockingTag(context.Background(), client, "JBR-1", "22-202", &productTag, &exceptions)

	assert.Equal(t, 1, requests)
	require.Len(t, exceptions, 1)
	assert.Contains(t, exceptions[0], productTag.Name)
}

func TestProductBlockingTagsHaveCompleteIds(t *testing.T) {
	t.Parallel()
	for product, tag := range productBlockingTags {
		assert.NotEmpty(t, tag.ID, "tag id missing for product %q", product)
		assert.NotEmpty(t, tag.Name, "tag name missing for product %q", product)
		assert.Equal(t, "Tag", tag.Type, "tag type wrong for product %q", product)
	}
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
