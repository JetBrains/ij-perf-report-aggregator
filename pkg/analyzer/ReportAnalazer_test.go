package analyzer

import (
	"log/slog"
	"testing"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOwnerFromBuildProperties(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		props string
		want  string
	}{
		{name: "present", props: `{"ultimate.codeowner":"PyCharm Exec Experts"}`, want: "PyCharm Exec Experts"},
		{name: "absent", props: `{"teamcity.build.branch":"master"}`, want: ""},
		{name: "empty properties", props: "", want: ""},
		{name: "empty object", props: "{}", want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ownerFromBuildProperties(model.ExtraData{TcBuildProperties: []byte(tc.props)})
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestBranchInferenceForJBRNumber(t *testing.T) {
	t.Parallel()
	r := &RunResult{}
	d := model.ExtraData{
		TcBuildProperties: []byte("{}"),
		TcBuildType:       "JBR_232_JBR17_Tests_Performance_DaCapo_macOS12aarch64Metal",
	}
	b, e := getBranch(r, d, "jbr", slog.Default())
	require.NoError(t, e)
	assert.Equal(t, "232_jbr17", b)
}

func TestBranchInferenceForJBRMaster(t *testing.T) {
	t.Parallel()
	r := &RunResult{}
	d := model.ExtraData{
		TcBuildProperties: []byte("{}"),
		TcBuildType:       "JBR_Master_JBR17_Tests_Performance_DaCapo_macOS12aarch64Metal",
	}
	b, e := getBranch(r, d, "jbr", slog.Default())
	require.NoError(t, e)
	assert.Equal(t, "master_jbr17", b)
}

func TestBranchInferenceForDevMain(t *testing.T) {
	t.Parallel()
	r := &RunResult{}
	d := model.ExtraData{
		TcBuildProperties: []byte("{}"),
		TcBuildType:       "JBR_Dev_Main_Tests_Performance_DaCapo_macOS12aarch64Metal",
	}
	b, e := getBranch(r, d, "jbr", slog.Default())
	require.NoError(t, e)
	assert.Equal(t, "dev_main", b)
}
