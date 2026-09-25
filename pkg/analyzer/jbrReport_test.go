package analyzer

import (
	"testing"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzePerfJbrReportProject(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		buildType string
		want      string
	}{
		{name: "config after Performance", buildType: "JBR_Dev_Main_Tests_Performance_DaCapo_macOS12aarch64Metal", want: "DaCapo_macOS12aarch64Metal"},
		{name: "last Performance wins", buildType: "JBR_Performance_Tests_Performance_DaCapo", want: "DaCapo"},
		{name: "no Performance", buildType: "JBR_Dev_Main_Tests_DaCapo", want: ""},
		{name: "ends with Performance", buildType: "JBR_Dev_Main_Tests_Performance", want: ""},
		{name: "only separator after Performance", buildType: "JBR_Dev_Main_Tests_Performance_", want: ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			runResult := &RunResult{ReportFileName: "results/report.txt", RawReport: []byte("time\t12,5\n")}
			ignore := analyzePerfJbrReport(runResult, model.ExtraData{TcBuildType: tc.buildType})
			assert.False(t, ignore)
			assert.Equal(t, tc.want, runResult.Report.Project)
		})
	}
}
