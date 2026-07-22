FROM buildpack-deps:curl AS builder

RUN cd /tmp && curl -L https://github.com/ClickHouse/ClickHouse/releases/download/v25.8.28.1-lts/clickhouse-common-static-25.8.28.1-amd64.tgz | tar xvz --strip-components=1
RUN cd /tmp && curl -L https://github.com/Altinity/clickhouse-backup/releases/download/v2.8.0/clickhouse-backup-linux-amd64.tar.gz | tar xz
RUN mkdir /var/lib/clickhouse

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=builder /tmp/usr/bin/clickhouse /usr/bin/clickhouse
COPY --from=builder /tmp/build/linux/amd64/clickhouse-backup /usr/bin/clickhouse-backup

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV TZ UTC

COPY users.xml /etc/clickhouse-server/users.xml
COPY --from=builder --chown=nonroot:nonroot /var/lib/clickhouse /var/lib/clickhouse

EXPOSE 9000 8123 9009
VOLUME /var/lib/clickhouse

ENTRYPOINT ["/usr/bin/entrypoint"]