package setting

import (
	"net/http"
	"strings"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

// Startup runs are reported under "<project>/measureStartup" and "<project>/warmup" — the same
// convention the startup dashboards select by (StartupProjectConfigurator). Every other project in
// the table is a regular performance test that never reports startup metrics, so alerting on it
// would only produce missing data.
var startupProjectSuffixes = []string{"/measureStartup"}

func fetchStartupProjects(backendUrl string, client *http.Client, settings detector.PerformanceSettings) ([]string, error) {
	projects, err := detector.FetchAllTests(backendUrl, client, settings)
	if err != nil {
		return nil, err
	}

	startupProjects := make([]string, 0, len(projects))
	for _, project := range projects {
		for _, suffix := range startupProjectSuffixes {
			if strings.HasSuffix(project, suffix) {
				startupProjects = append(startupProjects, project)
				break
			}
		}
	}
	return startupProjects, nil
}
