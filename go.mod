module github.com/JetBrains/ij-perf-report-aggregator

go 1.13

require (
	cloud.google.com/go v0.52.0 // indirect
	cloud.google.com/go/storage v1.5.0 // indirect
	facette.io/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/AlexAkulov/clickhouse-backup v0.5.1
	github.com/ClickHouse/clickhouse-go v1.3.13
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/VictoriaMetrics/fastcache v1.5.7
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/araddon/dateparse v0.0.0-20190622164848-0fb0a474d195
	github.com/asaskevich/govalidator v0.0.0-20200108200545-475eaeb16496
	github.com/aws/aws-sdk-go v1.28.9 // indirect
	github.com/cespare/xxhash v1.1.0
	github.com/deanishe/go-env v0.4.0
	github.com/develar/errors v0.9.0
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.9
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.9.8 // indirect
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-runewidth v0.0.8 // indirect
	github.com/mattn/go-sqlite3 v2.0.0+incompatible // indirect
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/nats-io/jwt v0.3.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.0 // indirect
	github.com/nats-io/nats.go v1.9.1
	github.com/olekukonko/tablewriter v0.0.4
	github.com/panjf2000/ants/v2 v2.2.2
	github.com/pierrec/lz4 v2.4.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/spaolacci/murmur3 v1.0.1-0.20190317074736-539464a789e9 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/tdewolff/minify/v2 v2.7.2
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/fastjson v1.4.5
	github.com/valyala/quicktemplate v1.4.1
	go.uber.org/atomic v1.5.1
	go.uber.org/multierr v1.4.0
	go.uber.org/zap v1.13.0
	golang.org/x/crypto v0.0.0-20200128174031-69ecbb4d6d5d // indirect
	golang.org/x/exp v0.0.0-20200119233911-0405dc783f0a // indirect
	golang.org/x/sys v0.0.0-20200124204421-9fbb57f87de9 // indirect
	golang.org/x/tools v0.0.0-20200129045341-207d3de1faaf // indirect
	google.golang.org/genproto v0.0.0-20200128133413-58ce757ed39b // indirect
	google.golang.org/grpc v1.27.0 // indirect
	gopkg.in/sakura-internet/go-rison.v3 v3.1.0
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/mattn/go-sqlite3 => github.com/mattn/go-sqlite3 v1.11.1-0.20191105054421-67c1376b46fb

replace github.com/AlexAkulov/clickhouse-backup => github.com/AlexAkulov/clickhouse-backup v0.5.2-0.20191205080651-6160f7f6d90f
