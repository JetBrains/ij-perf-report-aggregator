package meta

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
)

type YoutrackCreateIssueRequest struct {
	ProjectId      string        `json:"projectId"`
	AccidentId     string        `json:"accidentId"`
	BuildLink      string        `json:"buildLink"`
	ChangesLink    string        `json:"changesLink"`
	CustomFields   []CustomField `json:"customFields"`
	DashboardLink  string        `json:"dashboardLink"`
	AffectedMetric string        `json:"affectedMetric"`
	Delta          string        `json:"delta"`
	TestMethodName *string       `json:"testMethodName"`
}

type UploadAttachmentsToIssueRequest struct {
	IssueId                string                 `json:"issueId"`
	TeamCityAttachmentInfo TeamCityAttachmentInfo `json:"teamcityAttachmentInfo"`
	AffectedTest           string                 `json:"affectedTest"`
	ChartPng               *[]byte                `json:"chartPng"`
}

type GenerateDescriptionData struct {
	Kind           string
	AffectedTest   string
	AffectedMetric string
	Delta          string
	StackTrace     string
	BuildLink      string
	Changes        string
	DashboardLink  string
	TestHistoryUrl *string
}

type CreateIssueResponse struct {
	Issue      YoutrackIssue `json:"issue"`
	Exceptions []string      `json:"exceptions"`
}

var teamCityClient = NewTeamCityClient("https://buildserver.labs.intellij.net", os.Getenv("TEAMCITY_TOKEN"))
var youtrackClient = NewYoutrackClient("https://youtrack.jetbrains.com", os.Getenv("YOUTRACK_TOKEN"))

