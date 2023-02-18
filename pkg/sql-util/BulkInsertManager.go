package sql_util

import (
  "context"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "golang.org/x/sync/errgroup"
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

  group *errgroup.Group

  InsertContext context.Context

  dependencies []*BatchInsertManager
}

func NewBulkInsertManager(insertContext context.Context, db driver.Conn, insertSql string, insertWorkerCount int, logger *zap.Logger, ) (*BatchInsertManager, error) {
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

  logger.Debug("insert pool capacity", zap.Int("count", poolCapacity))

  errorGroup, ctx := errgroup.WithContext(insertContext)
  errorGroup.SetLimit(poolCapacity)
  manager := &BatchInsertManager{
    queuedItems:   0,
    Db:            db,
    InsertContext: ctx,
    insertSql:     insertSql,
    logger:        logger,

    group: errorGroup,

    // large inserts leads to large memory usage, so, insert by 2000 items
    BatchSize: 2000,
  }
  return manager, nil
}

func (t *BatchInsertManager) AddDependency(dependency *BatchInsertManager) {
  t.dependencies = append(t.dependencies, dependency)
}

func (t *BatchInsertManager) GetQueuedItemCount() int {
  return t.queuedItems
}

func (t *BatchInsertManager) ScheduleSendBatch() {
  batch := t.prepareForFlush()
  if batch != nil {
    t.group.Go(func() error {
      return t.sendBatch(batch)
    })
  }
}

func (t *BatchInsertManager) SendBatchNow() error {
  batch := t.prepareForFlush()
  if batch != nil {
    return t.sendBatch(batch)
  }
  return nil
}

func (t *BatchInsertManager) sendBatch(batch driver.Batch) error {
  for _, dependency := range t.dependencies {
    err := dependency.SendBatchNow()
    if err != nil {
      t.logger.Error("cannot commit dependency", zap.Error(err))
      return err
    }
  }

  t.logger.Info("send batch")
  err := batch.Send()
  if err != nil {
    t.logger.Error("cannot send batch", zap.Error(err))
    return err
  }
  return nil
}

func (t *BatchInsertManager) prepareForFlush() driver.Batch {
  t.mutex.Lock()
  defer t.mutex.Unlock()

  batch := t.batch
  if batch == nil {
    return nil
  }

  t.batch = nil
  queuedItems := t.queuedItems
  t.queuedItems = 0
  t.logger.Info("items scheduled to be sent", zap.Int("count", queuedItems))
  return batch
}

func (t *BatchInsertManager) PrepareForAppend() (driver.Batch, error) {
  t.mutex.Lock()
  defer t.mutex.Unlock()

  if t.queuedItems >= t.BatchSize {
    t.ScheduleSendBatch()
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
  // flush
  t.ScheduleSendBatch()
  err := t.group.Wait()

  if t.batch != nil {
    t.logger.Error("abort batch", zap.String("reason", "was not sent, but close is called"))
    abortErr := t.batch.Abort()
    t.batch = nil
    if err == nil {
      err = abortErr
    }
  }
  return err
}
