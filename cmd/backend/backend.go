package main

import (
	"log/slog"
	"os"

	"github.com/JetBrains/ij-perf-report-aggregator/pkg/server"
	"github.com/JetBrains/ij-perf-report-aggregator/pkg/util"
)

func main() {
	err := server.Serve(util.GetEnv("CLICKHOUSE", server.DefaultDbUrl), util.GetEnv("NATS", ""))
	if err != nil {
		slog.Error("error on starting backend", "error", err)
		os.Exit(1)
	}
}
