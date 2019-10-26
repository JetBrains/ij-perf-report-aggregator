module perf-stats/clickhouse-backup

go 1.13

replace github.com/JetBrains/ij-perf-report-aggregator/common => ../../pkg

replace github.com/AlexAkulov/clickhouse-backup => github.com/develar/clickhouse-backup v0.4.3-0.20191026115939-e08f51c1c381

require (
	github.com/AlexAkulov/clickhouse-backup v0.4.2
	github.com/JetBrains/ij-perf-report-aggregator/common v0.0.0-00010101000000-000000000000
	github.com/develar/errors v0.9.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/robfig/cron/v3 v3.0.0
	go.uber.org/zap v1.11.0
)
