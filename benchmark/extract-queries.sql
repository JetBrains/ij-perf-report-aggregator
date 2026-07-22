-- Extract the benchmark corpus from the LOCAL baseline server's query_log.
--
-- Capture flow: start the baseline server (`go run ./cmd/clickhouse` with the restored
-- data), run the backend (it defaults to 127.0.0.1:9000) and the frontend, then open
-- the dashboards to benchmark in the browser — every query they produce lands in the
-- local system.query_log. Then:
--
--   ~/clickhouse-26.6.2.81 client --port 9000 < benchmark/extract-queries.sql > benchmark/corpus.jsonl
--
-- One row per distinct (query shape, database): the backend sets a default database per
-- data source and sends unqualified table names, so the same SQL text against ij and
-- perfint are different queries — replaying needs both the text and the database.
-- The slowest concrete instance is kept as the replayable text.
SELECT
  concat(toString(normalized_query_hash), '_', current_database) AS id,
  current_database AS db,
  count() AS runs,
  round(median(query_duration_ms)) AS med_ms,
  max(query_duration_ms) AS max_ms,
  round(max(memory_usage) / 1048576) AS max_mem_mb,
  argMax(query, query_duration_ms) AS query
FROM system.query_log
WHERE type = 'QueryFinish'
  AND query_kind = 'Select'
  AND is_initial_query
  AND NOT has(databases, 'system')
  -- keep only application traffic (backend uses ch-go / clickhouse-go drivers);
  -- interactive clickhouse-client sessions (incl. benchmark replays) are not workload
  AND client_name NOT LIKE 'ClickHouse client%'
  AND log_comment NOT LIKE 'bench:%'
GROUP BY normalized_query_hash, current_database
ORDER BY sum(query_duration_ms) DESC
FORMAT JSONEachRow
