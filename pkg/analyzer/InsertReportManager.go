package analyzer

import (
	"context"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.deanishe.net/env"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// use `select distinct cast(machine, 'Uint16') as id, machine as name FROM report order by id` to get current enum values
// for now machine enum should be updated manually if a new machine will be added

var ErrMetricsCannotBeComputed = errors.New("metrics cannot be computed")

type RunResult struct {
	Product string

	Machine string

	BuildTime     time.Time
	GeneratedTime time.Time

	TcBuildId          int
	TcBuildType        string
	TcInstallerBuildId int

	RawReport []byte

	// maybe null
	Report         *model.Report
	ReportFileName string

	BuildC1 int
	BuildC2 int
	BuildC3 int

	TriggeredBy string
	BuildNumber string

	ExtraFieldData []interface{}

	branch string
}

type InsertReportManager struct {
	sql_util.InsertDataManager

	context                          context.Context
	MaxGeneratedTime                 time.Time
	IsCheckThatNotAlreadyAddedNeeded bool
	config                           DatabaseConfiguration
	nonMetricFieldCount              int
	insertInstallerManager           *InsertInstallerManager
	insertMetaManager                *InsertMetaManager
	TableName                        string
}

func NewInsertReportManager(ctx context.Context, db driver.Conn, metaDb *pgxpool.Pool, config DatabaseConfiguration, tableName string, insertWorkerCount int) (*InsertReportManager, error) {
	var err error

	if config.TableName != "" {
		tableName = config.TableName
	}

	metaFields := make([]string, 0, 16)
	metaFields = append(metaFields, "machine", "generated_time", "project", "tc_build_id", "branch")
	if config.HasBuildTypeField {
		metaFields = append(metaFields, "tc_build_type")
	}
	if config.HasProductField {
		metaFields = append(metaFields, "product")
	}
	if config.HasInstallerField {
		metaFields = append(metaFields, "build_time", "tc_installer_build_id", "build_c1", "build_c2", "build_c3")
	}
	if config.HasBuildNumber {
		metaFields = append(metaFields, "build_number")
	}
	metaFields = append(metaFields, "triggeredBy")
	var sb strings.Builder
	sb.WriteString("insert into ")
	sb.WriteString(tableName)
	sb.WriteString(" (")
	for i, field := range metaFields {
		if i != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(field)
	}
	config.insertStatementWriter(&sb)
	sb.WriteString(") values (")

	for i, n := 0, len(metaFields)+config.extraFieldCount; i < n; i++ {
		if i != 0 {
			sb.WriteRune(',')
		}
		sb.WriteRune('?')
	}
	sb.WriteRune(')')

	effectiveSql := sb.String()

	insertManager := sql_util.NewBatchInsertManager(ctx, db, effectiveSql, insertWorkerCount, slog.With("type", "report"))

	// large inserts leads to large memory usage, so, allow to override INSERT_BATCH_SIZE via env
	insertManager.BatchSize = env.GetInt("INSERT_BATCH_SIZE", 20_000)

	var installerManager *InsertInstallerManager
	if config.HasInstallerField || config.HasNoInstallerButHasChanges {
		installerManager, err = NewInstallerInsertManager(ctx, db)
		if err != nil {
			return nil, err
		}
	}

	var metaManager *InsertMetaManager
	if config.HasMetaDB && metaDb != nil {
		metaManager, err = NewInsertMetaManager(ctx, metaDb)
		if err != nil {
			return nil, err
		}
	}

	manager := &InsertReportManager{
		nonMetricFieldCount: len(metaFields),
		config:              config,
		TableName:           tableName,
		InsertDataManager: sql_util.InsertDataManager{
			InsertManager: insertManager,
		},

		context:                ctx,
		insertInstallerManager: installerManager,
		insertMetaManager:      metaManager,
	}

	if installerManager != nil {
		insertManager.AddDependency(installerManager.InsertManager)
	}
	return manager, nil
}

// checks that not duplicated, warn if metrics cannot be computed
func (t *InsertReportManager) Insert(runResult *RunResult) error {
	logger := slog.Default()
	if t.config.HasProductField {
		logger = logger.With("product", runResult.Product)
	}
	logger = logger.With(
		"db", t.config.DbName,
		"table", t.config.TableName,
		"file", runResult.ReportFileName,
	)

	// tc collector uses tc build id to avoid duplicates, so, IsCheckThatNotAlreadyAddedNeeded is set to false by default
	if t.IsCheckThatNotAlreadyAddedNeeded && !runResult.GeneratedTime.After(t.MaxGeneratedTime) {
		selectStatement := "select 1 from " + t.config.TableName + " where "
		if t.config.HasProductField {
			selectStatement = "product = '" + sql_util.StringEscaper.Replace(runResult.Product) + "' and "
		}
		selectStatement += "machine = '" + sql_util.StringEscaper.Replace(runResult.Machine) +
			"' and project = '" + sql_util.StringEscaper.Replace(runResult.Report.Project) +
			"' and generated_time = " + strconv.FormatInt(runResult.GeneratedTime.Unix(), 10)

		exists, err := t.CheckExists(t.InsertManager.Db.QueryRow(t.context, selectStatement))
		if err != nil {
			return err
		}

		if exists {
			logger.Debug("report already processed")
			return nil
		}
	}

	err := t.WriteMetrics(runResult.Product, runResult, runResult.branch, runResult.Report.Project, logger)
	if err != nil {
		if errors.Is(err, ErrMetricsCannotBeComputed) {
			logger.Warn(err.Error())
			return nil
		}
		return err
	}

	logger.Debug("new report added")
	return nil
}

//goland:noinspection SpellCheckingInspection
var projectIdToName = map[string]string{
	"Fleet": "fleet",

	// IJ simple project - project v3
	"Mc92Qmj3NY0xxdIiX9ayVbbEZ7s": "simple for IJ",
	// IJ simple project - project v2
	"73YWaW9bytiPDGuKvwNIYMK5CKI": "simple for IJ",

	// idea project (v2)
	"26hfTKDRtXpJ6U7ivgfKthtyU0A": "idea",
	// idea project (v3)
	"nC4MRRFMVYUSQLNIvPgDt+B3JqA": "idea",
	// idea project (v4)
	"Xplo4RZSHXIFu5elYBOiDDkwu20": "idea",

	// light edit
	"6hglkyp/cmAi7ntjrg7dHwd5NG4": "light edit (IJ)",
	"1PbxeQ044EEghMOG9hNEFee05kM": "light edit (IJ)",

	// restoring editors
	"/q9N7EHxr8F1NHjbNQnpqb0Q0fs": "restoring editors",
}

func (t *InsertReportManager) WriteMetrics(product string, row *RunResult, branch string, providedProject string, logger *slog.Logger) error {
	batch, err := t.InsertManager.PrepareForAppend()
	if err != nil {
		return err
	}

	project := row.Report.Project
	if project == "" {
		project = providedProject
	} else if t.config.HasInstallerField {
		customName := projectIdToName[row.Report.Project]
		if customName != "" {
			project = customName
		}
	}

	args := make([]interface{}, 0, t.nonMetricFieldCount+t.config.extraFieldCount)
	args = append(args, row.Machine, row.GeneratedTime, project, uint32(row.TcBuildId), branch)
	if t.config.HasBuildTypeField {
		args = append(args, row.TcBuildType)
	}

	if t.config.HasProductField {
		args = append(args, product)
	}
	if t.config.HasInstallerField {
		buildTimeUnix, err := getBuildTimeFromReport(row.Report)
		if err != nil {
			return err
		}

		if buildTimeUnix.IsZero() {
			buildTimeUnix = row.BuildTime
		}

		if strings.HasPrefix(row.Machine, "intellij-linux-hw-compile-hp-blade-") {
			return nil
		}
		args = append(args, buildTimeUnix, uint32(row.TcInstallerBuildId), uint8(row.BuildC1), uint16(row.BuildC2), uint16(row.BuildC3))
	}
	if t.config.HasBuildNumber {
		args = append(args, row.BuildNumber)
	}
	args = append(args, row.TriggeredBy)

	if t.config.DbName == "ij" || t.config.DbName == "ijDev" {
		err = ComputeIjMetrics(t.nonMetricFieldCount, row.Report, &args, logger)
		if err != nil {
			return err
		}
	}

	args = append(args, row.ExtraFieldData...)

	err = batch.Append(args...)
	if err != nil {
		return fmt.Errorf("cannot append: %w", err)
	}
	return nil
}
