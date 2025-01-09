package meta

import (
	"bytes"
	"context"
	"encoding/json"
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

// Build represents the root element of the XML structure
type Build struct {
	XMLName    xml.Name   `xml:"build"`
	BuildType  BuildType  `xml:"buildType"`
	Properties Properties `xml:"properties"`
}

// BuildType represents the buildType element
type BuildType struct {
	ID string `xml:"id,attr"`
}

// Properties represents the properties element
type Properties struct {
	Property []Property `xml:"property"`
}

// Property represents an individual property element
type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type BuildResponse struct {
	XMLName xml.Name `xml:"build"`
	WebURL  string   `xml:"webUrl,attr"`
}

type BuildInfo struct {
	BuildTypeId string `json:"buildTypeId"`
}

type Change struct {
	Version string `json:"version"`
}

type Changes struct {
	Change []Change `json:"change"`
}

func (client *TeamCityClient) getBuildType(ctx context.Context, buildID string) (string, error) {
	res, err := client.makeRequest(ctx, "/app/rest/builds/id:"+buildID, map[string]string{"Accept": "application/json"})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var build BuildInfo
	if err := json.NewDecoder(res.Body).Decode(&build); err != nil {
		return "", fmt.Errorf("failed to decode changes response: %w", err)
	}

	return build.BuildTypeId, nil
}

func (client *TeamCityClient) getChanges(ctx context.Context, buildID string) (*CommitRevisions, error) {
	res, err := client.makeRequest(ctx, "/app/rest/changes?locator=build:(id:"+buildID+")&count=10000", map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var changes Changes
	if err := json.NewDecoder(res.Body).Decode(&changes); err != nil {
		return nil, fmt.Errorf("failed to decode changes response: %w", err)
	}

	if len(changes.Change) == 0 {
		return nil, fmt.Errorf("no changes found for build %s", buildID)
	}

	revisions := &CommitRevisions{
		LastCommit:  changes.Change[0].Version,
		FirstCommit: changes.Change[len(changes.Change)-1].Version,
	}

	return revisions, nil
}

func (client *TeamCityClient) startBuild(ctx context.Context, buildId string, params map[string]string) (*string, error) {
	endpoint := "/app/rest/buildQueue"

	properties := Properties{
		Property: make([]Property, 0, len(params)),
	}

	for key, value := range params {
		properties.Property = append(properties.Property, Property{
			Name:  key,
			Value: value,
		})
	}

	build := Build{
		BuildType:  BuildType{ID: buildId},
		Properties: properties,
	}

	myUrl := fmt.Sprintf("%s%s", client.teamCityURL, endpoint)
	fmt.Printf("Requesting URL: %s\n", myUrl)

	buildXml, err := xml.Marshal(build)
	if err != nil {
		return nil, fmt.Errorf("error marshalling build %v: %w", build, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, myUrl, bytes.NewBuffer(buildXml))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.authToken)
	req.Header.Set("Content-Type", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	body := resp.Body
	all, err := io.ReadAll(body)
	defer body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var buildResponse BuildResponse
	err = xml.Unmarshal(all, &buildResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling XML: %w", err)
	}
	return &buildResponse.WebURL, nil
}
