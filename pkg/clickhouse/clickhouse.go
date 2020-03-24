package clickhouse

import (
  "context"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v6"
  "go.uber.org/zap"
  "os"
)

const MetaFileName = "meta.json"

type BaseBackupManager struct {
  Bucket      string
  Client      *minio.Client
  TaskContext context.Context
  LocalPath   string

  Logger *zap.Logger
}

func CreateBaseBackupManager(taskContext context.Context, logger *zap.Logger) (*BaseBackupManager, error) {
  client, err := minio.New(os.Getenv("S3_ENDPOINT"), os.Getenv("S3_ACCESS_KEY"), os.Getenv("S3_SECRET_KEY"), true)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return &BaseBackupManager{
    Bucket:      os.Getenv("S3_BUCKET"),
    Client:      client,
    TaskContext: taskContext,
    LocalPath:   getClickhouseDir(),
    Logger:      logger,
  }, nil
}

func getClickhouseDir() string {
  localPath := os.Getenv("CLICKHOUSE_DATA_PATH")
  if len(localPath) == 0 {
    localPath = "/var/lib/clickhouse"
  }
  return localPath
}

type MetaFile struct {
  RequiredBackup      string   `json:"requiredBackup"`
  EstimatedBackupSize int64    `json:"estimatedBackupSize"`
  Hardlinks           []string `json:"hardlinks"`
}
