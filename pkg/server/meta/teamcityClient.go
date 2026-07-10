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
	// FirstCommitAfterPreviousDot is the oldest uncovered commit — the first commit after the
	// previous dot. Empty when there is no gap. Callers can use it as the lower bound to analyse
	// across the gap.
	FirstCommitAfterPreviousDot string `json:"firstCommitAfterPreviousDot"`
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

func (client *TeamCityClient) getChangesBetweenBuilds(ctx context.Context, buildTypeID, branch, sinceBuildID, untilBuildID string) ([]changeRef, error) {
	branchDim := ""
	if branch != "" {
		branchDim = "branch:(name:" + url.QueryEscape(branch) + "),"
	}
	endpoint := fmt.Sprintf(
		"/app/rest/builds?locator=buildType:(id:%s),%ssinceBuild:(id:%s),untilBuild:(id:%s),count:10000&fields=build(changes($locator(count:10000),change(id,version)))",
		url.QueryEscape(buildTypeID), branchDim, url.QueryEscape(sinceBuildID), url.QueryEscape(untilBuildID))
	res, err := client.makeRequest(ctx, endpoint, map[string]string{"Accept": "application/json"})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var builds struct {
		Build []struct {
			Changes changesResponse `json:"changes"`
		} `json:"build"`
	}
	if err := json.NewDecoder(res.Body).Decode(&builds); err != nil {
		return nil, fmt.Errorf("failed to decode builds-with-changes response: %w", err)
	}

	changes := make([]changeRef, 0)
	for _, b := range builds.Build {
		changes = append(changes, b.Changes.Change...)
	}
	return changes, nil
}

func (client *TeamCityClient) getChangesGap(ctx context.Context, buildID, previousBuildID, currentFirstCommit string) (*ChangesGap, error) {
	if currentFirstCommit == "" {
		return &ChangesGap{Known: false}, nil
	}

	// Anchor the previous dot on the commit it was built on, not on its attributed changeset:
	// builds run out of commit order, so the previous dot's build is frequently attributed 0
	// changes even though it has a revision.
	prevRev, err := client.getBuildRevision(ctx, previousBuildID)
	if err != nil {
		return nil, err
	}
	if prevRev == "" {
		return &ChangesGap{Known: false}, nil
	}

	build, err := client.getBuildInfo(ctx, buildID)
	if err != nil {
		return nil, err
	}

	prevRevID, err := client.getChangeID(ctx, build.BuildTypeId, prevRev)
	if err != nil {
		return nil, err
	}
	currentFirstID, err := client.getChangeID(ctx, build.BuildTypeId, currentFirstCommit)
	if err != nil {
		return nil, err
	}
	if prevRevID == 0 || currentFirstID == 0 {
		return &ChangesGap{Known: false}, nil
	}

	gap := &ChangesGap{Known: true}
	// The current range starts at or before the previous dot (e.g. out-of-order builds), so
	// nothing is left uncovered.
	if currentFirstID <= prevRevID {
		return gap, nil
	}

	changes, err := client.getChangesBetweenBuilds(ctx, build.BuildTypeId, build.BranchName, previousBuildID, buildID)
	if err != nil {
		return nil, err
	}
	// Distinct commits strictly between the two dots' revisions. The change-id bounds discard the
	// current build's own commits and anything an out-of-order build in the window pulled in from
	// outside the range.
	uncovered := make(map[int64]struct{})
	var firstAfterID int64
	for _, c := range changes {
		if c.Id > prevRevID && c.Id < currentFirstID {
			uncovered[c.Id] = struct{}{}
			// The oldest uncovered commit (smallest change id) is the first commit after the
			// previous dot.
			if firstAfterID == 0 || c.Id < firstAfterID {
				firstAfterID = c.Id
				gap.FirstCommitAfterPreviousDot = c.Version
			}
		}
	}
	gap.GapCommitCount = len(uncovered)
	gap.HasGap = gap.GapCommitCount > 0
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
