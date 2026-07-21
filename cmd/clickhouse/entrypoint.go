package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	clickhousebackup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
	"github.com/nats-io/nats.go"
	"go.deanishe.net/env"
)

//go:embed config.xml
var clickhouseConfig []byte

func main() {
	clickhouseExecutable := "/usr/bin/clickhouse"

	isLocalRun := os.Getenv("KUBERNETES_SERVICE_HOST") == ""
	ctx := context.Background()
	if isLocalRun {
		clickhouseExecutable = "/Users/maxim.kolmakov/clickhouse"
		clickhousebackup.SetS3EnvForLocalRun(ctx)
	}

	bucket := getEnvOrFile("S3_BUCKET", "/etc/s3/bucket")
	s3AccessKey := getEnvOrFile("S3_ACCESS_KEY", "/etc/s3/accessKey")
	s3SecretKey := getEnvOrFile("S3_SECRET_KEY", "/etc/s3/secretKey")

	restoreData := os.Getenv("RESTORE_DB") == "true"

	configFile := "/var/lib/clickhouse/config.xml"
	if isLocalRun {
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		// the file must be named config.xml: clickhouse-backup reads object disk credentials
		// from the preprocessed config, which clickhouse names after the main config file
		configFile = filepath.Join(workingDir, "deployment", "ch-local", "config.xml")
	}

	if restoreData {
		// the data dir is about to be wiped — refuse if another server still uses it
		dialer := net.Dialer{Timeout: time.Second}
		if conn, err := dialer.DialContext(ctx, "tcp", "127.0.0.1:9000"); err == nil {
			_ = conn.Close()
			log.Fatal("another clickhouse-server is already running on 127.0.0.1:9000, stop it before restore")
		}

		err := prepareConfigAndDir(isLocalRun, bucket, s3AccessKey, s3SecretKey, configFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd := exec.CommandContext(ctx, clickhouseExecutable, "server", "--config-file="+configFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	if restoreData {
		defer func() {
			err = cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				log.Println(err)
				_ = cmd.Process.Kill()
			}

			err = cmd.Wait()
			if err != nil {
				log.Println(err)
			}
		}()

		err = restoreDb(ctx)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	go func() {
		// wait for clickhouse server start
		time.Sleep(10 * time.Second)
		requestClearCache()
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func prepareConfigAndDir(isLocalRun bool, bucket string, s3AccessKey string, s3SecretKey string, configFile string) error {
	chDir := "/var/lib/clickhouse"
	if isLocalRun {
		chDir = env.GetString("CLICKHOUSE_DATA_PATH", "/Volumes/data/ij-perf-db/clickhouse")
	}

	entries, err := os.ReadDir(chDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	for _, entry := range entries {
		err = os.RemoveAll(filepath.Join(chDir, entry.Name()))
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	if !isLocalRun {
		s3Url := "https://" + bucket + ".s3.eu-west-1.amazonaws.com/data/"
		log.Print("S3 URL: " + s3Url)

		s := strings.NewReplacer(
			"$S3_URL", s3Url,
			"$S3_ACCESS_KEY", s3AccessKey,
			"$S3_SECRET_KEY", s3SecretKey,
		).Replace(string(clickhouseConfig))

		// /etc is not writeable
		err = os.WriteFile(configFile, []byte(s), 0o666)
		if err != nil {
			return err
		}
	}
	return nil
}

func restoreDb(ctx context.Context) error {
	// wait a little bit for clickhouse start
	time.Sleep(4 * time.Second)

	attemptCount := 3
	var backupName string
	var err error
	for i := range attemptCount {
		backupName, err = clickhousebackup.LatestRemoteBackup(ctx)
		if err == nil {
			break
		}
		log.Println("cannot get latest remote backup", err)
		if i < attemptCount-1 {
			time.Sleep(time.Duration((i+1)*3) * time.Second)
		}
	}
	if err != nil {
		return err
	}

	err = clickhousebackup.Run(ctx, "restore_remote", backupName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	// the downloaded backup copy is not needed after restore
	err = clickhousebackup.Run(ctx, "delete", "local", backupName)
	if err != nil {
		log.Println("cannot delete local backup copy", err)
	}

	log.Println("DB is restored (backup=" + backupName + ")")
	return nil
}

func requestClearCache() {
	url := os.Getenv("NATS")
	if url == "" {
		url = "nats://nats:4222"
	}

	nc, err := nats.Connect(url, nats.Name("NATS Sample Publisher"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	err = nc.Publish("server.clearCache", []byte("clickhouse"))
	if err != nil {
		log.Fatal(err)
	}

	err = nc.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func getEnvOrFile(envName string, file string) string {
	v := os.Getenv(envName)
	if v == "" {
		b, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		return string(b)
	}
	return v
}
