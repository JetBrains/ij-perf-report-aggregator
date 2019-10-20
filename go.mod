module report-aggregator

go 1.13

require (
	facette.io/natsort v0.0.0-20181210072756-2cd4dd1e2dcb
	github.com/VictoriaMetrics/fastcache v1.5.1
	github.com/alecthomas/kingpin v2.2.6+incompatible
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/araddon/dateparse v0.0.0-20190622164848-0fb0a474d195
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/cespare/xxhash v1.1.0
	github.com/cespare/xxhash/v2 v2.1.0 // indirect
	github.com/develar/errors v0.9.0
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/jmoiron/sqlx v1.2.0
	github.com/json-iterator/go v1.1.7
	github.com/kr/pretty v0.1.0 // indirect
	github.com/kshvakov/clickhouse v1.3.12-0.20191009193648-01e778c59781
	github.com/magiconair/properties v1.8.1
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/panjf2000/ants/v2 v2.2.2
	github.com/pkg/errors v0.8.1
	github.com/rs/cors v1.7.0
	github.com/tdewolff/minify/v2 v2.5.2
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/quicktemplate v1.3.1
	go.uber.org/atomic v1.4.0
	go.uber.org/multierr v1.2.0 // indirect
	go.uber.org/zap v1.10.0
	google.golang.org/appengine v1.6.5 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

replace github.com/bvinc/go-sqlite-lite => github.com/develar/go-sqlite-lite v0.6.2-0.20191004054224-45b503c7be7f
