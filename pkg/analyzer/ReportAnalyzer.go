package analyzer

import (
  "context"
  errors2 "errors"
  "github.com/ClickHouse/clickhouse-go/v2"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/develar/errors"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "golang.org/x/sync/errgroup"
  "strconv"
  "strings"
  "time"
)

type ReportAnalyzer struct {
  config DatabaseConfiguration

  insertQueue chan *ReportInfo

  analyzeContext context.Context

  waitGroup *errgroup.Group

  InsertReportManager *InsertReportManager

  logger *zap.Logger
}

func CreateReportAnalyzer(
  db driver.Conn,
  config DatabaseConfiguration,
  parentContext context.Context,
  logger *zap.Logger,
) (*ReportAnalyzer, error) {
  group, analyzeContext := errgroup.WithContext(parentContext)
  group.SetLimit(1)

  insertReportManager, err := NewInsertReportManager(db, config, analyzeContext, "report", env.GetInt("INSERT_WORKER_COUNT", -1), logger)
  if err != nil {
    return nil, err
  }

  analyzer := &ReportAnalyzer{
    config:         config,
    insertQueue:    make(chan *ReportInfo, 1024),
    analyzeContext: analyzeContext,
    waitGroup:      group,

    InsertReportManager: insertReportManager,
    logger:              logger,
  }

  go func() {
    for {
      report, ok := <-analyzer.insertQueue
      if !ok {
        logger.Debug("analyze stopped", zap.String("reason", "insert queue is closed"))
        return
      }

      group.Go(func() error {
        return analyzer.insert(report)
      })
    }
  }()
  return analyzer, nil
}

func OpenDb(clickHouseUrl string, config DatabaseConfiguration) (driver.Conn, error) {
  // well, go-faster/ch is not so easy to use for such a generic case as our code (each column should be created in advance, no API to simply pass slice of any values)
  db, err := clickhouse.Open(&clickhouse.Options{
    Addr: []string{clickHouseUrl},
    Auth: clickhouse.Auth{
      Database: config.DbName,
    },
    DialTimeout:     10 * time.Second,
    ConnMaxLifetime: time.Hour,
    Settings: map[string]interface{}{
      // https://github.com/ClickHouse/ClickHouse/issues/2833
      // ZSTD 19+ is used, read/write timeout should be quite large (10 minutes)
      "send_timeout":    30_000,
      "receive_timeout": 3000,
    },
  })
  return db, err
}

func (t *ReportAnalyzer) Analyze(data []byte, extraData model.ExtraData) error {
  if t.analyzeContext.Err() != nil {
    return nil
  }

  var err error

  runResult := &RunResult{
    RawReport: data,

    TcBuildId:          getNullIfEmpty(extraData.TcBuildId),
    TcBuildType:        extraData.TcBuildType,
    TcInstallerBuildId: getNullIfEmpty(extraData.TcInstallerBuildId),
    ReportFileName:     extraData.ReportFile,
    TriggeredBy:        extraData.TriggeredBy,
  }

  if t.config.DbName == "jbr" {
    ignore := analyzePerfJbrReport(runResult, extraData)
    if ignore {
      //ignore empty report
      return nil
    }
  } else {
    err = ReadReport(runResult, t.config, t.logger)
  }
  projectId := t.config.DbName
  if err != nil {
    return err
  }

  if runResult.Report == nil {
    // ignore report
    return nil
  }

  if len(extraData.ProductCode) == 0 {
    extraData.ProductCode = runResult.Report.ProductCode
  }
  if len(extraData.BuildNumber) == 0 {
    extraData.BuildNumber = runResult.Report.Build
  }

  if len(extraData.Machine) == 0 {
    return errors.New("machine is not specified")
  }

  runResult.Product = extraData.ProductCode
  runResult.Machine = extraData.Machine

  if runResult.GeneratedTime.IsZero() {
    runResult.GeneratedTime, err = computeGeneratedTime(runResult.Report, extraData)
    if err != nil {
      if extraData.CurrentBuildTime.IsZero() {
        return err
      } else {
        runResult.GeneratedTime = extraData.CurrentBuildTime
      }
    }
  }

  if t.config.HasInstallerField {
    runResult.BuildTime = extraData.BuildTime

    if len(extraData.BuildNumber) == 0 {
      t.logger.Error("buildNumber is missed")
      return nil
    }

    buildComponents := strings.Split(extraData.BuildNumber, ".")
    if len(buildComponents) == 2 {
      buildComponents = append(buildComponents, "0")
    }

    runResult.BuildC1, runResult.BuildC2, runResult.BuildC3, err = splitBuildNumber(buildComponents)
    if err != nil {
      return err
    }
  }

  if t.analyzeContext.Err() != nil {
    return nil
  }

  if len(extraData.TcBuildProperties) == 0 {
    runResult.branch = "master"
  } else {
    runResult.branch, err = getBranch(runResult, extraData, projectId, t.logger)
    if err != nil {
      return err
    }
  }

  t.insertQueue <- &ReportInfo{
    extraData: extraData,
    runResult: runResult,
  }
  return nil
}

