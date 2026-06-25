package meta

import (
	"context"
	"io"
	"regexp"
	"slices"
	"strings"
	"sync"
)

var consecutiveHyphensRegex = regexp.MustCompile(`-+`)

// replicates sanitizing logic on publishing artifacts by driver
// see com.intellij.ide.starter.ci.teamcity.TeamCityClient#publishTeamCityArtifacts
func replaceSpecialCharactersWithHyphens(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == '.' || r == '/' || r == '\\':
			b.WriteRune(r)
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteRune('-')
		}
	}
	return strings.Trim(consecutiveHyphensRegex.ReplaceAllString(b.String(), "-"), "-")
}

type UploadAttachmentsRequest struct {
	TeamCityAttachmentInfo TeamCityAttachmentInfo `json:"teamcityAttachmentInfo"`
	ProjectName            string                 `json:"projectName"`
	TestType               string                 `json:"testType"`
	MethodName             *string                `json:"methodName,omitempty"`
}

type teamCityArtifact struct {
	BuildId      int
	ArtifactPath string
}

type UploadArtifact struct {
	BuildId       int
	FileName      string
	Body          io.Reader
	ContentLength int64
}

type UploadConfig struct {
	UploadArtifact func(ctx context.Context, artifact UploadArtifact) error
	OnError        func(buildId int, message string, err error)
	OnSuccess      func(buildId int, fileName string)
	SkipPostfix    bool
	BuildsToSkip   map[int]struct{}
}

type artifactCollector interface {
	getArtifactsPaths(params UploadAttachmentsRequest) []string
	checkArtifact(artifactName string) bool
}

type fleetStartupCollector struct{}

func (f fleetStartupCollector) getArtifactsPaths(UploadAttachmentsRequest) []string {
	return []string{""}
}

func (f fleetStartupCollector) checkArtifact(artifactName string) bool {
	return strings.HasSuffix(artifactName, "fleet.fahrplan.json")
}

type fleetPerfTestCollector struct{}

func (f fleetPerfTestCollector) getArtifactsPaths(params UploadAttachmentsRequest) []string {
	if params.MethodName == nil {
		// fallback to situations when MethodName is not provided
		return []string{""}
	}
	methodName := *params.MethodName
	artifactPath := strings.ReplaceAll(methodName, "#", "/")
	testMethodName := methodName
	if _, after, ok := strings.Cut(methodName, "#"); ok {
		testMethodName = after
	}
	return []string{"logs/" + artifactPath, "metrics/" + testMethodName}
}

func (f fleetPerfTestCollector) checkArtifact(artifactName string) bool {
	switch artifactName {
	// logs.zip is fallback for situation when MethodName is not provided and we just download logs.zip from the root
	case "logs.zip", "fsdaemon.log", "fleet.log", "spans.json", "fleet.test.json":
		return true
	default:
		return false
	}
}

type perfUnitTestCollector struct{}

func (f perfUnitTestCollector) getArtifactsPaths(params UploadAttachmentsRequest) []string {
	return []string{params.ProjectName}
}

func (f perfUnitTestCollector) checkArtifact(artifactName string) bool {
	switch artifactName {
	case "log.zip", "testlog.zip", "metrics.performance.json":
		return true
	default:
		return false
	}
}

type perfintCollector struct{}

func (f perfintCollector) getArtifactsPaths(params UploadAttachmentsRequest) []string {
	return []string{replaceSpecialCharactersWithHyphens(params.ProjectName)}
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
	name, ext := filename, ""
	if dot := strings.LastIndex(filename, "."); dot > 0 {
		name, ext = filename[:dot], filename[dot+1:]
	}

	nameParts := strings.Split(name, "-")
	prefix := nameParts[0]
	if slices.Contains(nameParts[1:], "frontend") {
		prefix += "-frontend"
	}

	result := prefix + "-" + suffix
	if ext != "" {
		result += "." + ext
	}
	return result
}

// ProcessAndUploadTeamcityArtifacts handles the common logic for fetching artifacts
// from TeamCity and uploading them via the provided config
func ProcessAndUploadTeamcityArtifacts(ctx context.Context, params UploadAttachmentsRequest, config UploadConfig) {
	var uploadWg sync.WaitGroup

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
		if _, skip := config.BuildsToSkip[buildId]; skip {
			continue
		}
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
	var attachmentPostfix string
	if !config.SkipPostfix {
		if index == 0 {
			attachmentPostfix = "current"
		} else {
			attachmentPostfix = "before"
		}
	}
	for _, testArtifactPath := range collector.getArtifactsPaths(params) {
		children, err := teamCityClient.getArtifactChildren(ctx, buildId, testArtifactPath)
		if err != nil {
			config.OnError(buildId, "Failed to get teamcity artifact children", err)
			continue
		}

		var filteredChildren []string
		for _, str := range children {
			if collector.checkArtifact(str) {
				filteredChildren = append(filteredChildren, str)
			}
		}

		for _, str := range filteredChildren {
			fileName := str
			if !config.SkipPostfix {
				fileName = getAttachmentName(str, attachmentPostfix)
			}
			artifact := teamCityArtifact{
				BuildId:      buildId,
				ArtifactPath: testArtifactPath + "/" + str,
			}
			uploadWg.Go(func() {
				resp, err := teamCityClient.getDownloadArtifactResponse(ctx, artifact.BuildId, artifact.ArtifactPath)
				if err != nil {
					config.OnError(artifact.BuildId, "Failed to download artifact from TeamCity", err)
					return
				}
				defer resp.Body.Close()

				err = config.UploadArtifact(ctx, UploadArtifact{
					BuildId:       artifact.BuildId,
					FileName:      fileName,
					Body:          resp.Body,
					ContentLength: resp.ContentLength,
				})
				if err != nil {
					config.OnError(artifact.BuildId, "Failed to upload attachment", err)
				} else {
					config.OnSuccess(artifact.BuildId, fileName)
				}
			})
		}
	}
}
