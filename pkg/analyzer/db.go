package analyzer

import (
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/develar/errors"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"report-aggregator/pkg/util"
	"strings"
	"time"
)

func prepareDatabaseFile(filePath string, logger *zap.Logger) error {
	dir := filepath.Dir(filePath)

	dirStat, err := os.Stat(dir)
	if err == nil && dirStat.IsDir() {
		// dir exists - check file and copy if needed (for backup purposes)
		_, err := os.Stat(filePath)
		if err == nil {
			err = createBackup(filePath, dir, logger)
			if err != nil {
				return err
			}
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
	fileBaseName := filepath.Base(filePath)
	backupFile, err := os.Create(filepath.Join(dirPath, strings.TrimSuffix(fileBaseName, filepath.Ext(fileBaseName))+"-backup-"+time.Now().Format(time.RFC822)+".sqlite"))
	if err != nil {
		return errors.WithStack(err)
	}

	defer util.Close(backupFile, logger)

	file, err := os.Open(filePath)
	if err != nil {
		return errors.WithStack(err)
	}

	defer util.Close(file, logger)

	_, err = io.Copy(file, backupFile)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func prepareDatabase(dbPath string, logger *zap.Logger) (*sqlite3.Conn, error) {
	db, err := sqlite3.Open(dbPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db.BusyTimeout(5 * time.Second)

	err = db.Exec(`
create table if not exists report
(
	id string not null
		constraint reports_pk
			primary key,
  generated_time int,
  metrics_version int,
	metrics string,
	raw_report string
);

create unique index if not exists reports_id_uindex
	on report (id);
`)

	if err != nil {
		util.Close(db, logger)
		return nil, err
	}

	return db, nil
}
