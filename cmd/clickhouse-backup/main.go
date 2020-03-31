package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v6"
  "github.com/nats-io/nats.go"
  "go.uber.org/zap"
  "math/rand"
  "os"
  "path/filepath"
  "runtime/debug"
  "time"
)

// example: if data collected each 3 hours, will be 8 backup per day, so, upload full backup at least once a day
const maxIncrementalBackupCount = 8

// RFC3339 is not suitable because of colon ':'
const timeFormat = "2006-01-02T15-04-05"

// If pod started first time (not only this container), then first backup will be not incremental,
// because clickhouse-restore renames directory and do not story copy in the backup dir. It is ok, since if the whole pod is restarted, then maybe clickhouse was upgraded.
// Easy to copy, but for now decided that better to do full backup in this case.
func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  rand.Seed(time.Now().UnixNano())

  err := start("nats://"+env.Get("NATS", "nats:4222"), logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}

type BackupManager struct {
  *clickhouse.BaseBackupManager

  backupParentDir string
}

func start(natsUrl string, logger *zap.Logger) error {
  logger.Info("started", zap.String("nats", natsUrl))
  nc, err := nats.Connect(natsUrl)
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
    BaseBackupManager: baseBackupManager,
    backupParentDir:   filepath.Join(baseBackupManager.ClickhouseDir, "backup"),
  }

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
    backupName := time.Now().UTC().Format(timeFormat)
    backupDir := filepath.Join(backupManager.backupParentDir, backupName)
    err = backupManager.backup(backupDir, backupName)
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))

      err = os.RemoveAll(backupDir)
      if err != nil {
        logger.Error("cannot remove", zap.Error(err))
      }
    }
  }

  return nil
}

func (t *BackupManager) backup(backupDir string, backupName string) (err error) {
  logger := t.Logger.With(zap.String("backup", backupName))

  defer func() {
    if r := recover(); r != nil {
      err = errors.New("panic: " + string(debug.Stack()))
    }
  }()

  diffFromPath, err := t.removeOldLocalBackups(t.backupParentDir, maxIncrementalBackupCount)
  if err != nil {
    return errors.WithStack(err)
  }

  if len(diffFromPath) != 0 {
    key := diffFromPath + ".tar"
    info, err := t.Client.StatObjectWithContext(t.TaskContext, t.Bucket, key, minio.StatObjectOptions{})
    if err != nil || info.Err != nil || info.Size == 0 {
      logger.Warn("incremental backup is not created because previous backup doesn't exist on remote side", zap.String("remoteBackupPath", key), zap.String("bucket", t.Bucket), zap.Any("endpoint", t.Client.EndpointURL()))
      diffFromPath = ""
    }
  }

    estimatedTarSize, err := t.createBackup(backupDir, logger)
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("upload", zap.String("backup", backupDir))
  task := Task{
    metadataDir:      filepath.Join(t.ClickhouseDir, "metadata"),
    backupDir:        backupDir,
    estimatedTarSize: estimatedTarSize,
    diffFromPath:     diffFromPath,
    logger:           logger,
  }

  err = t.upload(backupName+".tar", task)
  if err != nil {
    return errors.WithStack(err)
  }

  if t.TaskContext.Err() != nil {
    return nil
  }

  logger.Info("uploaded")

  return nil
}

func (t *BackupManager) createBackup(backupDir string, logger *zap.Logger) (int64, error) {
  _, err := os.Stat(backupDir)
  if err == nil || !os.IsNotExist(err) {
    return 0, errors.Errorf("backup '%s' already exists", backupDir)
  }

  logger.Info("create")
  estimatedTarSize, err := t.freezeAndMoveToBackupDir(logger, backupDir)
  if err != nil {
    return 0, errors.WithStack(err)
  }

  return estimatedTarSize, nil
}
