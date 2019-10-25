.PHONY: build build-mac build-win build-linux

# https://github.com/bvinc/go-sqlite-lite/issues/10#issuecomment-498539630

# https://github.com/valyala/quicktemplate
# go get -u github.com/valyala/quicktemplate/qtc
# go get -u github.com/go-bindata/go-bindata/...

assets:
	qtc -dir pkg/server
	qtc -dir pkg/teamcity

build: lint
	go mod tidy
	make build-mac
	make build-linux
	make build-windows

build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -tags "clz4 sqlite_json sqlite_stat4 sqlite_foreign_keys" -ldflags='-s -w' -o dist/mac/report-aggregator ./cmd/report-aggregator
	XZ_OPT=-9 tar -cJf dist/mac-report-aggregator.tar.xz dist/mac/report-aggregator

build-linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -tags "clz4 sqlite_json sqlite_stat4 sqlite_foreign_keys" -ldflags='-s -w' -o dist/linux/report-aggregator ./cmd/report-aggregator
	XZ_OPT=-9 tar -cJf dist/linux-report-aggregator.tar.xz dist/linux/report-aggregator

build-windows:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-w64-mingw32-gcc CXX=/usr/local/bin/x86_64-w64-mingw32-g++ go build -tags "clz4 sqlite_json sqlite_stat4 sqlite_foreign_keys" -ldflags='-s -w' -o dist/windows/report-aggregator.exe ./cmd/report-aggregator

lint:
	golangci-lint run

build-server:
	go build -tags "clz4 sqlite_json sqlite_stat4 sqlite_foreign_keys" -ldflags='-s -w' -o dist/server ./cmd/server

update-deps:
	GOPROXY=https://proxy.golang.org go get -u ./cmd/report-aggregator
	GOPROXY=https://proxy.golang.org go get -u ./cmd/server
	go mod tidy

# https://medium.com/@valyala/promql-tutorial-for-beginners-9ab455142085

#   -influxSkipSingleField {measurement}
#    	Uses {measurement} instead of `{measurement}{separator}{field_name}` for metic name if Influx line contains only a single field

# docker run -it --rm --name ij-perf-clickhouse-server --ulimit nofile=262144:262144 -p 9000:9000 -p 8123:8123 --volume=$HOME/ij-perf-db/clickhouse:/var/lib/clickhouse:delegated yandex/clickhouse-server:19.15.2.2

# docker run -it --rm --link ij-perf-clickhouse-server:clickhouse-server yandex/clickhouse-client --host clickhouse-server
# optimize table report

# SELECT partition, name, active FROM system.parts WHERE table = 'report'

# kubectl port-forward svc/clickhouse 2000:8123