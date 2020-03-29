package main

import (
  "archive/tar"
  "database/sql"
  _ "github.com/ClickHouse/clickhouse-go"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

func (t *BackupManager) freezeAndMoveToBackupDir(logger *zap.Logger, backupDir string) (int64, error) {
  db, err := sql.Open("clickhouse", "tcp://"+env.Get("CLICKHOUSE", "clickhouse:9000")+"?read_timeout=600&write_timeout=600&debug=0&compress=1&send_timeout=30000&receive_timeout=3000")
  if err != nil {
    return 0, errors.WithStack(err)
  }

  defer util.Close(db, logger)

  shadowPath := filepath.Join(t.ClickhouseDir, "shadow")

  dir, err := os.Open(shadowPath)
  if err != nil {
    if !os.IsNotExist(err) {
      return 0, errors.WithStack(err)
    }
  } else {
    names, err := dir.Readdirnames(-1)
    if err != nil {
      return 0, errors.WithStack(err)
    } else if len(names) > 0 {
      return 0, errors.Errorf("'%s' is not empty: %s", shadowPath, strings.Join(names, ","))
    }
  }

  dbName := "default"

  tables, err := t.getTables(db, dbName)
  if err != nil {
    return 0, errors.WithStack(err)
  }

  if len(tables) == 0 {
    return 0, errors.Errorf("no tables to backup")
  }

  logger.Debug("freeze tables", zap.Strings("tables", tables))

  backupShadowDir := filepath.Join(backupDir, dbName)
  err = os.MkdirAll(backupShadowDir, os.ModePerm)
  if err != nil {
    return 0, errors.WithStack(err)
  }

  shadowDir := filepath.Join(t.ClickhouseDir, "shadow")
  for index, table := range tables {
    logger.Info("freeze table", zap.String("table", table))
    _, err := db.Exec("alter table default.`" + table + "` freeze")
    if err != nil {
      return 0, errors.WithStack(err)
    }

    err = os.Rename(filepath.Join(shadowDir, strconv.Itoa(index+1), "data", dbName, table), filepath.Join(backupShadowDir, table))
    if err != nil {
      return 0, errors.WithStack(err)
    }
  }

  err = os.RemoveAll(shadowDir)
  if err != nil {
    return 0, errors.WithStack(err)
  }

  row := db.QueryRow("select sum(bytes) + (count() * 345) from system.parts where active")
  var estimatedTarSize int64
  err = row.Scan(&estimatedTarSize)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  return estimatedTarSize, nil
}

func (t *BackupManager) getTables(db *sql.DB, dbName string) ([]string, error) {
  var tables []string
  rows, err := db.QueryContext(t.TaskContext, "select name from system.tables where database = '"+dbName+"' and is_temporary = 0 and engine like '%MergeTree';")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  defer util.Close(rows, t.Logger)
  for rows.Next() {
    var table string
    err = rows.Scan(&table)
    if err != nil {
      return nil, errors.WithStack(err)
    }

    tables = append(tables, table)
  }
  err = rows.Err()
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return tables, nil
}

func writeMetadata(writer *tar.Writer, metadataDir string) error {
  dbName := "default"
  dbMetadataDir := filepath.Join(metadataDir, dbName)
  files, err := ioutil.ReadDir(dbMetadataDir)
  if err != nil {
    return errors.WithStack(err)
  }

  copyBuffer := make([]byte, 64*1024)
  for _, file := range files {
    name := file.Name()
    if file.Mode().IsRegular() && strings.HasSuffix(name, ".sql") {
      err = writeFile(filepath.Join(dbMetadataDir, name), "_metadata_/"+dbName+"/"+name, file, writer, copyBuffer)
      if err != nil {
        return errors.WithStack(err)
      }
    }
  }
  return nil
}
