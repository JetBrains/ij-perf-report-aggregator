package meta

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// LlmAnalysisMatch is a successful analysis whose guilty commits intersect the requested commits.
type LlmAnalysisMatch struct {
	Id             int      `json:"id"`
	CreatedAt      string   `json:"createdAt"`
	Project        string   `json:"project"`
	Metric         string   `json:"metric"`
	CurrentBuildId string   `json:"currentBuildId"`
	YtIssueId      *string  `json:"ytIssueId,omitempty"`
	DashboardLink  *string  `json:"dashboardLink,omitempty"`
	MatchedCommits []string `json:"matchedCommits"`
}

// LlmAnalysisDegradationMatch is a detected regression of the same or a related metric, on a chart without an analysis
// at that build, whose build contains a guilty commit.
type LlmAnalysisDegradationMatch struct {
	Project string `json:"project"`
	Metric  string `json:"metric"`
	BuildId string `json:"buildId"`
	// Machine is filled when the request names the analysed chart table (db and table) and the regression is in it,
	// so a link to that chart page can select the point; empty means the regression is in another table
	Machine        string   `json:"machine,omitempty"`
	Date           string   `json:"date"`
	MatchedCommits []string `json:"matchedCommits"`
	RangeSize      int      `json:"rangeSize"`
}

type LlmAnalysisMatches struct {
	// RangeSize is the number of commits in the build's range (buildId mode only): one matched commit
	// out of 2 is strong evidence, out of 500 a weak one.
	RangeSize int                `json:"rangeSize,omitempty"`
	Matches   []LlmAnalysisMatch `json:"matches"`
	// Degradations are filled in analysisId mode with degradations=true only, as they take seconds to check.
	Degradations []LlmAnalysisDegradationMatch `json:"degradations,omitempty"`
	// MachinesResolved tells that degradation machines were looked up in the requested chart table, so an empty machine
	// means the regression is in another table, not that the lookup failed
	MachinesResolved bool `json:"machinesResolved,omitempty"`
}

const (
	llmAnalysisMatchesLimit = 50
	// regressions detected this many days before/after the analysed build ran are checked; on prod data
	// this window keeps 94% of the matches found over 60 days
	degradationMatchDaysBefore = 14
	degradationMatchDaysAfter  = 7
	// metrics are related when analyses blamed the same commit for both at least this many times (distinct commits),
	// e.g. localInspections and localInspections#mean_value
	relatedMetricMinCommits = 2
	// how far back from the analysis start its run is looked up; older points fall back to the analysis date
	analysedRunMaxAgeDays = 90
)

// BuildCommitsMatch is the part of a build's commits that is among the requested commits.
type BuildCommitsMatch struct {
	RangeSize      int
	MatchedCommits []string
}

type BuildProject struct {
	BuildId string
	Project string
}

type BuildRun struct {
	Machine       string
	GeneratedTime time.Time
}

// BuildStore reads builds from ClickHouse. Commits are full SHA-1 hex strings.
type BuildStore interface {
	// BuildCommits returns the commits of the build's range, none for an unknown build.
	BuildCommits(ctx context.Context, buildId string) ([]string, error)
	// MatchBuildCommits returns, by build id, the builds that contain any of the commits; the others are absent.
	MatchBuildCommits(ctx context.Context, buildIds []string, commits []string) (map[string]BuildCommitsMatch, error)
	// BuildRuns returns the runs of the projects in the builds found in db.table; non-zero from and to bound
	// generated_time, which lets ClickHouse skip partitions.
	BuildRuns(ctx context.Context, db string, table string, buildIds []string, projects []string, from time.Time, to time.Time) (map[BuildProject]BuildRun, error)
}

