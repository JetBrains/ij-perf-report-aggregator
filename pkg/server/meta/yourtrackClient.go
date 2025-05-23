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
	Name  string      `json:"name"`
	Type  string      `json:"$type"`
	Value interface{} `json:"value"`
}

type Visibility struct {
	PermittedGroups []auth.YTUser `json:"permittedGroups"`
	PermittedUsers  []auth.YTUser `json:"permittedUsers"`
	Type            string        `json:"$type"`
}

type CreateIssueInfo struct {
	Summary      string          `json:"summary"`
	Description  string          `json:"description"`
	Project      YoutrackProject `json:"project"`
	Reporter     *auth.YTUser    `json:"reporter,omitempty"`
	Visibility   Visibility      `json:"visibility"`
	CustomFields []CustomField   `json:"customFields"`
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

	var jsonData map[string]interface{}
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

func (client *YoutrackClient) UploadAttachment(ctx context.Context, issueId string, file []byte, fileName string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("error creating form file: %w", err)
	}

	_, err = io.Copy(part, bytes.NewReader(file))
	if err != nil {
		return fmt.Errorf("error copying file data: %w", err)
	}

	writer.Close()

	err = client.waitIssueIsCreated(ctx, issueId)
	if err != nil {
		return fmt.Errorf("issue was not created: %w", err)
	}

	_, err = client.fetchFromYouTrack(ctx, fmt.Sprintf("/api/issues/%s/attachments", issueId), "POST", &body, map[string]string{"Content-Type": writer.FormDataContentType()})
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
