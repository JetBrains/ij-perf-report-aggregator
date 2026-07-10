package meta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LLMAnalysisRequest struct {
	Project             string                         `json:"project"`
	Metric              string                         `json:"metric"`
	CurrentBuildId      string                         `json:"currentBuildId"`
	PrevBuildId         string                         `json:"prevBuildId"`
	SpaceAttachments    SpaceUploadAttachmentsResponse `json:"spaceAttachments"`
	CurrentValue        *string                        `json:"currentValue,omitempty"`
	PreviousValue       *string                        `json:"previousValue,omitempty"`
	UserName            *string                        `json:"userName,omitempty"`
	FirstCommitRevision *string                        `json:"firstCommitRevision,omitempty"`
	LastCommitRevision  *string                        `json:"lastCommitRevision,omitempty"`
	TestMethodName      *string                        `json:"testMethodName,omitempty"`
	YtIssueId           *string                        `json:"ytIssueId,omitempty"`
	DashboardLink       *string                        `json:"dashboardLink,omitempty"`
}

type LlmAnalysisRun struct {
	Id         int    `json:"id"`
	CreatedAt  string `json:"createdAt"`
	RunBuildId string `json:"runBuildId"`
	State      string `json:"state"`
}

type LlmAnalysisDetails struct {
	LlmAnalysisRun

	Project             string   `json:"project"`
	Metric              string   `json:"metric"`
	CurrentBuildId      string   `json:"currentBuildId"`
	PrevBuildId         string   `json:"prevBuildId"`
	CurrentValue        *string  `json:"currentValue,omitempty"`
	PreviousValue       *string  `json:"previousValue,omitempty"`
	UserName            *string  `json:"userName,omitempty"`
	UserEmail           *string  `json:"userEmail,omitempty"`
	FirstCommitRevision *string  `json:"firstCommitRevision,omitempty"`
	LastCommitRevision  *string  `json:"lastCommitRevision,omitempty"`
	TestMethodName      *string  `json:"testMethodName,omitempty"`
	YtIssueId           *string  `json:"ytIssueId,omitempty"`
	DashboardLink       *string  `json:"dashboardLink,omitempty"`
	LlmGuiltyCommits    []string `json:"llmGuiltyCommits,omitempty"`
	LlmComment          *string  `json:"llmComment,omitempty"`
	TotalCostUsd        *float64 `json:"totalCostUsd,omitempty"`
}

type LlmAnalysisListItem struct {
	LlmAnalysisRun

	Project          string   `json:"project"`
	Metric           string   `json:"metric"`
	CurrentBuildId   string   `json:"currentBuildId"`
	PrevBuildId      string   `json:"prevBuildId"`
	CurrentValue     *string  `json:"currentValue,omitempty"`
	PreviousValue    *string  `json:"previousValue,omitempty"`
	UserName         *string  `json:"userName,omitempty"`
	UserEmail        *string  `json:"userEmail,omitempty"`
	YtIssueId        *string  `json:"ytIssueId,omitempty"`
	DashboardLink    *string  `json:"dashboardLink,omitempty"`
	LlmGuiltyCommits []string `json:"llmGuiltyCommits,omitempty"`
	TotalCostUsd     *float64 `json:"totalCostUsd,omitempty"`
	FeedbackCount    int      `json:"feedbackCount"`
	AvgRating        *float64 `json:"avgRating,omitempty"`
}

const (
	llmAnalysisListDefaultLimit = 100
	llmAnalysisListMaxLimit     = 1000
)

type LlmAnalysisState string

const (
	LlmAnalysisStateInProgress LlmAnalysisState = "in_progress"
	LlmAnalysisStateSuccess    LlmAnalysisState = "success"
	LlmAnalysisStateFailed     LlmAnalysisState = "failed"
)

func (s *LlmAnalysisState) UnmarshalText(text []byte) error {
	v := LlmAnalysisState(text)
	switch v {
	case LlmAnalysisStateInProgress, LlmAnalysisStateSuccess, LlmAnalysisStateFailed:
		*s = v
		return nil
	default:
		return fmt.Errorf("invalid llm analysis state: %q", string(text))
	}
}

