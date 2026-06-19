package degradation_detector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	dataQuery "github.com/JetBrains/ij-perf-report-aggregator/pkg/data-query"
)

func ExpandTestsByPattern(backendUrl string, client *http.Client, tests []string, baseSettings PerformanceSettings) []string {
	testsExpanded := make([]string, 0, len(tests)*5)
	for _, test := range tests {
		if strings.Contains(test, "%") {
			matchingTests, err := fetchTestsByPattern(backendUrl, client, baseSettings, test)
			if err != nil {
				slog.Error("error while fetching tests by pattern", "error", err, "pattern", test)
				continue
			}
			testsExpanded = append(testsExpanded, matchingTests...)
		} else {
			testsExpanded = append(testsExpanded, test)
		}
	}
	return testsExpanded
}

func fetchTestsByPattern(backendUrl string, client *http.Client, settings PerformanceSettings, pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	filters := []dataQuery.QueryFilter{
		{Field: "branch", Value: settings.Branch},
		{Field: "generated_time", Sql: ">subtractDays(now(),30)"},
		{Field: "triggeredBy", Value: ""},
	}
	if settings.Machine != "" {
		filters = append(filters, dataQuery.QueryFilter{Field: "machine", Value: settings.Machine, Operator: "like"})
	}
	if pattern != "" {
		filters = append(filters, dataQuery.QueryFilter{Field: "project", Value: pattern, Operator: "like"})
	}
	query := dataQuery.Query{
		Database: settings.Db,
		Table:    settings.Table,
		Fields:   []dataQuery.QueryDimension{{Name: "project", Sql: "distinct project"}},
		Flat:     true,
		Filters:  filters,
		Order:    []string{"project"},
	}
	response, err := getValuesFromServer(ctx, client, backendUrl, query)
	if err != nil {
		return nil, err
	}
	tests, err := extractValuesFromRequest(response)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

func FetchAllTests(backendUrl string, client *http.Client, settings PerformanceSettings) ([]string, error) {
	return fetchTestsByPattern(backendUrl, client, settings, "")
}

// FetchProjectOwners returns the full project->owner mapping for the given db/table
// from the meta database, in a single request. Projects without a resolved code
// owner are simply absent from the map.
func FetchProjectOwners(backendUrl string, client *http.Client, db, table string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	params := url.Values{}
	params.Set("db", db)
	params.Set("table", table)
	requestURL := backendUrl + "/api/meta/projectOwners?" + params.Encode()

	return getJSONMap(ctx, client, requestURL, "project owners")
}

// FetchCodeOwnerChannels returns a map from code-owner group name (and each alias) to its
// Slack channel, resolved by the backend via the CodeOwners service. Owners without a
// configured channel are absent from the map.
func FetchCodeOwnerChannels(backendUrl string, client *http.Client) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return getJSONMap(ctx, client, backendUrl+"/api/meta/codeOwnerChannels", "code-owner channels")
}

// getJSONMap issues a GET request to requestURL and decodes the JSON response body into a
// string->string map. description identifies the call in error messages (e.g. "project owners").
func getJSONMap(ctx context.Context, client *http.Client, requestURL, description string) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", description, err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send GET request for %s: %w", description, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get %s: %v", description, resp.Status)
	}
	result := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", description, err)
	}
	return result, nil
}

// FetchMetricNamesByPattern returns distinct measures.name values matching the given SQL LIKE
// pattern for the project in settings. Patterns without "%" are returned as-is without a query.
func FetchMetricNamesByPattern(backendUrl string, client *http.Client, settings PerformanceSettings, pattern string) ([]string, error) {
	if !strings.Contains(pattern, "%") {
		return []string{pattern}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	filters := []dataQuery.QueryFilter{
		{Field: "branch", Value: settings.Branch},
		{Field: "generated_time", Sql: ">subtractDays(now(),30)"},
		{Field: "project", Value: settings.Project},
		{Field: "measures.name", Value: pattern, Operator: "like"},
		{Field: "triggeredBy", Value: ""},
	}
	if settings.Machine != "" {
		filters = append(filters, dataQuery.QueryFilter{Field: "machine", Value: settings.Machine, Operator: "like"})
	}
	query := dataQuery.Query{
		Database: settings.Db,
		Table:    settings.Table,
		Fields:   []dataQuery.QueryDimension{{Name: "measures", SubName: "name", Sql: "distinct measures.name"}},
		Flat:     true,
		Filters:  filters,
		Order:    []string{"measures.name"},
	}
	response, err := getValuesFromServer(ctx, client, backendUrl, query)
	if err != nil {
		return nil, err
	}
	metrics, err := extractValuesFromRequest(response)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func extractValuesFromRequest(response []byte) ([]string, error) {
	var data [][]any

	err := json.Unmarshal(response, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	if len(data[0]) < 1 {
		return nil, errors.New("not enough data")
	}
	tests, err := SliceToSliceOfString(data[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert values: %w", err)
	}
	return tests, nil
}
