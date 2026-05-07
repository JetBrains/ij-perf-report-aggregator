package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LLMAnalysisRequest struct {
	Date                string  `json:"date"`
	Project             string  `json:"project"`
	Metric              string  `json:"metric"`
	CurrentBuildId      string  `json:"currentBuildId"`
	PrevBuildId         string  `json:"prevBuildId"`
	CurrentValue        *string `json:"currentValue,omitempty"`
	PreviousValue       *string `json:"previousValue,omitempty"`
	UserName            *string `json:"userName,omitempty"`
	FirstCommitRevision *string `json:"firstCommitRevision,omitempty"`
	LastCommitRevision  *string `json:"lastCommitRevision,omitempty"`
	TestMethodName      *string `json:"testMethodName,omitempty"`
}

type LlmAnalysisRun struct {
	Id         int    `json:"id"`
	Date       string `json:"date"`
	RunBuildId string `json:"runBuildId"`
	State      string `json:"state"`
}

type DegradationData struct {
	CommitRange   *CommitRange `json:"commitRange,omitempty"`
	TestName      *string      `json:"testName,omitempty"`
	Metric        *Metric      `json:"metric,omitempty"`
	UploadedFiles []string     `json:"uploadedFiles"`
}

type CommitRange struct {
	FromSha *string `json:"from_sha,omitempty"`
	ToSha   *string `json:"to_sha,omitempty"`
}

type Metric struct {
	Name        string  `json:"name"`
	ValueBefore *string `json:"valueBefore,omitempty"`
	ValueAfter  *string `json:"valueAfter,omitempty"`
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

		id, err := insertLlmAnalysisRow(request.Context(), metaDb, llmAnalysisRequest)
		if err != nil {
			http.Error(writer, "Failed to insert LLM analysis row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		/*degradationData := DegradationData{
			TestName: llmAnalysisRequest.TestMethodName,
			Metric: &Metric{
				Name:        llmAnalysisRequest.Metric,
				ValueBefore: llmAnalysisRequest.PreviousValue,
				ValueAfter:  llmAnalysisRequest.CurrentValue,
			},
			UploadedFiles: nil,
		}

		if llmAnalysisRequest.FirstCommitRevision != nil || llmAnalysisRequest.LastCommitRevision != nil {
			degradationData.CommitRange = &CommitRange{
				FromSha: llmAnalysisRequest.FirstCommitRevision,
				ToSha:   llmAnalysisRequest.LastCommitRevision,
			}
		}

		degradationDataJSON, err := json.Marshal(degradationData)
		if err != nil {
			http.Error(writer, "Failed to marshal degradation data: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = spacePackagesClient.UploadFile(
			request.Context(),
			"platform-test-automation",
			"performance-regression-llm-analysis",
			"analyses/1",
			"degradation-data.json",
			bytes.NewReader(degradationDataJSON),
		)
		if err != nil {
			slog.Error("failed to upload degradation data to space", "error", err)
		}*/

		buildParams := map[string]string{
			//"degradation.data":           string(degradationDataJSON),
			"llm.analysis.id": strconv.Itoa(id),
		}

		//TODO: remove staging after testing
		weburlPtr, err := teamCityClient.startBuild(request.Context(), "ijplatform_staging_PerformanceDegradationAnalyzer", buildParams)
		if err != nil {
			http.Error(writer, "Failed to start LLM analysis: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if weburlPtr != nil {
			weburl := *weburlPtr
			runBuildId := weburl[strings.LastIndex(weburl, "/")+1:]
			state := "queued"
			if err := updateLlmAnalysisRun(request.Context(), metaDb, id, LlmAnalysisRunUpdate{
				RunBuildId: &runBuildId,
				State:      &state,
			}); err != nil {
				http.Error(writer, "Failed to update LLM analysis run: "+err.Error(), http.StatusInternalServerError)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(writer).Encode(LlmAnalysisRun{
				Id:         id,
				Date:       llmAnalysisRequest.Date,
				RunBuildId: runBuildId,
				State:      state,
			}); err != nil {
				http.Error(writer, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(writer, "TC response doesn't have weburl", http.StatusInternalServerError)
		}
	}
}

type LlmAnalysisRunUpdate struct {
	RunBuildId       *string
	State            *string
	LlmGuiltyCommits *[]string
	LlmComment       *string
	UserRate         *bool
	UserComment      *string
}

func updateLlmAnalysisRun(ctx context.Context, metaDb *pgxpool.Pool, id int, u LlmAnalysisRunUpdate) error {
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
		add("state", *u.State)
	}
	if u.LlmGuiltyCommits != nil {
		add("llm_guilty_commits", *u.LlmGuiltyCommits)
	}
	if u.LlmComment != nil {
		add("llm_comment", *u.LlmComment)
	}
	if u.UserRate != nil {
		add("user_rate", *u.UserRate)
	}
	if u.UserComment != nil {
		add("user_comment", *u.UserComment)
	}

	if len(setClauses) == 0 {
		return nil
	}
	args = append(args, id)
	sql := fmt.Sprintf("UPDATE llm_analysis_runs SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "), len(args))

	if _, err := metaDb.Exec(ctx, sql, args...); err != nil {
		slog.Error("cannot execute update llm_analysis_runs query", "error", err, "id", id, "sql", sql)
		return err
	}
	return nil
}

func insertLlmAnalysisRow(ctx context.Context, metaDb *pgxpool.Pool, params LLMAnalysisRequest) (int, error) {
	var id int
	idRow := metaDb.QueryRow(ctx,
		"INSERT INTO llm_analysis_runs (date, project, metric, current_build_id, prev_build_id, current_value, previous_value, user_name, first_commit_revision, last_commit_revision, test_method_name) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		params.Date, params.Project, params.Metric, params.CurrentBuildId,
		params.PrevBuildId, params.CurrentValue, params.PreviousValue, params.UserName,
		params.FirstCommitRevision, params.LastCommitRevision, params.TestMethodName)
	if err := idRow.Scan(&id); err != nil {
		slog.Error("cannot execute insert llm_analysis_runs query", "error", err,
			"date", params.Date, "project", params.Project, "metric", params.Metric)
		return 0, err
	}
	return id, nil
}
