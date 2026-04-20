package meta

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type LLMAnalysisRequest struct {
	CurrentBuildId          string   `json:"currentBuildId"`
	AffectedMetric          string   `json:"affectedMetric"`
	CurrentValue            *string  `json:"currentValue"`
	PreviousValue           *string  `json:"previousValue"`
	TestMethodName          *string  `json:"testMethodName"`
	YoutrackIssueReadableId string   `json:"youtrackIssueReadableId"`
	YoutrackIssueId         string   `json:"youtrackIssueId"`
	SpaceUploadedFiles      []string `json:"spaceUploadedFiles"`
}

type DegradationData struct {
	CommitRange   *CommitRange `json:"commitRange,omitempty"`
	TestName      *string      `json:"testName,omitempty"`
	Metric        *Metric      `json:"metric,omitempty"`
	UploadedFiles []string     `json:"uploadedFiles"`
}

type CommitRange struct {
	FromSha string `json:"from_sha"`
	ToSha   string `json:"to_sha"`
}

type Metric struct {
	Name        string  `json:"name"`
	ValueBefore *string `json:"valueBefore,omitempty"`
	ValueAfter  *string `json:"valueAfter,omitempty"`
}

func CreatePostStartLlmAnalysis() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var llmAnalysisRequest LLMAnalysisRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		err := decoder.Decode(&llmAnalysisRequest)
		if err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		var commits *CommitRevisions
		commits, err = teamCityClient.getChanges(request.Context(), llmAnalysisRequest.CurrentBuildId)
		if err != nil {
			slog.Error("cannot get commits from build", "buildId", llmAnalysisRequest.CurrentBuildId, "error", err)
		}

		degradationData := DegradationData{
			TestName: llmAnalysisRequest.TestMethodName,
			Metric: &Metric{
				Name:        llmAnalysisRequest.AffectedMetric,
				ValueBefore: llmAnalysisRequest.PreviousValue,
				ValueAfter:  llmAnalysisRequest.CurrentValue,
			},
			UploadedFiles: llmAnalysisRequest.SpaceUploadedFiles,
		}

		if commits != nil {
			degradationData.CommitRange = &CommitRange{
				FromSha: commits.FirstCommit,
				ToSha:   commits.LastCommit,
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
			"analyses/"+llmAnalysisRequest.YoutrackIssueId,
			"degradation-data.json",
			degradationDataJSON,
		)
		if err != nil {
			slog.Error("failed to upload degradation data to space", "error", err)
		}

		buildParams := map[string]string{
			"degradation.data":           string(degradationDataJSON),
			"youtrack.issue.readable.id": llmAnalysisRequest.YoutrackIssueReadableId,
			"youtrack.issue.id":          llmAnalysisRequest.YoutrackIssueId,
		}

		weburlPtr, err := teamCityClient.startBuild(request.Context(), "ijplatform_master_PerformanceDegradationAnalyzer", buildParams)
		if err != nil {
			http.Error(writer, "Failed to start LLM analysis: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if weburlPtr != nil {
			byteSlice := []byte(*weburlPtr)
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
