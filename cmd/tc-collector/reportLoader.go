package main

import (
  "bytes"
  "encoding/base64"
  "encoding/hex"
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

  var err error

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
      t.logger.Error("cannot find installer build", zap.Int("buildId", build.Id))
      continue
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
    t.logger.Debug("load installer info", zap.Int("count", len(notLoadedInstallerBuildIds)), zap.Array("ids", zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
      for _, installerInfo := range notLoadedInstallerBuildIds {
        encoder.AppendInt(installerInfo.id)
      }
      return nil
    })))
    err = util.MapAsyncConcurrency(len(notLoadedInstallerBuildIds), networkRequestCount, func(index int) (f func() error, err error) {
      return func() error {
        installerInfo := notLoadedInstallerBuildIds[index]
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

  err = util.MapAsyncConcurrency(len(builds), networkRequestCount, func(taskIndex int) (f func() error, err error) {
    return func() error {
      build := builds[taskIndex]

      installerInfo := build.installerInfo
      if installerInfo == nil {
        // or already processed or cannot compute installer info
        return nil
      }

      dataList, err := t.downloadStartUpReports(*build)
      if err != nil {
        return err
      }

      if len(dataList) == 0 {
        t.logger.Error("cannot find any start-up report", zap.Int("id", build.Id))
        return nil
      }

      tcBuildProperties, err := t.downloadBuildProperties(*build)
      if err != nil {
        return err
      }

      for _, data := range dataList {
        err = t.reportAnalyzer.Analyze(data, model.ExtraData{
          Machine:            build.Agent.Name,
          TcBuildId:          build.Id,
          TcInstallerBuildId: installerInfo.id,
          BuildTime:          installerInfo.buildTime,
          TcBuildProperties:  tcBuildProperties,
          Changes:            installerInfo.changes,
        })
      }

      if err != nil {
        return err
      }
      return nil
    }, nil
  })
  if err != nil {
    return errors.WithStack(err)
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
