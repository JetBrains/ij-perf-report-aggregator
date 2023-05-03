package main

import (
  "context"
  "database/sql"
  "encoding/json"
  "errors"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
  sqlutil "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  e "github.com/develar/errors"
  "github.com/nats-io/nats.go"
  "go.uber.org/atomic"
  "go.uber.org/zap"
  "io"
  "net/http"
  "net/url"
  "os"
  "sort"
  "strings"
  "time"
)

const tcTimeFormat = "20060102T150405-0700"

type Collector struct {
  serverUrl string

  httpClient  *http.Client
  taskContext context.Context

  config analyzer.DatabaseConfiguration

  logger *zap.Logger

  tcSessionId atomic.String

  installerBuildIdToInfo map[int]*InstallerInfo
  buildIdToInfo          map[int]*BuildInfo
}

func doNotifyServer(natsUrl string, logger *zap.Logger) error {
  logger.Info("ask report aggregator server to clear cache")
  nc, err := nats.Connect("nats://" + natsUrl)
  if err != nil {
    return e.WithStack(err)
  }

  err = nc.Publish("server.clearCache", []byte("tcCollector"))
  if err != nil {
    return e.WithStack(err)
  }

  logger.Info("ask to backup db")
  err = nc.Publish("db.backup", []byte("tcCollector"))
  if err != nil {
    return e.WithStack(err)
  }

  // ensure that message is delivered, because app will be exited very soon
  err = nc.Flush()
  if err != nil {
    return e.WithStack(err)
  }
  return nil
}

func collectFromTeamCity(taskContext context.Context, clickHouseUrl string, tcUrl string, projectId string, buildConfigurationIds []string, initialSince time.Time, userSpecifiedSince time.Time, httpClient *http.Client, logger *zap.Logger, ) error {
  serverUrl := tcUrl + "/app/rest"

  serverBuildUrl, err := url.Parse(serverUrl + "/builds/")
  if err != nil {
    return err
  }

  config := analyzer.GetAnalyzer(projectId)

  db, err := analyzer.OpenDb(clickHouseUrl, config)
  if err != nil {
    return e.WithStack(err)
  }

  defer util.Close(db, logger)

  for _, buildTypeId := range buildConfigurationIds {
    if taskContext.Err() != nil {
      return e.WithStack(taskContext.Err())
    }

    err = collectBuildConfiguration(
      taskContext,
      httpClient,
      db,
      config,
      buildTypeId,
      serverUrl,
      serverBuildUrl,
      tcUrl,
      userSpecifiedSince,
      initialSince,
      logger.With(zap.String("buildTypeId", buildTypeId), zap.String("projectId", projectId)),
    )
    if err != nil {
      return err
    }
  }
  return err
}

func collectBuildConfiguration(
  taskContext context.Context,
  httpClient *http.Client,
  db driver.Conn,
  config analyzer.DatabaseConfiguration,
  buildTypeId string,
  serverUrl string,
  serverBuildUrl *url.URL,
  serverHost string,
  userSpecifiedSince time.Time,
  initialSince time.Time,
  logger *zap.Logger,
) error {
  q := serverBuildUrl.Query()
  locator := "buildType:(id:" + buildTypeId + "),defaultFilter:false,failedToStart:false,state:finished,canceled:false,count:500"

  since := userSpecifiedSince
  if since.IsZero() {
    since = initialSince
    if since.IsZero() {
      //goland:noinspection SqlResolve
      query := "select last_time from collector_state where build_type_id = '" + sqlutil.StringEscaper.Replace(buildTypeId) + "' order by last_time desc limit 1"
      err := db.QueryRow(taskContext, query).Scan(&since)
      if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return e.WithStack(err)
      }
    }
  }

  if !since.IsZero() {
    locator += ",sinceDate:" + since.Format(tcTimeFormat)
  }

  q.Set("locator", locator)
  q.Set("fields", buildTeamCityQuery())
  serverBuildUrl.RawQuery = q.Encode()

  logger.Info("collect", zap.String("buildTypeId", buildTypeId), zap.Time("since", since))

  reportExistenceChecker := &ReportExistenceChecker{}

  err := reportExistenceChecker.reset(taskContext, config.DbName, config.TableName, buildTypeId, db, since)
  if err != nil {
    return err
  }

  // TC returns from newest to oldest, but we need
  // 1) to insert in opposite order (less merge work for ClickHouse)
  // 2) set last collect state once the oldest chunk is committed, but it is possible only if the oldest will be inserted before newest (as we ask TC to returns since some date)
  var buildsToLoad [][]*Build

  collector := &Collector{
    serverUrl: serverUrl,

    httpClient:  httpClient,
    taskContext: taskContext,
    config:      config,

    installerBuildIdToInfo: make(map[int]*InstallerInfo),
    buildIdToInfo:          make(map[int]*BuildInfo),

    logger: logger,
  }

  buildList, err := collector.loadBuilds(serverBuildUrl.String())
  if err != nil {
    logger.Warn(err.Error())
    return nil
  }

  buildsToLoad = append(buildsToLoad, buildList.Builds)

  totalCount := len(buildList.Builds)
  nextHref := buildList.NextHref
  for len(buildList.NextHref) != 0 {
    if taskContext.Err() != nil {
      return e.WithStack(taskContext.Err())
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

    logger.Debug("load reports", zap.Int("chunk", i))

    lastBuildStartDate, err := time.Parse(tcTimeFormat, builds[len(builds)-1].StartDate)
    if err != nil {
      return e.WithStack(err)
    }

    reportAnalyzer, err := analyzer.CreateReportAnalyzer(taskContext, db, config, logger)
    if err != nil {
      return err
    }

    err = collector.loadReports(builds, reportExistenceChecker, reportAnalyzer)
    if err != nil {
      return err
    }

    logger.Debug("wait for analyze and insert", zap.Int("chunk", i))
    err = reportAnalyzer.WaitAnalyzeAndInsert()
    if err != nil {
      return err
    }

    // engine ReplacingMergeTree(last_time) is used, no need to delete old entry
    // set last collect time to 1 second after last build in chunk
    err = updateLastCollectTime(taskContext, buildTypeId, lastBuildStartDate.Add(1*time.Second), db)
    if err != nil {
      return err
    }
  }
  return nil
}

