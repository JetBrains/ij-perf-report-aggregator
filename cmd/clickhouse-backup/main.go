package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	clickhousebackup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/nats-io/nats.go"
)

func main() {
	err := start("nats://" + util.GetEnv("NATS", "nats:4222"))
	if err != nil {
		slog.Error("cannot start backup", "err", err)
		os.Exit(1)
	}
}

func start(natsUrl string) error {
	taskContext, cancel := util.CreateCommandContext()
	defer cancel()

	if os.Getenv("KUBERNETES_SERVICE_HOST") == "" {
		clickhousebackup.SetS3EnvForLocalRun(taskContext)
	}

	if util.GetEnvAsBool("DO_BACKUP", false) {
		return executeBackup(taskContext)
	}

	slog.Info("started", "nats", natsUrl)
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		return fmt.Errorf("cannot connect to nats: %w", err)
	}

	sub, err := nc.SubscribeSync("db.backup")
	if err != nil {
		return fmt.Errorf("cannot subscribe to db.backup: %w", err)
	}

	lastBackupTime := time.Time{}
	for taskContext.Err() == nil {
		_, err = sub.NextMsgWithContext(taskContext)
		if err != nil {
			contextError := taskContext.Err()
			if contextError != nil {
				slog.Info("cancelled", "reason", contextError)
				return nil
			}
			return fmt.Errorf("cannot receive message: %w", err)
		}

		if taskContext.Err() != nil {
			return nil
		}

		if time.Since(lastBackupTime) < 24*time.Hour {
			// do not create backups too often
			slog.Info("backup request skipped", "reason", "time threshold", "lastBackupTime", lastBackupTime)
			continue
		}

		slog.Info("backup requested")
		err = executeBackup(taskContext)
		if err != nil {
			slog.Error("cannot backup", "error", err)
		} else {
			lastBackupTime = time.Now()
		}
	}

	return nil
}

func executeBackup(taskContext context.Context) error {
	// the backup name is auto-generated from the current UTC time,
	// remote retention (BACKUPS_TO_KEEP_REMOTE) is applied after upload
	return clickhousebackup.Run(taskContext, "create_remote", "--delete-source")
}
