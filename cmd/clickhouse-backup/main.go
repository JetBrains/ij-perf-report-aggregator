package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Altinity/clickhouse-backup/pkg/backup"
	"github.com/Altinity/clickhouse-backup/pkg/status"
	clickhousebackup "github.com/JetBrains/ij-perf-report-aggregator/pkg/clickhouse-backup"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
	"github.com/nats-io/nats.go"
	"go.deanishe.net/env"
)

func main() {
	err := start("nats://" + env.Get("NATS", "nats:4222"))
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

	backuper := clickhousebackup.CreateBackuper()

	if env.GetBool("DO_BACKUP") {
		err := executeBackup(taskContext, backuper)
		return err
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
		err = executeBackup(taskContext, backuper)
		if err != nil {
			slog.Error("cannot backup", "error", err)
		} else {
			lastBackupTime = time.Now()
		}
	}

	return nil
}

func executeBackup(taskContext context.Context, backuper *backup.Backuper) error {
	backupName := backup.NewBackupName()
	logger := slog.With("backup", backupName)

	err := backuper.CreateBackup(backupName, "", nil, false, false, false, false, "unknown", status.NotFromAPI)
	if err != nil {
		return fmt.Errorf("cannot create backup: %w", err)
	}

	if taskContext.Err() != nil {
		return nil
	}

	logger.Info("upload")
	err = backuper.Upload(backupName, "", "", "", nil, false, false, status.NotFromAPI)
	if err != nil {
		return err
	}

	if taskContext.Err() != nil {
		return nil
	}

	logger.Info("uploaded")
	return nil
}
