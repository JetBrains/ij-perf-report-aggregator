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
	"sync"
	"time"
)

type ReportAnalyzer struct {
	input        chan *model.Report
	waitChannel  chan struct{}
	ErrorChannel chan error

	analyzeContext context.Context

	waitGroup sync.WaitGroup
	closeOnce sync.Once

	minifier        *minify.M
	db              *sqlite3.Conn
	insertStatement *sqlite3.Stmt
	hash            hash.Hash

	machine string

	logger *zap.Logger
}

func CreateReportAnalyzer(dbPath string, machine string, analyzeContext context.Context, logger *zap.Logger) (*ReportAnalyzer, error) {
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
		input:          make(chan *model.Report),
		analyzeContext: analyzeContext,
		waitChannel:    make(chan struct{}),
		ErrorChannel:   make(chan error),

		minifier: m,
		db:       db,
		hash:     sha1.New(),

		logger: logger,

		machine: machine,
	}

	analyzer.insertStatement, err = db.Prepare(`INSERT INTO report (id, machine, generated_time, metrics_version, metrics, raw_report) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		util.Close(db, logger)
		return nil, err
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

func readReport(data []byte) (*model.Report, error) {
	var report model.Report
	err := jsoniter.ConfigFastest.Unmarshal(data, &report)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &report, nil
}

func (t *ReportAnalyzer) Analyze(data []byte, lastGeneratedTime int64) error {
	if t.analyzeContext.Err() != nil {
		return nil
	}

	report, err := readReport(data)
	if err != nil {
		return err
	}

	if t.analyzeContext.Err() != nil {
		return nil
	}

	// normalize to compute consistent unique id
	report.RawData, err = t.minifier.Bytes("json", data)
	if err != nil {
		return err
	}

	if t.analyzeContext.Err() != nil {
		return nil
	}

	err = computeGeneratedTime(report, lastGeneratedTime)
	if err != nil {
		return err
	}

	t.input <- report
	return nil
}

func computeGeneratedTime(report *model.Report, generatedTimeFromLog int64) error {
	if report.Generated == "" {
		if generatedTimeFromLog == -1 {
			return errors.New("Generated time not in report and not provided explicitly")
		}
		report.GeneratedTime = generatedTimeFromLog
	} else {
		parsedTime, err := parseTime(report)
		if err != nil {
			return err
		}
		report.GeneratedTime = parsedTime.Unix()
	}
	return nil
}

func parseTime(report *model.Report) (*time.Time, error) {
	parsedTime, err := time.Parse(time.RFC1123Z, report.Generated)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC1123, report.Generated)
		if err != nil {
			parsedTime, err = time.Parse("Jan 2, 2006, 3:04:05 PM MST", report.Generated)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
	}
	return &parsedTime, nil
}

func (t *ReportAnalyzer) Close() error {
	t.closeOnce.Do(func() {
		close(t.input)
	})

	insertStatement := t.insertStatement
	if insertStatement != nil {
		util.Close(insertStatement, t.logger)
		insertStatement = nil
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

const metricsVersion = 1

func (t *ReportAnalyzer) doAnalyze(report *model.Report) error {
	t.waitGroup.Add(1)
	defer t.waitGroup.Done()

	t.hash.Reset()
	t.hash.Write(report.RawData)

	id := base64.RawURLEncoding.EncodeToString(t.hash.Sum(nil))

	isAlreadyProcessed, err := t.isReportAlreadyProcessed(id)
	if err != nil {
		return errors.WithStack(err)
	}

	logger := t.logger.With(zap.String("id", id), zap.String("generatedTime", time.Unix(report.GeneratedTime, 0).Format(time.RFC1123)))
	if isAlreadyProcessed {
		logger.Info("report already processed")
		return nil
	}

	metrics := t.computeMetrics(report, logger)
	if metrics == nil {
		return nil
	}

	serializedMetrics, err := jsoniter.ConfigFastest.Marshal(metrics)
	if err != nil {
		return errors.WithStack(err)
	}

	err = t.insertStatement.Exec(id, t.machine, report.GeneratedTime, metricsVersion, serializedMetrics, report.RawData)
	if err != nil {
		return errors.WithStack(err)
	}

	logger.Info("new report added")
	return nil
}

func (t *ReportAnalyzer) isReportAlreadyProcessed(id string) (bool, error) {
	stmt, err := t.db.Prepare(`SELECT metrics_version FROM report WHERE id = ?`, id)
	if err != nil {
		return false, errors.WithStack(err)
	}

	defer util.Close(stmt, t.logger)

	hasRow, err := stmt.Step()
	if err != nil {
		return false, errors.WithStack(err)
	}
	return hasRow, nil
}
