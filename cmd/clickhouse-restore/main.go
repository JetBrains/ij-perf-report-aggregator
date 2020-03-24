package main

import (
  "archive/tar"
  "encoding/json"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/cheggaaa/pb/v3"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v6"
  "go.uber.org/zap"
  "io"
  "io/ioutil"
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
  remoteFileInfo, err := backupManager.findBackup(bucket)
  if err != nil {
    return errors.WithStack(err)
  }

  err = backupManager.download(remoteFileInfo.Key, remoteFileInfo.Size)
  if err != nil {
    return errors.WithStack(err)
  }

  // rename shadow to data
  err = os.Rename(filepath.Join(backupManager.LocalPath, "shadow"), filepath.Join(backupManager.LocalPath, "data"))
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

func (t *BackupManager) findBackup(bucket string) (*minio.ObjectInfo, error) {
  var candidate *minio.ObjectInfo
  var lastModified time.Time
  for objectInfo := range t.Client.ListObjectsV2(bucket, "", false, t.TaskContext.Done()) {
    if objectInfo.Err != nil {
      return nil, errors.WithStack(objectInfo.Err)
    }

    if strings.HasSuffix(objectInfo.Key, ".tar") {
      if lastModified.Before(objectInfo.LastModified) {
        candidate = &objectInfo
        lastModified = objectInfo.LastModified
      }
    }
  }

  if candidate == nil {
    return nil, errors.Errorf("backup not found (endpoint=%s, bucket=%s)", t.Client.EndpointURL(), bucket)
  }
  return candidate, nil
}

func (t *BackupManager) download(filePath string, estimatedFileSize int64) error {
  err := os.MkdirAll(t.LocalPath, os.ModePerm)
  if err != nil {
    return errors.WithStack(err)
  }

  object, err := t.Client.GetObjectWithContext(t.TaskContext, t.Bucket, filePath, minio.GetObjectOptions{})
  if err != nil {
    return errors.WithStack(err)
  }

  //noinspection GoUnhandledErrorResult
  defer object.Close()

  copyBuffer := make([]byte, 32*1024)

  var proxyReader io.Reader
  if env.GetBool("DISABLE_PROGRESS", false) {
    proxyReader = object
  } else {
    bar := pb.Start64(estimatedFileSize)
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
      b, err := ioutil.ReadAll(tarReader)
      if err != nil {
        return errors.WithStack(err)
      }

      err = json.Unmarshal(b, &metafile)
      if err != nil {
        return errors.WithStack(err)
      }
      continue
    }

    err = decompressTarFile(tarReader, header, t.LocalPath, copyBuffer, createdDirs)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  if len(metafile.RequiredBackup) != 0 {
    t.Logger.Info("download required parts", zap.String("requiredBackup", metafile.RequiredBackup), zap.String("currentBackup", filePath))
    err = t.download(metafile.RequiredBackup, metafile.EstimatedBackupSize)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  for _, hardlink := range metafile.Hardlinks {
    newName := filepath.Join(t.LocalPath, hardlink)
    extractDir := filepath.Dir(newName)
    oldName := filepath.Join(filepath.Dir(t.LocalPath), metafile.RequiredBackup, hardlink)
    err = os.MkdirAll(extractDir, os.ModePerm)
    if err != nil {
      return errors.WithStack(err)
    }

    err = os.Link(oldName, newName)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  return nil
}

func decompressTarFile(tarReader *tar.Reader, header *tar.Header, rootDirectory string, copyBuffer []byte, createdDirs map[string]bool) error {
  outputFile := filepath.Join(rootDirectory, header.Name)

  var destinationDir string
  if header.Typeflag == tar.TypeDir {
    destinationDir = outputFile
  } else {
    destinationDir = filepath.Dir(outputFile)
  }

  if !createdDirs[destinationDir] {
    err := os.MkdirAll(destinationDir, os.ModePerm)
    if err != nil {
      return errors.WithStack(err)
    }

    createdDirs[destinationDir] = true
  }

  switch header.Typeflag {
  case tar.TypeDir:
    return nil
  case tar.TypeReg, tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
    return writeFile(tarReader, outputFile, copyBuffer)
  case tar.TypeSymlink:
    return os.Symlink(outputFile, header.Linkname)
  case tar.TypeLink:
    return os.Link(outputFile, header.Linkname)
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
