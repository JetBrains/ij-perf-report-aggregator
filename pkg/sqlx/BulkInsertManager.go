package sqlx

import (
  "database/sql"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "report-aggregator/pkg/util"
)

type BulkInsertManager struct {
  transaction     *sql.Tx
  InsertStatement *sql.Stmt
  Db              *sql.DB

  insertSql string

  queuedItems int

  logger *zap.Logger
}

func NewBulkInsertManager(db *sql.DB, insertSql string, logger *zap.Logger) *BulkInsertManager {
  return &BulkInsertManager{queuedItems: 0, Db: db, insertSql: insertSql, logger: logger}
}

func (t *BulkInsertManager) Commit() error {
  if t.transaction != nil {
    defer util.Close(t.InsertStatement, t.logger)

    err := t.transaction.Commit()
    if err != nil {
      return errors.WithStack(err)
    }

    t.logger.Info("items were inserted", zap.Int("count", t.queuedItems))
    t.transaction = nil
    t.InsertStatement = nil
    t.queuedItems = 0
  }
  return nil
}

func (t *BulkInsertManager) PrepareForInsert() error {
  if t.queuedItems >= 1000 {
    if t.transaction != nil {
      err := t.Commit()
      if err != nil {
        return err
      }
    }
  } else {
    t.queuedItems++
  }

  if t.InsertStatement == nil {
    var err error
    t.transaction, err = t.Db.Begin()
    if err != nil {
      return errors.WithStack(err)
    }

    t.InsertStatement, err = t.transaction.Prepare(t.insertSql)
    if err != nil {
      return errors.WithStack(err)
    }
  }
  return nil
}

func (t *BulkInsertManager) Close() error {
  if t.InsertStatement != nil {
    util.Close(t.InsertStatement, t.logger)
  }
  if t.transaction != nil {
    err := t.transaction.Rollback()
    if err != nil {
      t.logger.Error("cannot rollback", zap.Error(err))
    }
  }

  return nil
}
