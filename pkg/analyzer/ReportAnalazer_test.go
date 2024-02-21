package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/require"
  "log/slog"
  "testing"
)

func TestBranchInferenceForJBRNumber(t *testing.T) {
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
  r := &RunResult{}
  d := model.ExtraData{
    TcBuildProperties: []byte("{}"),
    TcBuildType:       "JBR_Dev_Main_Tests_Performance_DaCapo_macOS12aarch64Metal",
  }
  b, e := getBranch(r, d, "jbr", slog.Default())
  require.NoError(t, e)
  assert.Equal(t, "dev_main", b)
}
