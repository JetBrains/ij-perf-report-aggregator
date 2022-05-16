package sql_util

import (
  "context"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/develar/errors"
  "github.com/panjf2000/ants/v2"
  "go.uber.org/zap"
  "runtime"
  "sync"
)

type BatchInsertManager struct {
  mutex sync.Mutex

  batch driver.Batch
  Db    driver.Conn

  insertSql string

  BatchSize   int
  queuedItems int

  logger *zap.Logger

  waitGroup sync.WaitGroup
  pool      *ants.PoolWithFunc
  Error     error

  InsertContext context.Context

  dependencies []*BatchInsertManager
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

  logger.Info("insert pool capacity", zap.Int("count", poolCapacity))

  manager := &BatchInsertManager{
    queuedItems:   0,
    Db:            db,
    InsertContext: insertContext,
    insertSql:     insertSql,
    logger:        logger,

    // large inserts leads to large memory usage, so, insert by 2000 items
    BatchSize: 2000,
  }

  var err error
  manager.pool, err = ants.NewPoolWithFunc(poolCapacity, func(b interface{}) {
    manager.doSendBatch(b.(driver.Batch))
  })
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return manager, nil
}

func (t *BatchInsertManager) doSendBatch(batch driver.Batch) {
  defer t.waitGroup.Done()

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
    return
  }
}

func (t *BatchInsertManager) AddDependency(dependency *BatchInsertManager) {
  t.dependencies = append(t.dependencies, dependency)
}

func (t *BatchInsertManager) GetQueuedItemCount() int {
  return t.queuedItems
}

func (t *BatchInsertManager) SendAndWait() error {
  uncommittedCount := t.pool.Running()
  if uncommittedCount > 0 {
    t.logger.Info("waiting sending", zap.Int("batches", uncommittedCount))
  }

  err := t.ScheduleSendBatch()
  if err != nil {
    return err
  }

  t.waitGroup.Wait()
  return t.Error
}

func (t *BatchInsertManager) ScheduleSendBatch() error {
  batch, err := t.prepareForFlush()
  if err != nil {
    return err
  }
  if batch == nil {
    return nil
  }

  t.waitGroup.Add(1)
  return t.pool.Invoke(batch)
}

func (t *BatchInsertManager) prepareForFlush() (driver.Batch, error) {
  t.mutex.Lock()
  defer t.mutex.Unlock()

  batch := t.batch
  if batch == nil {
    return nil, nil
  }

  t.batch = nil
  queuedItems := t.queuedItems
  t.queuedItems = 0
  t.logger.Info("items scheduled to be sent", zap.Int("count", queuedItems))
  return batch, nil
}

func (t *BatchInsertManager) PrepareForAppend() (driver.Batch, error) {
  t.mutex.Lock()
  defer t.mutex.Unlock()

  if t.queuedItems >= t.BatchSize {
    err := t.ScheduleSendBatch()
    if err != nil {
      return nil, err
    }
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
  t.pool.Release()

  var err error
  if t.batch != nil {
    t.logger.Error("abort batch", zap.String("reason", "was not sent, but close is called"))
    err = t.batch.Abort()
    t.batch = nil
  }

  return err
}
