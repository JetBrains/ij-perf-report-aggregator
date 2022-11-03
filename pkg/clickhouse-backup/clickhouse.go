package clickhouse_backup

import (
  "bytes"
  "context"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
  "github.com/develar/errors"
  "github.com/minio/minio-go/v7"
  "github.com/minio/minio-go/v7/pkg/credentials"
  "github.com/valyala/fastjson"
  "go.deanishe.net/env"
  "go.uber.org/zap"
  "os"
  "os/exec"
  "path/filepath"
  "strings"
)

const MetaFileName = "meta.json.gz"
const InfoFileName = "info.json"

type TableInfo struct {
  Name         string `ch:"name" json:"name"`
  Uuid         string `ch:"uuid" json:"uuid"`
  MetadataPath string `ch:"metadata_path" json:"metadataPath"`
  Database     string `ch:"database" json:"db"`
}

type DbInfo struct {
  Name string `ch:"name" json:"name"`
}

type MappingInfo struct {
  Tables []TableInfo `ch:"tables" json:"tables"`
  Db     []DbInfo    `ch:"db" json:"db"`
}

type BackupManager struct {
  Bucket        string
  Client        *minio.Client
  TaskContext   context.Context
  ClickhouseDir string

  Logger *zap.Logger
}

func CreateBackupManager(taskContext context.Context, logger *zap.Logger) (*BackupManager, error) {
  // do not try to use doppler on K8S
  if len(os.Getenv("KUBERNETES_SERVICE_HOST")) == 0 {
    cmd := exec.Command("doppler", "secrets", "download", "--project", "s3", "--config", "prd", "--no-file")
    stdout, err := cmd.Output()
    if err != nil {
      logger.Warn("failed to use doppler to retrieve credentials", zap.Error(err))
    } else {
      excludePrefix := []byte("DOPPLER_")
      fastjson.MustParseBytes(stdout).GetObject().Visit(func(key []byte, v *fastjson.Value) {
        if !bytes.HasPrefix(key, excludePrefix) {
          err = os.Setenv(string(key), string(v.GetStringBytes()))
          if err != nil {
            logger.Fatal("cannot set env", zap.Error(err))
          }
        }
      })
    }
  }

  endpoint, err := util.GetEnvOrFile("S3_ENDPOINT", "/etc/s3/endpoint")
  if err != nil {
    if os.IsNotExist(err) {
      endpoint = "s3.amazonaws.com"
    } else {
      return nil, errors.WithStack(err)
    }
  }

  endpoint = strings.TrimSuffix(strings.TrimPrefix(endpoint, "https://"), "/")

  accessKey, err := util.GetEnvOrFile("S3_ACCESS_KEY", "/etc/s3/accessKey")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  secretKey, err := util.GetEnvOrFile("S3_SECRET_KEY", "/etc/s3/secretKey")
  if err != nil {
    return nil, errors.WithStack(err)
  }

  bucket, err := util.GetEnvOrFile("S3_BUCKET", "/etc/s3/bucket")
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

  return &BackupManager{
    Bucket:        bucket,
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

type MetaFile struct {
  RequiredBackup      string   `json:"requiredBackup"`
  EstimatedBackupSize int64    `json:"estimatedBackupSize"`
  Hardlinks           []string `json:"hardlinks"`
}
