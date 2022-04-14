module github.com/JetBrains/ij-perf-report-aggregator

go 1.18

replace github.com/go-faster/ch => github.com/develar/ch v0.27.1

//replace github.com/go-faster/ch => /Volumes/data/Documents/ch

require (
	facette.io/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/ClickHouse/clickhouse-go v1.5.4
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cheggaaa/pb/v3 v3.0.8
	github.com/develar/errors v0.9.0
	github.com/dgraph-io/ristretto v0.1.0
	github.com/fatih/color v1.13.0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jmoiron/sqlx v1.3.4
	github.com/json-iterator/go v1.1.12
	github.com/klauspost/cpuid/v2 v2.0.12 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/magiconair/properties v1.8.6
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/minio-go/v7 v7.0.24
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats-server/v2 v2.1.4 // indirect
	github.com/nats-io/nats.go v1.14.0
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/panjf2000/ants/v2 v2.4.8
	github.com/pierrec/lz4 v2.4.1+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.8.2
	github.com/rs/xid v1.4.0 // indirect
	github.com/segmentio/ksuid v1.0.4
	github.com/stretchr/testify v1.7.1
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/fastjson v1.6.3
	github.com/valyala/quicktemplate v1.7.0
	github.com/zeebo/xxh3 v1.0.2
	go.deanishe.net/env v0.5.1
	go.uber.org/atomic v1.9.0
	go.uber.org/multierr v1.8.0
	go.uber.org/zap v1.21.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.10
	google.golang.org/protobuf v1.24.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

require (
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dmarkham/enumer v1.5.5 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/go-faster/ch v0.26.0
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.5.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/hashicorp/go-version v1.4.0 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pascaldekloe/name v1.0.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/segmentio/asm v1.1.4 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	go.opentelemetry.io/otel v1.6.3 // indirect
	go.opentelemetry.io/otel/metric v0.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.6.3 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220106191415-9b9b3d81d5e3 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/xerrors v0.0.0-20220411194840-2f41105eb62f // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

require (
	github.com/jackc/puddle v1.2.2-0.20220404125616-4e959849469a
	github.com/sakura-internet/go-rison/v4 v4.0.0
)
