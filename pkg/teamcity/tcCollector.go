package teamcity

import (
  "context"
  "encoding/base64"
  "encoding/hex"
  "github.com/alecthomas/kingpin"
  "github.com/araddon/dateparse"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "go.uber.org/atomic"
  "go.uber.org/zap"
  "io/ioutil"
  "log"
  "net/http"
  "net/url"
  "os"
  "report-aggregator/pkg/analyzer"
  "report-aggregator/pkg/model"
  "report-aggregator/pkg/util"
  "strconv"
  "strings"
  "sync"
  "time"
)

const tcTimeFormat = "20060102T150405-0700"

// TC REST API: By default only builds from the default branch are returned (https://www.jetbrains.com/help/teamcity/rest-api.html#Build-Locator),
// so, no need to explicitly specify filter
func ConfigureCollectFromTeamCity(app *kingpin.Application, log *zap.Logger) {
  command := app.Command("collect-tc", "Collect reports from TeamCity.")
  buildTypeIds := command.Flag("build-type-id", "The TeamCity build type id.").Short('c').Required().Strings()
  clickHouseUrl := command.Flag("clickHouse", "The ClickHouse server URL.").Default("localhost:9000").String()
  sinceDate := command.Flag("since", "The date to force collecting since").String()

  command.Action(func(context *kingpin.ParseContext) error {
    var since time.Time
    if len(*sinceDate) > 0 {
      var err error
      since, err = dateparse.ParseStrict(*sinceDate)
      if err != nil {
        return errors.WithStack(err)
      }
    }
    return collectFromTeamCity(*clickHouseUrl, *buildTypeIds, since, log)
  })
}

func collectFromTeamCity(clickHouseUrl string, buildTypeIds []string, since time.Time, logger *zap.Logger) error {
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

  var httpClient = &http.Client{
    Timeout: 30 * time.Second,
  }

  serverHost := "https://buildserver.labs.intellij.net"
  collector := &Collector{
    serverUrl: serverHost + "/app/rest",

    reportAnalyzer: reportAnalyzer,

    httpClient:  httpClient,
    taskContext: taskContext,

    rwLock:                  &sync.RWMutex{},
    installerBuildToChanges: make(map[int]string),

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
    locator := "buildType:(id:" + buildTypeId + "),count:500"
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

func (t *Collector) loadInstallerChanges(installerBuildId int) (string, error) {
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
    return "", err
  }

  response, err := t.get(artifactUrl.String())
  if err != nil {
    return "", err
  }

  defer util.Close(response.Body, t.logger)

  if response.StatusCode > 300 {
    responseBody, _ := ioutil.ReadAll(response.Body)
    return "", errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  t.storeSessionIdCookie(response)

  var changeList ChangeList
  err = jsoniter.ConfigFastest.NewDecoder(response.Body).Decode(&changeList)
  if err != nil {
    return "", err
  }

  var sb strings.Builder
  for index, change := range changeList.List {
    if index != 0 {
      sb.WriteRune('\n')
    }

    bytes, err := hex.DecodeString(change.Version)
    if err != nil {
      return "", err
    }

    encoder := base64.RawStdEncoding
    buf := make([]byte, encoder.EncodedLen(len(bytes)))
    encoder.Encode(buf, bytes)
    sb.Write(buf)
  }

  result = sb.String()
  t.installerBuildToChanges[installerBuildId] = result
  return result, nil
}

func (t *Collector) loadReports(builds []Build) error {
  err := util.MapAsync(len(builds), func(taskIndex int) (f func() error, err error) {
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
  installerBuildToChanges map[int]string
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
