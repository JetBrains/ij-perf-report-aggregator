package main

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "io/ioutil"
)

type BuildList struct {
  Count    int    `json:"count"`
  Href     string `json:"href"`
  NextHref string `json:"nextHref"`

  Builds []Build `json:"build"`
}

type Build struct {
  Id     int    `json:"id"`
  Status string `json:"status"`
  Agent  Agent  `json:"agent"`

  ArtifactDependencies ArtifactDependencies `json:"artifact-dependencies"`
  Artifacts            Artifacts            `json:"artifacts"`
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

func (t *Collector) processBuilds(url string) (string, error) {
  t.logger.Info("request", zap.String("url", url))

  request, err := t.createRequest(url)
  if err != nil {
    return "", err
  }

  response, err := t.httpClient.Do(request)
  if err != nil {
    return "", err
  }

  defer util.Close(response.Body, t.logger)

  if response.StatusCode > 300 {
    responseBody, _ := ioutil.ReadAll(response.Body)
    return "", errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  var result BuildList
  err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&result)
  if err != nil {
    return "", err
  }

  err = t.loadReports(result.Builds)
  if err != nil {
    return "", err
  }
  return result.NextHref, nil
}

type ChangeList struct {
  List []Change `json:"change"`
}

type Change struct {
  Version string `json:"version"`
}
