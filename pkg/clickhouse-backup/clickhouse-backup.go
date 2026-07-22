package clickhouse_backup

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/valyala/fastjson"
)

// BackupsToKeepRemote keeps a month of daily restore points; self-contained backups are ~28 GB each
const BackupsToKeepRemote = 30

func BinaryPath() string {
	if p := os.Getenv("CLICKHOUSE_BACKUP_BIN"); p != "" {
		return p
	}
	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		return filepath.Join(home, "clickhouse-backup")
	}
	return "/usr/bin/clickhouse-backup"
}

// backupEnv returns the process environment for the clickhouse-backup binary,
// see https://github.com/Altinity/clickhouse-backup#default-config
func backupEnv() []string {
	overrides := map[string]string{
		"REMOTE_STORAGE":         "s3",
		"BACKUPS_TO_KEEP_REMOTE": strconv.Itoa(BackupsToKeepRemote),
		"CLICKHOUSE_HOST":        "127.0.0.1",
		// must be pinned: kubernetes service links inject CLICKHOUSE_PORT=tcp://<ip>:9000,
		// which breaks the binary's numeric port parsing
		"CLICKHOUSE_PORT":             "9000",
		"S3_ALLOW_MULTIPART_DOWNLOAD": "true",
		// the default is derived from the node CPU count and OOMs the sidecar's memory limit
		"UPLOAD_CONCURRENCY": "2",
		// v2 requires backups to live under a non-empty prefix, disjoint from object_disk_path;
		// legacy v1 backups at the bucket root can be targeted by explicitly setting S3_PATH=""
		"S3_PATH": envOrDefault("S3_PATH", "backup"),
		// where backups keep server-side copies of s3-disk (object disk) data — this is what
		// makes backups self-contained, unlike v1 which stored only the pointer files
		"S3_OBJECT_DISK_PATH": "object_disks",
		// local restores copy object-disk data from AWS into MinIO — different endpoints,
		// so server-side copy is impossible and the data must stream through the client;
		// in k8s source and destination are the same bucket and server-side copy is faster
		"ALLOW_OBJECT_DISK_STREAMING": envOrDefault("ALLOW_OBJECT_DISK_STREAMING",
			strconv.FormatBool(os.Getenv("KUBERNETES_SERVICE_HOST") == "")),
		"S3_ACCESS_KEY": util.GetEnvOrFileOrPanic("S3_ACCESS_KEY", "/etc/s3/accessKey"),
		"S3_SECRET_KEY": util.GetEnvOrFileOrPanic("S3_SECRET_KEY", "/etc/s3/secretKey"),
		"S3_BUCKET":     util.GetEnvOrFileOrPanic("S3_BUCKET", "/etc/s3/bucket"),
		"S3_REGION":     util.GetEnv("S3_REGION", "eu-west-1"),
	}

	result := make([]string, 0, len(os.Environ())+len(overrides))
	for _, kv := range os.Environ() {
		key, _, _ := strings.Cut(kv, "=")
		if _, ok := overrides[key]; !ok {
			result = append(result, kv)
		}
	}
	for key, value := range overrides {
		result = append(result, key+"="+value)
	}
	return result
}

func envOrDefault(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func Run(ctx context.Context, args ...string) error {
	cmd := exec.CommandContext(ctx, BinaryPath(), args...)
	cmd.Env = backupEnv()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("clickhouse-backup %s: %w", strings.Join(args, " "), err)
	}
	return nil
}

type RemoteBackup struct {
	BackupName  string `json:"BackupName"`
	Description string `json:"Description"`
}

// LatestRemoteBackup returns the name of the most recent remote backup.
// Broken backups have a zero creation date, so `latest` never resolves to one.
func LatestRemoteBackup(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, BinaryPath(), "list", "remote", "latest", "--format", "json")
	cmd.Env = backupEnv()
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("clickhouse-backup list remote: %w", err)
	}

	var backups []RemoteBackup
	err = json.Unmarshal(stdout, &backups)
	if err != nil {
		return "", fmt.Errorf("cannot parse backup list %q: %w", stdout, err)
	}
	if len(backups) == 0 {
		return "", errors.New("no remote backup")
	}
	backup := backups[len(backups)-1]
	if strings.HasPrefix(backup.Description, "broken") {
		return "", fmt.Errorf("latest remote backup %s is broken: %s", backup.BackupName, backup.Description)
	}
	return backup.BackupName, nil
}

func SetS3EnvForLocalRun(ctx context.Context) {
	cmd := exec.CommandContext(ctx, "doppler", "secrets", "download", "--project", "s3", "--config", "prd", "--no-file")
	stdout, err := cmd.Output()
	if err != nil {
		log.Println("failed to use doppler to retrieve credentials", err)
	} else {
		excludePrefix := []byte("DOPPLER_")
		fastjson.MustParseBytes(stdout).GetObject().Visit(func(key []byte, v *fastjson.Value) {
			// do not override explicitly set variables — allows pointing a local run
			// at MinIO (S3_ENDPOINT, S3_BUCKET, ...) without doppler clobbering it
			if !bytes.HasPrefix(key, excludePrefix) && os.Getenv(string(key)) == "" {
				err = os.Setenv(string(key), string(v.GetStringBytes()))
				if err != nil {
					log.Println("cannot set env", err)
				}
			}
		})
	}
}
