package main

import (
  "bytes"
  "context"
  "encoding/base64"
  "encoding/hex"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
  "golang.org/x/sync/errgroup"
  "io"
  "net/url"
  "runtime"
  "strconv"
  "strings"
  "time"
)

var networkRequestCount = runtime.NumCPU() + 1

func (t *Collector) loadReports(builds []*Build, reportExistenceChecker *ReportExistenceChecker, reportAnalyzer *analyzer.ReportAnalyzer) error {
  if networkRequestCount > 8 {
    networkRequestCount = 8
  }

  for index, build := range builds {
    if reportExistenceChecker.has(build.Id) {
      t.logger.Info("build already processed", zap.Int("id", build.Id), zap.String("startDate", build.StartDate))
      builds[index] = nil
    }
  }

  if t.config.HasInstallerField {
    err := t.loadInstallerInfo(builds, networkRequestCount)
    if err != nil {
      return err
    }
  }

  duration := time.Duration(len(builds)*120) * time.Second
  t.logger.Debug("load", zap.Int("timeout", int(duration.Seconds())))
  taskContextWithTimeout, cancel := context.WithTimeout(t.taskContext, duration)
  defer cancel()
  errGroup, loadContext := errgroup.WithContext(taskContextWithTimeout)

  errGroup.SetLimit(networkRequestCount)
  for _, build := range builds {
    if build == nil || build.Agent.Name == "Dead agent" {
      continue
    }
    errGroup.Go((func(build *Build) func() error {
      return func() error {
        t.logger.Info("processing build", zap.Int("id", build.Id))
        if t.config.HasInstallerField && build.installerInfo == nil {
          // or already processed or cannot compute installer info
          return nil
        }

        artifacts, err := t.downloadReports(*build, loadContext)
        if err != nil {
          return err
        }

        if len(artifacts) == 0 {
          t.logger.Error("cannot find any performance report", zap.Int("id", build.Id), zap.String("status", build.Status))
          return nil
        }

        tcBuildProperties, err := t.downloadBuildProperties(*build, loadContext)
        if err != nil {
          return err
        }
        if tcBuildProperties == nil {
          return nil
        }

        for _, artifact := range artifacts {
          if loadContext.Err() != nil {
            return nil
          }

          data := model.ExtraData{
            Machine:           build.Agent.Name,
            TcBuildId:         build.Id,
            TcBuildType:       build.Type,
            TcBuildProperties: tcBuildProperties,
            ReportFile:        artifact.path,
          }

          currentBuildTime, err := analyzer.ParseTime(build.StartDate)
          if err == nil {
            data.CurrentBuildTime = currentBuildTime
          }

          if build.Private && build.TriggeredBy.User != nil {
            data.TriggeredBy = build.TriggeredBy.User.Email
          }

          if t.config.HasInstallerField {
            installerInfo := build.installerInfo
            data.BuildTime = installerInfo.buildTime
            data.Changes = installerInfo.changes
            data.TcInstallerBuildId = installerInfo.id
          }

          err = reportAnalyzer.Analyze(artifact.data, data)
          if err != nil {
            if build.Status == "FAILURE" {
              t.logger.Warn("cannot parse performance report in the failed build", zap.Int("buildId", build.Id), zap.Error(err))
            } else {
              return err
            }
          }
        }
        return nil
      }
    })(build))
  }
  return errGroup.Wait()
}

func (t *Collector) loadInstallerInfo(builds []*Build, networkRequestCount int) error {
  var notLoadedInstallerBuildIds []*InstallerInfo
  for _, build := range builds {
    if build == nil {
      continue
    }

    id, buildTime, err := computeBuildDate(build)
    if err != nil {
      return err
    }

    if id == -1 {
      if t.config.HasInstallerField {
        t.logger.Error("cannot find installer build", zap.Int("buildId", build.Id))
      }
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

  if len(notLoadedInstallerBuildIds) == 0 {
    return nil
  }

  t.logger.Debug("load installer info", zap.Int("count", len(notLoadedInstallerBuildIds)), zap.Array("ids", zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
    for _, installerInfo := range notLoadedInstallerBuildIds {
      encoder.AppendInt(installerInfo.id)
    }
    return nil
  })))

  errGroup, loadContext := errgroup.WithContext(t.taskContext)
  errGroup.SetLimit(networkRequestCount)
  for _, installerInfo := range notLoadedInstallerBuildIds {
    if installerInfo.id == -1 {
      continue
    }

    errGroup.Go((func(installerInfo *InstallerInfo) func() error {
      return func() error {
        var err error
        installerInfo.changes, err = t.loadInstallerChanges(installerInfo.id, loadContext)
        return errors.WithStack(err)
      }
    })(installerInfo))
  }
  return errGroup.Wait()
}

func computeBuildDate(build *Build) (int, time.Time, error) {
  for _, dependencyBuild := range build.ArtifactDependencies.Builds {
    if strings.Contains(dependencyBuild.BuildTypeId, "Installer") {
      result, err := time.Parse(tcTimeFormat, dependencyBuild.FinishDate)
      if err != nil {
        return -1, time.Time{}, err
      }

      return dependencyBuild.Id, result, nil
    }
  }
  return -1, time.Time{}, nil
}

func (t *Collector) loadInstallerChanges(installerBuildId int, ctx context.Context) ([]string, error) {
  artifactUrl, err := url.Parse(t.serverUrl + "/changes?locator=build:(id:" + strconv.Itoa(installerBuildId) + ")&fields=change(version)")
  if err != nil {
    return nil, err
  }

  response, err := t.get(artifactUrl.String(), ctx)
  if err != nil {
    return nil, err
  }

  defer util.Close(response.Body, t.logger)

  if response.StatusCode > 300 {
    responseBody, _ := io.ReadAll(response.Body)
    return nil, errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  var changeList ChangeList
  err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&changeList)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  encoding := base64.RawStdEncoding

  var b bytes.Buffer
  result := make([]string, len(changeList.List))
  for index, change := range changeList.List {
    if strings.Contains(change.Version, " ") {
      //private build with custom change, format: 13 04 2022 12:14
      continue
    }
    data, err := hex.DecodeString(change.Version)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    b.Reset()

    buf := make([]byte, encoding.EncodedLen(len(data)))
    encoding.Encode(buf, data)
    result[index] = string(buf)
  }

  return result, nil
}
