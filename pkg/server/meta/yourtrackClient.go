package meta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server/auth"
)

type YoutrackClient struct {
	youTrackUrl   string
	youtrackToken string
}

type YoutrackProject struct {
	ID string `json:"id"`
}

type YoutrackIssue struct {
	ID         string `json:"id"`
	IDReadable string `json:"idReadable"`
}

type CustomFieldValue struct {
	Name string `json:"name"`
}

type CustomField struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name"`
	Type  string `json:"$type"`
	Value any    `json:"value"`
}

type Visibility struct {
	PermittedGroups []auth.YTUser `json:"permittedGroups"`
	PermittedUsers  []auth.YTUser `json:"permittedUsers"`
	Type            string        `json:"$type"`
}

type Tag struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Type string `json:"$type"`
}

type CreateIssueInfo struct {
	Summary      string          `json:"summary"`
	Description  string          `json:"description"`
	Project      YoutrackProject `json:"project"`
	Reporter     *auth.YTUser    `json:"reporter,omitempty"`
	Visibility   Visibility      `json:"visibility"`
	CustomFields []CustomField   `json:"customFields"`
	Tags         []Tag           `json:"tags,omitempty"`
}

func NewYoutrackClient(youTrackUrl, youtrackToken string) *YoutrackClient {
	return &YoutrackClient{
		youTrackUrl:   youTrackUrl,
		youtrackToken: youtrackToken,
	}
}

func (client *YoutrackClient) CreateIssue(ctx context.Context, info CreateIssueInfo) (*YoutrackIssue, error) {
	body, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("error marshaling info: %w", err)
	}

	slog.Info("Create issue data:", "info", string(body))

	var jsonData map[string]any
	if err := json.Unmarshal(body, &jsonData); err != nil {
		slog.Error("error unmarshalling JSON:", "error", err)
	} else {
		slog.Info("JSON is valid")
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	responseData, err := client.fetchFromYouTrack(ctx, "/api/issues?fields=id,idReadable", "POST", bytes.NewBuffer(body), headers)
	if err != nil {
		return nil, fmt.Errorf("error creating info: %w", err)
	}

	var issue YoutrackIssue
	if err := json.Unmarshal(responseData, &issue); err != nil {
		return nil, fmt.Errorf("error unmarshalling issue: %w", err)
	}

	return &issue, nil
}

func (client *YoutrackClient) SearchIssuesByLabel(ctx context.Context, label string) ([]YoutrackIssue, error) {
	encodedLabel := url.QueryEscape(fmt.Sprintf("{%s}", label))
	responseData, err := client.fetchFromYouTrack(ctx, "/api/issues?query=tag:"+encodedLabel, "GET", nil, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error fetching issues: %w", err)
	}

	var issues []YoutrackIssue
	if err := json.Unmarshal(responseData, &issues); err != nil {
		return nil, fmt.Errorf("error unmarshalling issues: %w", err)
	}

	return issues, nil
}

func (client *YoutrackClient) GetCustomFields(ctx context.Context, projectId string) ([]CustomField, error) {
	responseData, err := client.fetchFromYouTrack(ctx, fmt.Sprintf("/api/admin/projects/%s/customFields?fields=value(name)", projectId), "GET", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching custom fields: %w", err)
	}

	var customFields []CustomField
	if err := json.Unmarshal(responseData, &customFields); err != nil {
		return nil, fmt.Errorf("error unmarshalling custom fields: %w", err)
	}

	return customFields, nil
}

const maxAttachmentSize = 95 * 1024 * 1024 // 95 MB

func (client *YoutrackClient) UploadAttachment(ctx context.Context, issueId string, content io.Reader, fileName string, contentLength int64) error {
	if contentLength > maxAttachmentSize {
		return fmt.Errorf("YouTrack: attachment size of the file %s exceeds maximum allowed size of 95 MB", fileName)
	}

	err := client.waitIssueIsCreated(ctx, issueId)
	if err != nil {
		return fmt.Errorf("issue was not created: %w", err)
	}

	pr, pw := io.Pipe()
	defer pr.Close()

	writer := multipart.NewWriter(pw)
	contentType := writer.FormDataContentType()

	go func() {
		part, err := writer.CreateFormFile("file", fileName)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("error creating form file: %w", err))
			return
		}

		_, err = io.Copy(part, content)
		if err != nil {
			_ = pw.CloseWithError(fmt.Errorf("error copying file data: %w", err))
			return
		}

		if err := writer.Close(); err != nil {
			_ = pw.CloseWithError(fmt.Errorf("error closing multipart writer: %w", err))
			return
		}
		_ = pw.Close()
	}()

	_, err = client.fetchFromYouTrack(ctx, fmt.Sprintf("/api/issues/%s/attachments", issueId), "POST", pr, map[string]string{"Content-Type": contentType})
	if err != nil {
		return fmt.Errorf("error uploading attachment: %w", err)
	}

	return nil
}

func (client *YoutrackClient) fetchFromYouTrack(ctx context.Context, endpoint string, method string, body io.Reader, headers map[string]string) ([]byte, error) {
	youtrackUrl := fmt.Sprintf("%s%s", client.youTrackUrl, endpoint)
	log.Printf("Youtrack url: %s\n", youtrackUrl)

	req, err := http.NewRequestWithContext(ctx, method, youtrackUrl, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.youtrackToken)
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

func (client *YoutrackClient) waitIssueIsCreated(ctx context.Context, issueId string) error {
	var responseData []byte
	var err error
	var issue YoutrackIssue

	for range 5 {
		responseData, err = client.fetchFromYouTrack(ctx, fmt.Sprintf("/api/issues/%s?fields=id,idReadable", issueId), http.MethodGet, nil, map[string]string{
			"Content-Type": "application/json",
		})

		if err == nil {
			if err = json.Unmarshal(responseData, &issue); err == nil {
				break
			}
			err = fmt.Errorf("error unmarshalling issue: %w", err)
		} else {
			err = fmt.Errorf("error fetching from YouTrack: %w", err)
		}
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return fmt.Errorf("checking whether the issue is created failed after 5 retries: %w", err)
	}
	return nil
}

func CreatePostYoutrackUploadAttachments() http.HandlerFunc {
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

		var mu sync.Mutex
		ctx := request.Context()

		config := UploadConfig{
			UploadChartPng: func(ctx context.Context, chartData []byte) error {
				contentLength := int64(len(chartData))
				return youtrackClient.UploadAttachment(ctx, params.IssueId, bytes.NewReader(chartData), "dashboard.png", contentLength)
			},
			UploadArtifact: func(ctx context.Context, artifact UploadArtifact) error {
				return youtrackClient.UploadAttachment(ctx, params.IssueId, artifact.Body, artifact.FileName, artifact.ContentLength)
			},
			OnError: func(message string, err error) {
				slog.Error(message, "error", err)
				mu.Lock()
				defer mu.Unlock()
				response.Exceptions = append(response.Exceptions,
					fmt.Sprintf("Message: %s. Error: %s", message, err.Error()))
			},
			OnSuccess: func(fileName string) {
				mu.Lock()
				defer mu.Unlock()
				response.Uploads = append(response.Uploads, fileName)
			},
		}

		ProcessAndUploadArtifacts(ctx, params, config)

		if len(response.Exceptions) > 0 {
			if len(response.Uploads) > 0 {
				writer.WriteHeader(http.StatusMultiStatus)
			} else {
				writer.WriteHeader(http.StatusInternalServerError)
			}
		}
		_ = marshalAndWriteIssueResponse(writer, response)
	}
}
