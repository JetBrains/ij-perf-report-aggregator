package main

import (
	"fmt"
	"github.com/json-iterator/go"
	"io"
	"time"
)

type BuildList struct {
	Count    int    `json:"count"`
	Href     string `json:"href"`
	NextHref string `json:"nextHref"`

	Builds []*Build `json:"build"`
}

type Build struct {
	Id          int    `json:"id"`
	Type        string `json:"buildTypeId"`
	Status      string `json:"status"`
	Agent       Agent  `json:"agent"`
	FinishDate  string `json:"finishDate"`
	BuildNumber string `json:"number"`

	ArtifactDependencies ArtifactDependencies `json:"artifact-dependencies"`
	Artifacts            Artifacts            `json:"artifacts"`
	Private              bool                 `json:"personal"`
	TriggeredBy          *TriggeredBy         `json:"triggered"`

	installerInfo *InstallerInfo
	buildInfo     *BuildInfo
}

type TriggeredBy struct {
	User *User `json:"user"`
}

type User struct {
	Email string `json:"email"`
}

type InstallerInfo struct {
	id        int
	changes   []string
	buildTime time.Time
}

type BuildInfo struct {
	id      int
	changes []string
}

type Artifacts struct {
	File []Artifact `json:"file"`
}

type Artifact struct {
	Url      string    `json:"href"`
	Children Artifacts `json:"children"`
}

type ArtifactDependencies struct {
	Builds []ArtifactDependencyBuild `json:"build"`
}

type ArtifactDependencyBuild struct {
	Id          int    `json:"id"`
	BuildTypeId string `json:"buildTypeId"`
	FinishDate  string `json:"finishDate"`
}

type Agent struct {
	Name string `json:"name"`
}

func (t *Collector) loadBuilds(url string) (*BuildList, error) {
	t.logger.Info("request", "url", url)

	request, err := t.createRequest(t.taskContext, url)
	if err != nil {
		return nil, err
	}

	response, err := t.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to load builds: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode > 300 {
		responseBody, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("invalid response (%s): %s", response.Status, responseBody)
	}

	t.storeSessionIdCookie(response)

	var result BuildList
	err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse builds: %w", err)
	}
	return &result, nil
}

type ChangeList struct {
	List []Change `json:"change"`
}

type Change struct {
	Version string `json:"version"`
}
