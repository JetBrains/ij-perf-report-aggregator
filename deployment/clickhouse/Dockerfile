FROM buildpack-deps:curl AS builder

RUN cd /tmp && curl -L https://github.com/ClickHouse/ClickHouse/releases/download/v23.5.2.7-stable/clickhouse-common-static-23.5.2.7-amd64.tgz | tar xvz --strip-components=3
RUN mkdir /var/lib/clickhouse

FROM gcr.io/distroless/base-debian11:nonroot

COPY --from=builder /tmp/clickhouse /usr/bin/clickhouse

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV TZ UTC

COPY users.xml /etc/clickhouse-server/users.xml
COPY --from=builder --chown=nonroot:nonroot /var/lib/clickhouse /var/lib/clickhouse

EXPOSE 9000 8123 9009
VOLUME /var/lib/clickhouse

ENTRYPOINT ["/usr/bin/entrypoint"]