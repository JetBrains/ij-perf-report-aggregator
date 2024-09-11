package main

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"log"
)

func main() {
	taskContext, cancel := util.CreateCommandContext()
	defer cancel()

	err := moveToDisk(taskContext)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func moveToDisk(taskContext context.Context) error {
	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
	})
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	var result []struct {
		Partition string `ch:"partition"`
		Table     string `ch:"table"`
		Database  string `ch:"database"`
	}
	err = db.Select(taskContext, &result, "select distinct partition, database, table from system.parts where"+
		" database != 'system' and active = 1 and disk_name == 'default' and toYear(max_time) > 1970 and toYear(max_time) < 2022"+
		" order by database, table, name"+
		"")
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	for _, item := range result {
		log.Printf("move %s (%s.%s)", item.Partition, item.Database, item.Table)
		query := fmt.Sprintf("alter table %s.%s move partition %s to disk '%s'", item.Database, item.Table, item.Partition, "s3")
		// query := fmt.Sprintf("alter table %s.%s modify setting storage_policy='tiered'", item.Database, item.Table)
		println(query)
		err = db.Exec(taskContext, query)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}
