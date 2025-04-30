package meta

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/auth"
	"github.com/jackc/pgx/v5/pgxpool"
)

type YoutrackCreateIssueRequest struct {
	ProjectId      string        `json:"projectId"`
	AccidentId     string        `json:"accidentId"`
	TicketLabel    string        `json:"ticketLabel"`
	BuildLink      string        `json:"buildLink"`
	ChangesLink    string        `json:"changesLink"`
	CustomFields   []CustomField `json:"customFields"`
	DashboardLink  string        `json:"dashboardLink"`
	AffectedMetric string        `json:"affectedMetric"`
	Delta          string        `json:"delta"`
	TestMethodName *string       `json:"testMethodName"`
	TestType       string        `json:"testType"`
}

type UploadAttachmentsToIssueRequest struct {
	IssueId                string                 `json:"issueId"`
	TeamCityAttachmentInfo TeamCityAttachmentInfo `json:"teamcityAttachmentInfo"`
	AffectedTest           string                 `json:"affectedTest"`
	ChartPng               *[]byte                `json:"chartPng"`
	TestType               string                 `json:"testType"`
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
	TestMethod     *string
	TestType       string
}

type CreateIssueResponse struct {
	Issue      YoutrackIssue `json:"issue"`
	Exceptions []string      `json:"exceptions"`
}

var (
	teamCityClient = NewTeamCityClient("https://buildserver.labs.intellij.net", os.Getenv("TEAMCITY_TOKEN"))
	youtrackClient = NewYoutrackClient("https://youtrack.jetbrains.com", os.Getenv("YOUTRACK_TOKEN"))
	ytAuth         = auth.NewYTAuth("https://youtrack.jetbrains.com", os.Getenv("YOUTRACK_TOKEN"))
)

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

		testHistoryUrl := ""
		if params.TestMethodName != nil && lowerKind == "exception" {
			testHistoryUrl, err = teamCityClient.getTestHistoryUrl(request.Context(), *params.TestMethodName)
		}

		if err != nil {
			logError("cannot get test history link", err, &response.Exceptions)
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
			params.TestMethodName,
			params.TestType,
		}

		accessToken := request.Header.Get("X-Auth-Request-Access-Token")
		user, err := auth.FetchUserInfo(request.Context(), accessToken)
		if err != nil {
			slog.Warn("cannot fetch user info", "error", err)
		}

		userId, err := ytAuth.GetUser(request.Context(), user.Email)
		if err != nil {
			slog.Warn("error getting user id:", "error", err)
			userId = nil
		}
		issueInfo := CreateIssueInfo{
			Summary:     fmt.Sprintf("[%s] %s", lowerKind, params.TicketLabel),
			Description: generateDescription(descriptionData),
			Project:     YoutrackProject{ID: params.ProjectId},
			Reporter:    userId,
			Visibility: Visibility{
				PermittedGroups: []auth.YTUser{{ID: "10-3"}},
				PermittedUsers:  []auth.YTUser{{ID: "11-1539792"}},
				Type:            "LimitedVisibility",
			},
			CustomFields: []CustomField{
				{
					Name: "Type",
					Type: "SingleEnumIssueCustomField",
					Value: struct {
						Name string `json:"name"`
					}{
						Name: "Performance Problem",
					},
				},
			},
		}
		if params.ProjectId == "22-414" { // If project ID is KTIJ set subsystem as they require it
			subsystemsCustomField := CustomField{
				Name: "Subsystems",
				Type: "MultiOwnedIssueCustomField",
				Value: []CustomFieldValue{
					{Name: "IDE"},
				},
			}
			issueInfo.CustomFields = append(issueInfo.CustomFields, subsystemsCustomField)
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
			logError("unable to update accident reason", err, &response.Exceptions)
		}
	}
}

type artifactCollector interface {
	getArtifactsPath(params UploadAttachmentsToIssueRequest) string
	checkArtifact(artifactName string) bool
}

type fleetStartupCollector struct{}

func (f fleetStartupCollector) getArtifactsPath(UploadAttachmentsToIssueRequest) string {
	return ""
}

func (f fleetStartupCollector) checkArtifact(artifactName string) bool {
	return strings.HasSuffix(artifactName, "fleet.fahrplan.json")
}

type fleetPerfTestCollector struct{}

func (f fleetPerfTestCollector) getArtifactsPath(UploadAttachmentsToIssueRequest) string {
	return ""
}

func (f fleetPerfTestCollector) checkArtifact(artifactName string) bool {
	return artifactName == "logs.zip"
}

type perfUnitTestCollector struct{}

func (f perfUnitTestCollector) getArtifactsPath(params UploadAttachmentsToIssueRequest) string {
	return params.AffectedTest
}

func (f perfUnitTestCollector) checkArtifact(artifactName string) bool {
	return artifactName == "log.zip"
}

type perfintCollector struct{}

func (f perfintCollector) getArtifactsPath(params UploadAttachmentsToIssueRequest) string {
	return strings.ReplaceAll(params.AffectedTest, "_", "-")
}

