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

func (client *TeamCityClient) hasArtifacts(ctx context.Context, buildID string) (bool, error) {
	endpoint := fmt.Sprintf("/app/rest/builds/id:%s/artifacts/children/", buildID)
	resp, err := client.makeRequest(ctx, endpoint, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusNotFound {
				return false, nil
			}
		}
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var files Files
	err = xml.Unmarshal(body, &files)
	if err != nil {
		return false, err
	}

	return len(files.Files) > 0, nil
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

	children := make([]string, 0, len(files.Files))
	for _, child := range files.Files {
		children = append(children, child.Name)
	}

	return children, nil
}

func (client *TeamCityClient) getDownloadArtifactResponse(ctx context.Context, buildId int, filePath string) (*http.Response, error) {
	endpoint := fmt.Sprintf("/app/rest/builds/id:%d/artifacts/content/%s", buildId, filePath)
	resp, err := client.makeRequest(ctx, endpoint, nil)
	if err != nil {
		if resp != nil {
			resp.Body.Close()
		}
		return nil, err
	}
	return resp, nil
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
	Id      int64    `xml:"id,attr"`
	WebURL  string   `xml:"webUrl,attr"`
}

type BuildInfo struct {
	BuildTypeId string `json:"buildTypeId"`
	Number      string `json:"number"`
	BranchName  string `json:"branchName"`
	StartDate   string `json:"startDate"`
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

func (client *TeamCityClient) getBuildCounter(ctx context.Context, buildID string) (string, error) {
	res, err := client.makeRequest(ctx, "/app/rest/builds/id:"+buildID, map[string]string{"Accept": "application/json"})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var build BuildInfo
	if err := json.NewDecoder(res.Body).Decode(&build); err != nil {
		return "", fmt.Errorf("failed to decode build response: %w", err)
	}

	return build.Number, nil
}

type CommitRevisions struct {
	FirstCommit string `json:"firstCommit"`
	LastCommit  string `json:"lastCommit"`
}

func (client *TeamCityClient) getBuildInfo(ctx context.Context, buildID string) (*BuildInfo, error) {
	res, err := client.makeRequest(ctx, "/app/rest/builds/id:"+buildID, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var build BuildInfo
	if err := json.NewDecoder(res.Body).Decode(&build); err != nil {
		return nil, fmt.Errorf("failed to decode build info response: %w", err)
	}

	return &build, nil
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

// changeRef is a single VCS change as returned by the /changes endpoint with
// fields=...change(id,version). id is TeamCity's internal change id (used as the
// boundary for the sinceChange locator); version is the VCS revision (commit SHA).
type changeRef struct {
	Id      int64  `json:"id"`
	Version string `json:"version"`
}

type changesResponse struct {
	Count  int         `json:"count"`
	Change []changeRef `json:"change"`
}

// ChangesGap describes whether the current dot's bisect range covers all commits
// since the previous successful data point (dot) on the graph.
//
// When earlier builds failed/timed out and produced no metric, the commits they
// consumed exist in no dot at all. Such commits fall below the current build's
// range (which starts at currentFirstCommit), so a regression introduced by one
// of them would be unbisectable. This struct surfaces that.
type ChangesGap struct {
	// Known is false when the gap could not be determined (e.g. the previous dot's
	// revision can't be located in the configuration). Callers should not warn then.
	Known bool `json:"known"`
	// HasGap is true when there are commits between the previous dot and the start
	// of the current dot's bisect range that the range does not include.
	HasGap bool `json:"hasGap"`
	// GapCommitCount is the number of such uncovered commits.
	GapCommitCount int `json:"gapCommitCount"`
}

// getBuildRevision returns the VCS revision (commit SHA) the build was actually
// built on. Unlike the build's attributed changeset (getBuildChanges), this is
// always present even for builds TeamCity attributed 0 changes to — which happens
// routinely when builds of a configuration run out of commit order.
func (client *TeamCityClient) getBuildRevision(ctx context.Context, buildID string) (string, error) {
	endpoint := "/app/rest/builds/id:" + url.QueryEscape(buildID) + "?fields=revisions(revision(version))"
	res, err := client.makeRequest(ctx, endpoint, map[string]string{"Accept": "application/json"})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var build struct {
		Revisions struct {
			Revision []struct {
				Version string `json:"version"`
			} `json:"revision"`
		} `json:"revisions"`
	}
	if err := json.NewDecoder(res.Body).Decode(&build); err != nil {
		return "", fmt.Errorf("failed to decode build revision response: %w", err)
	}
	if len(build.Revisions.Revision) == 0 {
		return "", nil
	}
	return build.Revisions.Revision[0].Version, nil
}

// getChangeID resolves a VCS revision to TeamCity's internal change id within a
// build configuration, so it can be used as a sinceChange boundary.
func (client *TeamCityClient) getChangeID(ctx context.Context, buildTypeID, version string) (int64, error) {
	endpoint := fmt.Sprintf("/app/rest/changes?locator=buildType:(id:%s),version:%s&fields=change(id)",
		url.QueryEscape(buildTypeID), url.QueryEscape(version))
	res, err := client.makeRequest(ctx, endpoint, map[string]string{"Accept": "application/json"})
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	var changes changesResponse
	if err := json.NewDecoder(res.Body).Decode(&changes); err != nil {
		return 0, fmt.Errorf("failed to decode change-id response: %w", err)
	}
	if len(changes.Change) == 0 {
		return 0, nil
	}
	return changes.Change[0].Id, nil
}

func (client *TeamCityClient) getChangesSince(ctx context.Context, buildTypeID string, sinceChangeID int64) (*changesResponse, error) {
	endpoint := fmt.Sprintf("/app/rest/changes?locator=buildType:(id:%s),sinceChange:(id:%d),count:10000&fields=count,change(id,version)",
		url.QueryEscape(buildTypeID), sinceChangeID)
	res, err := client.makeRequest(ctx, endpoint, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var changes changesResponse
	if err := json.NewDecoder(res.Body).Decode(&changes); err != nil {
		return nil, fmt.Errorf("failed to decode changes-since response: %w", err)
	}
	return &changes, nil
}

// getChangesGap determines whether buildID includes all commits since the
// previous dot (previousBuildID). It compares the current build's own changeset
// against the full set of changes the build configuration accumulated since the
// previous dot; the difference is the commits consumed by builds in between that
// produced no data point.
// getChangesGap reports whether the current dot's bisect range covers all commits
// since the previous dot. currentFirstCommit is the oldest commit of the current
// range (the bisect's lower edge, sourced from the same ClickHouse installer
// changes the dialog builds the range from), so the check stays consistent with
// the range that is actually bisected — TeamCity's per-build changeset is not used,
// since out-of-order builds make it unreliable (often 0 changes).
func (client *TeamCityClient) getChangesGap(ctx context.Context, buildID, previousBuildID, currentFirstCommit string) (*ChangesGap, error) {
	if currentFirstCommit == "" {
		return &ChangesGap{Known: false}, nil
	}

	// Anchor the previous dot on the commit it was built on, not on its attributed
	// changeset: builds run out of commit order, so the previous dot's build is
	// frequently attributed 0 changes even though it has a revision.
	prevRev, err := client.getBuildRevision(ctx, previousBuildID)
	if err != nil {
		return nil, err
	}
	if prevRev == "" {
		return &ChangesGap{Known: false}, nil
	}

	buildTypeID, err := client.getBuildType(ctx, buildID)
	if err != nil {
		return nil, err
	}
	prevRevID, err := client.getChangeID(ctx, buildTypeID, prevRev)
	if err != nil {
		return nil, err
	}
	if prevRevID == 0 {
		return &ChangesGap{Known: false}, nil
	}

	// Every commit strictly after the previous dot's revision, newest-first.
	since, err := client.getChangesSince(ctx, buildTypeID, prevRevID)
	if err != nil {
		return nil, err
	}

	// Locate the current range's lower edge among the commits after the previous
	// dot. If it isn't there, the range reaches back to or before the previous dot
	// — no gap. Otherwise the commits older than it (but still after the previous
	// dot) are the uncovered gap.
	idx := -1
	for i, c := range since.Change {
		if c.Version == currentFirstCommit {
			idx = i
			break
		}
	}
	gap := &ChangesGap{Known: true}
	if idx != -1 {
		gap.GapCommitCount = len(since.Change) - 1 - idx
		gap.HasGap = gap.GapCommitCount > 0
	}
	return gap, nil
}

func (client *TeamCityClient) startBuild(ctx context.Context, buildId string, params map[string]string) (*BuildResponse, error) {
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
	req.Header.Set("Accept", "application/xml")

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
	return &buildResponse, nil
}
