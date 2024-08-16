package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

type YoutrackClient struct {
	youTrackUrl       string
	youtrackToken     string
	availableProjects []YoutrackProject
}

type YoutrackProject struct {
	ID string `json:"id"`
}

type YoutrackIssue struct {
	ID         string `json:"id"`
	IDReadable string `json:"idReadable"`
}

type CustomField struct {
	Name  string `json:"name"`
	Type  string `json:"$type"`
	Value struct {
		Name string `json:"name"`
	} `json:"value"`
}

type CreateIssueInfo struct {
	Summary      string          `json:"summary"`
	Description  string          `json:"description"`
	Project      YoutrackProject `json:"project"`
	CustomFields []CustomField   `json:"customFields"`
}

func NewYoutrackClient(youTrackUrl, youtrackToken string) *YoutrackClient {
	return &YoutrackClient{
		youTrackUrl:   youTrackUrl,
		youtrackToken: youtrackToken,
	}
}

func (client *YoutrackClient) fetchFromYouTrack(endpoint string, method string, body io.Reader, headers map[string]string) ([]byte, error) {
	youtrackUrl := fmt.Sprintf("%s%s", client.youTrackUrl, endpoint)
	log.Printf("Youtrack url: %s\n", youtrackUrl)

	req, err := http.NewRequest(method, youtrackUrl, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.youtrackToken))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (client *YoutrackClient) CreateIssue(info CreateIssueInfo) (*YoutrackIssue, error) {
	body, err := json.Marshal(info)
	if err != nil {
		return nil, fmt.Errorf("error marshaling info: %w", err)
	}

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	responseData, err := client.fetchFromYouTrack("/api/issues?fields=id,idReadable", "POST", bytes.NewBuffer(body), headers)
	if err != nil {
		return nil, fmt.Errorf("error creating info: %w", err)
	}

	var issue YoutrackIssue
	if err := json.Unmarshal(responseData, &issue); err != nil {
		return nil, fmt.Errorf("error unmarshalling issue: %w", err)
	}

	return &issue, nil
}

func (client *YoutrackClient) SearchIssuesByLabel(label string) ([]YoutrackIssue, error) {
	encodedLabel := url.QueryEscape(fmt.Sprintf("{%s}", label))
	responseData, err := client.fetchFromYouTrack(fmt.Sprintf("/api/issues?query=tag:%s", encodedLabel), "GET", nil, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, fmt.Errorf("error fetching issues: %w", err)
	}

	var issues []YoutrackIssue
	if err := json.Unmarshal(responseData, &issues); err != nil {
		return nil, fmt.Errorf("error unmarshalling issues: %w", err)
	}

	return issues, nil
}

func (client *YoutrackClient) GetCustomFields(projectId string) ([]CustomField, error) {
	responseData, err := client.fetchFromYouTrack(fmt.Sprintf("/api/admin/projects/%s/customFields?fields=value(name)", projectId), "GET", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching custom fields: %w", err)
	}

	var customFields []CustomField
	if err := json.Unmarshal(responseData, &customFields); err != nil {
		return nil, fmt.Errorf("error unmarshalling custom fields: %w", err)
	}

	return customFields, nil
}

func (client *YoutrackClient) GetProjects() ([]YoutrackProject, error) {
	if len(client.availableProjects) == 0 {
		responseData, err := client.fetchFromYouTrack("/api/admin/projects?fields=id,name", "GET", nil, map[string]string{"Content-Type": "application/json"})
		if err != nil {
			return nil, fmt.Errorf("error fetching projects: %w", err)
		}

		if err := json.Unmarshal(responseData, &client.availableProjects); err != nil {
			return nil, fmt.Errorf("error unmarshalling projects: %w", err)
		}
	}

	return client.availableProjects, nil
}

func (client *YoutrackClient) UploadAttachment(issueId string, file io.Reader, fileName string) error {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return fmt.Errorf("error creating form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error copying file data: %w", err)
	}

	writer.Close()

	_, err = client.fetchFromYouTrack(fmt.Sprintf("/api/issues/%s/attachments", issueId), "POST", &body, map[string]string{"Content-Type": writer.FormDataContentType()})
	if err != nil {
		return fmt.Errorf("error uploading attachment: %w", err)
	}

	return nil
}
