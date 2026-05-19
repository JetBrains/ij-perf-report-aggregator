package meta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"
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
		response := SpaceUploadAttachmentsResponse{
			Uploads:    map[int][]string{},
			Exceptions: map[int][]string{},
		}
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

		skip, alreadyUploaded := computeBuildsToSkip(ctx, metaDb, params.TeamCityAttachmentInfo, params.ProjectName)
		if len(alreadyUploaded) > 0 {
			response.Uploads = alreadyUploaded
		}

		config := UploadConfig{
			UploadArtifact: func(ctx context.Context, artifact UploadArtifact) error {
				folder := fmt.Sprintf("analyses/%d/%s", artifact.BuildId, params.ProjectName)
				return spacePackagesClient.UploadFile(ctx, "platform-test-automation", "performance-regression-llm-analysis", folder, artifact.FileName, artifact.Body)
			},
			OnError: func(buildId int, message string, err error) {
				slog.Error(message, "error", err)
				mu.Lock()
				defer mu.Unlock()
				response.Exceptions[buildId] = append(response.Exceptions[buildId],
					fmt.Sprintf("Message: %s. Error: %s", message, err.Error()))
			},
			OnSuccess: func(buildId int, fileName string) {
				mu.Lock()
				defer mu.Unlock()
				response.Uploads[buildId] = append(response.Uploads[buildId], fileName)
			},
			SkipPostfix:  true,
			BuildsToSkip: skip,
		}

		ProcessAndUploadTeamcityArtifacts(ctx, params.UploadAttachmentsRequest, config)

		if err := recordSpaceUploadedArtifacts(ctx, metaDb, params.ProjectName, response); err != nil {
			slog.Error("failed to record space uploaded artifacts", "error", err)
		}

		_ = marshalAndWriteIssueResponse(writer, response)
	}
}

func recordSpaceUploadedArtifacts(
	ctx context.Context,
	metaDb *pgxpool.Pool,
	project string,
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
		INSERT INTO space_uploaded_artifacts (build_id, project, uploaded_files, success)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (build_id, project) DO UPDATE
		SET uploaded_files = EXCLUDED.uploaded_files,
		    success        = EXCLUDED.success`

	for buildId := range buildIds {
		files := resp.Uploads[buildId]
		if files == nil {
			files = []string{}
		}
		success := len(resp.Exceptions[buildId]) == 0
		if _, err := metaDb.Exec(ctx, stmt, strconv.Itoa(buildId), project, files, success); err != nil {
			return fmt.Errorf("insert space_uploaded_artifacts for build %d project %s: %w", buildId, project, err)
		}
	}
	return nil
}

type TcBuildSpaceUpload struct {
	Files         []string
	SuccessStatus bool
}

func getSpaceUploads(ctx context.Context, metaDb *pgxpool.Pool, buildId int, project string) (*TcBuildSpaceUpload, error) {
	var upload TcBuildSpaceUpload
	err := metaDb.QueryRow(ctx,
		"SELECT success, uploaded_files FROM space_uploaded_artifacts WHERE build_id = $1 AND project = $2",
		strconv.Itoa(buildId), project).Scan(&upload.SuccessStatus, &upload.Files)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pgx.ErrNoRows
	}
	if err != nil {
		return nil, err
	}
	return &upload, nil
}

func computeBuildsToSkip(ctx context.Context, metaDb *pgxpool.Pool, info TeamCityAttachmentInfo, project string) (map[int]struct{}, map[int][]string) {
	buildIds := []int{info.CurrentBuildId}
	if info.PreviousBuildId != nil {
		buildIds = append(buildIds, *info.PreviousBuildId)
	}

	skip := make(map[int]struct{})
	alreadyUploaded := make(map[int][]string)
	for _, buildId := range buildIds {
		upload, err := getSpaceUploads(ctx, metaDb, buildId, project)
		if errors.Is(err, pgx.ErrNoRows) {
			continue
		}
		if err != nil {
			slog.Warn("failed to query previous space upload status", "buildId", buildId, "project", project, "error", err)
			continue
		}
		if upload.SuccessStatus {
			skip[buildId] = struct{}{}
			alreadyUploaded[buildId] = upload.Files
			slog.Info("skipping space upload for build (already uploaded successfully)", "buildId", buildId, "project", project, "files", upload.Files)
		}
	}
	return skip, alreadyUploaded
}
