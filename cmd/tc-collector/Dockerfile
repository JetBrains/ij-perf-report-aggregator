FROM golang:1.13-alpine3.10 AS builder

ENV GOPROXY=https://proxy.golang.org
ENV GO111MODULE=on

WORKDIR /project

# tzdata is required for clickhouse client
RUN apk add --update --no-cache tzdata

COPY pkg/go.mod ./pkg/
COPY pkg/go.sum ./pkg/

COPY cmd/tc-collector/go.mod ./cmd/app/
COPY cmd/tc-collector/go.sum ./cmd/app/
RUN cd ./cmd/app && go mod download

COPY cmd/tc-collector ./cmd/app
COPY pkg ./pkg
RUN cd ./cmd/app && go build -ldflags="-s -w -extldflags '-static'" -o /app .

FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app .

ENTRYPOINT ["/app"]