package main

import (
  "archive/tar"
  "fmt"
  _ "github.com/ClickHouse/clickhouse-go"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/segmentio/ksuid"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "path/filepath"
)

func (t *BackupManager) freezeAndMoveToBackupDir(db *sqlx.DB, table clickhouse.TableInfo, backupDir string, logger *zap.Logger) error {
  shadowDir := filepath.Join(t.ClickhouseDir, "shadow")
  dirName := ksuid.New().String()
  logger.Info("freeze table", zap.String("shadowDir", dirName))
  _, err := db.Exec(fmt.Sprintf("alter table `%s`.`%s` freeze with name '" + dirName + "'", table.Database, table.Name))
  if err != nil {
    return errors.WithStack(err)
  }

  tableShadowDir := filepath.Join(shadowDir, dirName)

  storeDir := filepath.Join(tableShadowDir, "store")
  storeDirs, err := ioutil.ReadDir(storeDir)
  if err != nil {
    return errors.WithStack(err)
  }

  for _, f := range storeDirs {
    if f.IsDir() {
      err = os.Rename(filepath.Join(storeDir, f.Name()), filepath.Join(backupDir, f.Name()))
      if err != nil {
        return errors.WithStack(err)
      }
    }
  }

  err = os.RemoveAll(tableShadowDir)
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

func (t *BackupManager) getTables(db *sqlx.DB, dbNames []string) ([]clickhouse.TableInfo, error) {
  t.Logger.Debug("getting tables")

  tables := make([]clickhouse.TableInfo, 0)
  var err error
  if len(dbNames) == 0 {
    err = db.SelectContext(t.TaskContext, &tables,
      "select name, uuid, database, metadata_path from system.tables where database != 'system' and is_temporary = 0 and engine like '%MergeTree' order by database, name;")
  } else {
    err = db.SelectContext(t.TaskContext, &tables,
      "select name, uuid, database, metadata_path from system.tables where database in (?) and is_temporary = 0 and engine like '%MergeTree' order by database, name;", dbNames)
  }
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return tables, nil
}

func writeMetadata(writer *tar.Writer, task Task) error {
  copyBuffer := make([]byte, 64*1024)

  for _, table := range task.tables {
    fullMetadataPath := filepath.Join(task.storeDir, table.MetadataPath)
    metaFileInfo, err := os.Stat(fullMetadataPath)
    if err != nil {
      return errors.WithStack(err)
    }

    err = writeFile(fullMetadataPath, "_metadata_/"+table.MetadataPath, metaFileInfo.Size(), writer, copyBuffer)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  for _, db := range task.db {
    dbMeta := filepath.Join(task.metadataDir, db.Name+".sql")
    dbMetaFileInfo, err := os.Stat(dbMeta)
    if err != nil {
      return errors.WithStack(err)
    }

    err = writeFile(dbMeta, "_metadata_/"+db.Name+".sql", dbMetaFileInfo.Size(), writer, copyBuffer)
    if err != nil {
      return errors.WithStack(err)
    }
  }
  return nil
}
