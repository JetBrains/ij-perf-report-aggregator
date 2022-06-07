package main

import (
  "context"
  "database/sql"
  e "errors"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  sqlutil "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
  "github.com/develar/errors"
  "github.com/nats-io/nats.go"
  "go.uber.org/atomic"
  "go.uber.org/zap"
  "net/http"
  "net/url"
  "os"
  "sort"
  "time"
)

const tcTimeFormat = "20060102T150405-0700"

type Collector struct {
  serverUrl string

  httpClient  *http.Client
  taskContext context.Context

  reportAnalyzer *analyzer.ReportAnalyzer

  logger *zap.Logger

  tcSessionId atomic.String

  installerBuildIdToInfo map[int]*InstallerInfo

  reportExistenceChecker *ReportExistenceChecker
}

var productCodeToBuildName = map[string]string{
  "IU": "Ultimate",
  "WS": "WebStorm",
  "PS": "PhpStorm",
  "DB": "DataGrip",
  "GO": "GoLand",
  "RM": "RubyMine",
}

func doNotifyServer(natsUrl string, logger *zap.Logger) error {
  logger.Info("ask report aggregator server to clear cache")
  nc, err := nats.Connect("nats://" + natsUrl)
  if err != nil {
    return errors.WithStack(err)
  }

  err = nc.Publish("server.clearCache", []byte("tcCollector"))
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("ask to backup db")
  err = nc.Publish("db.backup", []byte("tcCollector"))
  if err != nil {
    return errors.WithStack(err)
  }

  // ensure that message is delivered, because app will be exited very soon
  err = nc.Flush()
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

func collectFromTeamCity(
  clickHouseUrl string,
  tcUrl string,
  projectId string,
  buildConfigurationIds []string,
  initialSince time.Time,
  userSpecifiedSince time.Time,
  httpClient *http.Client, logger *zap.Logger,
  taskContext context.Context, cancel context.CancelFunc,
) error {
  reportAnalyzer, err := analyzer.CreateReportAnalyzer(clickHouseUrl, projectId, taskContext, logger, func() {
    logger.Debug("canceled by analyzer")
    cancel()
  })
  if err != nil {
    return err
  }

  serverHost := tcUrl
  collector := &Collector{
    serverUrl: serverHost + "/app/rest",

    reportAnalyzer: reportAnalyzer,

    httpClient:  httpClient,
    taskContext: taskContext,

    installerBuildIdToInfo: make(map[int]*InstallerInfo),

    logger:                 logger,
    reportExistenceChecker: &ReportExistenceChecker{},
  }

  serverUrl, err := url.Parse(collector.serverUrl + "/builds/")
  if err != nil {
    return err
  }

  for _, buildTypeId := range buildConfigurationIds {
    if taskContext.Err() != nil {
      return errors.WithStack(taskContext.Err())
    }

    q := serverUrl.Query()
    locator := "buildType:(id:" + buildTypeId + "),defaultFilter:false,failedToStart:false,state:finished,canceled:false,count:500"

    since := userSpecifiedSince
    if since.IsZero() {
      since = initialSince
      if since.IsZero() {
        query := "select last_time from collector_state where build_type_id = '" + sqlutil.StringEscaper.Replace(buildTypeId) + "' order by last_time desc limit 1"
        err = reportAnalyzer.InsertReportManager.InsertManager.Db.QueryRow(taskContext, query).Scan(&since)
        if err != nil && err != sql.ErrNoRows {
          return errors.WithStack(err)
        }
      }
    }

    if !since.IsZero() {
      locator += ",sinceDate:" + since.Format(tcTimeFormat)
    }

    q.Set("locator", locator)
    q.Set("fields", buildTeamCityQuery())
    serverUrl.RawQuery = q.Encode()

    logger.Info("collect", zap.String("buildTypeId", buildTypeId), zap.Time("since", since))

    err = collector.reportExistenceChecker.reset(projectId, buildTypeId, reportAnalyzer, taskContext, since)
    if err != nil {
      return err
    }

    // TC returns from newest to oldest, but we need
    // 1) to insert in opposite order (less merge work for ClickHouse)
    //2) set last collect state once the oldest chunk is committed, but it is possible only if the oldest will be inserted before newest (as we ask TC to returns since some date)
    var buildsToLoad [][]*Build

    buildList, err := collector.loadBuilds(serverUrl.String())
    if err != nil {
      return err
    }

    buildsToLoad = append(buildsToLoad, buildList.Builds)

    totalCount := len(buildList.Builds)
    nextHref := buildList.NextHref
    for len(buildList.NextHref) != 0 {
      if taskContext.Err() != nil {
        return errors.WithStack(taskContext.Err())
      }

      buildList, err = collector.loadBuilds(serverHost + nextHref)
      if err != nil {
        return err
      }

      nextHref = buildList.NextHref
      buildsToLoad = append(buildsToLoad, buildList.Builds)
      totalCount += len(buildList.Builds)
    }

    logger.Info("load reports", zap.Int("buildCount", totalCount), zap.String("buildTypeId", buildTypeId), zap.Time("since", since))

    for i := len(buildsToLoad) - 1; i >= 0; i-- {
      builds := buildsToLoad[i]
      if len(builds) == 0 {
        continue
      }

      sort.Slice(builds, func(i, j int) bool {
        return builds[i].Id < builds[j].Id
      })

      lastBuildStartDate, err := time.Parse(tcTimeFormat, builds[len(builds)-1].StartDate)
      if err != nil {
        return errors.WithStack(err)
      }

      err = collector.loadReports(builds)
      if err != nil {
        return err
      }

      select {
      case analyzeError := <-reportAnalyzer.WaitAnalyzer():
        if analyzeError != nil {
          return analyzeError
        }

        // engine ReplacingMergeTree(last_time) is used, no need to delete old entry
        // set last collect time to 1 second after last build in chunk
        err = updateLastCollectTime(buildTypeId, lastBuildStartDate.Add(1*time.Second), reportAnalyzer.InsertReportManager.InsertManager.Db, taskContext)
        if err != nil {
          return err
        }
      case <-taskContext.Done():
        return nil
      }
    }
  }

  err = reportAnalyzer.Close()
  return err
}

func buildTeamCityQuery() string {
  q := "file(href)"
  for i := 0; i < 3; i++ {
    q = "file(href,children(href," + q + "))"
  }
  return "count,href,nextHref,build(id,buildTypeId,startDate,status,agent(name),artifacts(" + q + "),artifact-dependencies(build(id,buildTypeId,finishDate)),personal,triggered(user(email)))"
}

func updateLastCollectTime(buildTypeId string, lastCollectTimeToSet time.Time, db driver.Conn, ctx context.Context) error {
  batch, err := db.PrepareBatch(ctx, "insert into collector_state values")
  if err != nil {
    return errors.WithStack(err)
  }

  err = batch.Append(buildTypeId, lastCollectTimeToSet)
  if err != nil {
    return errors.WithStack(err)
  }

  err = batch.Send()
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
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

func (t *Collector) get(url string, ctx context.Context) (*http.Response, error) {
  request, err := t.createRequest(url, ctx)
  if err != nil {
    return nil, err
  }

  response, err := t.httpClient.Do(request)
  if err != nil {
    if e.Is(err, context.Canceled) {
      return nil, err
    } else {
      return nil, errors.WithStack(err)
    }
  }
  return response, nil
}

func (t *Collector) createRequest(url string, ctx context.Context) (*http.Request, error) {
  request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  sessionId := t.tcSessionId.Load()
  if len(sessionId) != 0 {
    request.AddCookie(&http.Cookie{Name: "TCSESSIONID", Value: sessionId})
  }
  request.Header.Add("Authorization", "Bearer "+os.Getenv("TC_TOKEN"))

  request.Header.Add("Accept", "application/json")
  return request, nil
}
