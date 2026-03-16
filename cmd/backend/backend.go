package main

import (
	"log/slog"
	"os"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
	"github.com/joho/godotenv"
	"go.deanishe.net/env"
)

func main() {
	// Load .env file when present (local development).
	// Existing environment variables are never overwritten, so production
	// deployments that inject vars via the environment are unaffected.
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		slog.Warn("could not load .env file", "error", err)
	}

	err := server.Serve(env.Get("CLICKHOUSE", server.DefaultDbUrl), env.Get("NATS", ""))
	if err != nil {
		slog.Error("error on starting backend", "error", err)
		os.Exit(1)
	}
}
