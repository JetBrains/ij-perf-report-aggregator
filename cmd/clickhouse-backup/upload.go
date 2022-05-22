package main

import (
  "archive/tar"
  "bytes"
  "compress/gzip"
  "encoding/json"
  clickhousebackup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/cheggaaa/pb/v3"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v7"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "io"
  "log"
  "os"
  "path/filepath"
  "runtime/debug"
  "time"
)

// 345 - size of header in the UStar format
const estimatedTarHeaderSize = int64(345)

func (t *BackupManager) upload(remoteFilePath string, task Task) error {
  var bar *pb.ProgressBar
  putObjectOptions := minio.PutObjectOptions{
    ContentType: "application/x-tar",
    NumThreads:  env.GetUint("UPLOAD_WORKER_COUNT", 0),
    PartSize:    uint64(env.GetUint("UPLOAD_PART_SIZE", 0) * 1024 * 1024),
  }
  if !env.GetBool("DISABLE_PROGRESS", false) {
    bar = pb.Full.Start64(task.estimatedTarSize)
    bar.Set(pb.Bytes, true)
    bar.SetRefreshRate(time.Second)
    defer bar.Finish()
    putObjectOptions.Progress = &ProgressBarUpdater{bar: bar}
  }

  if t.TaskContext.Err() != nil {
    return nil
  }

  reader, writer := io.Pipe()
  //noinspection GoUnhandledErrorResult
  defer writer.Close()

  go func() {
    defer func() {
      if r := recover(); r != nil {
        log.Print(string(debug.Stack()))
        t.Logger.Error("recovered", zap.ByteString("stack", debug.Stack()))
      }
    }()

    err := createTar(writer, task, bar)
    _ = writer.CloseWithError(errors.WithStack(err))
  }()

  _, err := t.Client.PutObject(t.TaskContext, os.Getenv("S3_BUCKET"), remoteFilePath, reader, -1, putObjectOptions)
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}

type Task struct {
  metadataDir      string
  storeDir         string
  backupDir        string
  diffFromPath     string
  db               []clickhousebackup.DbInfo
  tables           []clickhousebackup.TableInfo
  estimatedTarSize int64

  logger *zap.Logger
}

func createTar(writer io.Writer, task Task, progressBar *pb.ProgressBar) (err error) {
  tarWriter := tar.NewWriter(writer)
  defer func() {
    closeErr := tarWriter.Close()
    if err == nil {
      err = closeErr
    }
  }()

  copyBuffer := make([]byte, 64*1024)

  err = writeMetadata(tarWriter, task)
  if err != nil {
    return errors.WithStack(err)
  }

  var hardlinks []string
  err = filepath.Walk(task.backupDir, func(filePath string, info os.FileInfo, err error) error {
    if err != nil {
      return errors.WithStack(err)
    }

    if info == nil || !info.Mode().IsRegular() {
      return nil
    }

    if task.diffFromPath != "" {
      relativePath := filePath[len(task.backupDir)+1:]
      diffFromFile, err := os.Stat(filepath.Join(task.diffFromPath, relativePath))
      if err == nil && os.SameFile(info, diffFromFile) {
        hardlinks = append(hardlinks, relativePath)

        if progressBar != nil {
          progressBar.Add64(info.Size() - estimatedTarHeaderSize)
        }
        return nil
      }
    }

    err = writeFile(filePath, filePath[len(task.backupDir)+1:], info.Size(), tarWriter, copyBuffer)
    if err != nil {
      return errors.WithStack(err)
    }

    return nil
  })
  if err != nil {
    return errors.WithStack(err)
  }

  if len(hardlinks) > 0 {
    metaContent, err := serializeMetaFile(task, hardlinks)
    if err != nil {
      return errors.WithStack(err)
    }

    err = writeHeader(tarWriter, clickhousebackup.MetaFileName, int64(len(metaContent)))
    if err != nil {
      return errors.WithStack(err)
    }
    _, err = tarWriter.Write(metaContent)
    if err != nil {
      return errors.WithStack(err)
    }
  }
  return nil
}

func writeFile(filePath string, name string, size int64, tarWriter *tar.Writer, copyBuffer []byte) error {
  err := writeHeader(tarWriter, name, size)
  if err != nil {
    return errors.WithStack(err)
  }

  file, err := os.Open(filePath)
  if err != nil {
    return errors.WithStack(err)
  }

  //noinspection GoUnhandledErrorResult
  defer file.Close()

  _, err = io.CopyBuffer(tarWriter, file, copyBuffer)
  if err != nil {
    return errors.WithStack(err)
  }
  return nil
}

func serializeMetaFile(task Task, hardlinks []string) ([]byte, error) {
  metafile := clickhousebackup.MetaFile{
    RequiredBackup:      filepath.Base(task.diffFromPath),
    EstimatedBackupSize: task.estimatedTarSize,
    Hardlinks:           hardlinks,
  }

  var b bytes.Buffer
  gzWriter, err := gzip.NewWriterLevel(&b, gzip.BestCompression)
  if err != nil {
    return nil, errors.WithStack(err)
  }
  defer util.Close(gzWriter, task.logger)

  err = json.NewEncoder(gzWriter).Encode(&metafile)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  err = gzWriter.Close()
  if err != nil {
    return nil, errors.WithStack(err)
  }
  return b.Bytes(), nil
}

// without OS-specific fields like gid/uid/Mode/ModTime
func writeHeader(tarWriter *tar.Writer, name string, size int64) error {
  return tarWriter.WriteHeader(&tar.Header{
    Typeflag: tar.TypeReg,
    Name:     name,
    Size:     size,
    Uid:      -1,
    Gid:      -1,
    Mode:     0644,
    Format:   tar.FormatPAX,
  })
}

type ProgressBarUpdater struct {
  ContentSize int64
  bar         *pb.ProgressBar
}

func (t *ProgressBarUpdater) Read(p []byte) (n int, err error) {
  t.bar.Add(len(p))
  return
}
