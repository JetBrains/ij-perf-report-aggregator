package meta

import (
	"context"
	"encoding/json"
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
	PrevBuildId         *string `json:"prevBuildId,omitempty"`
	CurrentValue        *string `json:"currentValue,omitempty"`
	PreviousValue       *string `json:"previousValue,omitempty"`
	UserName            *string `json:"userName,omitempty"`
	FirstCommitRevision *string `json:"firstCommitRevision,omitempty"`
	LastCommitRevision  *string `json:"lastCommitRevision,omitempty"`
	TestMethodName      *string `json:"testMethodName,omitempty"`
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

		weburlPtr, err := teamCityClient.startBuild(request.Context(), "ijplatform_master_PerformanceDegradationAnalyzer", buildParams)
		if err != nil {
			http.Error(writer, "Failed to start LLM analysis: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if weburlPtr != nil {
			weburl := *weburlPtr
			runBuildId := weburl[strings.LastIndex(weburl, "/")+1:]
			_ = updateLlmAnalysisRunBuildId(request.Context(), metaDb, id, runBuildId)

			byteSlice := []byte(weburl)
			_, err = writer.Write(byteSlice)
			if err != nil {
				http.Error(writer, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(writer, "TC response doesn't have weburl", http.StatusInternalServerError)
		}
	}
}

func updateLlmAnalysisRunBuildId(ctx context.Context, metaDb *pgxpool.Pool, id int, runBuildId string) error {
	_, err := metaDb.Exec(ctx,
		"UPDATE llm_analysis_runs SET run_build_id = $1 WHERE id = $2",
		runBuildId, id)
	if err != nil {
		slog.Error("cannot execute update llm_analysis_runs.run_build_id query", "error", err,
			"id", id, "runBuildId", runBuildId)
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