func getBranch(runResult *RunResult, extraData model.ExtraData, projectId string, logger *zap.Logger) (string, error) {
  parser := structParsers.Get()
  defer structParsers.Put(parser)

  props, err := parser.ParseBytes(extraData.TcBuildProperties)
  if err != nil {
    return "", errors.WithStack(err)
  }

  branch := string(props.GetStringBytes("vcsroot.branch"))
  if len(branch) != 0 && projectId != "fleet" {
    return strings.TrimPrefix(branch, "refs/heads/"), nil
  }

  if projectId == "ij" {
    logger.Error("cannot infer branch from TC properties", zap.ByteString("tcBuildProperties", extraData.TcBuildProperties))
    return "", errors.New("cannot infer branch from TC properties")
  } else {
    //goland:noinspection SpellCheckingInspection
    branch = string(props.GetStringBytes("teamcity.build.branch"))
    if len(branch) != 0 && branch != "<default>" {
      return branch, nil
    }
    var isMaster = props.GetStringBytes("vcsroot.ijplatform_master_IntelliJMonorepo.branch")
    if len(isMaster) == 0 {
      // we check that the property doesn't exist so it is not a master
      if runResult.BuildC3 == 0 {
        return strconv.Itoa(runResult.BuildC1), nil
      } else {
        // we have EAP branch
        return strconv.Itoa(runResult.BuildC1) + "." + strconv.Itoa(runResult.BuildC2), nil
      }
    } else {
      return "master", nil
    }
  }
}

type ReportInfo struct {
  extraData model.ExtraData

  runResult *RunResult
}

func computeGeneratedTime(report *model.Report, extraData model.ExtraData) (time.Time, error) {
  if report.Generated == "" {
    if extraData.LastGeneratedTime.IsZero() {
      return time.Time{}, errors.New("generated time not in report and not provided explicitly")
    }
    return extraData.LastGeneratedTime, nil
  } else {
    parsedTime, err := ParseTime(report.Generated)
    if err != nil {
      return time.Time{}, err
    }
    return parsedTime, nil
  }
}

func (t *ReportAnalyzer) CloseChannel() {
  close(t.insertQueue)
}

func (t *ReportAnalyzer) WaitAnalyzeAndInsert() error {
  // execute as part of the waitGroup because we use waitGroup context for InsertReportManager - after calling of Wait, context will be cancelled
  t.waitGroup.Go(func() error {
    // ensure that data is written
    t.logger.Debug("wait for insert")
    return t.InsertReportManager.InsertManager.Close()
  })
  return t.waitGroup.Wait()
}

func (t *ReportAnalyzer) insert(report *ReportInfo) error {
  runResult := report.runResult

  if report.extraData.TcInstallerBuildId > 0 {
    err := t.InsertReportManager.insertInstallerManager.Insert(report.extraData.TcInstallerBuildId, report.extraData.Changes)
    if err != nil {
      return err
    }
  }

  err := t.InsertReportManager.Insert(runResult)
  if err != nil {
    if errors2.Is(err, context.Canceled) {
      return err
    } else {
      return errors.WithMessagef(err, "cannot insert report (teamcityBuildId=%d, reportPath=%s)", report.extraData.TcBuildId, report.extraData.ReportFile)
    }
  }
  return nil
}

func getNullIfEmpty(v int) int {
  if v <= 0 {
    return 0
  } else {
    return v
  }
}

func splitBuildNumber(buildComponents []string) (int, int, int, error) {
  buildC1, err := strconv.Atoi(buildComponents[0])
  if err != nil {
    return 0, 0, 0, errors.WithStack(err)
  }
  buildC2, err := strconv.Atoi(buildComponents[1])
  if err != nil {
    return 0, 0, 0, errors.WithStack(err)
  }
  buildC3, err := strconv.Atoi(buildComponents[2])
  if err != nil {
    return 0, 0, 0, errors.WithStack(err)
  }
  return buildC1, buildC2, buildC3, nil
}
