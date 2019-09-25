package analyzer

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/develar/errors"
	"github.com/json-iterator/go"
	"github.com/mcuadros/go-version"
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
		input:          make(chan *model.Report),
		analyzeContext: analyzeContext,
		waitChannel:    make(chan struct{}),
		ErrorChannel:   make(chan error),

		minifier: m,
		db:       db,
		hash:     sha1.New(),

		logger: logger,
	}

	analyzer.insertStatement, err = db.Prepare(`INSERT INTO report (id, generated_time, metrics_version, metrics, raw_report) VALUES (?, ?, ?, ?, ?)`)
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

	err = t.insertStatement.Exec(id, report.GeneratedTime, metricsVersion, serializedMetrics, report.RawData)
	if err != nil {
		return errors.WithStack(err)
	}

	logger.Info("new report added")
	return nil
}

func (t *ReportAnalyzer) computeMetrics(report *model.Report, logger *zap.Logger) *model.Metrics {
	metrics := &model.Metrics{
		Bootstrap: -1,
		Splash:    -1,

		AppInitPreparation:       -1,
		AppInit:                  -1,
		PluginDescriptorsLoading: -1,

		AppComponentCreation:     -1,
		ProjectComponentCreation: -1,
		ModuleLoading:            -1,
	}

	if version.Compare(report.Version, "12", ">=") && len(report.TraceEvents) == 0 {
		logger.Warn("invalid report (due to opening second project?), report will be skipped")
		return nil
	}

	// v < 12: PluginDescriptorsLoading can be or in MainActivities, or in PrepareAppInitActivities

	for _, activity := range report.MainActivities {
		switch activity.Name {
		case "bootstrap":
			metrics.Bootstrap = activity.Duration

		case "app initialization preparation":
			metrics.AppInitPreparation = activity.Duration
		case "app initialization":
			metrics.AppInit = activity.Duration
		case "plugin descriptors loading":
			metrics.PluginDescriptorsLoading = activity.Duration

		case "app component creation":
			metrics.AppComponentCreation = activity.Duration
		case "project component creation":
			metrics.ProjectComponentCreation = activity.Duration
		case "module loading":
			metrics.ModuleLoading = activity.Duration
		}
	}

	if version.Compare(report.Version, "11", "<") {
		for _, activity := range report.PrepareAppInitActivities {
			switch activity.Name {
			case "plugin descriptors loading":
				metrics.PluginDescriptorsLoading = activity.Start
			case "splash initialization":
				metrics.Splash = activity.Start
			}
		}
	} else {
		for _, activity := range report.TraceEvents {
			if activity.Phase == "i" && (activity.Name == "splash" || activity.Name == "splash shown") {
				metrics.Splash = activity.Timestamp / 1000
			}
		}
	}

	if metrics.Bootstrap == -1 {
		logRequiredMetricNotFound(logger, "bootstrap")
		return nil
	}
	if metrics.PluginDescriptorsLoading == -1 {
		logRequiredMetricNotFound(logger, "pluginDescriptorsLoading")
		return nil
	}
	if metrics.AppComponentCreation == -1 {
		logRequiredMetricNotFound(logger, "AppComponentCreation")
		return nil
	}
	if metrics.ModuleLoading == -1 {
		logRequiredMetricNotFound(logger, "ModuleLoading")
		return nil
	}
	return metrics
}

func logRequiredMetricNotFound(logger *zap.Logger, metricName string) {
	logger.Error("metric is required, but not found, report will be skipped", zap.String("metric", metricName))
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
