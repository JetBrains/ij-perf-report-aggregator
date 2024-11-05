package meta

import (
	"encoding/json"
	"net/http"
)

type BisectRequest struct {
	TargetValue string `json:"targetValue"`
	Changes     string `json:"changes"`
	Direction   string `json:"direction"`
	Test        string `json:"test"`
	Metric      string `json:"metric"`
	Branch      string `json:"branch"`
	BuildType   string `json:"buildType"`
	ClassName   string `json:"className"`
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

		weburlPtr, err := teamCityClient.startBuild(request.Context(), "ijplatform_master_BisectChangeset", map[string]string{
			"target.bisect.direction":           bisectReq.Direction,
			"target.bisected.metric":            bisectReq.Metric,
			"target.bisected.simple.class":      bisectReq.ClassName,
			"target.bisected.test":              bisectReq.Test,
			"target.branch":                     bisectReq.Branch,
			"target.configuration.id":           bisectReq.BuildType,
			"target.git.commits":                bisectReq.Changes,
			"target.value.before.changed.point": bisectReq.TargetValue,
		})
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
