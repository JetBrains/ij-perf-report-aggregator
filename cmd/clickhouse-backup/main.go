package main

import (
  "context"
  "fmt"
  "github.com/AlexAkulov/clickhouse-backup/pkg/backup"
  "github.com/AlexAkulov/clickhouse-backup/pkg/status"
  clickhousebackup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/nats-io/nats.go"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "time"
)

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

func start(natsUrl string, logger *zap.Logger) error {
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  backuper := clickhousebackup.CreateBackuper()

  if env.GetBool("DO_BACKUP") {
    _, err := executeBackup(backuper, taskContext, 0, logger)
    return err
  }

  logger.Info("started", zap.String("nats", natsUrl))
  nc, err := nats.Connect(natsUrl)
  if err != nil {
    return errors.WithStack(err)
  }

  sub, err := nc.SubscribeSync("db.backup")
  if err != nil {
    return errors.WithStack(err)
  }

  lastBackupTime := time.Time{}
  backupCount := 0
  for taskContext.Err() == nil {
    _, err = sub.NextMsgWithContext(taskContext)
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

    if time.Now().Sub(lastBackupTime) < 4*time.Hour {
      // do not create backups too often
      logger.Info("backup request skipped", zap.String("reason", "time threshold"), zap.Time("lastBackupTime", lastBackupTime))
      return nil
    }

    logger.Info("backup requested")
    backupCount, err = executeBackup(backuper, taskContext, backupCount, logger)
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))
    } else {
      lastBackupTime = time.Now()
    }
  }

  return nil
}

func executeBackup(backuper *backup.Backuper, taskContext context.Context, backupCount int, logger *zap.Logger) (int, error) {
  backupName := backup.NewBackupName()
  logger = logger.With(zap.String("backup", backupName))

  diffFromRemote := ""
  if backupCount < clickhousebackup.MaxIncrementalBackupCount {
    remoteBackups, err := backuper.GetRemoteBackups(taskContext, true)
    if err != nil {
      logger.Error("cannot get remote backup list", zap.Error(err))
    } else if len(remoteBackups) > 0 {
      diffFromRemote = remoteBackups[len(remoteBackups)-1].BackupName
    }
  }

  err := backuper.CreateBackup(backupName, "", nil, false, false, false, "unknown", status.NotFromAPI)
  if err != nil {
    return backupCount, errors.WithStack(err)
  }

  if taskContext.Err() != nil {
    return backupCount, nil
  }

  logger.Info("upload", zap.String("diffFromRemote", diffFromRemote))
  err = backuper.Upload(backupName, "", diffFromRemote, "", nil, false, false, status.NotFromAPI)
  if err != nil {
    return backupCount, err
  }

  if taskContext.Err() != nil {
    return backupCount, nil
  }

  logger.Info("uploaded")
  return backupCount + 1, nil
}
