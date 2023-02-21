package sql_util

import (
  "database/sql"
  "errors"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  e "github.com/develar/errors"
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
  case !errors.Is(err, sql.ErrNoRows):
    return false, e.WithStack(err)
  default:
    return false, nil
  }
}
