module perf-stats/tc-collector

go 1.13

replace github.com/JetBrains/ij-perf-report-aggregator/common => ../../pkg

require (
	github.com/JetBrains/ij-perf-report-aggregator/common v0.0.0-00010101000000-000000000000
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/araddon/dateparse v0.0.0-20190622164848-0fb0a474d195
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/develar/errors v0.9.0
	github.com/json-iterator/go v1.1.7
	github.com/magiconair/properties v1.8.1
	github.com/nats-io/nats.go v1.8.1
	github.com/valyala/quicktemplate v1.3.1
	go.uber.org/atomic v1.4.0
	go.uber.org/zap v1.11.0
)
