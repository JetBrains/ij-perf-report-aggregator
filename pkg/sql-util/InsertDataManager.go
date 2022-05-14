package sql_util

import (
  "database/sql"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/develar/errors"
  "go.uber.org/zap"
)

type InsertDataManager struct {
  InsertManager *BatchInsertManager
  Logger        *zap.Logger
}

func (t *InsertDataManager) CheckExists(row driver.Row) (bool, error) {
  var fakeResult uint8
  err := row.Scan(&fakeResult)
  switch {
  case err == nil:
    return true, nil
  case err != sql.ErrNoRows:
    return false, errors.WithStack(err)
  default:
    return false, nil
  }
}
