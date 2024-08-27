package sql_util

import (
  "context"
  "fmt"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "golang.org/x/sync/errgroup"
  "log/slog"
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

  logger *slog.Logger

  group *errgroup.Group

  InsertContext context.Context

  dependencies []*BatchInsertManager
}

func NewBatchInsertManager(insertContext context.Context, db driver.Conn, insertSql string, insertWorkerCount int, logger *slog.Logger) *BatchInsertManager {
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

  logger.Info("insert pool capacity", "count", poolCapacity)

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
  return manager
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
      t.logger.Error("cannot commit dependency", "error", err)
      return err
    }
  }

  t.logger.Info("send batch")
  err := batch.Send()
  if err != nil {
    t.logger.Error("cannot send batch", "error", err)
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
  t.logger.Info("items scheduled to be sent", "count", queuedItems)
  return batch
}

func (t *BatchInsertManager) PrepareForAppend() (driver.Batch, error) {
  t.mutex.Lock()
  defer t.mutex.Unlock()

  if t.queuedItems >= t.BatchSize {
    t.mutex.Unlock()
    t.ScheduleSendBatch()
    t.mutex.Lock()
  } else {
    t.queuedItems++
  }

  if t.batch == nil {
    var err error
    t.batch, err = t.Db.PrepareBatch(t.InsertContext, t.insertSql, driver.WithReleaseConnection())
    if err != nil {
      return nil, fmt.Errorf("cannot prepare batch: %w", err)
    }
  }
  return t.batch, nil
}

func (t *BatchInsertManager) Close() error {
  // flush
  t.ScheduleSendBatch()
  err := t.group.Wait()

  if t.batch != nil {
    t.logger.Error("abort batch; was not sent, but close is called")
    abortErr := t.batch.Abort()
    t.batch = nil
    if err == nil {
      err = abortErr
    }
  }
  return err
}
