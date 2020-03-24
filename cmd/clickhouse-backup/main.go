package main

import (
  "archive/tar"
  "encoding/json"
  "fmt"
  "github.com/AlexAkulov/clickhouse-backup/pkg/chbackup"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/cheggaaa/pb/v3"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/kelseyhightower/envconfig"
  "github.com/minio/minio-go/v6"
  "github.com/nats-io/nats.go"
  "go.uber.org/zap"
  "io"
  "io/ioutil"
  "log"
  "os"
  "path/filepath"
  "runtime/debug"
  "sort"
  "strings"
  "time"
)

var lastLocalBackups []string

const backupSuffix = ".tar"

var incrementalBackupCounter = 0

const maxIncrementalBackupCount = 4
const backupsToKeepLocal = maxIncrementalBackupCount + 1

// 345 - size of header in the UStar format
const estimatedTarHeaderSize = int64(345)

func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  err := start("nats://"+env.Get("NATS", "nats:4222"), logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}

type BackupManager struct {
  *clickhouse.BaseBackupManager
}

func start(natsUrl string, logger *zap.Logger) error {
  logger.Info("started", zap.String("nats", natsUrl))
  nc, err := nats.Connect(natsUrl)
  if err != nil {
    return errors.WithStack(err)
  }

  config := chbackup.DefaultConfig()
  config.General.BackupsToKeepRemote = 10
  err = envconfig.Process("", config)
  if err != nil {
    return errors.WithStack(err)
  }

  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  sub, err := nc.SubscribeSync("db.backup")
  if err != nil {
    return errors.WithStack(err)
  }

  baseBackupManager, err := clickhouse.CreateBaseBackupManager(taskContext, logger)
  if err != nil {
    return errors.WithStack(err)
  }
  backupManager := &BackupManager{
    baseBackupManager,
  }

  backupParentDir := filepath.Join(backupManager.LocalPath, "backup")

  for taskContext.Err() == nil {
    _, err := sub.NextMsgWithContext(taskContext)
    if err != nil {
      contextError := taskContext.Err()
      if contextError != nil {
        logger.Info("cancelled", zap.NamedError("reason", contextError))
        return nil
      }
      return errors.WithStack(contextError)
    }

    if taskContext.Err() != nil {
      return nil
    }

    logger.Info("backup requested")

    err = backupManager.backup(config, backupParentDir)
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))
    }
  }

  return nil
}

func (t *BackupManager) backup(config *chbackup.Config, backupParentDir string) error {
  defer func() {
    if r := recover(); r != nil {
      log.Print(string(debug.Stack()))
      t.Logger.Error("recovered", zap.ByteString("stack", debug.Stack()))
    }
  }()

  // backupsToKeepLocal - 1 as after this function will be yet another backup
  err := t.removeOldLocalBackups(backupParentDir, backupsToKeepLocal - 1)
  if err != nil {
    return errors.WithStack(err)
  }

  backupName := time.Now().Format(chbackup.BackupTimeFormat)

  err = chbackup.CreateBackup(*config, backupName, "")
  if err != nil {
    return errors.WithStack(err)
  }

  var diffFrom string
  // first or each 4 backup is full
  if len(lastLocalBackups) < 1 || incrementalBackupCounter > maxIncrementalBackupCount {
    incrementalBackupCounter = 0
  } else {
    diffFrom = strings.TrimSuffix(lastLocalBackups[len(lastLocalBackups)-1], backupSuffix)
  }

  backupDir := filepath.Join(backupParentDir, backupName)

  var diffFromPath string
  if len(diffFrom) != 0 {
    diffFromPath = filepath.Join(backupParentDir, diffFrom)
  }

  //if diffFromPath != "" {
  //}

  t.Logger.Info("upload", zap.String("backup", backupDir))
  err = t.upload(backupDir, backupName + ".tar", diffFromPath)
  if err != nil {
   return errors.WithStack(err)
  }

  if t.TaskContext.Err() != nil {
    return nil
  }

  t.Logger.Info("uploaded", zap.String("backup", backupDir))

  incrementalBackupCounter++

  if len(lastLocalBackups) > 2 {
    lastLocalBackups = append(lastLocalBackups[1:], backupName+backupSuffix)
  } else {
    lastLocalBackups = append(lastLocalBackups, backupName+backupSuffix)
  }

  return nil
}

