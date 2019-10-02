package teamcity

import (
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "io/ioutil"
  "report-aggregator/pkg/util"
)

type BuildList struct {
  Count    int    `json:"count"`
  Href     string `json:"href"`
  NextHref string `json:"nextHref"`

  Builds []Build `json:"build"`
}

type Build struct {
  Id    int   `json:"id"`
  Agent Agent `json:"agent"`
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
