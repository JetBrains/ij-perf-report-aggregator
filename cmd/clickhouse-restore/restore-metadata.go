package main

import (
  "archive/tar"
  "compress/gzip"
  "encoding/json"
  clickhouse_backup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "go.uber.org/zap"
  "os"
  "path/filepath"
)

func readInfoMappingFile(tarReader *tar.Reader) (*clickhouse_backup.MappingInfo, error) {
  var info clickhouse_backup.MappingInfo
  err := json.NewDecoder(tarReader).Decode(&info)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return &info, nil
}

func readMetaFile(tarReader *tar.Reader, logger *zap.Logger) (clickhouse_backup.MetaFile, error) {
  var metafile clickhouse_backup.MetaFile

  gzipReader, err := gzip.NewReader(tarReader)
  if err != nil {
    return metafile, errors.WithStack(err)
  }

  defer util.Close(gzipReader, logger)

  err = json.NewDecoder(gzipReader).Decode(&metafile)
  if err != nil {
    return metafile, errors.WithStack(err)
  }
  return metafile, nil
}

// metadata/db/table symlink is not restored - not needed, clickhouse creates these symlinks automatically
func extractMetadata(clickhouseDir string, info *clickhouse_backup.MappingInfo, currentMetadataDir string) error {
  // move metadata to root
  metadataDir := filepath.Join(clickhouseDir, "metadata")
  err := os.RemoveAll(metadataDir)
  if err != nil {
    return errors.WithStack(err)
  }

  err = os.MkdirAll(metadataDir, os.ModePerm)
  if err != nil {
    return errors.WithStack(err)
  }

  for _, item := range info.Db {
    dbDir := filepath.Join(clickhouseDir, "data", item.Name)
    err = os.MkdirAll(dbDir, os.ModePerm)
    if err != nil {
      return errors.WithStack(err)
    }

    err = os.Rename(filepath.Join(currentMetadataDir, item.Name+".sql"), filepath.Join(metadataDir, item.Name+".sql"))
    if err != nil {
      return errors.WithStack(err)
    }
  }

  for _, table := range info.Tables {
    metadataPath := filepath.Join(clickhouseDir, "store", table.MetadataPath)

    err = os.MkdirAll(filepath.Dir(metadataPath), os.ModePerm)
    if err != nil {
      return errors.WithStack(err)
    }

    err = os.Rename(filepath.Join(currentMetadataDir, table.MetadataPath), metadataPath)
    if err != nil {
      return errors.WithStack(err)
    }

    source := filepath.Join(clickhouseDir, "store", table.Uuid[0:3], table.Uuid)
    // no file if all data stored on S3
    if _, err := os.Stat(source); err == nil {
      dbDir := filepath.Join(clickhouseDir, "data", table.Database)
      err = os.Symlink(source, filepath.Join(dbDir, table.Name))
      if err != nil && !os.IsExist(err) {
        return errors.WithStack(err)
      }
    }
  }

  err = os.RemoveAll(currentMetadataDir)
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}
