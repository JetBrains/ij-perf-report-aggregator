package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/analyzer"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/model"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"go.deanishe.net/env"
)

/*
1. run restore-backup RC
2. change `migrate/report.sql` as needed and execute.
*/
func main() {
	db := env.Get("DB")
	table := env.Get("TABLE")
	split := strings.Split(env.Get("DB"), "_")
	if len(split) > 1 {
		table = split[1]
	}
	if db == "" || table == "" {
		slog.Error("Missing db or table, don't forget to set env variables DB and/or TABLE")
		os.Exit(1)
	}
	err := transform("localhost:9000", db, table)
	if err != nil {
		slog.Error("trasnform failed", "error", err)
	}
}

type ReportRow struct {
	Product string `ch:"product"`
	Machine string `ch:"machine"`
	Branch  string `ch:"branch"`

	Project string `ch:"project"`

	GeneratedTime time.Time `ch:"generated_time"`
	BuildTime     time.Time `ch:"build_time"`

	RawReport string `ch:"raw_report"`

	TcBuildId          uint32 `ch:"tc_build_id"`
	TcInstallerBuildId uint32 `ch:"tc_installer_build_id"`
	TcBuildType        string `ch:"tc_build_type"`
	TriggeredBy        string `ch:"triggeredBy"`

	BuildC1 uint8  `ch:"build_c1"`
	BuildC2 uint16 `ch:"build_c2"`
	BuildC3 uint16 `ch:"build_c3"`

	MeasuresName  []string `ch:"measures.name"`
	MeasuresValue []int32  `ch:"measures.value"`
	MeasuresType  []string `ch:"measures.type"`
}

// set insertWorkerCount to 1 if not enough memory
const insertWorkerCount = 4

func transform(clickHouseUrl string, idName string, tableName string) error {
	slog.Info("start transforming", "db", idName)

	split := strings.Split(idName, "_")
	dbName := idName
	if len(split) > 1 {
		dbName = split[0]
	}

	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickHouseUrl},
		Auth: clickhouse.Auth{
			Database: dbName,
		},
		DialTimeout:     time.Second,
		ConnMaxLifetime: time.Hour,
		Settings: map[string]interface{}{
			// https://github.com/ClickHouse/ClickHouse/issues/2833
			// ZSTD 19+ is used, read/write timeout should be quite large (10 minutes)
			"send_timeout":     30_000,
			"receive_timeout":  3000,
			"max_memory_usage": 100000000000,
		},
	})
	if err != nil {
		return fmt.Errorf("cannot connect to clickhouse: %w", err)
	}

	defer util.Close(db)

	taskContext, cancel := util.CreateCommandContext()
	defer cancel()

	config := analyzer.GetAnalyzer(idName)

	config.TableName = tableName + "2"
	insertReportManager, err := analyzer.NewInsertReportManager(taskContext, db, nil, config, tableName+"2", insertWorkerCount)
	if err != nil {
		return err
	}

	insertManager := insertReportManager.InsertManager
	// we send batch in the end of each iteration
	insertManager.BatchSize = 50_000

	// the whole select result in memory - so, limit
	var minTime time.Time
	var maxTime time.Time
	// use something like (now() - toIntervalMonth(1)) to test the transformer on a fresh data
	err = db.QueryRow(taskContext, "select min(generated_time) as min, max(generated_time) as max from "+tableName).Scan(&minTime, &maxTime)
	if err != nil {
		return fmt.Errorf("cannot query min/max: %w", err)
	}

	slog.Info("time range", "start", minTime, "end", maxTime)

	// round to the start of the month
	minTime = time.Date(minTime.Year(), minTime.Month(), 1, 0, 0, 0, 0, minTime.Location())
	// round to the end of the month
	if maxTime.Month() == 12 {
		maxTime = time.Date(maxTime.Year()+1, 1, 1, 0, 0, 0, 0, maxTime.Location())
	} else {
		maxTime = time.Date(maxTime.Year(), maxTime.Month()+1, 1, 0, 0, 0, 0, maxTime.Location())
	}

	for current := minTime; current.Before(maxTime); {
		// 1 month
		next := current.AddDate(0, 1, 0)
		err = process(taskContext, db, config, current, next, insertReportManager, tableName)
		if err != nil {
			return err
		}

		current = next

		if insertManager.GetQueuedItemCount() > 10_000 {
			insertManager.ScheduleSendBatch()
		}
	}

	err = insertReportManager.InsertManager.Close()
	if err != nil {
		return err
	}

	slog.Info("transforming finished")
	return nil
}

func process(taskContext context.Context, db driver.Conn, config analyzer.DatabaseConfiguration, startTime time.Time, endTime time.Time, insertReportManager *analyzer.InsertReportManager, tableName string) error {
	slog.Info("process", "start", startTime, "end", endTime)
	// don't forget to update order clause if differs - better to insert data in an expected order

	var err error
	var rows driver.Rows
	if config.HasProductField {
		rows, err = db.Query(taskContext, `
      select product, machine, branch,
             generated_time, build_time, raw_report,
             tc_build_id, tc_installer_build_id,
             build_c1, build_c2, build_c3, project
      from report
      where generated_time >= $1 and generated_time < $2
      order by product, machine, branch, project, build_c1, build_c2, build_c3, build_time, generated_time
    `, startTime, endTime)
	} else {
		buildFields := ""
		if config.HasInstallerField {
			buildFields = "build_c1, build_c2, build_c3,"
		}
		installerFields := ""
		if config.HasInstallerField {
			installerFields = "tc_installer_build_id, " + buildFields
		}
		rows, err = db.Query(taskContext, `
      select machine, branch,
             generated_time, build_time, tc_build_id,`+installerFields+` project, measures.name, measures.value, measures.type, triggeredBy
      from `+tableName+`
      where generated_time >= $1 and generated_time < $2
      order by machine, branch, project, `+buildFields+` build_time, generated_time
    `, startTime, endTime)

	}
	if err != nil {
		return fmt.Errorf("cannot query: %w", err)
	}

	defer util.Close(rows)

	var row ReportRow

	for rows.Next() {
		err = rows.ScanStruct(&row)
		if err != nil {
			return fmt.Errorf("cannot scan: %w", err)
		}

		runResult := &analyzer.RunResult{
			Machine:       row.Machine,
			GeneratedTime: row.GeneratedTime,
			BuildTime:     row.BuildTime,
			TcBuildId:     int(row.TcBuildId),
		}

		if config.HasInstallerField {
			runResult.TcInstallerBuildId = int(row.TcInstallerBuildId)
			runResult.BuildC1 = int(row.BuildC1)
			runResult.BuildC2 = int(row.BuildC2)
			runResult.BuildC3 = int(row.BuildC3)
		}
		if config.HasRawReport {
			err = analyzer.ReadReport(runResult, config)
			if err != nil {
				return err
			}

			if runResult.Report == nil {
				// ignore report
				continue
			}
		}

		if config.DbName == "perfint" {
			runResult.Report = &model.Report{
				Project:   row.Project,
				BuildDate: row.BuildTime.Format("20060102T150405+0000"),
				Generated: row.GeneratedTime.Format("20060102T150405+0000"),
			}
			runResult.ExtraFieldData = []interface{}{row.MeasuresName, row.MeasuresValue, row.MeasuresType}
			runResult.TriggeredBy = row.TriggeredBy
			runResult.TcBuildType = row.TcBuildType
		}

		err = insertReportManager.WriteMetrics(row.Product, runResult, row.Branch, row.Project, slog.Default())
		if err != nil {
			return err
		}
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("cannot scan: %w", err)
	}

	return nil
}
