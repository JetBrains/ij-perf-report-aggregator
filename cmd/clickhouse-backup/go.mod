module perf-stats/clickhouse-backup

go 1.13

replace github.com/JetBrains/ij-perf-report-aggregator/common => ../../pkg

replace github.com/AlexAkulov/clickhouse-backup => /Volumes/data/Documents/clickhouse-backup

require (
	github.com/JetBrains/ij-perf-report-aggregator/common v0.0.0-00010101000000-000000000000
	github.com/robfig/cron/v3 v3.0.0
	go.uber.org/zap v1.11.0
)
