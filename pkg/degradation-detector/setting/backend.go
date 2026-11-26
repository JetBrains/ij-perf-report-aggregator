package setting

import (
	"net/http"
	"sync"

	detector "github.com/JetBrains/ij-perf-report-aggregator/pkg/degradation-detector"
)

// Backend is what the settings need to know about the collected data: which tests ran, which metrics they have and
// who owns them.
type Backend interface {
	// Tests returns the tests that ran on the branch and the machines in the last 30 days.
	Tests(db, table, branch, machine string) ([]string, error)
	// MetricNames returns the metrics of the test of the settings whose name matches the SQL LIKE pattern.
	MetricNames(test detector.PerformanceSettings, pattern string) ([]string, error)
	// ProjectOwners returns the code owner of each test that has one.
	ProjectOwners(db, table string) (map[string]string, error)
	OwnerGroups() ([]detector.CodeOwnerGroup, error)
}

type httpBackend struct {
	url    string
	client *http.Client

	mu    sync.Mutex
	tests map[[4]string]testsResult
}

type testsResult struct {
	tests []string
	err   error
}

func NewBackend(url string, client *http.Client) Backend {
	return &httpBackend{url: url, client: client, tests: map[[4]string]testsResult{}}
}

// Tests keeps the answer: several subscribers select tests of the same table, branch and machines.
func (b *httpBackend) Tests(db, table, branch, machine string) ([]string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	key := [4]string{db, table, branch, machine}
	result, ok := b.tests[key]
	if !ok {
		result.tests, result.err = detector.FetchAllTests(b.url, b.client, detector.PerformanceSettings{
			Db: db, Table: table, Branch: branch, Machine: machine,
		})
		b.tests[key] = result
	}
	return result.tests, result.err
}

func (b *httpBackend) MetricNames(test detector.PerformanceSettings, pattern string) ([]string, error) {
	return detector.FetchMetricNamesByPattern(b.url, b.client, test, pattern)
}

func (b *httpBackend) ProjectOwners(db, table string) (map[string]string, error) {
	return detector.FetchProjectOwners(b.url, b.client, db, table)
}

func (b *httpBackend) OwnerGroups() ([]detector.CodeOwnerGroup, error) {
	return detector.FetchCodeOwnerGroups(b.url, b.client)
}
