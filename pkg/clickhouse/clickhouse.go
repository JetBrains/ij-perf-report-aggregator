package clickhouse

import (
  "context"
  "github.com/deanishe/go-env"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v6"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "strings"
)

const MetaFileName = "meta.json.gz"

type BaseBackupManager struct {
  Bucket        string
  Client        *minio.Client
  TaskContext   context.Context
  ClickhouseDir string

  Logger *zap.Logger
}

func CreateBaseBackupManager(taskContext context.Context, logger *zap.Logger) (*BaseBackupManager, error) {
  endpoint, err := getEnvOrFile("S3_ENDPOINT", "/etc/s3/endpoint")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  endpoint = strings.TrimSuffix(strings.TrimPrefix(endpoint, "https://"), "/")

  accessKey, err := getEnvOrFile("S3_ACCESS_KEY", "/etc/s3/accessKey")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  secretKey, err := getEnvOrFile("S3_SECRET_KEY", "/etc/s3/secretKey")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  client, err := minio.New(endpoint, accessKey, secretKey, true)
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return &BaseBackupManager{
    Bucket:        os.Getenv("S3_BUCKET"),
    Client:        client,
    TaskContext:   taskContext,
    ClickhouseDir: env.GetString("CLICKHOUSE_DATA_PATH", "/var/lib/clickhouse"),
    Logger:        logger,
  }, nil
}

func getEnvOrFile(envName string, file string) (string, error) {
  v := os.Getenv(envName)
  if len(v) == 0 {
    b, err := ioutil.ReadFile(file)
    if err != nil {
      return "", errors.WithStack(err)
    }
    return string(b), err
  }
  return v, nil
}

type MetaFile struct {
  RequiredBackup      string   `json:"requiredBackup"`
  EstimatedBackupSize int64    `json:"estimatedBackupSize"`
  Hardlinks           []string `json:"hardlinks"`
}
