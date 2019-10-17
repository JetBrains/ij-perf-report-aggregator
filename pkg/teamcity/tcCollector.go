package teamcity

import (
  "context"
  "github.com/alecthomas/kingpin"
  "github.com/araddon/dateparse"
  "github.com/develar/errors"
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
  "time"
)

/*
collect-tc -c ijplatform_master_UltimateStartupPerfTestMac -c ijplatform_master_UltimateStartupPerfTestWindows -c ijplatform_master_UltimateStartupPerfTestLinux \
-c ijplatform_master_WebStormStartupPerfTestMac -c ijplatform_master_WebStormStartupPerfTestWindows -c ijplatform_master_WebStormStartupPerfTestLinux \
--db /Volumes/data/ij-perf-db/db.sqlite
*/

// TC REST API: By default only builds from the default branch are returned (https://www.jetbrains.com/help/teamcity/rest-api.html#Build-Locator),
// so, no need to explicitly specify filter
func ConfigureCollectFromTeamCity(app *kingpin.Application, log *zap.Logger) {
  command := app.Command("collect-tc", "Collect reports from TeamCity.")
  buildTypeIds := command.Flag("build-type-id", "The TeamCity build type id.").Short('c').Required().Strings()
  dbPath := command.Flag("db", "The output SQLite database file.").Short('o').Required().String()
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
    return collectFromTeamCity(*dbPath, *buildTypeIds, since, log)
  })
}

func collectFromTeamCity(dbPath string, buildTypeIds []string, since time.Time, logger *zap.Logger) error {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  reportAnalyzer, err := analyzer.CreateReportAnalyzer(dbPath, taskContext, logger)
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

    logger: logger,
  }

  serverUrl, err := url.Parse(collector.serverUrl + "/builds/")
  if err != nil {
    return err
  }

  if since.IsZero() {
    lastGeneratedTime, err := reportAnalyzer.GetLastGeneratedTime()
    if err != nil {
      return err
    }

    if lastGeneratedTime > 0 {
      since = time.Unix(lastGeneratedTime, -1)
    }
  }

  for _, buildTypeId := range buildTypeIds {
    q := serverUrl.Query()
    locator := "buildType:(id:" + buildTypeId + "),count:500"
    if !since.IsZero() {
      locator += ",sinceDate:" + since.Format("20060102T150405-0700")
    }
    q.Set("locator", locator)
    q.Set("fields", "count,href,nextHref,build(id,status,agent(name))")
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

func (t *Collector) loadReports(builds []Build) error {
  err := util.MapAsync(len(builds), func(taskIndex int) (f func() error, err error) {
    return func() error {
      build := builds[taskIndex]

      artifactUrl, err := url.Parse(t.serverUrl + "/builds/id:" + strconv.Itoa(build.Id) + "/artifacts/content/run/startup/startup-stats-startup.json")
      if err != nil {
        return err
      }

      request, err := t.createRequest(artifactUrl.String())
      if err != nil {
        return err
      }

      response, err := t.httpClient.Do(request)
      if err != nil {
        return err
      }

      defer util.Close(response.Body, t.logger)

      if response.StatusCode > 300 {
        if response.StatusCode == 404 && build.Status == "FAILURE" {
          t.logger.Warn("no report", zap.Int("id", build.Id), zap.String("status", build.Status))
          return nil
        }
        responseBody, _ := ioutil.ReadAll(response.Body)
        return errors.Errorf("Invalid response (%s): %s", response.Status, responseBody)
      }

      t.storeSessionIdCookie(response)

      // ReadAll is used because report not only required to be decoded, but also stored as is (after minification)
      data, err := ioutil.ReadAll(response.Body)
      if err != nil {
        return err
      }

      return t.reportAnalyzer.Analyze(data, model.ExtraData{
        Machine:   build.Agent.Name,
        TcBuildId: build.Id,
      })
    }, nil
  })

  if err != nil {
    return err
  }

  return nil
}

type Collector struct {
  serverUrl string

  httpClient  *http.Client
  taskContext context.Context

  reportAnalyzer *analyzer.ReportAnalyzer

  logger *zap.Logger

  tcSessionId atomic.String
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
