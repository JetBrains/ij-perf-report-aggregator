package analyzer

import (
  "context"
  _ "github.com/ClickHouse/clickhouse-go"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/valyala/fastjson"
  "go.deanishe.net/env"
  "go.uber.org/atomic"
  "go.uber.org/multierr"
  "go.uber.org/zap"
  "strconv"
  "strings"
  "sync"
)

type ReportAnalyzer struct {
  dbName   string
  analyzer CustomReportAnalyzer

  insertQueue chan *ReportInfo
  DoneChannel chan error

  firstError atomic.Value

  analyzeContext context.Context

  waitGroup sync.WaitGroup
  closeOnce sync.Once

  Db *sqlx.DB

  InsertReportManager *InsertReportManager

  logger *zap.Logger
}

func CreateReportAnalyzer(
  clickHouseUrl string,
  dbName string,
  reportAnalyzer CustomReportAnalyzer,
  analyzeContext context.Context,
  logger *zap.Logger,
  cancelOnError context.CancelFunc,
) (*ReportAnalyzer, error) {

  // https://github.com/ClickHouse/ClickHouse/issues/2833
  // ZSTD 19+ is used, read/write timeout should be quite large (10 minutes)
  db, err := sqlx.Open("clickhouse", "tcp://"+clickHouseUrl+"?read_timeout=600&write_timeout=600&debug=0&compress=1&send_timeout=30000&receive_timeout=3000&database="+dbName)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  insertReportManager, err := NewInsertReportManager(db, dbName, analyzeContext, "report", env.GetInt("INSERT_WORKER_COUNT", -1), logger)
  if err != nil {
    return nil, err
  }

  analyzer := &ReportAnalyzer{
    dbName:         dbName,
    analyzer:       reportAnalyzer,
    insertQueue:    make(chan *ReportInfo, 32),
    analyzeContext: analyzeContext,
    // buffered channel - do not block send if no one yet read
    DoneChannel: make(chan error, 1),

    waitGroup: sync.WaitGroup{},

    InsertReportManager: insertReportManager,

    Db: db,

    logger: logger,
  }

  go func() {
    for {
      select {
      case <-analyzeContext.Done():
        logger.Debug("analyze stopped", zap.String("reason", "context cancelled"))
        return

      case report, ok := <-analyzer.insertQueue:
        if !ok {
          return
        }

        err = analyzer.insert(report)
        if err != nil {
          logger.Error("analyze stopped", zap.String("reason", "error occurred"), zap.Error(err))
          analyzer.firstError.Store(err)
          // first, publish error
          analyzer.DoneChannel <- err
          cancelOnError()

          // do not process insertQueue anymore
          return
        }
      }
    }
  }()
  return analyzer, nil
}

func (t *ReportAnalyzer) Analyze(data []byte, extraData model.ExtraData) error {
  if t.analyzeContext.Err() != nil {
    return nil
  }

  var err error

  runResult := &RunResult{
    RawReport: data,

    TcBuildId:          getNullIfEmpty(extraData.TcBuildId),
    TcInstallerBuildId: getNullIfEmpty(extraData.TcInstallerBuildId),
    TcBuildProperties:  extraData.TcBuildProperties,
  }

  err = ReadReport(runResult, t.analyzer, t.logger)
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
  runResult.BuildTime = extraData.BuildTime

  runResult.GeneratedTime, err = computeGeneratedTime(runResult.Report, extraData)
  if err != nil {
    return err
  }

  if len(extraData.BuildNumber) == 0 && t.dbName != "ij" {
    // ignore report (remove later once fleet project will be collected)
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

  if t.analyzeContext.Err() != nil {
    return nil
  }

  if len(extraData.TcBuildProperties) == 0 {
    runResult.branch = "master"
  } else {
    //noinspection SpellCheckingInspection
    runResult.branch = fastjson.GetString(extraData.TcBuildProperties, "vcsroot.branch")
    if len(runResult.branch) == 0 {
      if t.dbName == "ij" {
        t.logger.Error("cannot infer branch from TC properties", zap.ByteString("tcBuildProperties", extraData.TcBuildProperties))
        return errors.New("cannot infer branch from TC properties")
      } else {
        runResult.branch = "master"
      }
    } else {
      runResult.branch = strings.TrimPrefix(runResult.branch, "refs/heads/")
    }
  }

  t.waitGroup.Add(1)
  t.insertQueue <- &ReportInfo{
    extraData: extraData,
    runResult: runResult,
  }
  return nil
}

type ReportInfo struct {
  extraData model.ExtraData

  runResult *RunResult
}

func computeGeneratedTime(report *model.Report, extraData model.ExtraData) (int64, error) {
  if report.Generated == "" {
    if extraData.LastGeneratedTime <= 0 {
      return -1, errors.New("generated time not in report and not provided explicitly")
    }
    return extraData.LastGeneratedTime, nil
  } else {
    parsedTime, err := parseTime(report.Generated)
    if err != nil {
      return -1, err
    }
    return parsedTime.Unix(), nil
  }
}

func (t *ReportAnalyzer) Close() error {
  t.closeOnce.Do(func() {
    close(t.insertQueue)
  })
  return errors.WithStack(multierr.Combine(t.InsertReportManager.insertInstallerManager.Close(), t.InsertReportManager.Close(), t.Db.Close()))
}

func (t *ReportAnalyzer) WaitAndCommit() <-chan error {
  go func() {
    t.waitGroup.Wait()

    if t.firstError.Load() != nil {
      return
    }

    err := t.InsertReportManager.InsertManager.CommitAndWait()
    t.DoneChannel <- err
  }()
  return t.DoneChannel
}

func (t *ReportAnalyzer) insert(report *ReportInfo) error {
  defer t.waitGroup.Done()

  runResult := report.runResult

  if report.extraData.TcInstallerBuildId > 0 {
    err := t.InsertReportManager.insertInstallerManager.Insert(report.extraData.TcInstallerBuildId, report.extraData.Changes)
    if err != nil {
      return err
    }
  }

  err := t.InsertReportManager.Insert(runResult)
  if err != nil {
    return errors.WithMessagef(err, "Cannot insert report (teamcityBuildId=%d, reportPath=%s)", report.extraData.TcBuildId, report.extraData.ReportFile)
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
