module perf-stats/report-aggregator

go 1.13

replace github.com/JetBrains/ij-perf-report-aggregator/common => ../../pkg

require (
	github.com/JetBrains/ij-perf-report-aggregator/common v0.0.0-00010101000000-000000000000
	github.com/alecthomas/kingpin v2.2.6+incompatible
	go.uber.org/zap v1.11.0
)
