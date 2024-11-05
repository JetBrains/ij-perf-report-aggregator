package sql_util

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type InsertDataManager struct {
	InsertManager *BatchInsertManager
}

func (t *InsertDataManager) CheckExists(row driver.Row) (bool, error) {
	var fakeResult uint8
	err := row.Scan(&fakeResult)

	switch {
	case err == nil:
		return true, nil
	case !errors.Is(err, sql.ErrNoRows):
		return false, fmt.Errorf("cannot check exists: %w", err)
	default:
		return false, nil
	}
}
