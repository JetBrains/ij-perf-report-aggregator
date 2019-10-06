FROM golang:1.13-alpine3.10 AS builder

ENV GOPROXY=https://proxy.golang.org

WORKDIR /project

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN apk add --update gcc libc-dev

COPY cmd/server ./cmd/server
COPY pkg ./pkg

RUN go build -ldflags="-s -w -extldflags '-static'" -o /report-aggregator ./cmd/server

FROM scratch

ENV SERVER_PORT=80

COPY --from=builder /report-aggregator .
EXPOSE 80
ENTRYPOINT ["/report-aggregator", "--db", "/ij-perf-db/db.sqlite"]