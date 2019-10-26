FROM golang:1.13-alpine3.10 AS builder

ENV GOPROXY=https://proxy.golang.org

WORKDIR /project

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/clickhouse-backup ./cmd/clickhouse-backup
COPY pkg ./pkg

RUN go build -ldflags="-s -w -extldflags '-static'" -o /app ./cmd/clickhouse-backup

FROM scratch

COPY --from=builder /app .

ENTRYPOINT ["/app"]