func buildTeamCityQuery() string {
  q := "file(href)"
  for i := 0; i < 3; i++ {
    q = "file(href,children(href," + q + "))"
  }
  return "count,href,nextHref,build(id,buildTypeId,number,startDate,status,agent(name),artifacts(" + q + "),artifact-dependencies(build(id,buildTypeId,finishDate)),personal,triggered(user(email)))"
}

func updateLastCollectTime(ctx context.Context, buildTypeId string, lastCollectTimeToSet time.Time, db driver.Conn) error {
  //goland:noinspection SqlResolve
  batch, err := db.PrepareBatch(ctx, "insert into collector_state values")
  if err != nil {
    return e.WithStack(err)
  }

  err = batch.Append(buildTypeId, lastCollectTimeToSet)
  if err != nil {
    return e.WithStack(err)
  }

  err = batch.Send()
  if err != nil {
    return e.WithStack(err)
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

func (t *Collector) get(ctx context.Context, url string) (*http.Response, error) {
  request, err := t.createRequest(ctx, url)
  if err != nil {
    return nil, err
  }

  response, err := t.httpClient.Do(request)
  if err != nil {
    if errors.Is(err, context.Canceled) {
      return nil, err
    }
    return nil, e.WithStack(err)
  }
  return response, nil
}

func (t *Collector) getSnapshots(ctx context.Context, configuration string) ([]string, error) {
  isComposite, err := t.isComposite(ctx, configuration)
  if err != nil {
    return nil, err
  }
  if !isComposite {
    return []string{configuration}, nil
  }
  configurations := make([]string, 0)
  err = t.getSnapshotsRecursive(ctx, configuration, &configurations)
  return configurations, err
}

func (t *Collector) getSnapshotsRecursive(ctx context.Context, configuration string, configurations *[]string) error {
  isComposite, err := t.isComposite(ctx, configuration)
  if err != nil {
    return nil
  }
  if strings.Contains(configuration, "Installers") {
    return nil
  }
  if !isComposite {
    *configurations = append(*configurations, configuration)
    return nil
  }

  response, err := t.get(ctx, t.serverUrl+"/buildTypes/"+configuration+"/snapshot-dependencies")
  if err != nil {
    return err
  }
  defer response.Body.Close()
  responseBody, _ := io.ReadAll(response.Body)
  if response.StatusCode > 300 {
    return e.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }

  type Dependency struct {
    Id string
  }
  type AllDependencies struct {
    Dependencies []Dependency `json:"snapshot-dependency"`
  }

  dependency := &AllDependencies{}
  err = json.Unmarshal(responseBody, dependency)
  if err != nil {
    return err
  }

  for _, dependency := range dependency.Dependencies {
    err = t.getSnapshotsRecursive(ctx, dependency.Id, configurations)
    if err != nil {
      t.logger.Warn(err.Error())
    }
  }
  return nil
}

func (t *Collector) isComposite(ctx context.Context, configuration string) (bool, error) {
  response, err := t.get(ctx, t.serverUrl+"/buildTypes/"+configuration+"/settings/buildConfigurationType")
  if err != nil {
    return false, err
  }
  defer response.Body.Close()
  responseBody, _ := io.ReadAll(response.Body)
  if response.StatusCode > 300 {
    return false, e.Errorf("Invalid response (%s): %s", response.Status, responseBody)
  }
  type BuildType struct {
    Name  string
    Value string
  }
  var buildType BuildType
  err = json.Unmarshal(responseBody, &buildType)
  return buildType.Value == "COMPOSITE", err
}

func (t *Collector) createRequest(ctx context.Context, url string) (*http.Request, error) {
  request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
  if err != nil {
    return nil, e.WithStack(err)
  }

  sessionId := t.tcSessionId.Load()
  if len(sessionId) != 0 {
    request.AddCookie(&http.Cookie{Name: "TCSESSIONID", Value: sessionId})
  }
  request.Header.Add("Authorization", "Bearer "+os.Getenv("TC_TOKEN"))

  request.Header.Add("Accept", "application/json")
  return request, nil
}
