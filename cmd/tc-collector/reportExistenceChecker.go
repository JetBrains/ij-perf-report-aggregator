package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"golang.org/x/tools/container/intsets"
)

type ReportExistenceChecker struct {
	ids intsets.Sparse
}

func (t *ReportExistenceChecker) reset(taskContext context.Context, dbName string, tableName string, db driver.Conn, since time.Time, buildType string) error {
	t.ids.Clear()

	table := "report"
	if tableName != "" {
		table = tableName
	}
	query := "select tc_build_id from " + table + " where generated_time > " + strconv.FormatInt(since.Unix(), 10) + " and tc_build_type = '" + buildType + "' order by tc_build_id"
	rows, err := db.Query(taskContext, query, since)
	if err != nil {
		return fmt.Errorf("cannot query db %s table %s: %w", dbName, tableName, err)
	}

	for rows.Next() {
		// clickhouse requires explicit type
		var id uint32
		err = rows.Scan(&id)
		if err != nil {
			return fmt.Errorf("cannot scan db %s table %s: %w", dbName, tableName, err)
		}

		t.ids.Insert(int(id))
	}

	if rows.Err() != nil {
		return fmt.Errorf("cannot scan db %s table %s: %w", dbName, tableName, rows.Err())
	}
	return nil
}

func (t *ReportExistenceChecker) has(id int) bool {
	return t.ids.Has(id)
}
