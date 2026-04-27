package meta

import (
	"fmt"
	"slices"
	"strings"
)

type UploadAttachmentsRequest struct {
	IssueId                string                 `json:"issueId"`
	TeamCityAttachmentInfo TeamCityAttachmentInfo `json:"teamcityAttachmentInfo"`
	AffectedTest           string                 `json:"affectedTest"`
	ChartPng               *[]byte                `json:"chartPng"`
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
