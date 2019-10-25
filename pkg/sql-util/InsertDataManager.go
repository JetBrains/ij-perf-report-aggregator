package sql_util

import (
  "database/sql"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/multierr"
  "go.uber.org/zap"
)

type InsertDataManager struct {
  Db *sqlx.DB

  InsertManager   *BulkInsertManager
  SelectStatement *sql.Stmt

  Logger *zap.Logger
}

func (t *InsertDataManager) GetUncommittedTransactionCount() int {
  return t.InsertManager.GetUncommittedTransactionCount()
}

func (t *InsertDataManager) CommitAndWait() error {
  err := t.InsertManager.Commit()
  if err != nil {
    return err
  }

  t.InsertManager.WaitGroup.Wait()
  return t.InsertManager.Error
}

func (t *InsertDataManager) Close() error {
  return errors.WithStack(multierr.Combine(t.InsertManager.Close(), t.SelectStatement.Close()))
}

func (t *InsertDataManager) CheckExists(row *sql.Row) (bool, error) {
  fakeResult := -1
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
