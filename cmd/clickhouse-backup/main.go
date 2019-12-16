package main

import (
  "flag"
  "fmt"
  "github.com/AlexAkulov/clickhouse-backup/pkg/chbackup"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/kelseyhightower/envconfig"
  "github.com/nats-io/nats.go"
  "go.uber.org/zap"
  "strings"
  "time"
)

var lastLocalBackups []string

const backupSuffix = ".tar"
var incrementalBackupCounter = 0
const maxIncrementalBackupCount = 4

func main() {
  logger := util.CreateLogger()
  defer func() {
    _ = logger.Sync()
  }()

  natsUrl := flag.String("nats", "", "The NATS URL.")
  flag.Parse()

  err := start(*natsUrl, logger)
  if err != nil {
    logger.Fatal(fmt.Sprintf("%+v", err))
  }
}

func start(natsUrl string, logger *zap.Logger) error {
  nc, err := nats.Connect(natsUrl)
  if err != nil {
    return errors.WithStack(err)
  }

  config := chbackup.DefaultConfig()
  config.General.DisableProgressBar = true
  config.General.BackupsToKeepLocal = maxIncrementalBackupCount + 1
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

    logger.Info("backup requested")

    err = backup(config)
    if err != nil {
      logger.Error("cannot backup", zap.Error(err))
    }
  }

  return nil
}

func backup(config *chbackup.Config) error {
  backupName := time.Now().Format(chbackup.BackupTimeFormat)

  err := chbackup.CreateBackup(*config, backupName, "")
  if err != nil {
    return errors.WithStack(err)
  }

  var diffFrom string
  // first or each 4 backup is full
  if len(lastLocalBackups) < 1 || incrementalBackupCounter > maxIncrementalBackupCount {
    incrementalBackupCounter = 0
    diffFrom = ""
  } else {
    diffFrom = strings.TrimSuffix(lastLocalBackups[len(lastLocalBackups)-1], backupSuffix)
  }

  err = chbackup.Upload(*config, backupName, diffFrom)
  if err != nil {
    return errors.WithStack(err)
  }

  incrementalBackupCounter++

  if len(lastLocalBackups) > 2 {
    lastLocalBackups = append(lastLocalBackups[1:], backupName+backupSuffix)
  } else {
    lastLocalBackups = append(lastLocalBackups, backupName+backupSuffix)
  }

  return nil
}