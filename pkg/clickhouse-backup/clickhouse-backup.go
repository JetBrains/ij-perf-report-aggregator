package clickhouse_backup

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/Altinity/clickhouse-backup/pkg/backup"
	"github.com/Altinity/clickhouse-backup/pkg/config"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/valyala/fastjson"
)

// MaxIncrementalBackupCount works like if data collected each 3 hours, will be 8 backup per day, so, upload full backup at least once a day
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

func SetS3EnvForLocalRun(ctx context.Context) {
	cmd := exec.CommandContext(ctx, "doppler", "secrets", "download", "--project", "s3", "--config", "prd", "--no-file")
	stdout, err := cmd.Output()
	if err != nil {
		log.Println("failed to use doppler to retrieve credentials", err)
	} else {
		excludePrefix := []byte("DOPPLER_")
		fastjson.MustParseBytes(stdout).GetObject().Visit(func(key []byte, v *fastjson.Value) {
			if !bytes.HasPrefix(key, excludePrefix) {
				err = os.Setenv(string(key), string(v.GetStringBytes()))
				if err != nil {
					log.Println("cannot set env", err)
				}
			}
		})
	}
}
