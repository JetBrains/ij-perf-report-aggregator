package main

import (
  "context"
  "fmt"
  "github.com/AlexAkulov/clickhouse-backup/pkg/backup"
  "github.com/AlexAkulov/clickhouse-backup/pkg/config"
  "github.com/AlexAkulov/clickhouse-backup/pkg/status"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/nats-io/nats.go"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "time"
)

// example: if data collected each 3 hours, will be 8 backup per day, so, upload full backup at least once a day
const maxIncrementalBackupCount = 4

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

  backupConfig := config.DefaultConfig()
  backupConfig.General.RemoteStorage = "s3"
  backupConfig.General.BackupsToKeepRemote = maxIncrementalBackupCount * 2
  backupConfig.S3.AccessKey = util.GetEnvOrFileOrPanic("S3_ACCESS_KEY", "/etc/s3/accessKey")
  backupConfig.S3.SecretKey = util.GetEnvOrFileOrPanic("S3_SECRET_KEY", "/etc/s3/secretKey")
  backupConfig.S3.Bucket = util.GetEnvOrFileOrPanic("S3_BUCKET", "/etc/s3/bucket")
  backupConfig.S3.Region = "eu-west-1"
  backupConfig.S3.AllowMultipartDownload = true
  backuper := backup.NewBackuper(backupConfig)

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

    logger.Info("backup requested")
    backupCount, err = executeBackup(backuper, taskContext, backupCount, logger)
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))
    }
  }

  return nil
}

func executeBackup(backuper *backup.Backuper, taskContext context.Context, backupCount int, logger *zap.Logger) (int, error) {
  backupName := time.Now().UTC().Format(backup.NewBackupName())
  logger = logger.With(zap.String("backup", backupName))

  diffFromRemote := ""
  if backupCount < maxIncrementalBackupCount {
    remoteBackups, err := backuper.GetRemoteBackups(taskContext, false)
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
