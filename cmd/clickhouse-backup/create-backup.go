package main

import (
  "archive/tar"
  "fmt"
  "github.com/ClickHouse/clickhouse-go/v2/lib/driver"
  clickhousebackup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "os"
  "path/filepath"
  "strconv"
  "strings"
  "time"
)

func inClause(names []string) string {
  return "'" + strings.Join(names, "', '") + "'"
}

func (t *BackupManager) freezeAndMoveToBackupDir(db driver.Conn, table clickhousebackup.TableInfo, backupDir string, logger *zap.Logger) error {
  shadowDir := filepath.Join(t.ClickhouseDir, "shadow")
  dirName := strconv.FormatInt(int64(os.Getpid()), 36) + "_" + strconv.FormatInt(time.Now().UnixNano(), 36)
  logger.Info("optimize table")
  err := db.Exec(t.TaskContext, fmt.Sprintf("optimize table `%s`.`%s`", table.Database, table.Name))
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("freeze table", zap.String("shadowDir", dirName))
  err = db.Exec(t.TaskContext, fmt.Sprintf("alter table `%s`.`%s` freeze with name '"+dirName+"'", table.Database, table.Name))
  if err != nil {
    return errors.WithStack(err)
  }

  tableShadowDir := filepath.Join(shadowDir, dirName)

  storeDir := filepath.Join(tableShadowDir, "store")
  storeDirs, err := os.ReadDir(storeDir)
  if err != nil && !os.IsNotExist(err) {
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

func (t *BackupManager) getTables(db driver.Conn, dbNames []string) ([]clickhousebackup.TableInfo, error) {
  t.Logger.Debug("getting tables")

  tables := make([]clickhousebackup.TableInfo, 0)
  var err error
  if len(dbNames) == 0 {
    err = db.Select(t.TaskContext, &tables,
      "select name, toString(uuid) as uuid, database, metadata_path from system.tables where database != 'system' and is_temporary = 0 and engine like '%MergeTree' order by database, name")
  } else {
    err = db.Select(t.TaskContext, &tables,
      "select name, toString(uuid) as uuid, database, metadata_path from system.tables where database in ("+inClause(dbNames)+") and is_temporary = 0 and engine like '%MergeTree' order by database, name;")
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
