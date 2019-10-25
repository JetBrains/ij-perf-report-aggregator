package analyzer

import (
  "context"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  _ "github.com/kshvakov/clickhouse"
  "github.com/tdewolff/minify/v2"
  "github.com/tdewolff/minify/v2/json"
  "github.com/valyala/fastjson"
  "go.uber.org/multierr"
  "go.uber.org/zap"
  "report-aggregator/pkg/model"
  "strconv"
  "strings"
  "sync"
  "time"
)

type ReportAnalyzer struct {
  input        chan *model.ReportInfo
  waitChannel  chan struct{}
  ErrorChannel chan error

  analyzeContext context.Context

  waitGroup sync.WaitGroup
  closeOnce sync.Once

  minifier *minify.M
  db       *sqlx.DB

  insertInstallerManager *InsertInstallerManager
  insertReportManager    *InsertReportManager

  logger *zap.Logger
}

func CreateReportAnalyzer(clickHouseUrl string, analyzeContext context.Context, logger *zap.Logger) (*ReportAnalyzer, error) {
  m := minify.New()
  m.AddFunc("json", json.Minify)

  // ZSTD 19 is used, read/write timeout should be quite large (10 minutes)
  db, err := sqlx.Open("clickhouse", "tcp://"+clickHouseUrl+"?read_timeout=600&write_timeout=600&compress=1")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  installerManager, err := NewInstallerInsertManager(db, logger)
  if err != nil {
    return nil, err
  }

  insertReportManager, err := NewInsertReportManager(db, analyzeContext, logger)
  if err != nil {
    return nil, err
  }

  analyzer := &ReportAnalyzer{
    input:          make(chan *model.ReportInfo),
    analyzeContext: analyzeContext,
    waitChannel:    make(chan struct{}),
    ErrorChannel:   make(chan error),

    insertInstallerManager: installerManager,
    insertReportManager:    insertReportManager,

    minifier: m,
    db:       db,

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
          // commit installer ids to ensure that correctly added reports have corresponding installer ids
          err2 := analyzer.insertInstallerManager.CommitAndWait()
          if err2 != nil {
            logger.Error("cannot commit installers", zap.Error(err2))
          }

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
  return errors.WithStack(multierr.Combine(t.insertInstallerManager.Close(), t.insertReportManager.Close(), t.db.Close()))
}

func (t *ReportAnalyzer) Done() <-chan struct{} {
  go func() {
    t.waitGroup.Wait()

    err := multierr.Combine(t.insertInstallerManager.CommitAndWait(), t.insertReportManager.CommitAndWait())
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

  logger := t.logger.With(zap.String("generatedTime", time.Unix(report.GeneratedTime, 0).Format(time.RFC1123)))

  durationMetrics, instantMetrics := ComputeMetrics(report.Report, logger)
  // or both null, or not - no need to check each one
  if durationMetrics == nil || instantMetrics == nil {
    logger.Warn("report skipped (metrics cannot be computed)")
    return nil
  }

  reportRow := &MetricResult{
    Product: report.ExtraData.ProductCode,
    Machine: report.ExtraData.Machine,

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

  reportRow.BuildTime, err = GetBuildTimeFromReport(report.Report)
  if err != nil {
    return err
  }
  if reportRow.BuildTime <= 0 {
    reportRow.BuildTime = report.ExtraData.BuildTime
  }

  if report.ExtraData.TcInstallerBuildId > 0 {
    err = t.insertInstallerManager.Insert(report.ExtraData.TcInstallerBuildId, report.ExtraData.Changes)
    if err != nil {
      return err
    }
  }

  err = t.insertReportManager.Insert(reportRow, branch)
  if err != nil {
    return errors.WithStack(err)
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
