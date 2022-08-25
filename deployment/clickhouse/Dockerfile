FROM golang:1.19.0 AS builder
WORKDIR /project/src
ADD entrypoint.go .
RUN GO111MODULE=off go get -d && GO111MODULE=off GOOS=linux go build -ldflags='-s -w' -trimpath -o /tmp/entrypoint .

RUN cd /tmp && curl -L https://github.com/ClickHouse/ClickHouse/releases/download/v22.8.2.11-lts/clickhouse-common-static-22.8.2.11-amd64.tgz | tar xvz --strip-components=3

RUN mkdir /var/lib/clickhouse

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=builder /tmp/entrypoint /usr/bin/entrypoint
COPY --from=builder /tmp/clickhouse /usr/bin/clickhouse

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV TZ UTC

COPY config.xml /etc/clickhouse-server/config.xml
COPY users.xml /etc/clickhouse-server/users.xml
COPY --from=builder --chown=nonroot:nonroot /var/lib/clickhouse /var/lib/clickhouse

EXPOSE 9000 8123 9009
VOLUME /var/lib/clickhouse

ENTRYPOINT ["/usr/bin/entrypoint"]