FROM golang:1.13-alpine3.10 AS builder

ENV GOPROXY=https://proxy.golang.org
ENV GO111MODULE=on

WORKDIR /project

COPY pkg/go.mod ./pkg/
COPY pkg/go.sum ./pkg/

COPY cmd/clickhouse-backup/go.mod ./cmd/app/
COPY cmd/clickhouse-backup/go.sum ./cmd/app/
RUN cd ./cmd/app && go mod download

COPY cmd/clickhouse-backup ./cmd/app
COPY pkg/util ./pkg/util
RUN cd ./cmd/app && go build -ldflags="-s -w" -o /app .

# from scratch doesn't work
FROM alpine:3.10

RUN apk add --no-cache --update ca-certificates tzdata

COPY --from=builder /app .
ENTRYPOINT ["/app"]