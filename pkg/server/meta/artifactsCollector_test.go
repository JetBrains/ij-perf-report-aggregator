package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceSpecialCharactersWithHyphens(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{name: "underscore", input: "my_project", want: "my-project"},
		{name: "collapse runs", input: "foo___bar", want: "foo-bar"},
		{name: "trim leading and trailing", input: "__foo__", want: "foo"},
		{name: "space becomes hyphen", input: "my project", want: "my-project"},
		{name: "preserve dot and slash", input: "foo.bar/baz", want: "foo.bar/baz"},
		{name: "preserve backslash", input: `foo\bar`, want: `foo\bar`},
		{name: "colon becomes hyphen", input: "foo:bar", want: "foo-bar"},
		{name: "mixed specials collapse", input: "a!@#b", want: "a-b"},
		{name: "digits preserved", input: "build123", want: "build123"},
		{name: "non-ascii letter becomes hyphen", input: "café_test", want: "caf-test"},
		{name: "cyrillic becomes hyphen", input: "тест_project", want: "project"},
		{name: "leading non-ascii trimmed", input: "ä-foo", want: "foo"},
		{name: "all specials trimmed to empty", input: "***", want: ""},
		{name: "empty input", input: "", want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := replaceSpecialCharactersWithHyphens(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPerfintCollectorGetArtifactsPaths(t *testing.T) {
	t.Parallel()
	c := perfintCollector{}

	cases := []struct {
		project string
		want    []string
	}{
		{project: "my_project", want: []string{"my-project"}},
		{project: "foo___bar", want: []string{"foo-bar"}},
		{project: "__foo__", want: []string{"foo"}},
		{project: "foo.bar/baz", want: []string{"foo.bar/baz"}},
	}

	for _, tc := range cases {
		t.Run(tc.project, func(t *testing.T) {
			t.Parallel()
			got := c.getArtifactsPaths(UploadAttachmentsRequest{ProjectName: tc.project})
			assert.Equal(t, tc.want, got)
		})
	}
}
