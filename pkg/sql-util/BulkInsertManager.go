package sql_util

import (
  "database/sql"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/panjf2000/ants/v2"
  "go.uber.org/zap"
  "report-aggregator/pkg/util"
  "runtime"
  "sync"
)

type BulkInsertManager struct {
  transaction     *sql.Tx
  insertStatement *sql.Stmt
  Db              *sqlx.DB

  insertSql string

  queuedItems int

  logger *zap.Logger

  WaitGroup sync.WaitGroup
  pool      *ants.Pool
  Error     error
}

func NewBulkInsertManager(db *sqlx.DB, insertSql string, logger *zap.Logger) (*BulkInsertManager, error) {
  // not enough RAM (if docker has access to 4 GB on a machine where there is only 16 GB)
  poolCapacity := runtime.NumCPU() - 4
  if poolCapacity < 2 {
    poolCapacity = 2
  }

  pool, err := ants.NewPool(poolCapacity)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return &BulkInsertManager{
    queuedItems: 0,
    Db:          db,
    insertSql:   insertSql,
    logger:      logger,

    pool: pool,
  }, nil
}

func (t *BulkInsertManager) GetUncommittedTransactionCount() int {
  return t.pool.Running()
}

func (t *BulkInsertManager) Commit() error {
  transaction := t.transaction
  if transaction == nil {
    return nil
  }

  insertStatement := t.insertStatement

  t.transaction = nil
  t.insertStatement = nil
  queuedItems := t.queuedItems
  t.queuedItems = 0

  t.logger.Info("add committing of insert transaction to queue", zap.Int("count", queuedItems))
  t.WaitGroup.Add(1)
  err := t.pool.Submit(func() {
    defer t.WaitGroup.Done()

    defer util.Close(insertStatement, t.logger)

    if t.Error != nil {
      t.rollbackTransaction(transaction)
      return
    }

    err := transaction.Commit()
    if err != nil {
      t.Error = errors.WithStack(err)
    }
    t.logger.Info("items were inserted", zap.Int("count", queuedItems))
  })

  if err != nil {
    t.WaitGroup.Done()
    return errors.WithStack(err)
  }
  return nil
}

func (t *BulkInsertManager) PrepareForInsert() (*sql.Stmt, error) {
  // large inserts leads to large memory usage, so, insert by 2000 items
  if t.queuedItems >= 2000 {
    if t.transaction != nil {
      err := t.Commit()
      if err != nil {
        return nil, err
      }
    }
  } else {
    t.queuedItems++
  }

  if t.insertStatement == nil {
    var err error
    t.transaction, err = t.Db.Begin()
    if err != nil {
      return nil, errors.WithStack(err)
    }

    t.insertStatement, err = t.transaction.Prepare(t.insertSql)
    if err != nil {
      return nil, errors.WithStack(err)
    }
  }
  return t.insertStatement, nil
}

func (t *BulkInsertManager) Close() error {
  t.pool.Release()

  var err error
  if t.insertStatement != nil {
    err = t.insertStatement.Close()
  }

  transaction := t.transaction
  if transaction != nil {
    t.transaction = nil
    t.rollbackTransaction(transaction)
  }

  return err
}

func (t *BulkInsertManager) rollbackTransaction(transaction *sql.Tx) {
  err := transaction.Rollback()
  if err != nil {
    t.logger.Error("cannot rollback", zap.Error(err))
  }
}
