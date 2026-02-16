package meta

import (
	"encoding/json"
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

		weburlPtr, err := teamCityClient.startBuild(request.Context(), "ijplatform_master_BisectChangesetOnSpace", buildParams)
		if err != nil {
			http.Error(writer, "Failed to start bisect: "+err.Error(), http.StatusInternalServerError)
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
