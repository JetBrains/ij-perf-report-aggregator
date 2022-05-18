package main

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
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
  Id        int    `json:"id"`
  Status    string `json:"status"`
  Agent     Agent  `json:"agent"`
  StartDate string `json:"startDate"`

  ArtifactDependencies ArtifactDependencies `json:"artifact-dependencies"`
  Artifacts            Artifacts            `json:"artifacts"`
  Private              bool                 `json:"personal"`
  TriggeredBy          *TriggeredBy         `json:"triggered"`

  installerInfo *InstallerInfo
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
  t.logger.Info("request", zap.String("url", url))

  request, err := t.createRequest(url, t.taskContext)
  if err != nil {
    return nil, err
  }

  response, err := t.httpClient.Do(request)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  defer util.Close(response.Body, t.logger)

  if response.StatusCode > 300 {
    responseBody, _ := io.ReadAll(response.Body)
    return nil, errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  var result BuildList
  err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&result)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &result, nil
}

type ChangeList struct {
  List []Change `json:"change"`
}

type Change struct {
  Version string `json:"version"`
}
