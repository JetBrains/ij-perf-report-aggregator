package analyzer

import (
  "context"
  "crypto/sha1"
  "encoding/base64"
  "github.com/bvinc/go-sqlite-lite/sqlite3"
  "github.com/develar/errors"
  "github.com/json-iterator/go"
  "github.com/tdewolff/minify/v2"
  "github.com/tdewolff/minify/v2/json"
  "go.uber.org/zap"
  "hash"
  "report-aggregator/pkg/model"
  "report-aggregator/pkg/util"
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
  db       *sqlite3.Conn

  insertReportStatement  *sqlite3.Stmt
  insertMachineStatement *sqlite3.Stmt

  hash hash.Hash

  logger *zap.Logger
}

func CreateReportAnalyzer(dbPath string, analyzeContext context.Context, logger *zap.Logger) (*ReportAnalyzer, error) {
  err := prepareDatabaseFile(dbPath, logger)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  db, err := prepareDatabase(dbPath, logger)
  if err != nil {
    return nil, err
  }

  m := minify.New()
  m.AddFunc("json", json.Minify)

  analyzer := &ReportAnalyzer{
    input:          make(chan *model.ReportInfo),
    analyzeContext: analyzeContext,
    waitChannel:    make(chan struct{}),
    ErrorChannel:   make(chan error),

    minifier: m,
    db:       db,
    hash:     sha1.New(),

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

        err := analyzer.doAnalyze(report)
        if err != nil {
          analyzer.ErrorChannel <- err
        }
      }
    }
  }()
  return analyzer, nil
}

func ReadReport(data []byte) (*model.Report, error) {
  var report model.Report
  err := jsoniter.ConfigFastest.Unmarshal(data, &report)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &report, nil
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

  reportInfo := &model.ReportInfo{Report:report, ExtraData:extraData}

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
    parsedTime, err := ParseTime(reportInfo.Report.Generated)
    if err != nil {
      return err
    }
    reportInfo.GeneratedTime = parsedTime.Unix()
  }
  return nil
}

func ParseTime(s string) (*time.Time, error) {
  parsedTime, err := time.Parse(time.RFC1123Z, s)
  if err != nil {
    parsedTime, err = time.Parse(time.RFC1123, s)
  }

  if err != nil {
    parsedTime, err = time.Parse("Jan 2, 2006, 3:04:05 PM MST", s)
  }

  if err != nil {
    parsedTime, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", s)
  }

  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &parsedTime, nil
}

func (t *ReportAnalyzer) Close() error {
  t.closeOnce.Do(func() {
    close(t.input)
  })

  statement := t.insertReportStatement
  if statement != nil {
    t.insertReportStatement = nil
    util.Close(statement, t.logger)
  }

  statement = t.insertMachineStatement
  if statement != nil {
    t.insertMachineStatement = nil
    util.Close(statement, t.logger)
  }

  db := t.db
  t.db = nil
  if db == nil {
    return nil
  }
  return errors.WithStack(db.Close())
}

func (t *ReportAnalyzer) Done() <-chan struct{} {
  go func() {
    t.waitGroup.Wait()
    close(t.waitChannel)
  }()
  return t.waitChannel
}

func (t *ReportAnalyzer) doAnalyze(report *model.ReportInfo) error {
  t.waitGroup.Add(1)
  defer t.waitGroup.Done()

  t.hash.Reset()
  _, err := t.hash.Write(report.RawData)
  if err != nil {
    return errors.WithStack(err)
  }

  id := base64.RawURLEncoding.EncodeToString(t.hash.Sum(nil))

  isReportAlreadyAdded, err := t.isReportAlreadyAdded(id)
  if err != nil {
    return errors.WithStack(err)
  }

  logger := t.logger.With(zap.String("id", id), zap.String("generatedTime", time.Unix(report.GeneratedTime, 0).Format(time.RFC1123)))
  if isReportAlreadyAdded {
    logger.Info("report already processed")
    return nil
  }

  durationMetrics, instantMetrics := ComputeMetrics(report.Report, logger)
  // or both null, or not - no need to check each one
  if durationMetrics == nil || instantMetrics == nil {
    logger.Warn("report skipped (metrics cannot be computed)")
    return nil
  }

  buildComponents := strings.Split(report.ExtraData.BuildNumber, ".")
  if len(buildComponents) == 2 {
    buildComponents = append(buildComponents, "0")
  }

  buildC1, buildC2, buildC3, err := splitBuildNumber(buildComponents)
  if err != nil {
    return err
  }

  statement := t.insertReportStatement
  if statement == nil {
    statement, err = t.db.Prepare(string(MustAsset("insert-report.sql")))
    if err != nil {
      return errors.WithStack(err)
    }

    t.insertReportStatement = statement
  }

  machineId, err := t.getMachineId(report.ExtraData.Machine)
  if err != nil {
    return errors.WithStack(err)
  }

  var buildId interface{}
  if report.ExtraData.TcBuildId == 0 {
    buildId = nil
  } else {
    buildId = report.ExtraData.TcBuildId
  }

  err = statement.Exec(id, machineId, report.ExtraData.ProductCode,
    report.GeneratedTime, buildId,
    buildC1, buildC2, buildC3,
    report.RawData)
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("new report added")
  return nil
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

// https://stackoverflow.com/questions/13244393/sqlite-insert-or-ignore-and-return-original-rowid
func (t *ReportAnalyzer) getMachineId(machineName string) (int, error) {
  statement := t.insertMachineStatement
  var err error
  if statement == nil {
    statement, err = t.db.Prepare("insert or ignore into machine(name) values(?)")
    if err != nil {
      return -1, errors.WithStack(err)
    }

    t.insertMachineStatement = statement
  }

  err = statement.Exec(machineName)
  if err != nil {
    return -1, errors.WithStack(err)
  }

  statement, err = t.db.Prepare("select ROWID from machine where name = ?")
  defer util.Close(statement, t.logger)

  if err != nil {
    return -1, errors.WithStack(err)
  }

  err = statement.BindString(machineName, 0)
  if err != nil {
    return -1, errors.WithStack(err)
  }

  _, err = statement.Step()
  if err != nil {
    return -1, errors.WithStack(err)
  }

  id, _, err := statement.ColumnInt(0)
  return id, err
}

func (t *ReportAnalyzer) isReportAlreadyAdded(id string) (bool, error) {
  statement, err := t.db.Prepare(`select count() from report where id = ? limit 1`, id)
  if err != nil {
    return false, errors.WithStack(err)
  }

  defer util.Close(statement, t.logger)

  hasRow, err := statement.Step()
  if err != nil {
    return false, errors.WithStack(err)
  }

  if hasRow {
    result, _, err := statement.ColumnInt(0)
    return result == 1, errors.WithStack(err)
  }

  return false, nil
}

func (t *ReportAnalyzer) GetLastGeneratedTime() (int64, error) {
  statement, err := t.db.Prepare(`select max(generated_time) from report`)
  if err != nil {
    return 1, errors.WithStack(err)
  }

  defer util.Close(statement, t.logger)

  hasRow, err := statement.Step()
  if err != nil {
    return -1, errors.WithStack(err)
  }

  if hasRow {
    result, _, err := statement.ColumnInt64(0)
    return result, errors.WithStack(err)
  }

  return -1, nil
}
