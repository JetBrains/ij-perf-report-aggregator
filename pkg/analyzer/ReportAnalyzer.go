package analyzer

import (
  "context"
  "crypto/sha1"
  "encoding/base64"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  _ "github.com/mattn/go-sqlite3"
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
  db       *sqlx.DB

  insertReportStatement  *sqlx.Stmt
  insertMachineStatement *sqlx.Stmt

  hash hash.Hash

  logger *zap.Logger
}

func CreateReportAnalyzer(dbPath string, analyzeContext context.Context, logger *zap.Logger) (*ReportAnalyzer, error) {
  err := prepareDatabaseFile(dbPath)
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
    //err = t.db.Exec("delete from report where id = ?", id)
    //if err != nil {
    //  return errors.WithStack(err)
    //}
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
    statement, err = t.db.Preparex(string(MustAsset("insert-report.sql")))
    if err != nil {
      return errors.WithStack(err)
    }

    t.insertReportStatement = statement
  }

  machineId, err := t.getMachineId(report.ExtraData.Machine)
  if err != nil {
    return errors.WithStack(err)
  }

  buildTimeFromReport, err := GetBuildTimeFromReport(report.Report)
  if err != nil {
    return errors.WithStack(err)
  }

  if report.ExtraData.TcInstallerBuildId > 0 {
    err = t.insertInstallerIdIfMissed(&report.ExtraData)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  _, err = statement.ExecContext(t.analyzeContext, id, machineId, report.ExtraData.ProductCode,
    report.GeneratedTime, chooseNotEmpty(buildTimeFromReport, report.ExtraData.BuildTime),
    getNullIfEmpty(report.ExtraData.TcBuildId), getNullIfEmpty(report.ExtraData.TcInstallerBuildId), report.ExtraData.TcBuildProperties,
    buildC1, buildC2, buildC3,
    report.RawData)
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("new report added")
  return nil
}

func chooseNotEmpty(v1 int64, v2 int64) int64 {
  if v1 <= 0 {
    return v2
  } else {
    return v1
  }
}

func getNullIfEmpty(v int) interface{} {
  if v <= 0 {
    return nil
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

// https://stackoverflow.com/questions/13244393/sqlite-insert-or-ignore-and-return-original-rowid
func (t *ReportAnalyzer) getMachineId(machineName string) (int, error) {
  insertMachineStatement := t.insertMachineStatement
  var err error
  if insertMachineStatement == nil {
    insertMachineStatement, err = t.db.Preparex("insert or ignore into machine(name) values(?)")
    if err != nil {
      return -1, errors.WithStack(err)
    }

    t.insertMachineStatement = insertMachineStatement
  }

  _, err = insertMachineStatement.Exec(machineName)
  if err != nil {
    return -1, errors.WithStack(err)
  }

  statement, err := t.db.Preparex("select ROWID from machine where name = ?")
  defer util.Close(statement, t.logger)

  if err != nil {
    return -1, errors.WithStack(err)
  }

  var id int
  err = statement.Get(&id, machineName)
  if err != nil {
    return -1, errors.WithStack(err)
  }
  return id, err
}

func (t *ReportAnalyzer) insertInstallerIdIfMissed(extraData *model.ExtraData) error {
  _, err := t.db.Exec("insert or ignore into installer(id, changes) values(?, ?)", extraData.TcInstallerBuildId, extraData.Changes)
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

func (t *ReportAnalyzer) isReportAlreadyAdded(id string) (bool, error) {
  var result int
  err := t.db.Get(&result, `select count() from report where id = ?`, id)
  if err != nil {
    return false, errors.WithStack(err)
  }
  return result == 1, nil
}

func (t *ReportAnalyzer) GetLastGeneratedTime() (int64, error) {
  var result int64
  err := t.db.Get(&result, `select max(generated_time) from report`)
  if err != nil {
    return -1, errors.WithStack(err)
  }
  return result, nil
}