// CreateGetLlmAnalysisMatches finds analyses that blamed a commit which could also explain another degradation.
// With buildId, the commits are the ones of that build (a chart point not analysed yet).
// With analysisId, they are the guilty commits of that analysis (other charts that commit was blamed for); with
// degradations=true, detected regressions of other charts at builds containing them are returned instead.
func CreateGetLlmAnalysisMatches(metaDb *pgxpool.Pool, builds BuildStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		query := request.URL.Query()
		result := LlmAnalysisMatches{Matches: []LlmAnalysisMatch{}}
		var commits []string
		excludeId := 0
		var analysis analysisRef
		switch {
		case query.Get("buildId") != "":
			buildCommits, err := builds.BuildCommits(ctx, query.Get("buildId"))
			if err != nil {
				slog.Error("unable to get build commits", "buildId", query.Get("buildId"), "error", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			commits = buildCommits
			result.RangeSize = len(buildCommits)
		case query.Get("analysisId") != "":
			id, err := strconv.Atoi(query.Get("analysisId"))
			if err != nil || id <= 0 {
				http.Error(writer, "invalid analysisId", http.StatusBadRequest)
				return
			}
			err = metaDb.QueryRow(ctx, "SELECT COALESCE(llm_guilty_commits, '{}'), project, metric, current_build_id, created_at FROM analyses WHERE id = $1", id).
				Scan(&commits, &analysis.project, &analysis.metric, &analysis.currentBuildId, &analysis.date)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				slog.Error("unable to select analysis guilty commits", "error", err, "id", id)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			excludeId = id
		default:
			http.Error(writer, "buildId or analysisId is required", http.StatusBadRequest)
			return
		}

		switch {
		case len(commits) == 0:
		case excludeId != 0 && query.Get("degradations") == "true":
			degradations, machinesResolved, err := findDegradationMatches(ctx, metaDb, builds, analysis, commits, query.Get("db"), query.Get("table"))
			if err != nil {
				slog.Error("unable to find degradation matches", "error", err, "analysisId", excludeId)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			result.Degradations = degradations
			result.MachinesResolved = machinesResolved
		default:
			matches, err := findAnalysisMatches(ctx, metaDb, commits, excludeId, analysis)
			if err != nil {
				slog.Error("unable to find analysis matches", "error", err)
				writer.WriteHeader(http.StatusInternalServerError)
				return
			}
			result.Matches = matches
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(result); err != nil {
			slog.Error("unable to write analysis matches response", "error", err)
		}
	}
}

// findAnalysisMatches skips the analysis excludeId and other runs on its chart and build, they are not other charts
func findAnalysisMatches(ctx context.Context, metaDb *pgxpool.Pool, commits []string, excludeId int, exclude analysisRef) ([]LlmAnalysisMatch, error) {
	const sql = "SELECT id, created_at, project, metric, current_build_id, yt_issue_id, dashboard_link, " +
		"ARRAY(SELECT unnest(llm_guilty_commits::text[]) INTERSECT SELECT unnest($1::text[])) " +
		"FROM analyses WHERE state = 'success' AND llm_guilty_commits::text[] && $1::text[] AND id <> $2 " +
		"AND (project, metric, current_build_id) IS DISTINCT FROM ($4, $5, $6) " +
		"ORDER BY id DESC LIMIT $3"
	rows, err := metaDb.Query(ctx, sql, commits, excludeId, llmAnalysisMatchesLimit, exclude.project, exclude.metric, exclude.currentBuildId)
	if err != nil {
		return nil, err
	}
	matches, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (LlmAnalysisMatch, error) {
		var m LlmAnalysisMatch
		var createdAt time.Time
		if err := row.Scan(&m.Id, &createdAt, &m.Project, &m.Metric, &m.CurrentBuildId, &m.YtIssueId, &m.DashboardLink, &m.MatchedCommits); err != nil {
			return LlmAnalysisMatch{}, err
		}
		m.CreatedAt = createdAt.Format(time.RFC3339)
		return m, nil
	})
	if matches == nil {
		matches = []LlmAnalysisMatch{}
	}
	return matches, err
}

type analysisRef struct {
	project        string
	metric         string
	currentBuildId string
	// when the analysed build ran if known, else when the analysis was started
	date time.Time
}

// findDegradationMatches returns regressions of the analysis metric, or a metric related to it, on charts without a
// successful or running analysis at that build (a failed one explains nothing), whose build contains one of the commits. db and table are the analysed chart table, if known.
// build_number holds the TeamCity build id for all detector settings except perfint and fleet (IDE build number),
// those are skipped as their commits can't be resolved.
func findDegradationMatches(ctx context.Context, metaDb *pgxpool.Pool, builds BuildStore, analysis analysisRef, commits []string, db string, table string) ([]LlmAnalysisDegradationMatch, bool, error) {
	knownTable := db != "" && table != ""
	if knownTable {
		// the analysed run happened before the analysis was started, mostly days before
		runs, err := builds.BuildRuns(ctx, db, table, []string{analysis.currentBuildId}, []string{analysis.project},
			analysis.date.AddDate(0, 0, -analysedRunMaxAgeDays), analysis.date.AddDate(0, 0, 1))
		if err != nil {
			// the analysis date is a good enough window center
			slog.Warn("unable to look up the analysed run date", "error", err, "db", db, "table", table)
		} else if run, ok := runs[BuildProject{analysis.currentBuildId, analysis.project}]; ok {
			analysis.date = run.GeneratedTime
		}
	}

	// affected_test is project/metric; DISTINCT ON keeps one row per chart and build (an inferred regression a user
	// also reported), with the longest metric when one metric is a suffix of another
	const sql = "WITH blamed AS (SELECT DISTINCT metric, unnest(llm_guilty_commits) AS c FROM analyses WHERE state = 'success'), " +
		"metrics AS (SELECT $1::text AS metric UNION " +
		"SELECT b2.metric FROM blamed b1 JOIN blamed b2 ON b1.c = b2.c AND b1.metric <> b2.metric " +
		"WHERE b1.metric = $1 GROUP BY b2.metric HAVING count(DISTINCT b1.c) >= $5) " +
		"SELECT DISTINCT ON (ac.affected_test, ac.build_number) left(ac.affected_test, length(ac.affected_test) - length(m.metric) - 1), m.metric, ac.build_number, ac.date " +
		"FROM accidents ac JOIN metrics m ON right(ac.affected_test, length(m.metric) + 1) = '/' || m.metric " +
		"WHERE ac.kind IN ('InferredRegression', 'Regression') AND ac.date BETWEEN $2::date - $3::int AND $2::date + $4::int AND ac.build_number ~ '^[0-9]+$' " +
		"AND NOT EXISTS (SELECT 1 FROM analyses a WHERE a.current_build_id = ac.build_number AND a.project || '/' || a.metric = ac.affected_test " +
		"AND a.state IN ('success', 'in_progress')) " +
		"ORDER BY ac.affected_test, ac.build_number, length(m.metric) DESC"
	rows, err := metaDb.Query(ctx, sql, analysis.metric, analysis.date, degradationMatchDaysBefore, degradationMatchDaysAfter, relatedMetricMinCommits)
	if err != nil {
		return nil, false, err
	}
	candidates, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (LlmAnalysisDegradationMatch, error) {
		var m LlmAnalysisDegradationMatch
		var date time.Time
		if err := row.Scan(&m.Project, &m.Metric, &m.BuildId, &date); err != nil {
			return LlmAnalysisDegradationMatch{}, err
		}
		m.Date = date.Format(time.DateOnly)
		return m, nil
	})
	if err != nil {
		return nil, false, err
	}

	buildIds := make([]string, 0, len(candidates))
	for _, c := range candidates {
		buildIds = append(buildIds, c.BuildId)
	}
	matched, err := builds.MatchBuildCommits(ctx, buildIds, commits)
	if err != nil {
		return nil, false, err
	}
	var result []LlmAnalysisDegradationMatch
	for _, c := range candidates {
		if build, ok := matched[c.BuildId]; ok {
			c.MatchedCommits = build.MatchedCommits
			c.RangeSize = build.RangeSize
			result = append(result, c)
		}
	}

	// the smaller the range, the stronger the evidence
	slices.SortFunc(result, func(a, b LlmAnalysisDegradationMatch) int {
		return cmp.Or(cmp.Compare(a.RangeSize, b.RangeSize), cmp.Compare(a.Project, b.Project), cmp.Compare(a.Metric, b.Metric))
	})
	if len(result) > llmAnalysisMatchesLimit {
		result = result[:llmAnalysisMatchesLimit]
	}
	machinesResolved := false
	if knownTable && len(result) > 0 {
		if err := fillMachines(ctx, builds, db, table, result); err != nil {
			// the list is still useful, links just can't select the point
			slog.Warn("unable to look up degradation machines", "error", err, "db", db, "table", table)
		} else {
			machinesResolved = true
		}
	}
	return result, machinesResolved, nil
}

func fillMachines(ctx context.Context, builds BuildStore, db string, table string, degradations []LlmAnalysisDegradationMatch) error {
	buildIds := make([]string, 0, len(degradations))
	projects := make([]string, 0, len(degradations))
	var from, to time.Time
	for _, d := range degradations {
		buildIds = append(buildIds, d.BuildId)
		projects = append(projects, d.Project)
		date, err := time.Parse(time.DateOnly, d.Date)
		if err != nil {
			return err
		}
		if from.IsZero() || date.Before(from) {
			from = date
		}
		if date.After(to) {
			to = date
		}
	}
	// a regression is dated by its run, a couple of days around covers runs finishing after midnight or timezone shifts
	runs, err := builds.BuildRuns(ctx, db, table, buildIds, projects, from.AddDate(0, 0, -2), to.AddDate(0, 0, 2))
	if err != nil {
		return err
	}
	for i := range degradations {
		degradations[i].Machine = runs[BuildProject{degradations[i].BuildId, degradations[i].Project}].Machine
	}
	return nil
}
