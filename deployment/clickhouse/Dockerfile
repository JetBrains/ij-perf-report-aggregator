FROM golang:1.19.3 AS builder
WORKDIR /project/src
ADD config.xml .
ADD entrypoint.go .
RUN GO111MODULE=off go get -d && GO111MODULE=off GOOS=linux go build -ldflags='-s -w' -trimpath -o /tmp/entrypoint .

RUN cd /tmp && curl -L https://github.com/ClickHouse/ClickHouse/releases/download/v22.10.2.11-stable/clickhouse-common-static-22.10.2.11-amd64.tgz | tar xvz --strip-components=3
RUN cd /tmp && curl -L https://github.com/AlexAkulov/clickhouse-backup/releases/download/v2.1.2/clickhouse-backup-linux-amd64.tar.gz | tar xvz --strip-components=3

RUN mkdir /var/lib/clickhouse

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=builder /tmp/entrypoint /usr/bin/entrypoint
COPY --from=builder /tmp/clickhouse /usr/bin/clickhouse
COPY --from=builder /tmp/clickhouse-backup /usr/bin/clickhouse-backup

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV TZ UTC

COPY users.xml /etc/clickhouse-server/users.xml
COPY --from=builder --chown=nonroot:nonroot /var/lib/clickhouse /var/lib/clickhouse

EXPOSE 9000 8123 9009
VOLUME /var/lib/clickhouse

ENTRYPOINT ["/usr/bin/entrypoint"]