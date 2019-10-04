FROM golang:1.13-alpine3.10 AS builder

ENV GOPROXY=https://proxy.golang.org
ENV GO111MODULE=on

WORKDIR /project

COPY go.mod .
COPY go.sum .
RUN go mod download

RUN apk update && apk add gcc libc-dev make

COPY main.go .
COPY pkg ./pkg

RUN go build -ldflags='-s -w' -o /report-aggregator ./

FROM alpine:3.10

ENV SERVER_PORT=80

COPY db/ij-perf-db ./ij-perf-db
COPY --from=builder /report-aggregator .
CMD ["/report-aggregator", "serve", "--db", "/ij-perf-db/db.sqlite"]