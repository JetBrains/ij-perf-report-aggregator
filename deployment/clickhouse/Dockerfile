FROM golang:1.16 AS builder
WORKDIR /project/src
ADD https://raw.githubusercontent.com/nats-io/nats.go/7192159fcab301e10f8f9c6d21b57857d501c490/examples/nats-pub/main.go .
RUN GO111MODULE=off go get -d && GO111MODULE=off GOOS=linux go build -ldflags='-s -w' -o /tmp/nats-pub .

ARG CH_VERSION

RUN cd /tmp && curl -L https://repo.clickhouse.tech/tgz/stable/clickhouse-common-static-$CH_VERSION.tgz | tar xvz --strip-components=4 ./clickhouse-common-static-$CH_VERSION/usr/bin/clickhouse

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive

# clickhouse requires locales, tzdata
# official Debian and Ubuntu images automatically run apt-get clean, so explicit invocation is not required
RUN apt-get update -qq && apt-get install -qq --no-install-recommends locales tzdata ca-certificates && rm -rf /var/lib/apt /var/lib/dpkg /var/cache/debconf /var/cache/apt /var/log/dpkg.log /var/log/apt && \
  locale-gen en_US.UTF-8

COPY --from=builder /tmp/nats-pub /usr/bin/nats-pub
COPY --from=builder /tmp/clickhouse /usr/bin/clickhouse

ENV LANG en_US.UTF-8
ENV LANGUAGE en_US:en
ENV LC_ALL en_US.UTF-8
ENV TZ UTC

COPY entrypoint.sh /entrypoint.sh
COPY config.xml /etc/clickhouse-server/config.xml
COPY users.xml /etc/clickhouse-server/users.xml

EXPOSE 9000 8123 9009
VOLUME /var/lib/clickhouse

ENTRYPOINT ["/entrypoint.sh"]