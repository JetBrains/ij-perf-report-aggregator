package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/nats-io/nats.go"
  "go.uber.org/zap"
  "log"
  "os"
  "path/filepath"
  "runtime/debug"
  "strings"
  "time"
)

const backupSuffix = ".tar"

var incrementalBackupCounter = 0

const maxIncrementalBackupCount = 4
const backupsToKeepLocal = maxIncrementalBackupCount + 1

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
    err = backupManager.backup()
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))
    }
  }

  return nil
}

func (t *BackupManager) backup() error {
  defer func() {
    if r := recover(); r != nil {
      log.Print(string(debug.Stack()))
      t.Logger.Error("recovered", zap.ByteString("stack", debug.Stack()))
    }
  }()

  lastLocalBackups, err := t.removeOldLocalBackups(t.backupParentDir, backupsToKeepLocal)
  if err != nil {
    return errors.WithStack(err)
  }

  backupName := time.Now().UTC().Format("2006-01-02T15-04-05")
  logger := t.Logger.With(zap.String("backup", backupName))

  backupDir := filepath.Join(t.backupParentDir, backupName)
  estimatedTarSize, extraEntries, err := t.createBackup(backupDir, logger)
  if err != nil {
    return errors.WithStack(err)
  }

  var diffFrom string
  // first or each 4 backup is full
  if len(lastLocalBackups) == 0 || incrementalBackupCounter > maxIncrementalBackupCount {
    incrementalBackupCounter = 0
  } else {
    diffFrom = strings.TrimSuffix(lastLocalBackups[len(lastLocalBackups)-1], backupSuffix)
  }

  var diffFromPath string
  if len(diffFrom) != 0 {
    diffFromPath = filepath.Join(t.backupParentDir, diffFrom)
  }

  logger.Info("upload", zap.String("backup", backupDir))
  task := Task{
    backupDir:        backupDir,
    estimatedTarSize: estimatedTarSize,
    diffFromPath:     diffFromPath,
    extraEntries:     extraEntries,
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

  incrementalBackupCounter++
  return nil
}

func (t *BackupManager) createBackup(backupDir string, logger *zap.Logger) (int64, []FileEntry, error) {
  _, err := os.Stat(backupDir)
  if err == nil || !os.IsNotExist(err) {
    return 0, nil, errors.Errorf("backup '%s' already exists", backupDir)
  }

  logger.Info("create")
  estimatedTarSize, err := t.freezeAndMoveToBackupDir(logger, backupDir)
  if err != nil {
    return 0, nil, errors.WithStack(err)
  }

  // copy metadata to memory to ensure that it is not changed during compression
  logger.Debug("copy metadata")
  extraEntries, err := t.collectMetadata()
  if err != nil {
    return 0, nil, errors.WithStack(err)
  }

  for _, entry := range extraEntries {
    estimatedTarSize += int64(len(entry.data)) + estimatedTarHeaderSize
  }

  return estimatedTarSize, extraEntries, nil
}