func (f perfintCollector) checkArtifact(artifactName string) bool {
	prefixes := []string{"logs-", "snapshots-"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(artifactName, prefix) {
			return true
		}
	}
	return false
}

func getArtifactCollector(testType string) artifactCollector {
	switch testType {
	case "fleet":
		return fleetStartupCollector{}
	case "intellij_dev", "intellij":
		return perfintCollector{}
	case "fleet_perf":
		return fleetPerfTestCollector{}
	case "perfUnitTests":
		return perfUnitTestCollector{}
	default:
		return nil
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

		collector := getArtifactCollector(params.TestType)

		if collector != nil {
			var wg sync.WaitGroup
			wg.Add(len(builds))

			for index, buildId := range builds {
				go func(index int, buildId int) {
					defer wg.Done()

					testArtifactPath := collector.getArtifactsPath(params)

					children, err := teamCityClient.getArtifactChildren(request.Context(), buildId, testArtifactPath)
					if err != nil {
						slog.Error("Failed to get teamcity artifact children", "error", err)
						errCh <- err
						return
					}

					var filteredChildren []string

					for _, str := range children {
						if collector.checkArtifact(str) {
							filteredChildren = append(filteredChildren, str)
						}
					}

					var attachmentPostfix string
					if index == 0 {
						attachmentPostfix = "current"
					} else {
						attachmentPostfix = "before"
					}

					var childWg sync.WaitGroup
					childWg.Add(len(filteredChildren))
					for _, str := range filteredChildren {
						go func(artifactName string) {
							defer childWg.Done()
							artifact, err := teamCityClient.downloadArtifact(request.Context(), buildId, testArtifactPath+"/"+artifactName)
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
		}
		close(errCh)

		for err := range errCh {
			if err != nil {
				exceptions.Exceptions = append(exceptions.Exceptions, err.Error())
			}
		}

		if len(exceptions.Exceptions) > 0 {
			writer.WriteHeader(http.StatusInternalServerError)
			_ = marshalAndWriteIssueResponse(writer, exceptions)
			return
		}

		_ = marshalAndWriteIssueResponse(writer, exceptions)
	}
}

func generateDescription(generateDescriptorData GenerateDescriptionData) string {
	var parts []string

	// Metric
	if generateDescriptorData.AffectedMetric != "" && generateDescriptorData.Delta != "" {
		parts = append(parts, fmt.Sprintf("**Metric:**\n%s (Delta: %s)", generateDescriptorData.AffectedMetric, generateDescriptorData.Delta))
	}

	// Test
	if generateDescriptorData.AffectedTest != "" {
		parts = append(parts, "**Test:**\n"+generateDescriptorData.AffectedTest)
	}

	// Test method
	if generateDescriptorData.TestMethod != nil && *generateDescriptorData.TestMethod != "" {
		parts = append(parts, "**Test method name:**\n"+*generateDescriptorData.TestMethod)
	}

	// Build
	if generateDescriptorData.BuildLink != "" {
		parts = append(parts, fmt.Sprintf("**Build:**\n[build link](%s)", generateDescriptorData.BuildLink))
	}

	// Changes in space
	if generateDescriptorData.Changes != "" {
		parts = append(parts, fmt.Sprintf("**Changes in space:**\n[space link](%s)", generateDescriptorData.Changes))
	}

	// Idea logs and snapshots
	if generateDescriptorData.TestType == "intellij" || generateDescriptorData.TestType == "intellij_dev" {
		logs := "**Idea logs, screenshots, thread dumps etc:**\nCurrent: [logs-current.zip](logs-current.zip)"
		snapshots := "**Snapshots:**\nCurrent: [snapshots-current.zip](snapshots-current.zip)"
		if generateDescriptorData.Kind != "exception" {
			logs += "\nBefore: [logs-before.zip](logs-before.zip)"
			snapshots += "\nBefore: [snapshots-before.zip](snapshots-before.zip)"
		}
		parts = append(parts, logs, snapshots)
	}

	if generateDescriptorData.TestType == "perfUnitTests" {
		snapshots := "**Snapshots:**\nCurrent: [log-current.zip](log-current.zip)"
		snapshots += "\nBefore: [log-before.zip](log-before.zip)"
		parts = append(parts, snapshots)
	}

	// Dashboard
	if generateDescriptorData.DashboardLink != "" {
		parts = append(parts, fmt.Sprintf("**Chart:**\n[link to test chart](%s)", generateDescriptorData.DashboardLink), "![](dashboard.png)")
	}

	// Stacktrace or test history
	if generateDescriptorData.Kind == "exception" {
		if generateDescriptorData.StackTrace != "" {
			parts = append(parts, fmt.Sprintf("**Stacktrace:**\n```%s```", generateDescriptorData.StackTrace))
		}
	} else {
		if generateDescriptorData.TestHistoryUrl != nil && *generateDescriptorData.TestHistoryUrl != "" {
			parts = append(parts, fmt.Sprintf("**Test history:**\n[test history link](%s)", *generateDescriptorData.TestHistoryUrl))
		}
	}

	description := strings.Join(parts, "\n\n")
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

func logError(message string, err error, exceptions *[]string) {
	slog.Error(message, "error", err)
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
