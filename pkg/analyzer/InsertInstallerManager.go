package analyzer

import (
  "context"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "golang.org/x/tools/container/intsets"
  "strconv"
)

type InsertInstallerManager struct {
  sql_util.InsertDataManager

  maxId       uint32
  insertedIds intsets.Sparse
}

func NewInstallerInsertManager(db driver.Conn, insertContext context.Context, logger *zap.Logger) (*InsertInstallerManager, error) {
  //noinspection GrazieInspection
  insertManager, err := sql_util.NewBulkInsertManager(db, insertContext, "insert into installer", 1, logger.Named("installer"))
  if err != nil {
    return nil, errors.WithStack(err)
  }

  manager := &InsertInstallerManager{
    InsertDataManager: sql_util.InsertDataManager{
      InsertManager: insertManager,

      Logger: logger,
    },

    insertedIds: intsets.Sparse{},
  }

  //noinspection SqlResolve
  err = db.QueryRow(insertContext, "select max(id) from installer").Scan(&manager.maxId)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return manager, nil
}

func (t *InsertInstallerManager) Insert(id int, changes []string) error {
  if t.insertedIds.Has(id) {
    return nil
  }

  if id <= int(t.maxId) {
    exists, err := t.CheckExists(t.InsertManager.Db.QueryRow(t.InsertManager.InsertContext, "select 1 from installer where id = "+strconv.Itoa(id)+" limit 1"))
    if err != nil {
      return err
    }

    if exists {
      return nil
    }
  }

  batch, err := t.InsertManager.PrepareForAppend()
  if err != nil {
    return err
  }

  err = batch.Append(uint32(id), changes)
  if err != nil {
    return errors.WithStack(err)
  }

  t.insertedIds.Insert(id)
  return nil
}
