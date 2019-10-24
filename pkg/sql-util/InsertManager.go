package sql_util

import (
  "database/sql"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/multierr"
  "go.uber.org/zap"
)

type InstallerManager struct {
  maxId int
  db    *sqlx.DB

  insertManager   *BulkInsertManager
  selectStatement *sql.Stmt

  insertedIds map[int]bool

  logger *zap.Logger
}

func NewInstallerManager(db *sqlx.DB, logger *zap.Logger) (*InstallerManager, error) {
  selectStatement, err := db.Prepare("select 1 from installer where id = ? limit 1")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  //noinspection GrazieInspection
  insertManager, err := NewBulkInsertManager(db, "insert into installer(id, changes) values(?, ?)", logger)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  manager := &InstallerManager{
    db: db,

    selectStatement: selectStatement,
    insertManager:   insertManager,

    insertedIds: make(map[int]bool),

    logger: logger,
  }

  //noinspection SqlResolve
  err = db.QueryRow("select max(id) from installer").Scan(&manager.maxId)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return manager, nil
}

func (t *InstallerManager) Commit() error {
  return t.insertManager.Commit()
}

func (t *InstallerManager) Close() error {
  return errors.WithStack(multierr.Combine(t.insertManager.Close(), t.selectStatement.Close()))
}

func (t *InstallerManager) Insert(id int, changes string) error {
  if t.insertedIds[id] {
    return nil
  }

  if id <= t.maxId {
    fakeResult := -1
    err := t.selectStatement.QueryRow(id).Scan(&fakeResult)
    if err == nil {
      return nil
    } else if err != sql.ErrNoRows {
      return errors.WithStack(err)
    }
  }

  statement, err := t.insertManager.PrepareForInsert()
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
