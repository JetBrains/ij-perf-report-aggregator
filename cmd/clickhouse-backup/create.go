package main

import (
  "archive/tar"
  "fmt"
  _ "github.com/ClickHouse/clickhouse-go"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/segmentio/ksuid"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
)

func (t *BackupManager) freezeAndMoveToBackupDir(db *sqlx.DB, dbName string, backupDir string, logger *zap.Logger) (int64, error) {
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
  dirPrefix := ksuid.New().String()
  for _, table := range tables {
    // todo How do clickhouse transform table to dir name?
    dirName := dirPrefix + "_" + table
    logger.Info("freeze table", zap.String("table", table), zap.String("shadowDir", dirName), zap.String("db", dbName))
    _, err := db.Exec(fmt.Sprintf("alter table `%s`.`%s` freeze with name '" + dirName + "'", dbName, table))
    if err != nil {
      return 0, errors.WithStack(err)
    }

    tableShadowDir := filepath.Join(shadowDir, dirName)
    err = os.Rename(filepath.Join(tableShadowDir, "data", dbName, table), filepath.Join(backupShadowDir, table))
    if err != nil {
      return 0, errors.WithStack(err)
    }

    err = os.RemoveAll(tableShadowDir)
    if err != nil {
      return 0, errors.WithStack(err)
    }
  }
  var estimatedTarSize int64
  err = db.GetContext(t.TaskContext, &estimatedTarSize, "select sum(bytes_on_disk) + (count() * 345) from system.parts where active and database = ?", dbName)
  if err != nil {
    return 0, errors.WithStack(err)
  }
  return estimatedTarSize, nil
}

func (t *BackupManager) getTables(db *sqlx.DB, dbName string) ([]string, error) {
  var tables []string
  err := db.SelectContext(t.TaskContext, &tables, "select name from system.tables where database = ? and is_temporary = 0 and engine like '%MergeTree';", dbName)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return tables, nil
}

func writeMetadata(writer *tar.Writer, dbName string, metadataDir string) error {
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
