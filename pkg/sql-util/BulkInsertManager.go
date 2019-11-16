package sql_util

import (
  "context"
  "database/sql"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/panjf2000/ants/v2"
  "go.uber.org/zap"
  "runtime"
  "sync"
)

type BulkInsertManager struct {
  transaction     *sql.Tx
  insertStatement *sql.Stmt
  Db              *sqlx.DB

  insertSql string

  BatchSize   int
  queuedItems int

  logger *zap.Logger

  waitGroup sync.WaitGroup
  pool      *ants.Pool
  Error     error

  InsertContext context.Context

  dependencies []*BulkInsertManager
}

func NewBulkInsertManager(db *sqlx.DB, insertContext context.Context, insertSql string, insertWorkerCount int, logger *zap.Logger) (*BulkInsertManager, error) {
  poolCapacity := insertWorkerCount
  if insertWorkerCount == -1 {
    // not enough RAM (if docker has access to 4 GB on a machine where there is only 16 GB)
    poolCapacity = runtime.NumCPU() - 4
    if poolCapacity < 2 {
      poolCapacity = 1
    }
  } else if poolCapacity > 99 {
    poolCapacity = 99
  }

  logger.Info("insert pool capacity", zap.Int("count", poolCapacity))

  pool, err := ants.NewPool(poolCapacity)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return &BulkInsertManager{
    queuedItems:   0,
    Db:            db,
    InsertContext: insertContext,
    insertSql:     insertSql,
    logger:        logger,

    // large inserts leads to large memory usage, so, insert by 2000 items
    BatchSize: 2000,

    pool: pool,
  }, nil
}

func (t *BulkInsertManager) AddDependency(dependency *BulkInsertManager) {
  t.dependencies = append(t.dependencies, dependency)
}

func (t *BulkInsertManager) GetUncommittedTransactionCount() int {
  return t.pool.Running()
}

func (t *BulkInsertManager) CommitAndWait() error {
  uncommittedTransactionCount := t.GetUncommittedTransactionCount()
  if uncommittedTransactionCount > 0 {
    t.logger.Info("waiting inserting", zap.Int("transactions", uncommittedTransactionCount))
  }

  err := t.Commit()
  if err != nil {
    return err
  }

  t.waitGroup.Wait()
  return t.Error
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

  var err error
  if t.pool.Cap() == 1 {
    t.logger.Info("commit", zap.Int("count", queuedItems))
    err = t.doInsert(insertStatement, transaction, queuedItems)
    if err != nil {
      return err
    }
  } else {
    t.logger.Info("add committing of insert transaction to queue", zap.Int("count", queuedItems))
    t.waitGroup.Add(1)
    err = t.pool.Submit(func() {
      defer t.waitGroup.Done()
      err := t.doInsert(insertStatement, transaction, queuedItems)
      if err != nil {
        t.Error = err
      }
    })

    if err != nil {
      t.waitGroup.Done()
      return err
    }
  }
  return nil
}

func (t *BulkInsertManager) doInsert(insertStatement *sql.Stmt, transaction *sql.Tx, queuedItems int) error {
  defer util.Close(insertStatement, t.logger)

  if t.Error != nil {
    t.logger.Error("rollback transaction", zap.String("reason", "previous transaction failed to be committed"), zap.NamedError("prevError", t.Error))
    t.rollbackTransaction(transaction)
    return nil
  }

  for _, dependency := range t.dependencies {
    err := dependency.CommitAndWait()
    if err != nil {
      t.Error = err
      t.logger.Info("cannot commit dependency", zap.Error(err))
      return nil
    }
  }

  err := transaction.Commit()
  if err != nil {
    t.logger.Info("cannot commit", zap.Error(err))
    return errors.WithStack(err)
  }

  t.logger.Info("items were inserted", zap.Int("count", queuedItems))
  return nil
}

func (t *BulkInsertManager) PrepareForInsert() (*sql.Stmt, error) {
  if t.queuedItems >= t.BatchSize {
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
    t.logger.Error("rollback transaction", zap.String("reason", "was not committed, but close is called"))
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
