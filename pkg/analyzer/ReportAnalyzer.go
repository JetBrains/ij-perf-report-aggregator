package analyzer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.deanishe.net/env"
)

type ReportAnalyzer struct {
	config DatabaseConfiguration

	insertQueue chan *ReportInfo

	analyzeContext context.Context

	waitGroup sync.WaitGroup
	cancel    func()
	errOnce   sync.Once
	err       error

	InsertReportManager *InsertReportManager

	logger *slog.Logger
}

func CreateReportAnalyzer(parentContext context.Context, db driver.Conn, metaDb *pgxpool.Pool, config DatabaseConfiguration, logger *slog.Logger) (*ReportAnalyzer, error) {
	insertReportManager, err := NewInsertReportManager(parentContext, db, metaDb, config, "report", env.GetInt("INSERT_WORKER_COUNT", -1))
	if err != nil {
		return nil, err
	}

	analyzeContext, cancel := context.WithCancel(parentContext)

	analyzer := &ReportAnalyzer{
		config:         config,
		insertQueue:    make(chan *ReportInfo, 1024),
		analyzeContext: analyzeContext,
		cancel:         cancel,

		InsertReportManager: insertReportManager,
		logger:              logger,
	}

	go func() {
		for {
			report, ok := <-analyzer.insertQueue
			if !ok {
				logger.Debug("analyze stopped; insert queue is closed")
				return
			}

			analyzer.invokeInsert(report, cancel)
		}
	}()
	return analyzer, nil
}

func (t *ReportAnalyzer) invokeInsert(report *ReportInfo, cancel context.CancelFunc) {
	defer t.waitGroup.Done()
	err := t.insert(report)
	if err != nil {
		t.errOnce.Do(func() {
			t.err = err
			cancel()
		})
	}
}

