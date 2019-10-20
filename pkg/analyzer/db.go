package analyzer

import (
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "go.uber.org/zap"
  "os"
  "path/filepath"
  "report-aggregator/pkg/util"
)

// sqlite can be used as document DB, index can be created for JSON (see https://news.ycombinator.com/item?id=19278019)

func prepareDatabaseFile(filePath string) error {
  dir := filepath.Dir(filePath)

  dirStat, err := os.Stat(dir)
  if err == nil && dirStat.IsDir() {
    // dir exists - check file and copy if needed (for backup purposes)
    return nil
  } else {
    err := os.MkdirAll(dir, 0777)
    if err != nil {
      return errors.WithStack(err)
    }

    // mode for create doesn't work because of umask
    err = os.Chmod(dir, 0700)
    if err != nil {
      return errors.WithStack(err)
    }
  }
  return nil
}

const toolDbVersion = 5

func prepareDatabase(dbPath string, logger *zap.Logger) (*sqlx.DB, error) {
  db, err := sqlx.Open("sqlite3", "file:"+dbPath)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  isPrepared := false

  defer func() {
    if !isPrepared {
      util.Close(db, logger)
    }
  }()

  var dbVersion int
  err = db.Get(&dbVersion, "PRAGMA user_version")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  switch {
  case dbVersion == 0:
    _, err = db.Exec(string(MustAsset("create-db.sql")))
    if err != nil {
      return nil, errors.WithStack(err)
    }

  case dbVersion < 5:
    return nil, errors.Errorf("Migration from db version %d is not implemented", dbVersion)

  case dbVersion > toolDbVersion:
    return nil, errors.Errorf("Database version %d is not supported (tool is outdated)", dbVersion)
  }

  isPrepared = true
  return db, nil
}
