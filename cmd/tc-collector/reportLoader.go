package main

import (
  "bytes"
  "context"
  "encoding/base64"
  "encoding/hex"
  e "errors"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
  "io/ioutil"
  "net/url"
  "runtime"
  "strconv"
  "strings"
  "time"
)

func (t *Collector) loadReports(builds []*Build) error {
  networkRequestCount := runtime.NumCPU() + 1
  if networkRequestCount > 8 {
    networkRequestCount = 8
  }

  var notLoadedInstallerBuildIds []*InstallerInfo
  for _, build := range builds {
    if t.reportExistenceChecker.has(build.Id) {
      t.logger.Info("build already processed", zap.Int("id", build.Id), zap.String("startDate", build.StartDate))
      continue
    }

    id, buildTime, err := computeBuildDate(build)
    if err != nil {
      return errors.WithStack(err)
    }

    if id == -1 {
      t.logger.Warn("cannot find installer build", zap.Int("buildId", build.Id))
    }

    installerInfo := t.installerBuildIdToInfo[id]
    if installerInfo == nil {
      installerInfo = &InstallerInfo{
        id:        id,
        buildTime: buildTime,
      }
      notLoadedInstallerBuildIds = append(notLoadedInstallerBuildIds, installerInfo)
      t.installerBuildIdToInfo[id] = installerInfo
    }
    build.installerInfo = installerInfo
  }

  if len(notLoadedInstallerBuildIds) > 0 {
    if notLoadedInstallerBuildIds[0].id != -1 {
      t.logger.Debug("load installer info", zap.Int("count", len(notLoadedInstallerBuildIds)), zap.Array("ids", zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
        for _, installerInfo := range notLoadedInstallerBuildIds {
          encoder.AppendInt(installerInfo.id)
        }
        return nil
      })))
      err := util.MapAsyncConcurrency(len(notLoadedInstallerBuildIds), networkRequestCount, func(index int) (func() error, error) {
        return func() error {
          installerInfo := notLoadedInstallerBuildIds[index]
          var err error
          installerInfo.changes, err = t.loadInstallerChanges(installerInfo.id)
          if err != nil {
            return errors.WithStack(err)
          }
          return nil
        }, nil
      })

      if err != nil {
        return errors.WithStack(err)
      }
    }
  }

  err := util.MapAsyncConcurrency(len(builds), networkRequestCount, func(taskIndex int) (f func() error, err error) {
    return func() error {
      build := builds[taskIndex]
      if build.Agent.Name == "Dead agent" {
        return nil
      }

      installerInfo := build.installerInfo
      if installerInfo == nil {
        // or already processed or cannot compute installer info
        return nil
      }

      artifacts, err := t.downloadStartUpReports(*build, t.taskContext)
      if err != nil {
        return err
      }

      if len(artifacts) == 0 {
        t.logger.Error("cannot find any performance report", zap.Int("id", build.Id), zap.String("status", build.Status))
        return nil
      }

      tcBuildProperties, err := t.downloadBuildProperties(*build, t.taskContext)
      if err != nil {
        return err
      }

      for _, artifact := range artifacts {
        if t.taskContext.Err() != nil {
          return nil
        }

        err = t.reportAnalyzer.Analyze(artifact.data, model.ExtraData{
          Machine:            build.Agent.Name,
          TcBuildId:          build.Id,
          TcInstallerBuildId: installerInfo.id,
          BuildTime:          installerInfo.buildTime,
          TcBuildProperties:  tcBuildProperties,
          Changes:            installerInfo.changes,
          ReportFile:         artifact.path,
        })
        if err != nil {
          if build.Status == "FAILURE" {
            t.logger.Warn("cannot parse performance report in the failed build", zap.Int("buildId", build.Id), zap.Error(err))
          } else {
            return err
          }
        }
      }
      return nil
    }, nil
  })
  if err != nil {
    if e.Is(err, context.Canceled) {
      return err
    } else {
      return errors.WithStack(err)
    }
  }
  return nil
}

func computeBuildDate(build *Build) (int, int64, error) {
  for _, dependencyBuild := range build.ArtifactDependencies.Builds {
    if strings.Contains(dependencyBuild.BuildTypeId, "Installer") {
      result, err := time.Parse(tcTimeFormat, dependencyBuild.FinishDate)
      if err != nil {
        return -1, -1, err
      }

      return dependencyBuild.Id, result.Unix(), nil
    }
  }
  return -1, -1, nil
}

func (t *Collector) loadInstallerChanges(installerBuildId int) ([][]byte, error) {
  artifactUrl, err := url.Parse(t.serverUrl + "/changes?locator=build:(id:" + strconv.Itoa(installerBuildId) + ")&fields=change(version)")
  if err != nil {
    return nil, err
  }

  response, err := t.get(artifactUrl.String(), t.taskContext)
  if err != nil {
    return nil, err
  }

  defer util.Close(response.Body, t.logger)

  if response.StatusCode > 300 {
    responseBody, _ := ioutil.ReadAll(response.Body)
    return nil, errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  var changeList ChangeList
  err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&changeList)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  var b bytes.Buffer
  result := make([][]byte, len(changeList.List))
  for index, change := range changeList.List {
    data, err := hex.DecodeString(change.Version)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    b.Reset()

    encoding := base64.RawStdEncoding
    buf := make([]byte, encoding.EncodedLen(len(data)))
    encoding.Encode(buf, data)
    result[index] = buf
  }

  return result, nil
}
