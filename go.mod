module github.com/JetBrains/ij-perf-report-aggregator

go 1.13

require (
	facette.io/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/AlexAkulov/clickhouse-backup v0.5.1
	github.com/ClickHouse/clickhouse-go v1.3.12
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/VictoriaMetrics/fastcache v1.5.4
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/araddon/dateparse v0.0.0-20190622164848-0fb0a474d195
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/aws/aws-sdk-go v1.25.48 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/deanishe/go-env v0.4.0
	github.com/develar/errors v0.9.0
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.8
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.9.4 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/mattn/go-sqlite3 v2.0.0+incompatible // indirect
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nats-io/jwt v0.3.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.0 // indirect
	github.com/nats-io/nats.go v1.9.1
	github.com/olekukonko/tablewriter v0.0.3
	github.com/panjf2000/ants/v2 v2.2.2
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron/v3 v3.0.0
	github.com/rs/cors v1.7.0
	github.com/spaolacci/murmur3 v1.0.1-0.20190317074736-539464a789e9 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/tdewolff/minify/v2 v2.6.1
	github.com/tdewolff/parse/v2 v2.4.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/fastjson v1.4.1
	github.com/valyala/quicktemplate v1.4.1
	go.uber.org/atomic v1.5.1
	go.uber.org/multierr v1.4.0
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20191202143827-86a70503ff7e // indirect
	golang.org/x/exp v0.0.0-20191129062945-2f5052295587 // indirect
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/net v0.0.0-20191204025024-5ee1b9f4859a // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sys v0.0.0-20191204072324-ce4227a45e2e // indirect
	golang.org/x/tools v0.0.0-20191205060818-73c7173a9f7d // indirect
	google.golang.org/genproto v0.0.0-20191203220235-3fa9dbf08042 // indirect
	gopkg.in/sakura-internet/go-rison.v3 v3.1.0
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.11.1-0.20191105054421-67c1376b46fb

replace github.com/AlexAkulov/clickhouse-backup => github.com/develar/clickhouse-backup v0.4.3-0.20191119202234-e508b07c24c7
