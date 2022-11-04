.PHONY: build build-mac build-win build-linux

# https://github.com/valyala/quicktemplate
# go get -u github.com/valyala/quicktemplate/qtc

# kubectl rollout status deployment.apps/clickhouse

assets:
	qtc -dir pkg/server
	qtc -dir pkg/tc-properties

lint:
	golangci-lint run
	eslint . --ext .js,.vue,.ts

build-server:
	go build -tags -ldflags='-s -w' -o dist/server ./cmd/backend

build-tc-collector:
	go build -tags -ldflags='-s -w' -o dist/tc-collector ./cmd/tc-collector

build-transform:
	go build -tags -ldflags='-s -w' -o dist/transformer ./cmd/transform

update-deps:
	touch node_modules/go.mod
	go get -u ./...
	go mod tidy
	unlink node_modules/go.mod
	#cd dashboard && npm-check-updates -u && pnpm i && pnpm update

# docker run -it --rm --name ij-perf-clickhouse-server --ulimit nofile=262144:262144 -p 9000:9000 -p 8123:8123 --volume=$HOME/ij-perf-db/clickhouse:/var/lib/clickhouse:delegated yandex/clickhouse-server:20.3.4.10

# docker run -it --rm --link ij-perf-clickhouse-server:clickhouse-server yandex/clickhouse-client:20.3.4.10 --host clickhouse-server
# optimize table report
# optimize table report final

# select partition, name, active from system.parts where table = 'report'

# kubectl port-forward svc/clickhouse 2000:8123
# kubectl port-forward svc/clickhouse 9900:9000

# install clickhouse for macOS: brew install clickhouse

# clickhouse server -C ~/Documents/report-aggregator/deployment/ch-local-config.xml
# clickhouse client -h 127.0.0.1 -d ij

# nats-server
# /Volumes/data/nats-pub db.backup "test"

# MINIO_ROOT_USER=minio MINIO_ROOT_PASSWORD=minio123 minio server --console-address ":9001" --address "127.0.0.1:9002" ~/ij-perf-db/s3

# doppler run --project s3 --config prd -- rclone sync --size-only --fast-list --progress /Volumes/data/ij-perf-db/s3/ij-perf/data :s3,region=eu-west-1,provider=AWS,env_auth:eks-eu-west-1-idea-ij-perf-data-zznrqycixv/data

# doppler run --project s3 --config prd -- aws s3 sync /Volumes/data/ij-perf-db/s3/ij-perf/data s3://eks-eu-west-1-idea-ij-perf-data-zznrqycixv/data/ --region=eu-west-1
# doppler run --project s3 --config prd -- aws s3 cp s3://eks-eu-west-1-idea-ij-perf-data-zznrqycixv/data/ryb/ihvlrwxuhfddotlaamluhjpvqzdfw --region=eu-west-1 f


# BACKUP_NAME=$(date -u +%Y-%m-%dT%H-%M-%S)
# clickhouse-backup create $BACKUP_NAME
# REMOTE_STORAGE=s3 S3_ALLOW_MULTIPART_DOWNLOAD=true BACKUPS_TO_KEEP_REMOTE=8 doppler run --project s3 --config prd -- clickhouse-backup upload $BACKUP_NAME

# REMOTE_STORAGE=s3 S3_ALLOW_MULTIPART_DOWNLOAD=true BACKUPS_TO_KEEP_REMOTE=8 doppler run --project s3 --config prd -- clickhouse-backup restore_remote