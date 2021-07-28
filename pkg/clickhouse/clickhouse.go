package clickhouse

import (
  "context"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v7"
  "github.com/minio/minio-go/v7/pkg/credentials"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
)

const MetaFileName = "meta.json.gz"
const InfoFileName = "info.json"

type TableInfo struct {
  Name         string `json:"name"`
  Uuid         string `json:"uuid"`
  MetadataPath string `db:"metadata_path" json:"metadataPath"`
  Database     string `json:"db"`
}

type DbInfo struct {
  Name string `json:"name"`
  Uuid string `json:"uuid"`
}

type MappingInfo struct {
  Tables []TableInfo `json:"tables"`
  Db     []DbInfo    `json:"db"`
}

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

  client, err := minio.New(endpoint, &minio.Options{
    Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
    Secure: !strings.HasPrefix(endpoint, "127.0.0.1:"),
  })
  if err != nil {
    return nil, errors.WithStack(err)
  }

  return &BaseBackupManager{
    Bucket:        os.Getenv("S3_BUCKET"),
    Client:        client,
    TaskContext:   taskContext,
    ClickhouseDir: GetClickhouseDir(),
    Logger:        logger,
  }, nil
}

func GetClickhouseDir() string {
  s := env.GetString("CLICKHOUSE_DATA_PATH", "/var/lib/clickhouse")
  if strings.HasPrefix(s, "~/") {
    homeDir, _ := os.UserHomeDir()
    return filepath.Join(homeDir, s[2:])
  } else {
    return s
  }
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
