package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type BisectRequest struct {
	TargetValue     string `json:"targetValue"`
	BuildId         string `json:"buildId"`
	Changes         string `json:"changes"`
	Mode            string `json:"mode"`
	Requester       string `json:"requester"`
	Direction       string `json:"direction"`
	Test            string `json:"test"`
	Metric          string `json:"metric"`
	BuildType       string `json:"buildType"`
	TestPatterns    string `json:"testPatterns"`
	ErrorMessage    string `json:"errorMessage"`
	ExcludedCommits string `json:"excludedCommits"`
	JpsCompilation  string `json:"jpsCompilation"`
	DashboardLink   string `json:"dashboardLink"`
	YtIssueId       string `json:"ytIssueId"`
}

// https://youtrack.jetbrains.com/articles/IJPL-A-201/Bisecting-integration-tests-on-TC
func generateParamsForPerfRun(bisectReq BisectRequest) map[string]string {
	return map[string]string{
		"target.bisect.direction":             bisectReq.Direction,
		"target.bisected.metric":              bisectReq.Metric,
		"target.intellij.build.test.patterns": bisectReq.TestPatterns,
		"target.bisected.test":                bisectReq.Test,
		"target.configuration.id":             bisectReq.BuildType,
		"target.build.id":                     bisectReq.BuildId,
		"target.git.commits":                  bisectReq.Changes,
		"target.mode":                         bisectReq.Mode,
		"target.executor.description":         bisectReq.Requester,
		"target.value.before.changed.point":   bisectReq.TargetValue,
		"target.perf.messages.mode":           "yes",
		"target.is.bisect.run":                "true",
		"target.commits.to.exclude":           bisectReq.ExcludedCommits,
		"target.jps.compile":                  bisectReq.JpsCompilation,
		"target.ij.perf.url":                  bisectReq.DashboardLink,
	}
}

// https://youtrack.jetbrains.com/articles/IJPL-A-201/Bisecting-integration-tests-on-TC
func generateParamsForFunctionalRun(bisectReq BisectRequest) map[string]string {
	return map[string]string{
		"target.intellij.build.test.patterns":   bisectReq.TestPatterns,
		"target.configuration.id":               bisectReq.BuildType,
		"target.build.id":                       bisectReq.BuildId,
		"target.git.commits":                    bisectReq.Changes,
		"target.mode":                           bisectReq.Mode,
		"target.executor.description":           bisectReq.Requester,
		"env.BISECT_FUNCTIONAL_FAILURE_MESSAGE": bisectReq.ErrorMessage,
		"target.perf.messages.mode":             "no",
		"target.is.bisect.run":                  "true",
		"target.commits.to.exclude":             bisectReq.ExcludedCommits,
		"target.jps.compile":                    bisectReq.JpsCompilation,
	}
}

// The bisect build sends its result notification to this issue on top of its usual notification.
const youtrackIssueBuildParam = "target.youtrack.issue"

func addYoutrackIssueParam(buildParams map[string]string, bisectReq BisectRequest) {
	if bisectReq.YtIssueId != "" {
		buildParams[youtrackIssueBuildParam] = bisectReq.YtIssueId
	}
}

func HandleGetTeamCityChanges() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildID := r.URL.Query().Get("buildId")
		if buildID == "" {
			http.Error(w, "buildId parameter is required", http.StatusBadRequest)
			return
		}

		revisions, err := teamCityClient.getChanges(r.Context(), buildID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(revisions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HandleGetTeamCityChangesGap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildID := r.URL.Query().Get("buildId")
		previousBuildID := r.URL.Query().Get("previousBuildId")
		currentFirstCommit := r.URL.Query().Get("currentFirstCommit")
		if buildID == "" || previousBuildID == "" || currentFirstCommit == "" {
			http.Error(w, "buildId, previousBuildId and currentFirstCommit parameters are required", http.StatusBadRequest)
			return
		}

		gap, err := teamCityClient.getChangesGap(r.Context(), buildID, previousBuildID, currentFirstCommit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(gap)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HandleGetTeamCityBuildType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildID := r.URL.Query().Get("buildId")
		if buildID == "" {
			http.Error(w, "buildId parameter is required", http.StatusBadRequest)
			return
		}

		revisions, err := teamCityClient.getBuildType(r.Context(), buildID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(revisions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HandleGetTeamCityBuildCounter() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildID := r.URL.Query().Get("buildId")
		if buildID == "" {
			http.Error(w, "buildId parameter is required", http.StatusBadRequest)
			return
		}

		counter, err := teamCityClient.getBuildCounter(r.Context(), buildID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(counter))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HandleGetTeamCityArtifactsExist() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildID := r.URL.Query().Get("buildId")
		if buildID == "" {
			http.Error(w, "buildId parameter is required", http.StatusBadRequest)
			return
		}

		hasArtifacts, err := teamCityClient.hasArtifacts(r.Context(), buildID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]bool{"hasArtifacts": hasArtifacts})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func HandleGetTeamCityBuildInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildID := r.URL.Query().Get("buildId")
		if buildID == "" {
			http.Error(w, "buildId parameter is required", http.StatusBadRequest)
			return
		}

		buildInfo, err := teamCityClient.getBuildInfo(r.Context(), buildID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(buildInfo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreatePostStartBisect() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var bisectReq BisectRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		err := decoder.Decode(&bisectReq)
		if err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		var buildParams map[string]string
		if bisectReq.ErrorMessage != "" {
			buildParams = generateParamsForFunctionalRun(bisectReq)
		} else {
			buildParams = generateParamsForPerfRun(bisectReq)
		}
		addYoutrackIssueParam(buildParams, bisectReq)

		buildResp, err := teamCityClient.startBuild(request.Context(), "ijplatform_master_BisectChangesetOnSpace", buildParams)
		if err != nil {
			http.Error(writer, "Failed to start bisect: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if buildResp != nil && buildResp.WebURL != "" {
			commentStartedBisectOnYoutrackIssue(request.Context(), bisectReq.YtIssueId, buildResp.WebURL)
			_, err = writer.Write([]byte(buildResp.WebURL))
			if err != nil {
				http.Error(writer, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(writer, "TC response doesn't have weburl", http.StatusInternalServerError)
		}
	}
}

// commentStartedBisectOnYoutrackIssue records the started bisect on the linked issue, so the ticket
// carries the link to the run before the bisect itself reports the result there.
//
// Runs detached from the request: the bisect is already started, so its response must not wait for
// YouTrack, and the issue is usually created seconds earlier and may not be readable yet.
// Best-effort: all failures are logged and swallowed.
func commentStartedBisectOnYoutrackIssue(ctx context.Context, issueID string, buildURL string) {
	if issueID == "" {
		return
	}
	ctx = context.WithoutCancel(ctx)
	go func() {
		if err := youtrackClient.waitIssueIsCreated(ctx, issueID); err != nil {
			slog.Warn("linked YouTrack issue is not readable, skipping the started-bisect comment",
				"issueId", issueID, "buildUrl", buildURL, "error", err)
			return
		}
		comment := fmt.Sprintf("Bisect started from IJ Perf: %s\n\nThe result will be posted here once the bisect finishes.", buildURL)
		if err := youtrackClient.AddComment(ctx, issueID, comment); err != nil {
			slog.Warn("cannot comment the started bisect on the linked YouTrack issue",
				"issueId", issueID, "buildUrl", buildURL, "error", err)
		}
	}()
}
