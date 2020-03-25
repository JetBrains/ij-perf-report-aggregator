package main

import (
  "archive/tar"
  "bytes"
  "compress/gzip"
  "encoding/json"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/cheggaaa/pb/v3"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/djherbis/buffer"
  "github.com/djherbis/nio"
  "github.com/minio/minio-go/v6"
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
  var progressBarUpdater *ProgressBarUpdater
  if !env.GetBool("DISABLE_PROGRESS", false) {
    bar := pb.Full.Start64(task.estimatedTarSize)
    bar.Set(pb.Bytes, true)
    bar.SetRefreshRate(time.Second)
    defer bar.Finish()
    progressBarUpdater = &ProgressBarUpdater{bar: bar}
  }

  if t.TaskContext.Err() != nil {
    return nil
  }

  buf := buffer.New(4 * 1024 * 1024)
  reader, writer := nio.Pipe(buf)
  //noinspection GoUnhandledErrorResult
  defer writer.Close()

  go func() {
    defer func() {
      if r := recover(); r != nil {
        log.Print(string(debug.Stack()))
        t.Logger.Error("recovered", zap.ByteString("stack", debug.Stack()))
      }
    }()

    err := createTar(writer, task, progressBarUpdater)
    _ = writer.CloseWithError(errors.WithStack(err))
    return
  }()

  putObjectOptions := minio.PutObjectOptions{Progress: progressBarUpdater, ContentType: "application/x-tar"}
  _, err := t.Client.PutObjectWithContext(t.TaskContext, os.Getenv("S3_BUCKET"), remoteFilePath, reader, -1, putObjectOptions)
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}

type Task struct {
  backupDir        string
  diffFromPath     string
  estimatedTarSize int64
  extraEntries     []FileEntry

  logger *zap.Logger
}

func createTar(writer *nio.PipeWriter, task Task, progressBarUpdater *ProgressBarUpdater) error {
  tarWriter := tar.NewWriter(writer)
  defer func() {
    err := tarWriter.Close()
    if err != nil {
      _ = writer.CloseWithError(errors.WithStack(err))
    }
  }()

  copyBuffer := make([]byte, 32*1024)

  var hardlinks []string
  err := filepath.Walk(task.backupDir, func(filePath string, info os.FileInfo, err error) error {
    if info == nil || !info.Mode().IsRegular() {
      return nil
    }

    if task.diffFromPath != "" {
      relativePath := filePath[len(task.backupDir)+1:]
      diffFromFile, err := os.Stat(filepath.Join(task.diffFromPath, relativePath))
      if err == nil && os.SameFile(info, diffFromFile) {
        hardlinks = append(hardlinks, relativePath)

        if progressBarUpdater != nil {
          task.estimatedTarSize -= info.Size() - estimatedTarHeaderSize
          progressBarUpdater.bar.SetTotal(task.estimatedTarSize)
        }
        return nil
      }
    }

    err = writeHeader(tarWriter, filePath[len(task.backupDir)+1:], info.Size())
    if err != nil {
      return errors.WithStack(err)
    }

    file, err := os.Open(filePath)
    if err != nil {
      return errors.WithStack(err)
    }

    _, err = io.CopyBuffer(tarWriter, file, copyBuffer)
    _ = file.Close()
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

    err = writeHeader(tarWriter, clickhouse.MetaFileName, int64(len(metaContent)))
    if err != nil {
      return errors.WithStack(err)
    }
    _, err = tarWriter.Write(metaContent)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  for _, entry := range task.extraEntries {
    err := writeHeader(tarWriter, entry.name, int64(len(entry.data)))
    if err != nil {
      return errors.WithStack(err)
    }

    _, err = tarWriter.Write(entry.data)
    if err != nil {
      return errors.WithStack(err)
    }
  }
  return nil
}

func serializeMetaFile(task Task, hardlinks []string) ([]byte, error) {
  metafile := clickhouse.MetaFile{
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
