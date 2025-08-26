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
)

func main() {
	sourceTable := "perfintDev.idea"
	targetTable := "perfintDev.idea2"

	err := copyTable("localhost:9000", sourceTable, targetTable)
	if err != nil {
		slog.Error("copy failed", "error", err)
		os.Exit(1)
	}
}

func parseTableName(fullTableName string) string {
	parts := strings.Split(fullTableName, ".")
	if len(parts) == 2 {
		return parts[0]
	}
	return "default"
}

func copyTable(clickHouseUrl string, sourceTable string, targetTable string) error {
	slog.Info("start copying", "source", sourceTable, "target", targetTable)

	sourceDb := parseTableName(sourceTable)

	db, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickHouseUrl},
		Auth: clickhouse.Auth{
			Database: sourceDb,
		},
		DialTimeout:     time.Second,
		ConnMaxLifetime: time.Hour,
		Settings: map[string]interface{}{
			"send_timeout":     30_000,
			"receive_timeout":  3000,
			"max_memory_usage": 100000000000,
		},
	})
	if err != nil {
		return fmt.Errorf("cannot connect to clickhouse: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get time range for processing
	var minTime, maxTime time.Time
	query := "SELECT min(generated_time) as min, max(generated_time) as max FROM " + sourceTable
	err = db.QueryRow(ctx, query).Scan(&minTime, &maxTime)
	if err != nil {
		return fmt.Errorf("cannot query min/max: %w", err)
	}

	slog.Info("time range", "start", minTime, "end", maxTime)

	// Round to month boundaries
	minTime = time.Date(minTime.Year(), minTime.Month(), 1, 0, 0, 0, 0, minTime.Location())
	if maxTime.Month() == 12 {
		maxTime = time.Date(maxTime.Year()+1, 1, 1, 0, 0, 0, 0, maxTime.Location())
	} else {
		maxTime = time.Date(maxTime.Year(), maxTime.Month()+1, 1, 0, 0, 0, 0, maxTime.Location())
	}

	// Process data month by month
	totalCopied := uint64(0)
	for current := minTime; current.Before(maxTime); {
		next := current.AddDate(0, 1, 0)

		copied, err := processMonthWithDailyBatches(ctx, db, sourceTable, targetTable, current, next)
		if err != nil {
			return fmt.Errorf("failed to process range %v to %v: %w", current, next, err)
		}

		totalCopied += copied
		current = next
	}

	slog.Info("copying finished", "total_rows_copied", totalCopied)
	return nil
}

func processMonthWithDailyBatches(ctx context.Context, db driver.Conn, sourceTable, targetTable string, startTime, endTime time.Time) (uint64, error) {
	slog.Info("processing month with daily batches", "start", startTime, "end", endTime)

	totalCopied := uint64(0)
	batchStart := startTime
	batchNum := 1

	// Process in daily batches
	for batchStart.Before(endTime) {
		batchEnd := batchStart.AddDate(0, 0, 1) // Add 1 day
		if batchEnd.After(endTime) {
			batchEnd = endTime
		}

		// Check row count for this day
		var batchCount uint64
		countQuery := fmt.Sprintf(`
			SELECT count(*) 
			FROM %s 
			WHERE generated_time >= ? AND generated_time < ?
		`, sourceTable)

		err := db.QueryRow(ctx, countQuery, batchStart, batchEnd).Scan(&batchCount)
		if err != nil {
			return totalCopied, fmt.Errorf("cannot count batch rows: %w", err)
		}

		if batchCount > 0 {
			slog.Info("copying daily batch",
				"batch_number", batchNum,
				"date", batchStart.Format("2006-01-02"),
				"rows", batchCount)

			// Copy this day's data
			insertQuery := fmt.Sprintf(`
				INSERT INTO %s 
				SELECT * FROM %s 
				WHERE generated_time >= ? AND generated_time < ?
			`, targetTable, sourceTable)

			err = db.Exec(ctx, insertQuery, batchStart, batchEnd)
			if err != nil {
				return totalCopied, fmt.Errorf("cannot execute INSERT SELECT for batch %d (date %s): %w",
					batchNum, batchStart.Format("2006-01-02"), err)
			}

			totalCopied += batchCount
			slog.Info("daily batch copied successfully",
				"batch_number", batchNum,
				"date", batchStart.Format("2006-01-02"),
				"rows", batchCount,
				"total_copied", totalCopied)
		}

		batchStart = batchEnd
		batchNum++
	}

	slog.Info("month processed successfully", "total_rows", totalCopied)
	return totalCopied, nil
}
