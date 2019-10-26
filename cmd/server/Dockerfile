FROM golang:1.13-alpine3.10 AS builder

ENV GOPROXY=https://proxy.golang.org
ENV GO111MODULE=on

WORKDIR /project

# tzdata is required for clickhouse client
RUN apk add --update --no-cache gcc libc-dev tzdata

COPY pkg/go.mod ./pkg/
COPY pkg/go.sum ./pkg/

COPY cmd/server/go.mod ./cmd/app/
COPY cmd/server/go.sum ./cmd/app/
RUN cd ./cmd/app && go mod download

COPY cmd/server ./cmd/server
COPY pkg ./pkg
RUN cd ./cmd/server && go build -tags clz4 -ldflags="-s -w -extldflags '-static'" -o /app .

FROM scratch

ENV SERVER_PORT=80
EXPOSE 80

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app .

ENTRYPOINT ["/app", "--db", "clickhouse:9000"]