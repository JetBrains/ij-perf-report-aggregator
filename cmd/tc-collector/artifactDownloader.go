package main

import (
  "compress/gzip"
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/tc-properties"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "io/ioutil"
  "net/url"
  "path"
  "strconv"
  "strings"
)

func (t *Collector) downloadStartUpReports(build Build) ([][]byte, error) {
  var result [][]byte
  err := t.findAndDownloadStartUpReports(build, build.Artifacts, &result)
  if err != nil {
    return nil, err
  }
  return result, nil
}

func (t *Collector) findAndDownloadStartUpReports(build Build, artifact Artifacts, result *[][]byte) error {
  for _, artifact := range artifact.File {
    if strings.HasSuffix(artifact.Url, ".json") {
      name := path.Base(artifact.Url)
      if strings.HasPrefix(name, "startup-stats") || strings.HasSuffix(name, ".performance.json") {
        artifactUrlString := t.serverUrl + strings.Replace(strings.TrimPrefix(artifact.Url, "/app/rest"), "/artifacts/metadata/", "/artifacts/content/", 1)
        report, err := t.downloadStartUpReport(build, artifactUrlString)
        if err != nil {
          return err
        }

        *result = append(*result, report)
        continue
      }
    }

    err := t.findAndDownloadStartUpReports(build, artifact.Children, result)
    if err != nil {
      return err
    }
  }

  return nil
}

func (t *Collector) downloadStartUpReport(build Build, artifactUrlString string) ([]byte, error) {
  artifactUrl, err := url.Parse(artifactUrlString)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  response, err := t.get(artifactUrl.String())
  if err != nil {
    if err == context.Canceled {
      return nil, err
    } else {
      return nil, errors.WithStack(err)
    }
  }

  defer util.Close(response.Body, t.logger)

  if response.StatusCode > 300 {
    if response.StatusCode == 404 && build.Status == "FAILURE" {
      t.logger.Warn("no report", zap.Int("id", build.Id), zap.String("status", build.Status))
      return nil, err
    }
    responseBody, _ := ioutil.ReadAll(response.Body)
    return nil, errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  // ReadAll is used because report not only required to be decoded, but also stored as is (after minification)
  data, err := ioutil.ReadAll(response.Body)
  if err != nil {
    return nil, err
  }
  return data, nil
}

func (t *Collector) downloadBuildProperties(build Build) ([]byte, error) {
  artifactUrl, err := url.Parse(t.serverUrl + "/builds/id:" + strconv.Itoa(build.Id) + "/artifacts/content/.teamcity/properties/build.start.properties.gz")
  if err != nil {
    return nil, err
  }

  response, err := t.get(artifactUrl.String())
  if err != nil {
    return nil, err
  }

  defer util.Close(response.Body, t.logger)

  if response.StatusCode > 300 {
    responseBody, _ := ioutil.ReadAll(response.Body)
    return nil, errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  gzipReader, err := gzip.NewReader(response.Body)
  if err != nil {
    return nil, err
  }

  data, err := ioutil.ReadAll(gzipReader)
  if err != nil {
    return nil, err
  }

  return tc_properties.ReadProperties(data)
}
