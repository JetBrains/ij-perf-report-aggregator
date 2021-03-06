package main

import (
  "fmt"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/jmoiron/sqlx"
  "github.com/minio/minio-go/v7"
  "github.com/nats-io/nats.go"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "os"
  "path/filepath"
  "runtime/debug"
  "strings"
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
  taskContext, cancel := util.CreateCommandContext()
  defer cancel()

  baseBackupManager, err := clickhouse.CreateBaseBackupManager(taskContext, logger)
  if err != nil {
    return errors.WithStack(err)
  }
  backupManager := &BackupManager{
    BaseBackupManager: baseBackupManager,
    backupParentDir:   filepath.Join(baseBackupManager.ClickhouseDir, "backup"),
  }

  if env.GetBool("DO_BACKUP") {
    return executeBackup(backupManager)
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
    _ = executeBackup(backupManager)
  }

  return nil
}

func executeBackup(backupManager *BackupManager) error {
  backupName := time.Now().UTC().Format(timeFormat)
  backupDir := filepath.Join(backupManager.backupParentDir, backupName)
  err := backupManager.backup(backupDir, backupName)
  if err != nil {
    backupManager.Logger.Error("cannot backup", zap.Error(err))

    err = os.RemoveAll(backupDir)
    if err != nil {
      backupManager.Logger.Error("cannot remove", zap.Error(err))
    }
  }
  return err
}

func (t *BackupManager) backup(backupDir string, backupName string) (err error) {
  logger := t.Logger.With(zap.String("backup", backupName))

  defer func() {
    if r := recover(); r != nil {
      err = errors.New("panic: " + string(debug.Stack()))
    }
  }()

  localBackups, err := t.collectLocalBackups(t.backupParentDir)
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Debug("local backups", zap.Strings("names", localBackups), zap.Int("maxIncrementalBackupCount", maxIncrementalBackupCount))
  var diffFromPath string
  // incremental backup only if limit is not exceed - if exceeded it means that full backup should be created
  removeLocalBackups := len(localBackups) > maxIncrementalBackupCount
  if len(localBackups) > 0 && !removeLocalBackups {
    removeLocalBackups = false

    diffFromName := localBackups[len(localBackups)-1]
    key := diffFromName + ".tar"
    info, err := t.Client.StatObject(t.TaskContext, t.Bucket, key, minio.StatObjectOptions{})
    if err != nil || info.Err != nil || info.Size == 0 {
      logger.Warn("incremental backup is not created because previous backup doesn't exist on remote side", zap.String("remoteBackupPath", key), zap.String("bucket", t.Bucket), zap.Any("endpoint", t.Client.EndpointURL()))
    } else {
      diffFromPath = filepath.Join(t.backupParentDir, diffFromName)
    }
  }

  task := Task{
    metadataDir:      filepath.Join(t.ClickhouseDir, "metadata"),
    backupDir:        backupDir,
    diffFromPath:     diffFromPath,
    estimatedTarSize: 0,
    logger:           logger,
  }

  err = t.createBackup(&task)
  if err != nil {
    return errors.WithStack(err)
  }

  logger.Info("upload", zap.String("backup", backupDir))
  err = t.upload(backupName+".tar", task)
  if err != nil {
    return errors.WithStack(err)
  }

  if t.TaskContext.Err() != nil {
    return nil
  }

  logger.Info("uploaded")

  if removeLocalBackups {
    for _, name := range localBackups {
      t.Logger.Info("remove old local backup", zap.String("backup", name))
      err = os.RemoveAll(filepath.Join(t.backupParentDir, name))
      if err != nil {
        logger.Error("cannot remove local backup", zap.String("fileName", name), zap.Error(err))
        break
      }
    }
  }

  return nil
}

func getDatabaseNamesToBackup() []string {
  var result []string
  for _, s := range strings.Split(env.GetString("DB"), ",") {
    name := strings.TrimSpace(s)
    if len(name) > 0 {
      result = append(result, name)
    }
  }
  return result
}

func (t *BackupManager) createBackup(task *Task) error {
  backupDir := task.backupDir
  _, err := os.Stat(backupDir)
  if err == nil || !os.IsNotExist(err) {
    return errors.Errorf("backup '%s' already exists", backupDir)
  }

  db, err := sqlx.Open("clickhouse", "tcp://"+env.Get("CLICKHOUSE", "clickhouse:9000")+"?read_timeout=600&write_timeout=600&debug=0&compress=1&send_timeout=30000&receive_timeout=3000")
  if err != nil {
    return errors.WithStack(err)
  }

  defer util.Close(db, task.logger)

  task.dbNames = getDatabaseNamesToBackup()
  if len(task.dbNames) == 0 {
    // select non-empty db
    err = db.SelectContext(t.TaskContext, &task.dbNames, "select distinct database from system.tables where database != 'system'")
    if err != nil {
      return errors.WithStack(err)
    }
  }

  task.logger.Info("create")
  for _, dbName := range task.dbNames {
    estimatedTarSize, err := t.freezeAndMoveToBackupDir(db, dbName, backupDir, task.logger.With(zap.String("db", dbName)))
    if err != nil {
      return errors.WithStack(err)
    }

    task.estimatedTarSize += estimatedTarSize
  }

  return nil
}
