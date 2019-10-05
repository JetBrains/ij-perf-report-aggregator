.PHONY: build build-mac build-win build-linux

# https://github.com/bvinc/go-sqlite-lite/issues/10#issuecomment-498539630

# https://github.com/valyala/quicktemplate
# go get -u github.com/valyala/quicktemplate/qtc
# go get -u github.com/go-bindata/go-bindata/...

assets:
	qtc -dir pkg/server
	qtc -dir pkg/analyzer
	go-bindata -o ./pkg/analyzer/sqlScript.go -pkg analyzer -prefix ./pkg/analyzer/sql ./pkg/analyzer/sql

build: lint
	go mod tidy
	make build-mac
	make build-linux
	make build-windows

build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -ldflags='-s -w' -o dist/mac/report-aggregator ./
	XZ_OPT=-9 tar -cJf dist/mac-report-aggregator.tar.xz dist/mac/report-aggregator

build-linux:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ go build -ldflags='-s -w' -o dist/linux/report-aggregator ./
	XZ_OPT=-9 tar -cJf dist/linux-report-aggregator.tar.xz dist/linux/report-aggregator

build-windows:
	env GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=/usr/local/bin/x86_64-w64-mingw32-gcc CXX=/usr/local/bin/x86_64-w64-mingw32-g++ go build -ldflags='-s -w' -o dist/windows/report-aggregator.exe ./

lint:
	golangci-lint run

docker:
	docker build -t docker-registry.labs.intellij.net/idea/report-aggregator:latest .
	# docker run -it --rm -p 9044:80 docker-registry.labs.intellij.net/idea/report-aggregator:latest

	# kubectl config set-context --current --namespace=idea
	# kubectl apply -f deployment.yaml
	# kubectl apply -f ingress.yaml
	# zstd -19 --long /Volumes/data/ij-perf-db/db.sqlite  -o ff.zstd

update-deps:
	GOPROXY=https://proxy.golang.org go get -u
	go mod tidy

update-computed-metrics:
	go build -ldflags='-s -w' -o dist/report-aggregator
	./dist/report-aggregator update-computed-metrics --db /Volumes/data/ij-perf-db/db.sqlite

# https://medium.com/@valyala/promql-tutorial-for-beginners-9ab455142085

#   -influxSkipSingleField {measurement}
#    	Uses {measurement} instead of `{measurement}{separator}{field_name}` for metic name if Influx line contains only a single field

# docker run -it --rm -v ~/ij-perf-db/victoria-metrics-data:/victoria-metrics-data -p 8428:8428 victoriametrics/victoria-metrics:v1.28.0-beta5 -retentionPeriod 120