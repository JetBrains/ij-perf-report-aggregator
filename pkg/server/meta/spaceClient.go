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
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
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

type SpaceUploadAttachmentsRequest struct {
	UploadAttachmentsRequest
}

type SpaceUploadAttachmentsResponse struct {
	Uploads    map[int][]string `json:"uploads"`
	Exceptions map[int][]string `json:"exceptions"`
}

func CreatePostSpaceUploadAttachments(metaDb *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var response SpaceUploadAttachmentsResponse
		var params SpaceUploadAttachmentsRequest
		decoder := json.NewDecoder(request.Body)
		defer request.Body.Close()
		err := decoder.Decode(&params)
		if err != nil {
			http.Error(writer, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}

		var mu sync.Mutex
		ctx := request.Context()

		config := UploadConfig{
			UploadChartPng: func(ctx context.Context, chartData []byte) error {
				folder := fmt.Sprintf("analyses/%d", params.TeamCityAttachmentInfo.CurrentBuildId)
				return spacePackagesClient.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", folder, "dashboard.png", bytes.NewReader(chartData))
			},
			UploadArtifact: func(ctx context.Context, artifact UploadArtifact) error {
				folder := fmt.Sprintf("analyses/%d", artifact.BuildId)
				return spacePackagesClient.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", folder, artifact.FileName, artifact.Body)
			},
			OnError: func(buildId int, message string, err error) {
				slog.Error(message, "error", err)
				mu.Lock()
				defer mu.Unlock()
				if response.Exceptions == nil {
					response.Exceptions = make(map[int][]string)
				}
				response.Exceptions[buildId] = append(response.Exceptions[buildId],
					fmt.Sprintf("Message: %s. Error: %s", message, err.Error()))
			},
			OnSuccess: func(buildId int, fileName string) {
				mu.Lock()
				defer mu.Unlock()
				if response.Uploads == nil {
					response.Uploads = make(map[int][]string)
				}
				response.Uploads[buildId] = append(response.Uploads[buildId], fileName)
			},
			SkipPostfix: true,
		}

		ProcessAndUploadArtifacts(ctx, params.UploadAttachmentsRequest, config)

		if err := recordSpaceUploadedArtifacts(ctx, metaDb, response); err != nil {
			slog.Error("failed to record space uploaded artifacts", "error", err)
		}

		_ = marshalAndWriteIssueResponse(writer, response)
	}
}

func recordSpaceUploadedArtifacts(
	ctx context.Context,
	metaDb *pgxpool.Pool,
	resp SpaceUploadAttachmentsResponse,
) error {
	buildIds := make(map[int]struct{}, len(resp.Uploads)+len(resp.Exceptions))
	for id := range resp.Uploads {
		buildIds[id] = struct{}{}
	}
	for id := range resp.Exceptions {
		buildIds[id] = struct{}{}
	}

	const stmt = `
		INSERT INTO space_uploaded_artifacts (build_id, uploaded_files, success)
		VALUES ($1, $2, $3)
		ON CONFLICT (build_id) DO UPDATE
		SET uploaded_files = EXCLUDED.uploaded_files,
		    success        = EXCLUDED.success`

	for buildId := range buildIds {
		files := resp.Uploads[buildId]
		if files == nil {
			files = []string{}
		}
		success := len(resp.Exceptions[buildId]) == 0
		if _, err := metaDb.Exec(ctx, stmt, strconv.Itoa(buildId), files, success); err != nil {
			return fmt.Errorf("insert space_uploaded_artifacts for build %d: %w", buildId, err)
		}
	}
	return nil
}