func OpenDb(clickHouseUrl string, config DatabaseConfiguration) (driver.Conn, *pgxpool.Pool, error) {
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
	if err != nil {
		return nil, nil, fmt.Errorf("cannot connect to clickhouse: %w", err)
	}

	metaDb, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	metaDb.Config().MaxConns = 10
	if err != nil {
		return nil, nil, fmt.Errorf("cannot create pool: %w", err)
	}

	return db, metaDb, err
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

	if isBisectRun(extraData, t.logger) {
		return nil
	}

	switch t.config.DbName {
	case "jbr":
		ignore := analyzePerfJbrReport(runResult, extraData)
		if ignore {
			// ignore empty report
			return nil
		}
	case "bazel":
		ignore := analyzePerfBazelReport(runResult)
		if ignore {
			// ignore empty report
			return nil
		}
	case "qodana":
		ignore := analyzeQodanaReport(runResult, extraData)
		if ignore {
			// ignore empty report
			return nil
		}
	default:
		err = ReadReport(runResult, t.config)

	}
	projectId := t.config.DbName
	if err != nil {
		return err
	}

	if runResult.Report == nil {
		// ignore report
		return nil
	}

	if extraData.ProductCode == "" {
		extraData.ProductCode = runResult.Report.ProductCode
	}
	if extraData.BuildNumber == "" {
		extraData.BuildNumber = runResult.Report.Build
	}

	if extraData.Machine == "" {
		return errors.New("machine is not specified")
	}

	runResult.Product = extraData.ProductCode
	runResult.Machine = extraData.Machine

	if runResult.GeneratedTime.IsZero() {
		runResult.GeneratedTime, err = computeGeneratedTime(runResult.Report, extraData)
		if err != nil {
			if extraData.CurrentBuildTime.IsZero() {
				return err
			}
			runResult.GeneratedTime = extraData.CurrentBuildTime
		}
	}

	if t.config.HasInstallerField {
		runResult.BuildTime = extraData.BuildTime

		if extraData.BuildNumber == "" {
			t.logger.Error("buildNumber is missed")
			return nil
		}

		buildComponents := strings.Split(extraData.BuildNumber, ".")
		if len(buildComponents) == 2 {
			buildComponents = append(buildComponents, "0")
		}

		runResult.BuildC1, runResult.BuildC2, runResult.BuildC3, err = splitBuildNumber(buildComponents)
		if err != nil {
			// we might get 231.snapshot build numbers, that is more or less fine and we need such build anyway
			t.logger.Error(err.Error())
		}
	}

	if t.config.HasBuildNumber {
		runResult.BuildNumber = extraData.TcBuildNumber
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

	t.waitGroup.Add(1)
	t.insertQueue <- &ReportInfo{
		extraData: extraData,
		runResult: runResult,
	}
	return nil
}

func isBisectRun(extraData model.ExtraData, logger *slog.Logger) bool {
	parser := parserPool.Get()
	defer parserPool.Put(parser)

	props, err := parser.ParseBytes(extraData.TcBuildProperties)
	if err != nil {
		logger.Warn("failed to parse build properties", "error", err)
		return false
	}
	return props.GetBool("env.IS_BISECT_RUN")
}

func getBranch(runResult *RunResult, extraData model.ExtraData, projectId string, logger *slog.Logger) (string, error) {
	parser := parserPool.Get()
	defer parserPool.Put(parser)

	props, err := parser.ParseBytes(extraData.TcBuildProperties)
	if err != nil {
		return "", fmt.Errorf("failed to parse build properties: %w", err)
	}

	if projectId == "mlEvaluation" {
		return "master", nil
	}
	if projectId == "jbr" {
		splitId := strings.SplitN(extraData.TcBuildType, "_", 4)
		if len(splitId) == 4 {
			jbrBranch := strings.ToLower(splitId[1]) + "_" + strings.ToLower(splitId[2])
			return jbrBranch, nil
		}
		logger.Error("format of JBR project is unexpected", "teamcity.project.id", extraData.TcBuildType)
		return "", errors.New("cannot infer branch from JBR project id")
	}
	if projectId == "qodana" {
		qodanaImage := string(props.GetStringBytes("image"))
		lastSlash := strings.LastIndex(qodanaImage, "/")

		if lastSlash >= 0 {
			return qodanaImage[lastSlash+1:], nil
		}
		logger.Warn("No slash found in string")
	}

	//goland:noinspection SpellCheckingInspection
	branch := string(props.GetStringBytes("teamcity.build.branch"))
	if branch != "" && branch != "<default>" {
		return branch, nil
	}
	branchInt := props.GetInt("teamcity.build.branch")
	if branchInt != 0 {
		return strconv.Itoa(branchInt), nil
	}
	isMaster := props.GetStringBytes("vcsroot.ijplatform_master_IntelliJMonorepo.branch")
	if len(isMaster) == 0 {
		// we check that the property doesn't exist so it is not a master
		if runResult.BuildC3 == 0 {
			return strconv.Itoa(runResult.BuildC1), nil
		}
		// we have EAP branch
		return strconv.Itoa(runResult.BuildC1) + "." + strconv.Itoa(runResult.BuildC2), nil
	}
	return "master", nil
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
	}
	parsedTime, err := ParseTime(report.Generated)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

func (t *ReportAnalyzer) WaitAnalyzeAndInsert() error {
	t.logger.Debug("wait for analyze")
	t.waitGroup.Wait()
	t.cancel()
	if t.err != nil {
		return t.err
	}

	close(t.insertQueue)

	t.logger.Debug("wait for insert")
	err := t.InsertReportManager.InsertManager.Close()
	if err != nil {
		return err
	}

	return nil
}

func (t *ReportAnalyzer) insert(report *ReportInfo) error {
	runResult := report.runResult

	if report.extraData.TcInstallerBuildId > 0 {
		err := t.InsertReportManager.insertInstallerManager.Insert(report.extraData.TcInstallerBuildId, report.extraData.Changes)
		if err != nil {
			return err
		}
	} else if report.extraData.Changes != nil {
		err := t.InsertReportManager.insertInstallerManager.Insert(report.extraData.TcBuildId, report.extraData.Changes)
		if err != nil {
			return err
		}
	}

	if t.InsertReportManager.insertMetaManager != nil {
		r := report.runResult.Report
		err := t.InsertReportManager.insertMetaManager.InsertProjectDescription(r.Project, runResult.branch, r.ProjectURL, r.MethodName, r.ProjectDescription)
		if err != nil {
			t.logger.Warn("cannot insert project description", "error", err)
		}
	}

	err := t.InsertReportManager.Insert(runResult)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return err
		}
		return fmt.Errorf("cannot insert report (teamcityBuildId=%d, reportPath=%s): %w", report.extraData.TcBuildId, report.extraData.ReportFile, err)
	}
	return nil
}

func getNullIfEmpty(v int) int {
	if v <= 0 {
		return 0
	}
	return v
}

func splitBuildNumber(buildComponents []string) (int, int, int, error) {
	buildC1, err := strconv.Atoi(buildComponents[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot parse build number: %w", err)
	}
	buildC2, err := strconv.Atoi(buildComponents[1])
	if err != nil {
		return buildC1, 0, 0, nil
	}
	buildC3, err := strconv.Atoi(buildComponents[2])
	if err != nil {
		return buildC1, buildC2, 0, nil
	}
	return buildC1, buildC2, buildC3, nil
}
