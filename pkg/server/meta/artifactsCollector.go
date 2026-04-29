package meta

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
)

type UploadAttachmentsRequest struct {
	IssueId                string                 `json:"issueId"`
	TeamCityAttachmentInfo TeamCityAttachmentInfo `json:"teamcityAttachmentInfo"`
	AffectedTest           string                 `json:"affectedTest"`
	ChartPng               []byte                 `json:"chartPng"`
	TestType               string                 `json:"testType"`
}

type UploadAttachmentsResponse struct {
	Uploads    []string `json:"uploads"`
	Exceptions []string `json:"exceptions"`
}

type teamCityArtifact struct {
	BuildId      int
	ArtifactPath string
	FileName     string
}

type UploadArtifact struct {
	FileName      string
	Body          io.Reader
	ContentLength int64
}

type UploadConfig struct {
	UploadChartPng func(ctx context.Context, chartData []byte) error
	UploadArtifact func(ctx context.Context, artifact UploadArtifact) error
	OnError        func(message string, err error)
	OnSuccess      func(fileName string)
}

type artifactCollector interface {
	getArtifactsPath(params UploadAttachmentsRequest) string
	checkArtifact(artifactName string) bool
}

type fleetStartupCollector struct{}

func (f fleetStartupCollector) getArtifactsPath(UploadAttachmentsRequest) string {
	return ""
}

func (f fleetStartupCollector) checkArtifact(artifactName string) bool {
	return strings.HasSuffix(artifactName, "fleet.fahrplan.json")
}

type fleetPerfTestCollector struct{}

func (f fleetPerfTestCollector) getArtifactsPath(UploadAttachmentsRequest) string {
	return ""
}

func (f fleetPerfTestCollector) checkArtifact(artifactName string) bool {
	return artifactName == "logs.zip"
}

type perfUnitTestCollector struct{}

func (f perfUnitTestCollector) getArtifactsPath(params UploadAttachmentsRequest) string {
	return params.AffectedTest
}

func (f perfUnitTestCollector) checkArtifact(artifactName string) bool {
	return artifactName == "log.zip"
}

type perfintCollector struct{}

func (f perfintCollector) getArtifactsPath(params UploadAttachmentsRequest) string {
	return strings.ReplaceAll(params.AffectedTest, "_", "-")
}

func (f perfintCollector) checkArtifact(artifactName string) bool {
	prefixes := []string{"logs-", "snapshots-"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(artifactName, prefix) {
			return true
		}
	}
	if artifactName == "metrics.performance.json" {
		return true
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

func getAttachmentName(filename, suffix string) string {
	// Handle metrics.performance.json specially
	if strings.HasPrefix(filename, "metrics.performance") && strings.HasSuffix(filename, ".json") {
		return fmt.Sprintf("metrics.performance-%s.json", suffix)
	}

	parts := strings.Split(filename, ".")
	if len(parts) != 2 {
		return filename
	}

	nameWithoutExt := parts[0]
	ext := parts[1]

	nameParts := strings.Split(nameWithoutExt, "-")

	prefix := nameParts[0]
	if slices.Contains(nameParts[1:], "frontend") {
		prefix += "-frontend"
	}

	updatedName := prefix + "-" + suffix
	return fmt.Sprintf("%s.%s", updatedName, ext)
}

// ProcessAndUploadArtifacts handles the common logic for fetching artifacts
// from TeamCity and uploading them via the provided config
func ProcessAndUploadArtifacts(ctx context.Context, params UploadAttachmentsRequest, config UploadConfig) {
	var uploadWg sync.WaitGroup

	if params.ChartPng != nil && config.UploadChartPng != nil {
		uploadWg.Go(func() {
			err := config.UploadChartPng(ctx, params.ChartPng)
			if err != nil {
				config.OnError("Failed to upload chart PNG", err)
			} else {
				config.OnSuccess("dashboard.png")
			}
		})
	}

	builds := []int{params.TeamCityAttachmentInfo.CurrentBuildId}
	if params.TeamCityAttachmentInfo.PreviousBuildId != nil {
		builds = append(builds, *params.TeamCityAttachmentInfo.PreviousBuildId)
	}

	collector := getArtifactCollector(params.TestType)
	if collector == nil {
		uploadWg.Wait()
		return
	}

	var wg sync.WaitGroup
	for index, buildId := range builds {
		wg.Go(func() {
			processBuildsArtifacts(ctx, buildId, index, params, collector, config, &uploadWg)
		})
	}
	wg.Wait()
	uploadWg.Wait()
}

func processBuildsArtifacts(
	ctx context.Context,
	buildId int,
	index int,
	params UploadAttachmentsRequest,
	collector artifactCollector,
	config UploadConfig,
	uploadWg *sync.WaitGroup,
) {
	testArtifactPath := collector.getArtifactsPath(params)

	children, err := teamCityClient.getArtifactChildren(ctx, buildId, testArtifactPath)
	if err != nil {
		config.OnError("Failed to get teamcity artifact children", err)
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

	for _, str := range filteredChildren {
		artifact := teamCityArtifact{
			BuildId:      buildId,
			ArtifactPath: testArtifactPath + "/" + str,
			FileName:     getAttachmentName(str, attachmentPostfix),
		}
		uploadWg.Go(func() {
			resp, err := teamCityClient.getDownloadArtifactResponse(ctx, artifact.BuildId, artifact.ArtifactPath)
			if err != nil {
				config.OnError("Failed to download artifact from TeamCity", err)
				return
			}
			defer resp.Body.Close()

			err = config.UploadArtifact(ctx, UploadArtifact{
				FileName:      artifact.FileName,
				Body:          resp.Body,
				ContentLength: resp.ContentLength,
			})
			if err != nil {
				config.OnError("Failed to upload attachment", err)
			} else {
				config.OnSuccess(artifact.FileName)
			}
		})
	}
}
