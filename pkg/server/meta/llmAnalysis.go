package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LLMAnalysisRequest struct {
	Project             string                          `json:"project"`
	Metric              string                          `json:"metric"`
	CurrentBuildId      string                          `json:"currentBuildId"`
	PrevBuildId         string                          `json:"prevBuildId"`
	CurrentValue        *string                         `json:"currentValue,omitempty"`
	PreviousValue       *string                         `json:"previousValue,omitempty"`
	UserName            *string                         `json:"userName,omitempty"`
	FirstCommitRevision *string                         `json:"firstCommitRevision,omitempty"`
	LastCommitRevision  *string                         `json:"lastCommitRevision,omitempty"`
	TestMethodName      *string                         `json:"testMethodName,omitempty"`
	YtIssueId           *string                         `json:"ytIssueId,omitempty"`
	SpaceAttachments    *SpaceUploadAttachmentsResponse `json:"spaceAttachments,omitempty"`
}

type LlmAnalysisRun struct {
	Id         int    `json:"id"`
	CreatedAt  string `json:"createdAt"`
	RunBuildId string `json:"runBuildId"`
	State      string `json:"state"`
}

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

type LlmAnalysisRunUpdate struct {
	Id               int               `json:"id"`
	RunBuildId       *string           `json:"runBuildId,omitempty"`
	State            *LlmAnalysisState `json:"state,omitempty"`
	LlmGuiltyCommits *[]string         `json:"llmGuiltyCommits,omitempty"`
	LlmComment       *string           `json:"llmComment,omitempty"`
	TotalCostUsd     *float64          `json:"totalCostUsd,omitempty"`
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

		id, createdAt, err := insertLlmAnalysisRow(request.Context(), metaDb, llmAnalysisRequest)
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
		if email := request.Header.Get("X-Auth-Request-Email"); email != "" {
			buildParams["user.email"] = email
		}

		weburlPtr, err := teamCityClient.startBuild(request.Context(), "ijplatform_master_PerformanceDegradationAnalyzer", buildParams)
		if err != nil || weburlPtr == nil {
			markLlmAnalysisFailed(request.Context(), metaDb, id)
			if err != nil {
				http.Error(writer, "Failed to start LLM analysis: "+err.Error(), http.StatusInternalServerError)
			} else {
				http.Error(writer, "TC response doesn't have weburl", http.StatusInternalServerError)
			}
			return
		}

		weburl := *weburlPtr
		runBuildId := weburl[strings.LastIndex(weburl, "/")+1:]
		if err := updateLlmAnalysisRun(request.Context(), metaDb, LlmAnalysisRunUpdate{
			Id:         id,
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
		project := query.Get("project")
		metric := query.Get("metric")
		currentBuildId := query.Get("currentBuildId")
		if project == "" || metric == "" || currentBuildId == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		rows, err := metaDb.Query(request.Context(),
			"SELECT id, created_at, run_build_id, state FROM analyses "+
				"WHERE project = $1 AND metric = $2 AND current_build_id = $3 "+
				"ORDER BY id DESC",
			project, metric, currentBuildId)
		if err != nil {
			slog.Error("unable to execute select analyses query", "error", err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer rows.Close()

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

func CreatePostUpdateLlmAnalysisRun(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var u LlmAnalysisRunUpdate
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		if err := decoder.Decode(&u); err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if u.Id <= 0 {
			http.Error(writer, "id is required", http.StatusBadRequest)
			return
		}
		if err := updateLlmAnalysisRun(request.Context(), metaDb, u); err != nil {
			http.Error(writer, "Failed to update LLM analysis run: "+err.Error(), http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusNoContent)
	}
}

func markLlmAnalysisFailed(ctx context.Context, metaDb *pgxpool.Pool, id int) {
	state := LlmAnalysisStateFailed
	if err := updateLlmAnalysisRun(ctx, metaDb, LlmAnalysisRunUpdate{Id: id, State: &state}); err != nil {
		slog.Error("cannot mark llm_analysis_run as failed", "error", err, "id", id)
	}
}

func updateLlmAnalysisRun(ctx context.Context, metaDb *pgxpool.Pool, u LlmAnalysisRunUpdate) error {
	var setClauses []string
	var args []any
	add := func(column string, value any) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, len(args)+1))
		args = append(args, value)
	}
	if u.RunBuildId != nil {
		add("run_build_id", *u.RunBuildId)
	}
	if u.State != nil {
		add("state", string(*u.State))
	}
	if u.LlmGuiltyCommits != nil {
		add("llm_guilty_commits", *u.LlmGuiltyCommits)
	}
	if u.LlmComment != nil {
		add("llm_comment", *u.LlmComment)
	}
	if u.TotalCostUsd != nil {
		add("total_cost_usd", *u.TotalCostUsd)
	}

	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, u.Id)
	sql := fmt.Sprintf("UPDATE analyses SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), len(args))

	if _, err := metaDb.Exec(ctx, sql, args...); err != nil {
		slog.Error("cannot execute update analyses query", "error", err, "id", u.Id, "sql", sql)
		return err
	}
	return nil
}

func insertLlmAnalysisRow(ctx context.Context, metaDb *pgxpool.Pool, params LLMAnalysisRequest) (int, time.Time, error) {
	var id int
	var createdAt time.Time
	idRow := metaDb.QueryRow(ctx,
		"INSERT INTO analyses (project, metric, current_build_id, prev_build_id, current_value, previous_value, user_name, first_commit_revision, last_commit_revision, test_method_name, yt_issue_id) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at",
		params.Project, params.Metric, params.CurrentBuildId,
		params.PrevBuildId, params.CurrentValue, params.PreviousValue, params.UserName,
		params.FirstCommitRevision, params.LastCommitRevision, params.TestMethodName, params.YtIssueId)
	if err := idRow.Scan(&id, &createdAt); err != nil {
		slog.Error("cannot execute insert analyses query", "error", err,
			"project", params.Project, "metric", params.Metric)
		return 0, time.Time{}, err
	}
	return id, createdAt, nil
}
