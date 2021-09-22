.PHONY: build build-mac build-win build-linux

# https://github.com/bvinc/go-sqlite-lite/issues/10#issuecomment-498539630

# https://github.com/valyala/quicktemplate
# go get -u github.com/valyala/quicktemplate/qtc

# kubectl rollout status deployment.apps/clickhouse

assets:
	qtc -dir pkg/server
	qtc -dir pkg/tc-properties

lint:
	golangci-lint run

build-server:
	go build -tags -ldflags='-s -w' -o dist/server ./cmd/backend

build-monitor:
	go build -tags -ldflags='-s -w' -o dist/monitor ./cmd/monitor

build-tc-collector:
	go build -tags -ldflags='-s -w' -o dist/tc-collector ./cmd/tc-collector

build-transform:
	go build -tags -ldflags='-s -w' -o dist/transformer ./cmd/transform

update-deps:
	touch node_modules/go.mod
	go get -d -u ./...
	go mod tidy
	unlink node_modules/go.mod
	#cd dashboard && npm-check-updates -u && pnpm i && pnpm update

# docker run -it --rm --name ij-perf-clickhouse-server --ulimit nofile=262144:262144 -p 9000:9000 -p 8123:8123 --volume=$HOME/ij-perf-db/clickhouse:/var/lib/clickhouse:delegated yandex/clickhouse-server:20.3.4.10

# docker run -it --rm --link ij-perf-clickhouse-server:clickhouse-server yandex/clickhouse-client:20.3.4.10 --host clickhouse-server
# optimize table report

# select partition, name, active from system.parts where table = 'report'

# kubectl port-forward svc/clickhouse 2000:8123
# kubectl port-forward svc/clickhouse 9900:9000

# download clickhouse for macOS: curl -O 'https://clickhouse-builds.s3.yandex.net/14587/83b4a620c80e01fedfc9605e522e3c40eebf9784/clickhouse_build_check/clang-10-darwin_relwithdebuginfo_none_bundled_unsplitted_disable_False_binary/clickhouse' && chmod a+x ./clickhouse
# see https://github.com/ClickHouse/ClickHouse/issues/235#issuecomment-548781417

# clickhouse server -C ~/Documents/report-aggregator/deployment/ch-local-config.xml
# clickhouse client -h 127.0.0.1 -d ij

# nats-server
# /Volumes/data/nats-pub db.backup "test"

# minio server minio --console-address ":9001" --address "127.0.0.1:9002"