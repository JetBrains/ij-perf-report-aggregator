package meta

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TeamCityClient struct {
	teamCityURL string
	authToken   string
}

type Test struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
	Href string `xml:"href,attr"`
}

type Tests struct {
	Count int    `xml:"count,attr"`
	Test  []Test `xml:"test"`
}

type TeamCityAttachmentInfo struct {
	CurrentBuildId  int  `json:"currentBuildId"`
	PreviousBuildId *int `json:"previousBuildId"`
}

type Files struct {
	XMLName xml.Name `xml:"files"`
	Files   []File   `xml:"file"`
}

type File struct {
	Name string `xml:"name,attr"`
}

func NewTeamCityClient(teamCityURL, authToken string) *TeamCityClient {
	return &TeamCityClient{
		teamCityURL: teamCityURL,
		authToken:   authToken,
	}
}

func (client *TeamCityClient) makeRequest(ctx context.Context, endpoint string, headers map[string]string) (*http.Response, error) {
	myUrl := fmt.Sprintf("%s%s", client.teamCityURL, endpoint)
	fmt.Printf("Requesting URL: %s\n", myUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, myUrl, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+client.authToken)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, fmt.Errorf("request failed: %s", resp.Status)
	}

	return resp, nil
}

func (client *TeamCityClient) getArtifactChildren(ctx context.Context, buildId int, testName string) ([]string, error) {
	endpoint := fmt.Sprintf("/app/rest/builds/id:%d/artifacts/children/%s", buildId, testName)
	resp, err := client.makeRequest(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var files Files
	err = xml.Unmarshal(body, &files)
	if err != nil {
		return nil, err
	}

	var children []string
	for _, child := range files.Files {
		children = append(children, child.Name)
	}

	return children, nil
}

func (client *TeamCityClient) downloadArtifact(ctx context.Context, buildId int, filePath string) ([]byte, error) {
	endpoint := fmt.Sprintf("/app/rest/builds/id:%d/artifacts/content/%s", buildId, filePath)
	resp, err := client.makeRequest(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *TeamCityClient) getTestHistoryUrl(ctx context.Context, testName string) (string, error) {
	endpoint := "/app/rest/tests?locator=name:" + url.QueryEscape(testName)
	resp, err := client.makeRequest(ctx, endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var tests Tests
	err = xml.Unmarshal(body, &tests)
	if err != nil {
		return "", err
	}

	if len(tests.Test) == 0 {
		return "", fmt.Errorf("test not found: %s", testName)
	}

	return fmt.Sprintf("%s/test/%s/?currentProjectId=ijplatform", client.teamCityURL, tests.Test[0].ID), nil
}