type LlmAnalysisRunPatch struct {
	RunBuildId       *string           `json:"runBuildId,omitempty"`
	State            *LlmAnalysisState `json:"state,omitempty"`
	LlmGuiltyCommits *[]string         `json:"llmGuiltyCommits,omitempty"`
	LlmComment       *string           `json:"llmComment,omitempty"`
	TotalCostUsd     *float64          `json:"totalCostUsd,omitempty"`
	YtIssueId        *string           `json:"ytIssueId,omitempty"`
}

var sha1HexRegex = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)

func validateLlmGuiltyCommits(commits []string) error {
	for i, c := range commits {
		if !sha1HexRegex.MatchString(c) {
			return fmt.Errorf("llmGuiltyCommits[%d] is not a 40-char hex SHA: %q", i, c)
		}
	}
	return nil
}

func CreatePostStartLlmAnalysis(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var llmAnalysisRequest LLMAnalysisRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		err := decoder.Decode(&llmAnalysisRequest)
		if err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		userEmail := request.Header.Get("X-Auth-Request-Email")

		fallBackToTCCommitsRangeIfMissing(request.Context(), &llmAnalysisRequest)
		extendRangeAcrossGap(request.Context(), &llmAnalysisRequest)

		id, createdAt, err := insertLlmAnalysisRow(request.Context(), metaDb, llmAnalysisRequest, userEmail)
		if err != nil {
			http.Error(writer, "Failed to insert LLM analysis row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		analysisData, err := json.Marshal(llmAnalysisRequest)
		if err != nil {
			http.Error(writer, "Failed to marshal analysis data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		buildParams := map[string]string{
			"llm.analysis.id": strconv.Itoa(id),
			"analysis.data":   string(analysisData),
		}
		if userEmail != "" {
			buildParams["user.email"] = userEmail
		}
		if llmAnalysisRequest.DashboardLink != nil && *llmAnalysisRequest.DashboardLink != "" {
			if link := buildDashboardLink(*llmAnalysisRequest.DashboardLink, llmAnalysisRequest.CurrentBuildId, id); link != "" {
				buildParams["dashboard.link"] = link
			}
		}

		buildResp, err := teamCityClient.startBuild(request.Context(), "ijplatform_master_PerformanceDegradationAnalyzer", buildParams)
		if err != nil || buildResp == nil || buildResp.Id == 0 {
			markLlmAnalysisFailed(request.Context(), metaDb, id)
			if err != nil {
				http.Error(writer, "Failed to start LLM analysis: "+err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(writer, "TC response doesn't have a build id", http.StatusInternalServerError)
			}
			return
		}

		runBuildId := strconv.FormatInt(buildResp.Id, 10)
		if err := updateLlmAnalysisRun(request.Context(), metaDb, id, LlmAnalysisRunPatch{
			RunBuildId: &runBuildId,
		}); err != nil {
			http.Error(writer, "Failed to update LLM analysis run: "+err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(LlmAnalysisRun{
			Id:         id,
			CreatedAt:  createdAt.Format(time.RFC3339),
			RunBuildId: runBuildId,
			State:      string(LlmAnalysisStateInProgress),
		}); err != nil {
			http.Error(writer, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateGetLlmAnalysisRuns(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		var whereClauses []string
		var args []any
		addFilter := func(column, value string) {
			if value == "" {
				return
			}
			args = append(args, value)
			whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", column, len(args)))
		}
		addFilter("project", query.Get("project"))
		addFilter("metric", query.Get("metric"))
		addFilter("current_build_id", query.Get("currentBuildId"))

		if len(whereClauses) == 0 {
			http.Error(writer, "at least one filter is required", http.StatusBadRequest)
			return
		}

		sql := "SELECT id, created_at, run_build_id, state FROM analyses WHERE " +
			strings.Join(whereClauses, " AND ") + " ORDER BY id DESC"

		rows, err := metaDb.Query(request.Context(), sql, args...)
		if err != nil {
			slog.Error("unable to execute select analyses query", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		runs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (LlmAnalysisRun, error) {
			var run LlmAnalysisRun
			var createdAt time.Time
			var runBuildId *string
			if err := row.Scan(&run.Id, &createdAt, &runBuildId, &run.State); err != nil {
				return LlmAnalysisRun{}, err
			}
			run.CreatedAt = createdAt.Format(time.RFC3339)
			if runBuildId != nil {
				run.RunBuildId = *runBuildId
			}
			return run, nil
		})
		if err != nil {
			slog.Error("unable to collect analyses rows", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if runs == nil {
			runs = []LlmAnalysisRun{}
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(runs); err != nil {
			slog.Error("unable to write analyses response", "error", err)
		}
	}
}

func CreateGetLlmAnalysisList(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		limit := llmAnalysisListDefaultLimit
		if v := query.Get("limit"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
				limit = parsed
			}
		}
		if limit > llmAnalysisListMaxLimit {
			limit = llmAnalysisListMaxLimit
		}
		offset := 0
		if v := query.Get("offset"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
				offset = parsed
			}
		}

		const sql = "SELECT a.id, a.created_at, a.run_build_id, a.state, a.project, a.metric, " +
			"a.current_build_id, a.prev_build_id, a.current_value, a.previous_value, " +
			"a.user_name, a.user_email, a.yt_issue_id, a.dashboard_link, a.llm_guilty_commits, a.total_cost_usd, " +
			"COALESCE(f.cnt, 0) AS feedback_count, f.avg_rating " +
			"FROM analyses a " +
			"LEFT JOIN (" +
			"SELECT analysis_id, COUNT(*) AS cnt, AVG(rate)::float8 AS avg_rating " +
			"FROM analysis_feedback GROUP BY analysis_id" +
			") f ON f.analysis_id = a.id " +
			"ORDER BY a.id DESC LIMIT $1 OFFSET $2"

		rows, err := metaDb.Query(request.Context(), sql, limit, offset)
		if err != nil {
			slog.Error("unable to execute select analyses list query", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		items, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (LlmAnalysisListItem, error) {
			var item LlmAnalysisListItem
			var createdAt time.Time
			var runBuildId *string
			if err := row.Scan(
				&item.Id,
				&createdAt,
				&runBuildId,
				&item.State,
				&item.Project,
				&item.Metric,
				&item.CurrentBuildId,
				&item.PrevBuildId,
				&item.CurrentValue,
				&item.PreviousValue,
				&item.UserName,
				&item.UserEmail,
				&item.YtIssueId,
				&item.DashboardLink,
				&item.LlmGuiltyCommits,
				&item.TotalCostUsd,
				&item.FeedbackCount,
				&item.AvgRating,
			); err != nil {
				return LlmAnalysisListItem{}, err
			}
			item.CreatedAt = createdAt.Format(time.RFC3339)
			if runBuildId != nil {
				item.RunBuildId = *runBuildId
			}
			return item, nil
		})
		if err != nil {
			slog.Error("unable to collect analyses list rows", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if items == nil {
			items = []LlmAnalysisListItem{}
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(items); err != nil {
			slog.Error("unable to write analyses list response", "error", err)
		}
	}
}

func CreateGetLlmAnalysisById(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil || id <= 0 {
			http.Error(writer, "invalid id", http.StatusBadRequest)
			return
		}

		details, err := getAnalysisById(request.Context(), metaDb, id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(writer, "analysis not found", http.StatusNotFound)
				return
			}
			slog.Error("unable to execute select analysis by id query", "error", err, "id", id)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(writer).Encode(details); err != nil {
			slog.Error("unable to write analysis details response", "error", err, "id", id)
		}
	}
}

func getAnalysisById(ctx context.Context, metaDb *pgxpool.Pool, id int) (*LlmAnalysisDetails, error) {
	const sql = "SELECT id, created_at, project, metric, current_build_id, prev_build_id, " +
		"current_value, previous_value, user_name, user_email, " +
		"first_commit_revision, last_commit_revision, test_method_name, run_build_id, yt_issue_id, " +
		"dashboard_link, state, llm_guilty_commits, llm_comment, total_cost_usd " +
		"FROM analyses WHERE id = $1"

	var (
		details    LlmAnalysisDetails
		createdAt  time.Time
		runBuildId *string
	)
	if err := metaDb.QueryRow(ctx, sql, id).Scan(
		&details.Id,
		&createdAt,
		&details.Project,
		&details.Metric,
		&details.CurrentBuildId,
		&details.PrevBuildId,
		&details.CurrentValue,
		&details.PreviousValue,
		&details.UserName,
		&details.UserEmail,
		&details.FirstCommitRevision,
		&details.LastCommitRevision,
		&details.TestMethodName,
		&runBuildId,
		&details.YtIssueId,
		&details.DashboardLink,
		&details.State,
		&details.LlmGuiltyCommits,
		&details.LlmComment,
		&details.TotalCostUsd,
	); err != nil {
		return nil, err
	}
	details.CreatedAt = createdAt.Format(time.RFC3339)
	if runBuildId != nil {
		details.RunBuildId = *runBuildId
	}
	return &details, nil
}

func CreatePatchLlmAnalysisRun(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(request, "id"))
		if err != nil || id <= 0 {
			http.Error(writer, "invalid id", http.StatusBadRequest)
			return
		}
		var patch LlmAnalysisRunPatch
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		if err := decoder.Decode(&patch); err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if patch.LlmGuiltyCommits != nil {
			if err := validateLlmGuiltyCommits(*patch.LlmGuiltyCommits); err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}
		}
		if err := updateLlmAnalysisRun(request.Context(), metaDb, id, patch); err != nil {
			http.Error(writer, "Failed to update LLM analysis run: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Once the job reports a terminal result back to ij-perf, mark the linked ticket as analysed.
		// Best-effort: a tagging failure must not fail the request.
		if patch.State != nil && (*patch.State == LlmAnalysisStateSuccess) {
			tagAnalysedYtIssue(request.Context(), metaDb, id)
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

// tagAnalysedYtIssue adds the analysed-by-ij-perf tag to the YouTrack issue linked to the analysis,
// if any. Best-effort: all failures are logged and swallowed.
func tagAnalysedYtIssue(ctx context.Context, metaDb *pgxpool.Pool, id int) {
	var ytIssueId *string
	if err := metaDb.QueryRow(ctx, "SELECT yt_issue_id FROM analyses WHERE id = $1", id).Scan(&ytIssueId); err != nil {
		slog.Warn("cannot read yt_issue_id for analysed-by-ij-perf tagging", "id", id, "error", err)
		return
	}
	if ytIssueId == nil || *ytIssueId == "" {
		return
	}
	if err := youtrackClient.AddTag(ctx, *ytIssueId, analysedByIjPerfTag); err != nil {
		slog.Warn("cannot tag linked YouTrack issue as analysed-by-ij-perf", "issueId", *ytIssueId, "error", err)
	}
}

// buildDashboardLink parses the FE-supplied link, strips scheme/host (defense in depth in case the
// FE still sends a full URL), appends point/analysis query params, and returns the absolute URL
// (dashboardBaseURL + "/path?query"). Returns "" if the input cannot be parsed or has no path.
const dashboardBaseURL = "https://ij-perf.labs.jb.gg"

func buildDashboardLink(rawLink string, currentBuildId string, analysisId int) string {
	u, err := url.Parse(rawLink)
	if err != nil {
		slog.Warn("invalid dashboardLink, skipping dashboard.link build param",
			"error", err, "dashboardLink", rawLink)
		return ""
	}
	u.Scheme = ""
	u.Host = ""
	u.User = nil
	if u.Path == "" || !strings.HasPrefix(u.Path, "/") {
		slog.Warn("dashboardLink has no absolute path, skipping dashboard.link build param",
			"dashboardLink", rawLink)
		return ""
	}
	q := u.Query()
	q.Set("point", currentBuildId)
	q.Set("analysis", strconv.Itoa(analysisId))
	u.RawQuery = q.Encode()
	return dashboardBaseURL + u.RequestURI()
}

// fallBackToTCCommitsRangeIfMissing populates FirstCommitRevision / LastCommitRevision from TC's build changes
// when the FE omitted them. Best-effort: any TC error is logged and the request proceeds unchanged.
func fallBackToTCCommitsRangeIfMissing(ctx context.Context, req *LLMAnalysisRequest) {
	firstMissing := req.FirstCommitRevision == nil
	lastMissing := req.LastCommitRevision == nil
	if !firstMissing && !lastMissing {
		return
	}
	commits, err := teamCityClient.getChanges(ctx, req.CurrentBuildId)
	if err != nil {
		slog.Warn("cannot get commits from build for revision fallback",
			"buildId", req.CurrentBuildId, "error", err)
		return
	}
	if commits == nil {
		return
	}
	if firstMissing && commits.FirstCommit != "" {
		req.FirstCommitRevision = &commits.FirstCommit
	}
	if lastMissing && commits.LastCommit != "" {
		req.LastCommitRevision = &commits.LastCommit
	}
}

// extendRangeAcrossGap moves the analysis range's first commit back to the first commit after the
// previous dot when builds failed/timed out in between, so the analysis covers the commits those
// builds consumed. Normally there is no gap and the current build's own first commit is left as
// is. Best-effort: on any TeamCity error the range is unchanged.
func extendRangeAcrossGap(ctx context.Context, req *LLMAnalysisRequest) {
	if req.FirstCommitRevision == nil || *req.FirstCommitRevision == "" || req.PrevBuildId == "" || req.CurrentBuildId == "" {
		return
	}
	gap, err := teamCityClient.getChangesGap(ctx, req.CurrentBuildId, req.PrevBuildId, *req.FirstCommitRevision)
	if err != nil {
		slog.Warn("cannot determine changes gap for LLM analysis range",
			"currentBuildId", req.CurrentBuildId, "prevBuildId", req.PrevBuildId, "error", err)
		return
	}
	if gap == nil || !gap.HasGap || gap.FirstCommitAfterPreviousDot == "" {
		return
	}
	req.FirstCommitRevision = &gap.FirstCommitAfterPreviousDot
}

func markLlmAnalysisFailed(ctx context.Context, metaDb *pgxpool.Pool, id int) {
	state := LlmAnalysisStateFailed
	if err := updateLlmAnalysisRun(ctx, metaDb, id, LlmAnalysisRunPatch{State: &state}); err != nil {
		slog.Error("cannot mark llm_analysis_run as failed", "error", err, "id", id)
	}
}

func updateLlmAnalysisRun(ctx context.Context, metaDb *pgxpool.Pool, id int, patch LlmAnalysisRunPatch) error {
	var setClauses []string
	var args []any
	add := func(column string, value any) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, len(args)+1))
		args = append(args, value)
	}
	if patch.RunBuildId != nil {
		add("run_build_id", *patch.RunBuildId)
	}
	if patch.State != nil {
		add("state", string(*patch.State))
	}
	if patch.LlmGuiltyCommits != nil {
		add("llm_guilty_commits", *patch.LlmGuiltyCommits)
	}
	if patch.LlmComment != nil {
		add("llm_comment", *patch.LlmComment)
	}
	if patch.TotalCostUsd != nil {
		add("total_cost_usd", *patch.TotalCostUsd)
	}
	if patch.YtIssueId != nil {
		add("yt_issue_id", *patch.YtIssueId)
	}

	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, id)
	sql := fmt.Sprintf("UPDATE analyses SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), len(args))

	if _, err := metaDb.Exec(ctx, sql, args...); err != nil {
		slog.Error("cannot execute update analyses query", "error", err, "id", id, "sql", sql)
		return err
	}
	return nil
}

func insertLlmAnalysisRow(ctx context.Context, metaDb *pgxpool.Pool, params LLMAnalysisRequest, userEmail string) (int, time.Time, error) {
	var id int
	var createdAt time.Time
	var userEmailArg *string
	if userEmail != "" {
		userEmailArg = &userEmail
	}
	idRow := metaDb.QueryRow(ctx,
		"INSERT INTO analyses (project, metric, current_build_id, prev_build_id, current_value, previous_value, user_name, user_email, first_commit_revision, last_commit_revision, test_method_name, yt_issue_id, dashboard_link) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, created_at",
		params.Project, params.Metric, params.CurrentBuildId,
		params.PrevBuildId, params.CurrentValue, params.PreviousValue, params.UserName, userEmailArg,
		params.FirstCommitRevision, params.LastCommitRevision, params.TestMethodName, params.YtIssueId, params.DashboardLink)
	if err := idRow.Scan(&id, &createdAt); err != nil {
		slog.Error("cannot execute insert analyses query", "error", err,
			"project", params.Project, "metric", params.Metric)
		return 0, time.Time{}, err
	}
	return id, createdAt, nil
}
