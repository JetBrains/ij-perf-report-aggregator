package main

import (
  "archive/tar"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/cheggaaa/pb/v3"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v7"
  "go.deanishe.net/env"
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

  err := restoreMain(logger)
  if err != nil {
    log.Fatalf("%+v", err)
  }
}

func restoreMain(logger *zap.Logger) error {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  clickhouseDir := clickhouse.GetClickhouseDir()
  dataDir := filepath.Join(clickhouseDir, "store")
  _, err := os.Stat(dataDir)
  if err == nil {
    if !env.GetBool("REMOVE_OLD_DATA_DIR", false) {
      return errors.Errorf("data directory \"%s\" already exists (set env REMOVE_OLD_DATA_DIR=true to force removing)", dataDir)
    }

    err = os.RemoveAll(dataDir)
    if err != nil && !os.IsNotExist(err) {
      return errors.WithStack(err)
    }
  }

  err = os.MkdirAll(clickhouseDir, os.ModePerm)
  if err != nil {
    return errors.WithStack(err)
  }

  if len(os.Args) > 1 {
    filePath := os.Args[1]
    if filePath == "local" {
      filePath, err = findLocalBackup()
      if err != nil {
        return err
      }
    }

    file, err := os.Open(filePath)
    if err != nil {
      return errors.WithStack(err)
    }

    defer util.Close(file, logger)
    return restore(filePath, dataDir, true, file, clickhouseDir, nil, logger, nil)
  } else {
    baseBackupManager, err := clickhouse.CreateBaseBackupManager(taskContext, logger)
    if err != nil {
      return errors.WithStack(err)
    }

    backupManager := &BackupManager{
      baseBackupManager,
    }

    remoteFile, err := backupManager.findBackup(util.GetEnvOrPanic("S3_BUCKET"))
    if err != nil {
      return err
    }

    var bar *pb.ProgressBar
    if !env.GetBool("DISABLE_PROGRESS", false) {
      bar = pb.Full.New(0)
      bar.Set(pb.Bytes, true)
      bar.SetRefreshRate(time.Second)
      defer bar.Finish()
    }

    return backupManager.download(remoteFile, dataDir, true, bar)
  }
}

type BackupManager struct {
  *clickhouse.BaseBackupManager
}

// Reader it's a wrapper for given reader, but with progress handle
type ProxyReader struct {
  io.Reader
  bar *pb.ProgressBar
}

// Read reads bytes from wrapped reader and add amount of bytes to progress bar
func (r *ProxyReader) Read(p []byte) (n int, err error) {
  n, err = r.Reader.Read(p)
  r.bar.Add(n)
  return
}

// Close the wrapped reader when it implements io.Closer
func (r *ProxyReader) Close() (err error) {
  if closer, ok := r.Reader.(io.Closer); ok {
    return closer.Close()
  }
  return
}

func restore(
  file string,
  outputRootDirectory string,
  isExtractMetadataNeeded bool,
  r io.ReadCloser,
  clickhouseDir string,
  backupManager *BackupManager,
  logger *zap.Logger,
  bar *pb.ProgressBar,
) error {
  copyBuffer := make([]byte, 32*1024)
  createdDirs := make(map[string]bool)

  if bar != nil {
    r = &ProxyReader{
      Reader: r,
      bar:    bar,
    }
  }

  tarReader := tar.NewReader(r)
  var metafile clickhouse.MetaFile
  var info *clickhouse.MappingInfo
  for {
    header, err := tarReader.Next()

    if err == io.EOF {
      break
    } else if err != nil {
      return errors.WithStack(err)
    }

    if header.Name == clickhouse.MetaFileName {
      metafile, err = readMetaFile(tarReader, logger)
      if err != nil {
        return err
      }
      continue
    } else if header.Name == clickhouse.InfoFileName {
      info, err = readInfoMappingFile(tarReader)
      if err != nil {
        return err
      }
      continue
    }

    err = decompressTarFile(tarReader, header, outputRootDirectory, copyBuffer, createdDirs)
    if err != nil {
      return err
    }
  }
  err := r.Close()
  if err != nil {
    return err
  }

  logger.Debug("move metadata", zap.String("backup", file), zap.String("outputRootDirectory", outputRootDirectory))
  currentMetadataDir := filepath.Join(outputRootDirectory, "_metadata_")
  if isExtractMetadataNeeded {
    err = extractMetadata(clickhouseDir, info, currentMetadataDir)
    if err != nil {
      return err
    }
  } else {
    err = os.RemoveAll(currentMetadataDir)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  if backupManager == nil || len(metafile.RequiredBackup) == 0 {
    return nil
  }

  logger.Info("download required parts", zap.String("requiredBackup", metafile.RequiredBackup), zap.String("currentBackup", file))
  previousBackupDir := filepath.Join(clickhouseDir, "backup", metafile.RequiredBackup)
  err = backupManager.download(metafile.RequiredBackup+".tar", previousBackupDir, false, bar)
  if err != nil {
    return errors.WithStack(err)
  }

  for _, hardlink := range metafile.Hardlinks {
    newName := filepath.Join(outputRootDirectory, hardlink)
    extractDir := filepath.Dir(newName)
    if !createdDirs[extractDir] {
      err = os.MkdirAll(extractDir, os.ModePerm)
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

func findLocalBackup() (string, error) {
  homeDir, err := os.UserHomeDir()
  if err != nil {
    return "", errors.WithStack(err)
  }

  var candidate string
  var lastModified time.Time
  downloadDir := filepath.Join(homeDir, "Downloads")
  files, err := ioutil.ReadDir(downloadDir)
  if err != nil {
    return "", errors.WithStack(err)
  }
  for _, file := range files {
    if strings.HasSuffix(file.Name(), ".tar") {
      if lastModified.Before(file.ModTime()) {
        candidate = file.Name()
        lastModified = file.ModTime()
      }
    }
  }

  if len(candidate) == 0 {
    return "", errors.Errorf("local backup not found (downloadDir=%s)", downloadDir)
  }
  return filepath.Join(downloadDir, candidate), nil
}

func (t *BackupManager) findBackup(bucket string) (string, error) {
  var candidate string
  var lastModified time.Time
  for objectInfo := range t.Client.ListObjects(t.TaskContext, bucket, minio.ListObjectsOptions{Recursive: false, WithMetadata: true}) {
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

func (t *BackupManager) download(file string, outputRootDirectory string, extractMetadata bool, bar *pb.ProgressBar) error {
  object, err := t.Client.GetObject(t.TaskContext, t.Bucket, file, minio.GetObjectOptions{})
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(object, t.Logger)

  if bar != nil {
    objectInfo, err := t.Client.StatObject(t.TaskContext, t.Bucket, file, minio.StatObjectOptions{})
    if err != nil {
      return errors.WithStack(err)
    }
    bar.AddTotal(objectInfo.Size)
    if !bar.IsStarted() {
      bar.Start()
    }
  }

  return restore(file, outputRootDirectory, extractMetadata, object, t.ClickhouseDir, t, t.Logger, bar)
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
