#!/bin/bash
set -ex

(sleep 10 && /usr/bin/nats-pub -s nats server.clearCache clickhouse) &
exec /usr/bin/clickhouse server --config-file=/etc/clickhouse-server/config.xml