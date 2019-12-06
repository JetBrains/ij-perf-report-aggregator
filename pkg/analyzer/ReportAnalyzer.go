package analyzer

import (
  "context"
  _ "github.com/ClickHouse/clickhouse-go"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/tdewolff/minify/v2"
  "github.com/tdewolff/minify/v2/json"
  "github.com/valyala/fastjson"
  "go.uber.org/multierr"
  "go.uber.org/zap"
  "strconv"
  "strings"
  "sync"
)

type ReportAnalyzer struct {
  input        chan *model.ReportInfo
  waitChannel  chan struct{}
  ErrorChannel chan error

  analyzeContext context.Context

  waitGroup sync.WaitGroup
  closeOnce sync.Once

  minifier *minify.M
  Db       *sqlx.DB

  insertReportManager *InsertReportManager

  logger *zap.Logger
}

func CreateReportAnalyzer(clickHouseUrl string, analyzeContext context.Context, logger *zap.Logger) (*ReportAnalyzer, error) {
  m := minify.New()
  m.AddFunc("json", json.Minify)

  // https://github.com/ClickHouse/ClickHouse/issues/2833
  // ZSTD 19 is used, read/write timeout should be quite large (10 minutes)
  db, err := sqlx.Open("clickhouse", "tcp://"+clickHouseUrl+"?read_timeout=600&write_timeout=600&debug=0  &compress=1&send_timeout=30000&receive_timeout=3000")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  insertReportManager, err := NewInsertReportManager(db, analyzeContext, "report", env.GetInt("INSERT_WORKER_COUNT", -1), logger)
  if err != nil {
    return nil, err
  }

  analyzer := &ReportAnalyzer{
    input:          make(chan *model.ReportInfo, 32),
    analyzeContext: analyzeContext,
    waitChannel:    make(chan struct{}),
    ErrorChannel:   make(chan error),

    insertReportManager:    insertReportManager,

    minifier: m,
    Db:       db,

    logger: logger,
  }

  go func() {
    for {
      select {
      case <-analyzeContext.Done():
        logger.Debug("analyze stopped")
        return

      case report, ok := <-analyzer.input:
        if !ok {
          return
        }

        err := analyzer.insert(report)
        if err != nil {
          analyzer.ErrorChannel <- err
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

  report, err := ReadReport(data)
  if err != nil {
    return err
  }

  if len(extraData.ProductCode) == 0 {
    extraData.ProductCode = report.ProductCode
  }
  if len(extraData.BuildNumber) == 0 {
    extraData.BuildNumber = report.Build
  }

  if len(extraData.Machine) == 0 {
    return errors.New("machine is not specified")
  }

  reportInfo := &model.ReportInfo{Report: report, ExtraData: extraData}

  if t.analyzeContext.Err() != nil {
    return nil
  }

  // normalize to compute consistent unique id
  reportInfo.RawData, err = t.minifier.Bytes("json", data)
  if err != nil {
    return err
  }

  if t.analyzeContext.Err() != nil {
    return nil
  }

  err = computeGeneratedTime(reportInfo, extraData)
  if err != nil {
    return err
  }

  t.input <- reportInfo
  return nil
}

func computeGeneratedTime(reportInfo *model.ReportInfo, extraData model.ExtraData) error {
  if reportInfo.Report.Generated == "" {
    if extraData.LastGeneratedTime <= 0 {
      return errors.New("generated time not in report and not provided explicitly")
    }
    reportInfo.GeneratedTime = extraData.LastGeneratedTime
  } else {
    parsedTime, err := parseTime(reportInfo.Report.Generated)
    if err != nil {
      return err
    }
    reportInfo.GeneratedTime = parsedTime.Unix()
  }
  return nil
}

func (t *ReportAnalyzer) Close() error {
  t.closeOnce.Do(func() {
    close(t.input)
  })
  return errors.WithStack(multierr.Combine(t.insertReportManager.insertInstallerManager.Close(), t.insertReportManager.Close(), t.Db.Close()))
}

func (t *ReportAnalyzer) Wait() {
  t.waitGroup.Wait()
}

func (t *ReportAnalyzer) Done() <-chan struct{} {
  go func() {
    t.waitGroup.Wait()

    err := t.insertReportManager.InsertManager.CommitAndWait()
    if err != nil {
      t.ErrorChannel <- err
    }

    close(t.waitChannel)
  }()
  return t.waitChannel
}

func (t *ReportAnalyzer) insert(report *model.ReportInfo) error {
  t.waitGroup.Add(1)
  defer t.waitGroup.Done()

  reportRow := &MetricResult{
    Product: report.ExtraData.ProductCode,
    Machine: report.ExtraData.Machine,

    BuildTime: report.ExtraData.BuildTime,

    GeneratedTime:      report.GeneratedTime,
    TcBuildId:          getNullIfEmpty(report.ExtraData.TcBuildId),
    TcInstallerBuildId: getNullIfEmpty(report.ExtraData.TcInstallerBuildId),
    TcBuildProperties:  report.ExtraData.TcBuildProperties,

    RawReport: report.RawData,
  }

  var branch string
  if len(report.ExtraData.TcBuildProperties) == 0 {
    branch = "master"
  } else {
    //noinspection SpellCheckingInspection
    branch = fastjson.GetString(report.ExtraData.TcBuildProperties, "vcsroot.branch")
    if len(branch) == 0 {
      t.logger.Error("cannot infer branch from TC properties", zap.ByteString("tcBuildProperties", report.ExtraData.TcBuildProperties))
      return errors.New("cannot infer branch from TC properties")
    }
  }

  buildComponents := strings.Split(report.ExtraData.BuildNumber, ".")
  if len(buildComponents) == 2 {
    buildComponents = append(buildComponents, "0")
  }

  var err error
  reportRow.BuildC1, reportRow.BuildC2, reportRow.BuildC3, err = splitBuildNumber(buildComponents)
  if err != nil {
    return err
  }

  if report.ExtraData.TcInstallerBuildId > 0 {
    err = t.insertReportManager.insertInstallerManager.Insert(report.ExtraData.TcInstallerBuildId, report.ExtraData.Changes)
    if err != nil {
      return err
    }
  }

  providedProject := report.Report.Project
  //if len(providedProject) == 0 {
  //  providedProject = "/q9N7EHxr8F1NHjbNQnpqb0Q0fs"
  //}

  err = t.insertReportManager.Insert(reportRow, branch, providedProject)
  if err != nil {
    return err
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

func (t *ReportAnalyzer) GetLastGeneratedTime() int64 {
  return t.insertReportManager.MaxGeneratedTime
}
