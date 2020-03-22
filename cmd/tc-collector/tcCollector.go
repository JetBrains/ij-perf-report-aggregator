package main

import (
  "context"
  "database/sql"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
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

var ProductCodeToBuildName = map[string]string{
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

func collectFromTeamCity(clickHouseUrl string, tcUrl string, buildTypeIds []string, userSpecifiedSince time.Time, httpClient *http.Client, logger *zap.Logger) error {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  reportAnalyzer, err := analyzer.CreateReportAnalyzer(clickHouseUrl, taskContext, logger, func() {
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

    logger: logger,
    reportExistenceChecker: &ReportExistenceChecker{},
  }

  serverUrl, err := url.Parse(collector.serverUrl + "/builds/")
  if err != nil {
    return err
  }

  for _, buildTypeId := range buildTypeIds {
    if taskContext.Err() != nil {
      return errors.WithStack(taskContext.Err())
    }

    q := serverUrl.Query()
    locator := "buildType:(id:" + buildTypeId + "),count:500"

    since := userSpecifiedSince
    if userSpecifiedSince.IsZero() {
      err = reportAnalyzer.Db.QueryRow("select last_time from collector_state where build_type_id = ? order by last_time desc limit 1", buildTypeId).Scan(&since)
      if err != nil && err != sql.ErrNoRows {
        return errors.WithStack(err)
      }
    }

    if !since.IsZero() {
      locator += ",sinceDate:" + since.Format(tcTimeFormat)
    }

    q.Set("locator", locator)
    q.Set("fields", "count,href,nextHref,build(id,startDate,status,agent(name),artifacts(file(children(file(children(file(href)))))),artifact-dependencies(build(id,buildTypeId,finishDate)))")
    serverUrl.RawQuery = q.Encode()

    logger.Info("collect", zap.String("buildTypeId", buildTypeId), zap.Time("since", since))

    err := collector.reportExistenceChecker.reset(buildTypeId, reportAnalyzer, taskContext, since)
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

      // commit each chunk to ensure that if later will be some error, we don't start from the beginning
      select {
      case analyzeError := <-reportAnalyzer.WaitAndCommit():
        if analyzeError != nil {
          return analyzeError
        }

        // engine ReplacingMergeTree(last_time) is used, no need to delete old entry
        // set last collect time to 1 second after last build in chunk
        err = updateLastCollectTime(buildTypeId, lastBuildStartDate.Add(1*time.Second), reportAnalyzer.Db, logger)
        if err != nil {
          return err
        }

        // break select and continue to process next build type chunk
        break

      case <-taskContext.Done():
        return nil
      }
    }
  }

  return nil
}

func updateLastCollectTime(buildTypeId string, lastCollectTimeToSet time.Time, db *sqlx.DB, logger *zap.Logger) error {
  tx, _ := db.Begin()
  stmt, err := tx.Prepare("insert into collector_state values (?, ?)")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(stmt, logger)

  _, err = stmt.Exec(buildTypeId, lastCollectTimeToSet)
  if err != nil {
    return errors.WithStack(err)
  }

  err = tx.Commit()
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

func (t *Collector) get(url string) (*http.Response, error) {
  request, err := t.createRequest(url)
  if err != nil {
    return nil, err
  }

  response, err := t.httpClient.Do(request)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return response, nil
}

func (t *Collector) createRequest(url string) (*http.Request, error) {
  request, err := http.NewRequestWithContext(t.taskContext, "GET", url, nil)
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
