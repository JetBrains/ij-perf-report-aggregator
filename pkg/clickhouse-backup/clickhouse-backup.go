package clickhouse_backup

import (
  "github.com/AlexAkulov/clickhouse-backup/pkg/backup"
  "github.com/AlexAkulov/clickhouse-backup/pkg/config"
  "github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
)

// example: if data collected each 3 hours, will be 8 backup per day, so, upload full backup at least once a day
const MaxIncrementalBackupCount = 4

func CreateBackuper() *backup.Backuper {
  backupConfig := config.DefaultConfig()
  backupConfig.General.RemoteStorage = "s3"
  backupConfig.General.BackupsToKeepRemote = MaxIncrementalBackupCount * 32
  backupConfig.ClickHouse.Host = "127.0.0.1"
  backupConfig.S3.AccessKey = util.GetEnvOrFileOrPanic("S3_ACCESS_KEY", "/etc/s3/accessKey")
  backupConfig.S3.SecretKey = util.GetEnvOrFileOrPanic("S3_SECRET_KEY", "/etc/s3/secretKey")
  backupConfig.S3.Bucket = util.GetEnvOrFileOrPanic("S3_BUCKET", "/etc/s3/bucket")
  backupConfig.S3.Region = "eu-west-1"
  backupConfig.S3.AllowMultipartDownload = true
  backuper := backup.NewBackuper(backupConfig)
  return backuper
}
