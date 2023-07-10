package main

import (
  "compress/gzip"
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/tc-properties"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "io"
  "net/http"
  "net/url"
  "path"
  "strconv"
  "strings"
)

type ArtifactItem struct {
  data []byte
  path string
}

func (t *Collector) downloadReports(ctx context.Context, build Build) ([]ArtifactItem, error) {
  var result []ArtifactItem
  err := t.findAndDownloadStartUpReports(ctx, build, build.Artifacts.File, &result)
  if err != nil {
    return nil, err
  }
  return result, nil
}

func (t *Collector) findAndDownloadStartUpReports(ctx context.Context, build Build, artifacts []Artifact, result *[]ArtifactItem) error {
  for _, artifact := range artifacts {
    name := path.Base(artifact.Url)
    if strings.HasSuffix(artifact.Url, ".json") && strings.HasPrefix(name, "startup-stats") ||
      strings.HasSuffix(name, ".performance.json") ||
      strings.HasSuffix(artifact.Url, ".json") && (strings.Contains(artifact.Url, "metrics") && name != "action.invoked.json" && name != "spans.json") ||
      t.config.DbName == "jbr" && strings.HasSuffix(name, ".txt") ||
      t.config.DbName == "bazel" && name == "metrics.txt" ||
      t.config.DbName == "qodana" && name == "open-telemetry.json" {
      artifactUrlString := t.serverUrl + strings.Replace(strings.TrimPrefix(artifact.Url, "/app/rest"), "/artifacts/metadata/", "/artifacts/content/", 1)
      report, err := t.downloadStartUpReport(ctx, build, artifactUrlString)
      if err != nil {
        return err
      }

      *result = append(*result, ArtifactItem{
        data: report,
        path: artifactUrlString,
      })
      continue

    }

    err := t.findAndDownloadStartUpReports(ctx, build, artifact.Children.File, result)
    if err != nil {
      return err
    }
  }

  return nil
}

func (t *Collector) downloadStartUpReport(ctx context.Context, build Build, artifactUrlString string) ([]byte, error) {
  artifactUrl, err := url.Parse(artifactUrlString)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  response, err := t.get(ctx, artifactUrl.String())
  if err != nil {
    return nil, err
  }

  defer response.Body.Close()

  if response.StatusCode > 300 {
    if response.StatusCode == http.StatusNotFound && build.Status == "FAILURE" {
      t.logger.Warn("no report", zap.Int("id", build.Id), zap.String("status", build.Status))
      return nil, err
    }
    responseBody, _ := io.ReadAll(response.Body)
    return nil, errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  // ReadAll is used because report not only required to be decoded, but also stored as is (after minification)
  data, err := io.ReadAll(response.Body)
  if err != nil {
    return nil, err
  }
  return data, nil
}

func (t *Collector) downloadBuildProperties(ctx context.Context, build Build) ([]byte, error) {
  artifactUrl, err := url.Parse(t.serverUrl + "/builds/id:" + strconv.Itoa(build.Id) + "/artifacts/content/.teamcity/properties/build.start.properties.gz")
  if err != nil {
    return nil, err
  }

  response, err := t.get(ctx, artifactUrl.String())
  if err != nil {
    return nil, err
  }

  defer response.Body.Close()

  if response.StatusCode > 300 {
    if response.StatusCode == http.StatusNotFound {
      t.logger.Warn("build.start.properties not found", zap.String("url", artifactUrl.String()))
      return nil, nil
    }

    responseBody, _ := io.ReadAll(response.Body)
    return nil, errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  gzipReader, err := gzip.NewReader(response.Body)
  if err != nil {
    return nil, err
  }

  data, err := io.ReadAll(gzipReader)
  if err != nil {
    return nil, err
  }

  return tc_properties.ReadProperties(data)
}
