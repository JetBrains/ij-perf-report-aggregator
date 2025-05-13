package meta

import (
	"encoding/json"
	"net/http"
)

type BisectRequest struct {
	TargetValue  string `json:"targetValue"`
	BuildId      string `json:"buildId"`
	Requester    string `json:"requester"`
	Direction    string `json:"direction"`
	Test         string `json:"test"`
	Metric       string `json:"metric"`
	BuildType    string `json:"buildType"`
	ClassName    string `json:"className"`
	ErrorMessage string `json:"errorMessage"`
}

// https://youtrack.jetbrains.com/articles/IJPL-A-201/Bisecting-integration-tests-on-TC
func generateParamsForPerfRun(bisectReq BisectRequest) map[string]string {
	return map[string]string{
		"target.bisect.direction":           bisectReq.Direction,
		"target.bisected.metric":            bisectReq.Metric,
		"target.bisected.simple.class":      bisectReq.ClassName,
		"target.bisected.test":              bisectReq.Test,
		"target.configuration.id":           bisectReq.BuildType,
		"target.build.id":                   bisectReq.BuildId,
		"target.executor.description":       bisectReq.Requester,
		"target.value.before.changed.point": bisectReq.TargetValue,
		"target.perf.messages.mode":         "yes",
		"target.is.bisect.run":              "true",
	}
}

// https://youtrack.jetbrains.com/articles/IJPL-A-201/Bisecting-integration-tests-on-TC
func generateParamsForFunctionalRun(bisectReq BisectRequest) map[string]string {
	return map[string]string{
		"target.bisected.simple.class":          bisectReq.ClassName,
		"target.configuration.id":               bisectReq.BuildType,
		"target.build.id":                       bisectReq.BuildId,
		"target.executor.description":           bisectReq.Requester,
		"env.BISECT_FUNCTIONAL_FAILURE_MESSAGE": bisectReq.ErrorMessage,
		"target.perf.messages.mode":             "no",
		"target.is.bisect.run":                  "true",
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
