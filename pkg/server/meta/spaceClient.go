package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
)

var spacePackagesClient = NewSpacePackagesClient("https://packages.jetbrains.team", os.Getenv("SPACE_TOKEN"))

type SpacePackagesClient struct {
	spaceUrl   string
	spaceToken string
}

func NewSpacePackagesClient(spaceUrl, spaceToken string) *SpacePackagesClient {
	return &SpacePackagesClient{
		spaceUrl:   spaceUrl,
		spaceToken: spaceToken,
	}
}

func (client *SpacePackagesClient) UploadFile(ctx context.Context, project string, packageName string, remoteFolder string, fileName string, body io.Reader) error {
	endpoint := fmt.Sprintf("/files/p/%s/%s/%s/%s", project, packageName, remoteFolder, fileName)

	_, err := client.doRequest(ctx, endpoint, http.MethodPut, body, nil)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}

	return nil
}

func (client *SpacePackagesClient) doRequest(ctx context.Context, endpoint string, method string, body io.Reader, headers map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%s%s", client.spaceUrl, endpoint)
	slog.Info("Space request", "url", url, "method", method)

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.spaceToken)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("request failed with status: %s. Body: %s", resp.Status, bodyString)
	}

	return bodyBytes, nil
}

func CreatePostSpaceUploadAttachments() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var response UploadAttachmentsResponse
		body := request.Body
		all, err := io.ReadAll(body)
		if err != nil {
			handleError(writer, "cannot read body", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, response)
			return
		}
		defer body.Close()

		var params UploadAttachmentsRequest
		if err = json.Unmarshal(all, &params); err != nil {
			handleError(writer, "cannot unmarshal parameters", err, &response.Exceptions)
			_ = marshalAndWriteIssueResponse(writer, response)
			return
		}

		var uploadsMu sync.Mutex
		addSuccess := func(fileName string) {
			uploadsMu.Lock()
			defer uploadsMu.Unlock()
			response.Uploads = append(response.Uploads, fileName)
		}

		logError := func(message string, err error) {
			slog.Error(message, "error", err)
		}

		var uploadWg sync.WaitGroup
		ctx := request.Context()

		remoteFolder := "analyses/" + params.IssueId

		if params.ChartPng != nil {
			uploadWg.Go(func() {
				err := spacePackagesClient.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", remoteFolder, "dashboard.png", bytes.NewReader(*params.ChartPng))
				if err != nil {
					logError("Failed to upload chart PNG to Space", err)
				} else {
					addSuccess("dashboard.png")
				}
			})
		}

		builds := []int{params.TeamCityAttachmentInfo.CurrentBuildId}
		if params.TeamCityAttachmentInfo.PreviousBuildId != nil {
			builds = append(builds, *params.TeamCityAttachmentInfo.PreviousBuildId)
		}

		collector := getArtifactCollector(params.TestType)

		if collector != nil {
			var wg sync.WaitGroup
			for index, buildId := range builds {
				wg.Go(func() {
					testArtifactPath := collector.getArtifactsPath(params)

					children, err := teamCityClient.getArtifactChildren(ctx, buildId, testArtifactPath)
					if err != nil {
						logError("Failed to get teamcity artifact children", err)
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
								logError("Failed to download artifact from TeamCity", err)
								return
							}
							defer resp.Body.Close()

							err = spacePackagesClient.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", remoteFolder, artifact.FileName, resp.Body)
							if err != nil {
								logError("Failed to upload attachment to Space", err)
							} else {
								addSuccess(artifact.FileName)
							}
						})
					}
				})
			}
			wg.Wait()
		}

		uploadWg.Wait()

		_ = marshalAndWriteIssueResponse(writer, response)
	}
}
