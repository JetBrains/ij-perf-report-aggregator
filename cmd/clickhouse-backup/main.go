package main

import (
  "context"
  "fmt"
  "github.com/Altinity/clickhouse-backup/pkg/backup"
  "github.com/Altinity/clickhouse-backup/pkg/status"
  clickhousebackup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/nats-io/nats.go"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "os"
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

  if len(os.Getenv("KUBERNETES_SERVICE_HOST")) == 0 {
    clickhousebackup.SetS3EnvForLocalRun()
  }

  backuper := clickhousebackup.CreateBackuper()

  if env.GetBool("DO_BACKUP") {
    err := executeBackup(taskContext, backuper, logger)
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

    if time.Since(lastBackupTime) < 24*time.Hour {
      // do not create backups too often
      logger.Info("backup request skipped", zap.String("reason", "time threshold"), zap.Time("lastBackupTime", lastBackupTime))
      continue
    }

    logger.Info("backup requested")
    err = executeBackup(taskContext, backuper, logger)
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))
    } else {
      lastBackupTime = time.Now()
    }
  }

  return nil
}

func executeBackup(taskContext context.Context, backuper *backup.Backuper, logger *zap.Logger) error {
  backupName := backup.NewBackupName()
  logger = logger.With(zap.String("backup", backupName))

  err := backuper.CreateBackup(backupName, "", nil, false, false, false, false, "unknown", status.NotFromAPI)
  if err != nil {
    return errors.WithStack(err)
  }

  if taskContext.Err() != nil {
    return nil
  }

  logger.Info("upload")
  err = backuper.Upload(backupName, "", "", "", nil, false, false, status.NotFromAPI)
  if err != nil {
    return err
  }

  if taskContext.Err() != nil {
    return nil
  }

  logger.Info("uploaded")
  return nil
}
