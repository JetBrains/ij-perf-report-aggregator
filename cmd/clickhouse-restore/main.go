package main

import (
  "archive/tar"
  "compress/gzip"
  "encoding/json"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/cheggaaa/pb/v3"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v6"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "io"
  "log"
  "os"
  "path/filepath"
  "strings"
  "time"
)

func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  err := restore(os.Getenv("S3_BUCKET"), logger)
  if err != nil {
    log.Fatalf("%+v", err)
  }
}

type BackupManager struct {
  *clickhouse.BaseBackupManager
}

func restore(bucket string, logger *zap.Logger) error {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  baseBackupManager, err := clickhouse.CreateBaseBackupManager(taskContext, logger)
  if err != nil {
    return errors.WithStack(err)
  }
  backupManager := &BackupManager{
    baseBackupManager,
  }
  remoteFile, err := backupManager.findBackup(bucket)
  if err != nil {
    return errors.WithStack(err)
  }

  err = os.MkdirAll(backupManager.ClickhouseDir, os.ModePerm)
  if err != nil {
    return errors.WithStack(err)
  }

  err = backupManager.download(remoteFile, filepath.Join(backupManager.ClickhouseDir, "data"), true)
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

func (t *BackupManager) findBackup(bucket string) (string, error) {
  var candidate string
  var lastModified time.Time
  for objectInfo := range t.Client.ListObjectsV2(bucket, "", false, t.TaskContext.Done()) {
    if objectInfo.Err != nil {
      return "", errors.WithStack(objectInfo.Err)
    }

    if strings.HasSuffix(objectInfo.Key, ".tar") {
      if lastModified.Before(objectInfo.LastModified) {
        candidate = objectInfo.Key
        lastModified = objectInfo.LastModified
      }
    }
  }

  if len(candidate) == 0 {
    return "", errors.Errorf("backup not found (endpoint=%s, bucket=%s)", t.Client.EndpointURL(), bucket)
  }
  return candidate, nil
}

func (t *BackupManager) download(file string, outputRootDirectory string, extractMetadata bool) error {
  object, err := t.Client.GetObjectWithContext(t.TaskContext, t.Bucket, file, minio.GetObjectOptions{})
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(object, t.Logger)

  copyBuffer := make([]byte, 32*1024)

  var proxyReader io.Reader
  if env.GetBool("DISABLE_PROGRESS", false) {
    proxyReader = object
  } else {
    objectInfo, err := t.Client.StatObjectWithContext(t.TaskContext, t.Bucket, file, minio.StatObjectOptions{})
    if err != nil {
      return errors.WithStack(err)
    }
    bar := pb.Full.Start64(objectInfo.Size)
    bar.SetRefreshRate(time.Second)
    proxyReader = bar.NewProxyReader(object)
    defer bar.Finish()
  }

  createdDirs := make(map[string]bool)

  tarReader := tar.NewReader(proxyReader)
  var metafile clickhouse.MetaFile
  for {
    header, err := tarReader.Next()
    if err == io.EOF {
      break
    } else if err != nil {
      return errors.WithStack(err)
    }

    if header.Name == clickhouse.MetaFileName {
      metafile, err = t.readMetaFile(tarReader)
      if err != nil {
        return errors.WithStack(err)
      }
      continue
    }

    err = decompressTarFile(tarReader, header, outputRootDirectory, copyBuffer, createdDirs)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  t.Logger.Debug("move metadata", zap.String("backup", file), zap.String("outputRootDirectory", outputRootDirectory))
  currentMetadataDir := filepath.Join(outputRootDirectory, "_metadata_")
  if extractMetadata {
    // move metadata to root
    metadataDir := filepath.Join(t.ClickhouseDir, "metadata")
    err = os.RemoveAll(metadataDir)
    if err != nil {
      return errors.WithStack(err)
    }
    err = os.Rename(currentMetadataDir, metadataDir)
    if err != nil {
      return errors.WithStack(err)
    }
  } else {
    err = os.RemoveAll(currentMetadataDir)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  if len(metafile.RequiredBackup) == 0 {
    return nil
  }

  t.Logger.Info("download required parts", zap.String("requiredBackup", metafile.RequiredBackup), zap.String("currentBackup", file))
  previousBackupDir := filepath.Join(t.ClickhouseDir, "backup", metafile.RequiredBackup)
  err = t.download(metafile.RequiredBackup+".tar", previousBackupDir, false)
  if err != nil {
    return errors.WithStack(err)
  }

  for _, hardlink := range metafile.Hardlinks {
    newName := filepath.Join(outputRootDirectory, hardlink)
    extractDir := filepath.Dir(newName)
    if !createdDirs[extractDir] {
      err := os.MkdirAll(extractDir, os.ModePerm)
      if err != nil {
        return errors.WithStack(err)
      }

      createdDirs[extractDir] = true
    }

    err = os.Link(filepath.Join(previousBackupDir, hardlink), newName)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  return nil
}

func (t *BackupManager) readMetaFile(tarReader *tar.Reader) (clickhouse.MetaFile, error) {
  var metafile clickhouse.MetaFile

  gzipReader, err := gzip.NewReader(tarReader)
  if err != nil {
    return metafile, errors.WithStack(err)
  }

  defer util.Close(gzipReader, t.Logger)

  err = json.NewDecoder(gzipReader).Decode(&metafile)
  if err != nil {
    return metafile, errors.WithStack(err)
  }
  return metafile, nil
}

func decompressTarFile(tarReader *tar.Reader, header *tar.Header, outputRootDirectory string, copyBuffer []byte, createdDirs map[string]bool) error {
  if header.Typeflag == tar.TypeDir {
    // archive doesn't contain directory entries, and even if exists, do not create empty directories
    return nil
  }

  file := filepath.Join(outputRootDirectory, header.Name)
  dir := filepath.Dir(file)

  if !createdDirs[dir] {
    err := os.MkdirAll(dir, os.ModePerm)
    if err != nil {
      return errors.WithStack(err)
    }

    createdDirs[dir] = true
  }

  switch header.Typeflag {
  case tar.TypeReg, tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
    return writeFile(tarReader, file, copyBuffer)
  case tar.TypeSymlink:
    return os.Symlink(file, header.Linkname)
  case tar.TypeLink:
    return os.Link(file, header.Linkname)
  default:
    return errors.Errorf("%s: unknown type flag: %c", header.Name, header.Typeflag)
  }
}

// don't restore file permissions - all files have regular perm and to the sake of speed, avoid chmod
func writeFile(source io.Reader, to string, buffer []byte) error {
  destinationFile, err := os.Create(to)
  if err != nil {
    return errors.WithStack(err)
  }

  _, err = io.CopyBuffer(destinationFile, source, buffer)
  if err != nil {
    _ = destinationFile.Close()
    return errors.WithStack(err)
  }

  err = destinationFile.Close()
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}