func CreatePostCreateIssueByAccident(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response := CreateIssueResponse{}

		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			handleError(writer, "cannot read body", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		defer body.Close()

		var params YoutrackCreateIssueRequest
		if err = json.Unmarshal(all, &params); err != nil {
			handleError(writer, "cannot unmarshal parameters", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		relatedAccident, err := getAccidentById(request.Context(), metaDb, params.AccidentId)
		lowerKind := strings.ToLower(relatedAccident.Kind)

		if err != nil {
			handleError(writer, "cannot get accident", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		affectedTest := relatedAccident.AffectedTest
		affectedMetric := params.AffectedMetric

		if strings.HasSuffix(affectedTest, affectedMetric) {
			affectedTest = strings.TrimSuffix(affectedTest, "/"+affectedMetric)
		}

		var testHistoryUrl string
		if params.TestMethodName != nil {
			testHistoryUrl, err = teamCityClient.getTestHistoryUrl(request.Context(), *params.TestMethodName)
		}

		if err != nil {
			handleError(writer, "cannot get test history link", err, &response.Exceptions)
		}

		descriptionData := GenerateDescriptionData{
			lowerKind,
			affectedTest,
			affectedMetric,
			params.Delta,
			relatedAccident.Stacktrace,
			params.BuildLink,
			params.ChangesLink,
			params.DashboardLink,
			&testHistoryUrl,
		}

		issueInfo := CreateIssueInfo{
			Summary:      fmt.Sprintf("[%s] %s", lowerKind, relatedAccident.Reason),
			Description:  generateDescription(descriptionData),
			Project:      YoutrackProject{ID: params.ProjectId},
			CustomFields: params.CustomFields,
		}

		issue, err := youtrackClient.CreateIssue(request.Context(), issueInfo)

		if err != nil {
			handleError(writer, "failed to create issue", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, &response)
			return
		}

		response.Issue = *issue

		err = marshalAndWriteIssueResponse(writer, &response)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		relatedAccident.Reason = fmt.Sprintf("%s: %s", issue.IDReadable, relatedAccident.Reason)
		err = updateAccidentReason(request.Context(), metaDb, relatedAccident)
		if err != nil {
			handleError(writer, "unable to update accident reason", err, &response.Exceptions)
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func CreatePostUploadAttachmentsToIssue() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		type Exceptions struct {
			Exceptions []string `json:"exceptions"`
		}

		var exceptions Exceptions
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			handleError(writer, "cannot read body", err, &exceptions.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			return
		}

		defer body.Close()

		var params UploadAttachmentsToIssueRequest
		if err = json.Unmarshal(all, &params); err != nil {
			handleError(writer, "cannot unmarshal parameters", err, &exceptions.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			return
		}

		errCh := make(chan error, 10)

		builds := []int{params.TeamCityAttachmentInfo.CurrentBuildId}
		if params.TeamCityAttachmentInfo.PreviousBuildId != nil {
			builds = append(builds, *params.TeamCityAttachmentInfo.PreviousBuildId)
		}

		if params.ChartPng != nil {
			err := youtrackClient.UploadAttachment(request.Context(), params.IssueId, *params.ChartPng, "dashboard.png")
			if err != nil {
				slog.Error("Failed to upload dashboard attachment to youtrack", "error", err)
				errCh <- err
			}
		}

		var wg sync.WaitGroup
		wg.Add(len(builds))

		for index, buildId := range builds {
			go func(index int, buildId int) {
				defer wg.Done()

				testArtifactPath := strings.ReplaceAll(params.AffectedTest, "_", "-")
				children, err := teamCityClient.getArtifactChildren(request.Context(), buildId, params.TeamCityAttachmentInfo.BuildTypeId, testArtifactPath)
				if err != nil {
					slog.Error("Failed to get teamcity artifact children", "error", err)
					errCh <- err
					return
				}

				var filteredChildren []string

				for _, str := range children {
					for _, keyword := range []string{"logs-", "snapshots-"} {
						if strings.Contains(str, keyword) {
							filteredChildren = append(filteredChildren, str)
						}
					}
				}

				var childWg sync.WaitGroup
				childWg.Add(len(filteredChildren))

				var attachmentPostfix string

				if index == 0 {
					attachmentPostfix = "current"
				} else {
					attachmentPostfix = "before"
				}

				for _, str := range filteredChildren {
					go func(artifactName string) {
						defer childWg.Done()
						artifact, err := teamCityClient.downloadArtifact(request.Context(), params.TeamCityAttachmentInfo.BuildTypeId, buildId, testArtifactPath+"/"+artifactName)
						if err != nil {
							slog.Error("Failed to download artefacts form teamcity", "error", err)
							errCh <- err
							return
						}

						attachmentName := getAttachmentName(artifactName, attachmentPostfix)
						err = youtrackClient.UploadAttachment(request.Context(), params.IssueId, artifact, attachmentName)
						if err != nil {
							slog.Error("Failed to upload attachment to youtrack", "error", err)
							errCh <- err
							return
						}
					}(str)
				}
				childWg.Wait()
			}(index, buildId)
		}
		wg.Wait()
		close(errCh)

		for err := range errCh {
			if err != nil {
				exceptions.Exceptions = append(exceptions.Exceptions, err.Error())
			}
		}

		if len(exceptions.Exceptions) > 0 {
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = marshalAndWriteIssueResponse(writer, exceptions)
		writer.WriteHeader(http.StatusOK)
	}
}

func generateDescription(generateDescriptorData GenerateDescriptionData) string {
	affectedTest := "**Affected test:**\n" + generateDescriptorData.AffectedTest
	affectedMetric := fmt.Sprintf("**Affected metric:**\n%s (Delta: %s)", generateDescriptorData.AffectedMetric, generateDescriptorData.Delta)
	build := fmt.Sprintf("**Build:**\n[build link](%s)", generateDescriptorData.BuildLink)
	changes := fmt.Sprintf("**Changes in space:**\n[space link](%s)", generateDescriptorData.Changes)
	logs := "**Idea logs, screenshots, thread dumps etc:**\n Current: [logs-current.zip](logs-current.zip)"
	snapshots := "**Snapshots:**\n Current: [snapshots-current.zip](snapshots-current.zip)"
	dashboard := fmt.Sprintf("**Dashboard:**\n[dashboard link](%s)", generateDescriptorData.DashboardLink)
	dashboardPng := "![](dashboard.png)"
	stacktrace := fmt.Sprintf("**Stacktrace:**\n```%s```", generateDescriptorData.StackTrace)
	var testHistory string
	if generateDescriptorData.TestHistoryUrl != nil {
		testHistory = fmt.Sprintf("**Test history:**\n[test history link](%s)", *generateDescriptorData.TestHistoryUrl)
	} else {
		testHistory = ""
	}
	var description string
	if generateDescriptorData.Kind != "exception" {
		logs += "\n Before: [logs-before.zip](logs-before.zip)"
		snapshots += "\n Before: [snapshots-before.zip](snapshots-before.zip)"
		description = strings.Join([]string{affectedTest, affectedMetric, build, changes, logs, snapshots, dashboard, dashboardPng}, "\n\n")
	} else {
		description = strings.Join([]string{affectedTest, testHistory, build, changes, logs, snapshots, stacktrace}, "\n\n")
	}
	return description
}

func getAttachmentName(filename, suffix string) string {
	parts := strings.Split(filename, ".")
	if len(parts) != 2 {
		return filename
	}

	nameWithoutExt := parts[0]
	ext := parts[1]

	nameParts := strings.Split(nameWithoutExt, "-")

	updatedName := nameParts[0] + "-" + suffix
	return fmt.Sprintf("%s.%s", updatedName, ext)
}

func handleError(writer http.ResponseWriter, message string, err error, exceptions *[]string) {
	slog.Error(message, "error", err)
	writer.WriteHeader(http.StatusInternalServerError)
	*exceptions = append(*exceptions, fmt.Sprintf("Message: %s. Error: %s", message, err.Error()))
}

func marshalAndWriteIssueResponse(writer http.ResponseWriter, response interface{}) error {
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		slog.Error("cannot marshal response", "error", err)
		return err
	}
	_, err = writer.Write(jsonBytes)
	if err != nil {
		slog.Error("cannot write response", "error", err)
		return err
	}

	return nil
}
