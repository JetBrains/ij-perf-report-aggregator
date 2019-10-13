package analyzer

import (
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/develar/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
  "report-aggregator/pkg/sqlx"
  "report-aggregator/pkg/util"
	"strings"
	"time"
)

// sqlite can be used as document DB, index can be created for JSON (see https://news.ycombinator.com/item?id=19278019)

func prepareDatabaseFile(filePath string, logger *zap.Logger) error {
	dir := filepath.Dir(filePath)

	dirStat, err := os.Stat(dir)
	if err == nil && dirStat.IsDir() {
		// dir exists - check file and copy if needed (for backup purposes)
		err = createBackup(filePath, dir, logger)
		if err != nil {
			return errors.WithStack(err)
		}
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

func createBackup(filePath string, dirPath string, logger *zap.Logger) error {
	oldConnection, err := sqlite3.Open(filePath, sqlite3.OPEN_READWRITE)
	if err != nil {
		sqlErr, ok := err.(*sqlite3.Error)
		if ok && sqlErr.Code() == sqlite3.CANTOPEN {
			// file is new, no need to backup
			return nil
		}
		return errors.WithStack(err)
	}

	defer util.Close(oldConnection, logger)

	fileBaseName := filepath.Base(filePath)
	newFilePath := filepath.Join(dirPath, strings.TrimSuffix(fileBaseName, filepath.Ext(fileBaseName))+"-backup-"+time.Now().Format("Jan-_2_15-04-05")+".sqlite")

	newConnection, err := sqlite3.Open(newFilePath)
	if err != nil {
		return errors.WithStack(err)
	}

	defer util.Close(newConnection, logger)

	backup, err := oldConnection.Backup("main", newConnection, "main")
	if err != nil {
		return errors.WithStack(err)
	}

	defer util.Close(backup, logger)

	err = backup.Step(-1)
	if err != nil && err != io.EOF {
		return errors.WithStack(err)
	}

	return nil
}

const toolDbVersion = 4

func prepareDatabase(dbPath string, logger *zap.Logger) (*sqlite3.Conn, error) {
	db, err := sqlite3.Open(dbPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	isPrepared := false

	defer func() {
		if !isPrepared {
			util.Close(db, logger)
		}
	}()

	db.BusyTimeout(5 * time.Second)

	dbVersion, err := sqlx.GetInt(db, "PRAGMA user_version", logger)
	if err != nil {
		return nil, err
	}

  if dbVersion == 0 {
    err = db.Exec(string(MustAsset("create-db.sql")))
    if err != nil {
      return nil, err
    }
  } else if dbVersion == 1 {
    return nil, errors.New("Migration from db version 1 is not possible (product code and build number cannot be inferred from old reports)")
  } else if dbVersion < 4 {
    return nil, errors.Errorf("Migration from db version %d is not implemented", dbVersion)
  } else if dbVersion > toolDbVersion {
    return nil, errors.Errorf("Database version %d is not supported (tool is outdated)", dbVersion)
  }

  if dbVersion <= 4 {
    err := db.Exec("create index generated_time on report (generated_time)")
    if err != nil {
      return nil, errors.WithStack(err)
    }
  }

	isPrepared = true
	return db, nil
}