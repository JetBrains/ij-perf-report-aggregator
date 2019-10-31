package analyzer

import (
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/sql-util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"
)

type InsertInstallerManager struct {
  sql_util.InsertDataManager

  maxId       int
  insertedIds map[int]bool
}

func NewInstallerInsertManager(db *sqlx.DB, logger *zap.Logger) (*InsertInstallerManager, error) {
  selectStatement, err := db.Prepare("select 1 from installer where id = ? limit 1")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  //noinspection GrazieInspection
  insertManager, err := sql_util.NewBulkInsertManager(db, "insert into installer(id, changes) values(?, ?)", logger)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  manager := &InsertInstallerManager{
    InsertDataManager: sql_util.InsertDataManager{
      Db: db,

      SelectStatement: selectStatement,
      InsertManager:   insertManager,

      Logger: logger,
    },

    insertedIds: make(map[int]bool),
  }

  //noinspection SqlResolve
  err = db.QueryRow("select max(id) from installer").Scan(&manager.maxId)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return manager, nil
}

func (t *InsertInstallerManager) Insert(id int, changes string) error {
  if t.insertedIds[id] {
    return nil
  }

  if id <= t.maxId {
    exists, err := t.CheckExists(t.SelectStatement.QueryRow(id))
    if err != nil {
      return err
    }

    if exists {
      return nil
    }
  }

  statement, err := t.InsertManager.PrepareForInsert()
  if err != nil {
    return errors.WithStack(err)
  }

  _, err = statement.Exec(id, changes)
  if err != nil {
    return errors.WithStack(err)
  }

  t.insertedIds[id] = true
  return nil
}
