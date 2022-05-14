package sql_util

import (
  "context"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "runtime"
  "sync"
)

type BatchInsertManager struct {
  batch driver.Batch
  Db    driver.Conn

  insertSql string

  BatchSize   int
  queuedItems int

  logger *zap.Logger

  waitGroup sync.WaitGroup
  Error     error

  InsertContext context.Context

  dependencies []*BatchInsertManager
  batches      chan driver.Batch
}

func NewBulkInsertManager(
  db driver.Conn,
  insertContext context.Context,
  insertSql string,
  insertWorkerCount int,
  logger *zap.Logger,
) (*BatchInsertManager, error) {
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

  // send batches to clickhouse in order, don't use pool here
  batches := make(chan driver.Batch)

  logger.Info("insert pool capacity", zap.Int("count", poolCapacity))

  manager := &BatchInsertManager{
    queuedItems:   0,
    Db:            db,
    InsertContext: insertContext,
    insertSql:     insertSql,
    logger:        logger,

    // large inserts leads to large memory usage, so, insert by 2000 items
    BatchSize: 2000,

    batches: batches,
  }

  for i := 0; i < poolCapacity; i++ {
    go func() {
      manager.processQueue(batches)
    }()
  }

  return manager, nil
}

func (t *BatchInsertManager) processQueue(batches chan driver.Batch) {
  for {
    select {
    case batch, ok := <-batches:
      if !ok {
        return
      }

      if t.Error != nil {
        return
      }

      for _, dependency := range t.dependencies {
        err := dependency.SendAndWait()
        if err != nil {
          t.logger.Error("cannot commit dependency", zap.Error(err))
          t.Error = err
          t.waitGroup.Done()
          return
        }
      }

      t.logger.Info("send batch")
      err := batch.Send()
      if err != nil {
        t.logger.Error("cannot send batch", zap.Error(err))
        t.Error = err
        t.waitGroup.Done()
        return
      }

    case <-t.InsertContext.Done():
      return
    }
  }
}

func (t *BatchInsertManager) AddDependency(dependency *BatchInsertManager) {
  t.dependencies = append(t.dependencies, dependency)
}

func (t *BatchInsertManager) GetUncommittedBatchCount() int {
  return len(t.batches)
}

func (t *BatchInsertManager) SendAndWait() error {
  uncommittedCount := t.GetUncommittedBatchCount()
  if uncommittedCount > 0 {
    t.logger.Info("waiting sending", zap.Int("batches", uncommittedCount))
  }

  t.Flush()
  t.waitGroup.Wait()
  return t.Error
}

func (t *BatchInsertManager) Flush() {
  batch := t.batch
  if batch == nil {
    return
  }

  t.batch = nil
  queuedItems := t.queuedItems
  t.queuedItems = 0

  t.logger.Info("items scheduled to be sent", zap.Int("count", queuedItems))
  t.batches <- batch
}

func (t *BatchInsertManager) PrepareForAppend() (driver.Batch, error) {
  if t.queuedItems >= t.BatchSize {
    t.Flush()
  } else {
    t.queuedItems++
  }

  if t.batch == nil {
    var err error
    t.batch, err = t.Db.PrepareBatch(t.InsertContext, t.insertSql)
    if err != nil {
      return nil, errors.WithStack(err)
    }
  }
  return t.batch, nil
}

func (t *BatchInsertManager) Close() error {
  close(t.batches)

  var err error
  if t.batch != nil {
    t.logger.Error("abort batch", zap.String("reason", "was not sent, but close is called"))
    err = t.batch.Abort()
    t.batch = nil
  }

  return err
}