func (t *BackupManager) removeOldLocalBackups(backupParentDir string, backupsToKeepLocal int) error {
  files, err := ioutil.ReadDir(backupParentDir)
  if err != nil {
    return errors.WithStack(err)
  }

  counter := len(files) - backupsToKeepLocal
  if counter <= 0 {
    return nil
  }

  sort.SliceStable(files, func(i, j int) bool {
    return files[i].ModTime().Before(files[j].ModTime())
  })

  for _, f := range files {
    counter--
    if counter == 0 {
      return nil
    }

    if !f.IsDir() {
      continue
    }

    t.Logger.Info("remove old local backup", zap.String("backup", f.Name()))
    err = os.RemoveAll(filepath.Join(backupParentDir, f.Name()))
    if err != nil {
      return errors.WithStack(err)
    }
  }

  return nil
}

func (t *BackupManager) upload(localDir, remoteFilePath, diffFromPath string) error {
  filesToPack, hardlinks, estimatedSize, err := collectFiles(localDir, diffFromPath)

  var progressBarUpdater *ProgressBarUpdater
  if !env.GetBool("DISABLE_PROGRESS", false) {
    if err != nil {
      return errors.WithStack(err)
    }

    bar := pb.Start64(estimatedSize)
    bar.Set(pb.Bytes, true)
    defer bar.Finish()
    progressBarUpdater = &ProgressBarUpdater{bar: bar}
  }

  if t.TaskContext.Err() != nil {
    return nil
  }

  var metaContent []byte
  if len(hardlinks) > 0 {
    metafile := clickhouse.MetaFile{
      RequiredBackup:      filepath.Base(diffFromPath),
      EstimatedBackupSize: estimatedSize,
      Hardlinks:           hardlinks,
    }
    metaContent, err = json.MarshalIndent(&metafile, "", "  ")
    if err != nil {
      return errors.WithStack(err)
    }
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

    err := createTar(writer, metaContent, err, filesToPack, localDir)
    _ = writer.CloseWithError(errors.WithStack(err))
    return
  }()

  putObjectOptions := minio.PutObjectOptions{Progress: progressBarUpdater, ContentType: "application/x-tar"}
  _, err = t.Client.PutObjectWithContext(t.TaskContext, os.Getenv("S3_BUCKET"), remoteFilePath, reader, -1, putObjectOptions)
  if err != nil {
    return errors.WithStack(err)
  }

  return nil
}

func createTar(writer *io.PipeWriter, metaContent []byte, err error, filesToPack []*tar.Header, localPath string) error {
  tarWriter := tar.NewWriter(writer)
  defer func() {
    err := tarWriter.Close()
    if err != nil {
      _ = writer.CloseWithError(errors.WithStack(err))
    }
  }()

  copyBuffer := make([]byte, 32*1024)

  if len(metaContent) > 0 {
    err = tarWriter.WriteHeader(&tar.Header{
      Typeflag: tar.TypeReg,
      Name:     clickhouse.MetaFileName,
      Size:     int64(len(metaContent)),
      Mode:     0644,
      Uid:      0,
      Gid:      0,
    })
    if err != nil {
      return errors.WithStack(err)
    }

    _, err = tarWriter.Write(metaContent)
    if err != nil {
      return errors.WithStack(err)
    }
  }

  for _, header := range filesToPack {
    filePath := header.Name
    header.Name = filePath[len(localPath)+1:]

    err = tarWriter.WriteHeader(header)
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
  }
  return nil
}

func collectFiles(localDir string, diffFromPath string) ([]*tar.Header, []string, int64, error) {
  var estimatedSize int64
  var filesToPack []*tar.Header
  var hardlinks []string
  err := filepath.Walk(localDir, func(filePath string, info os.FileInfo, err error) error {
    if info == nil || !info.Mode().IsRegular() {
      return nil
    }

    if diffFromPath != "" {
      relativePath := filePath[len(localDir)+1:]
      diffFromFile, err := os.Stat(filepath.Join(diffFromPath, relativePath))
      if err == nil && os.SameFile(info, diffFromFile) {
        hardlinks = append(hardlinks, relativePath)
        return nil
      }
    }

    estimatedSize += info.Size() + estimatedTarHeaderSize
    header, err := tar.FileInfoHeader(info, "")
    if err != nil {
      return errors.WithStack(err)
    }

    header.Name = filePath
    filesToPack = append(filesToPack, header)
    return nil
  })
  return filesToPack, hardlinks, estimatedSize, err
}

type ProgressBarUpdater struct {
  ContentSize int64
  bar         *pb.ProgressBar
}

func (t *ProgressBarUpdater) Read(p []byte) (n int, err error) {
  t.bar.Add(len(p))
  return
}
