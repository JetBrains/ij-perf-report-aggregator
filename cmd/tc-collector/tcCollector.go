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
  "github.com/nats-io/nats.go"
  "go.uber.org/atomic"
  "go.uber.org/zap"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
  "runtime"
  "sort"
  "strconv"
  "strings"
  "sync"
  "time"
)

const tcTimeFormat = "20060102T150405-0700"

func doNotifyServer(logger *zap.Logger) error {
  logger.Info("ask report aggregator server to clear cache")
  nc, err := nats.Connect("nats://nats:4222")
  if err != nil {
    return errors.WithStack(err)
  }

  err = nc.Publish("server.clearCache", []byte("tcCollector"))
  if err != nil {
    return errors.WithStack(err)
  }

  // ensure that message is delivered, because app will be exited very soon
  err = nc.Flush()
  if err != nil {
    return errors.WithStack(err)
  }
  return err
}

func collectFromTeamCity(clickHouseUrl string, tcUrl string, buildTypeIds []string, since time.Time, httpClient *http.Client, logger *zap.Logger) error {
  //memProfiler := profile.Start(profile.MemProfile)
  //defer func() {
  //  memProfiler.Stop()
  //}()

  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  reportAnalyzer, err := analyzer.CreateReportAnalyzer(clickHouseUrl, taskContext, logger)
  if err != nil {
    return err
  }

  go func() {
    err = <-reportAnalyzer.ErrorChannel
    cancel()

    if err != nil {
      // zap doesn't print with newlines
      log.Printf("%+v", err)
      logger.Error("cannot analyze", zap.Error(err))
    }
  }()

  serverHost := tcUrl
  collector := &Collector{
    serverUrl: serverHost + "/app/rest",

    reportAnalyzer: reportAnalyzer,

    httpClient:  httpClient,
    taskContext: taskContext,

    rwLock:                  &sync.RWMutex{},
    installerBuildToChanges: make(map[int][]byte),

    logger: logger,
  }

  serverUrl, err := url.Parse(collector.serverUrl + "/builds/")
  if err != nil {
    return err
  }

  if since.IsZero() {
    lastGeneratedTime := reportAnalyzer.GetLastGeneratedTime()
    if lastGeneratedTime > 0 {
      since = time.Unix(lastGeneratedTime, -1)
    }
  }

  for _, buildTypeId := range buildTypeIds {
    q := serverUrl.Query()
    locator := "buildType:(id:" + buildTypeId + "),count:1000"

    if !since.IsZero() {
      locator += ",sinceDate:" + since.Format(tcTimeFormat)
    }

    q.Set("locator", locator)
    q.Set("fields", "count,href,nextHref,build(id,status,agent(name),artifact-dependencies(build(id,buildTypeId,finishDate)))")
    serverUrl.RawQuery = q.Encode()

    nextHref, err := collector.processBuilds(serverUrl.String())
    if err != nil {
      return err
    }

    //memProfiler.Stop()
    //os.Exit(0)

    for len(nextHref) != 0 {
      nextHref, err = collector.processBuilds(serverHost + nextHref)
      if err != nil {
        return err
      }
    }
  }

  select {
  case analyzeError := <-reportAnalyzer.ErrorChannel:
    cancel()
    return analyzeError

  case <-reportAnalyzer.Done():
    cancel()
    return nil

  case <-taskContext.Done():
    return nil
  }
}

func (t *Collector) loadInstallerChanges(installerBuildId int) ([]byte, error) {
  t.rwLock.RLock()
  result := t.installerBuildToChanges[installerBuildId]
  t.rwLock.RUnlock()

  if len(result) != 0 {
    return result, nil
  }

  t.rwLock.Lock()
  defer t.rwLock.Unlock()

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
    return nil, err
  }

  var b bytes.Buffer
  for index, change := range changeList.List {
    if index != 0 {
      b.WriteByte('\n')
    }

    data, err := hex.DecodeString(change.Version)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    base64Encoder := base64.NewEncoder(base64.RawStdEncoding, &b)
    _, err = base64Encoder.Write(data)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    err = base64Encoder.Close()
    if err != nil {
      return nil, errors.WithStack(err)
    }
  }

  result = b.Bytes()
  t.installerBuildToChanges[installerBuildId] = result
  return result, nil
}

func (t *Collector) loadReports(builds []Build) error {
  sort.Slice(builds, func(i, j int) bool {
    return builds[i].Id < builds[j].Id
  })

  networkRequestCount := runtime.NumCPU() + 1
  if networkRequestCount > 8 {
    networkRequestCount = 8
  }

  err := util.MapAsyncConcurrency(len(builds), networkRequestCount, func(taskIndex int) (f func() error, err error) {
    return func() error {
      build := builds[taskIndex]

      data, err := t.downloadStartUpReport(build)
      if err != nil {
        return err
      }

      if data == nil {
        return nil
      }

      installerBuildId, buildTime, err := computeBuildDate(build)
      if err != nil {
        return err
      }

      tcBuildProperties, err := t.downloadBuildProperties(build)
      if err != nil {
        return err
      }

      changes, err := t.loadInstallerChanges(installerBuildId)
      if err != nil {
        return err
      }

      return t.reportAnalyzer.Analyze(data, model.ExtraData{
        Machine:            build.Agent.Name,
        TcBuildId:          build.Id,
        TcInstallerBuildId: installerBuildId,
        BuildTime:          buildTime,
        TcBuildProperties:  tcBuildProperties,
        Changes:            changes,
      })
    }, nil
  })
  if err != nil {
    return err
  }

  return nil
}

func computeBuildDate(build Build) (int, int64, error) {
  for _, dependencyBuild := range build.ArtifactDependencies.Builds {
    if strings.HasSuffix(dependencyBuild.BuildTypeId, "_Installers") {
      parseFinishData, err := time.Parse(tcTimeFormat, dependencyBuild.FinishDate)
      if err != nil {
        return -1, -1, err
      }

      return dependencyBuild.Id, parseFinishData.Unix(), nil
    }
  }
  return -1, -1, errors.Errorf("cannot find installer build: %+v", build)
}

type Collector struct {
  serverUrl string

  httpClient  *http.Client
  taskContext context.Context

  reportAnalyzer *analyzer.ReportAnalyzer

  logger *zap.Logger

  tcSessionId atomic.String

  rwLock                  *sync.RWMutex
  installerBuildToChanges map[int][]byte
}

func getTcSessionIdCookie(cookies []*http.Cookie) string {
  for _, cookie := range cookies {
    if cookie.Name == "TCSESSIONID" {
      return cookie.Value
    }
  }
  return ""
}

func (t *Collector) storeSessionIdCookie(response *http.Response) {
  cookie := getTcSessionIdCookie(response.Cookies())
  // TC doesn't set cookie if it was already set for request
  if len(cookie) > 0 {
    t.tcSessionId.Store(cookie)
  }
}

func (t *Collector) get(url string) (*http.Response, error) {
  request, err := t.createRequest(url)
  if err != nil {
    return nil, err
  }

  return t.httpClient.Do(request)
}

func (t *Collector) createRequest(url string) (*http.Request, error) {
  request, err := http.NewRequestWithContext(t.taskContext, "GET", url, nil)
  if err != nil {
    return nil, err
  }

  sessionId := t.tcSessionId.Load()
  if len(sessionId) != 0 {
    request.AddCookie(&http.Cookie{Name: "TCSESSIONID", Value: sessionId})
  }
  request.Header.Add("Authorization", "Bearer "+os.Getenv("TC_TOKEN"))

  request.Header.Add("Accept", "application/json")
  return request, nil
}
