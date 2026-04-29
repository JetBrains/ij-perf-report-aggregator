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
		var params UploadAttachmentsRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		err := decoder.Decode(&params)
		if err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		var uploadsMu sync.Mutex
		ctx := request.Context()
		remoteFolder := "analyses/" + params.IssueId

		config := UploadConfig{
			UploadChartPng: func(ctx context.Context, chartData []byte) error {
				return spacePackagesClient.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", remoteFolder, "dashboard.png", bytes.NewReader(chartData))
			},
			UploadArtifact: func(ctx context.Context, artifact UploadArtifact) error {
				return spacePackagesClient.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", remoteFolder, artifact.FileName, artifact.Body)
			},
			OnError: func(message string, err error) {
				slog.Error(message, "error", err)
			},
			OnSuccess: func(fileName string) {
				uploadsMu.Lock()
				defer uploadsMu.Unlock()
				response.Uploads = append(response.Uploads, fileName)
			},
		}

		ProcessAndUploadArtifacts(ctx, params, config)
		_ = marshalAndWriteIssueResponse(writer, response)
	}
}
