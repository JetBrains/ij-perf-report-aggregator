#!/usr/bin/env bash
# Compare two local ClickHouse servers on a corpus of real queries.
#
# Usage:  benchmark/run.sh benchmark/corpus.jsonl benchmark/results [iterations]
#
# corpus.jsonl comes from extract-queries.sql (fields: id, db, query, plus stats).
# Servers are expected on 127.0.0.1:$A_PORT (baseline) and 127.0.0.1:$B_PORT (candidate),
# both restored from the SAME backup. Iterations are interleaved A/B (order alternates
# per round) so environmental drift hits both sides equally.
#
# All measurements are taken server-side from each server's own system.query_log;
# runs are matched via log_comment (normalized_query_hash is not comparable across
# server versions). A query that fails on a server is recorded and skipped there —
# failures are findings, not errors.
set -euo pipefail

CORPUS=${1:?usage: run.sh corpus.jsonl results-dir [iterations]}
OUT=${2:?usage: run.sh corpus.jsonl results-dir [iterations]}
ITER=${3:-16}
CH=${CH:-$HOME/clickhouse-26.6.2.81} # used as plain client and clickhouse-local; works against both servers
A_PORT=${A_PORT:-9000}
B_PORT=${B_PORT:-9010}
MAX_S=${MAX_S:-120}

QDIR="$OUT/queries"
mkdir -p "$QDIR"
rm -f "$QDIR"/*.sql "$OUT"/failed-* "$OUT"/errors-*.log "$OUT"/stats-*.jsonl "$OUT"/manifest.tsv

# unpack corpus into one .sql file per query (debuggable, individually re-runnable)
# plus a manifest carrying the default database each query must run under
python3 - "$CORPUS" "$QDIR" "$OUT/manifest.tsv" <<'EOF'
import json, os, sys
with open(sys.argv[3], "w") as manifest:
    for line in open(sys.argv[1]):
        line = line.strip()
        if line:
            r = json.loads(line)
            with open(os.path.join(sys.argv[2], f"{r['id']}.sql"), "w") as f:
                f.write(r["query"])
            manifest.write(f"{r['id']}\t{r['db']}\n")
EOF

ids=()
dbs=()
while IFS=$'\t' read -r qid db; do
  ids+=("$qid")
  dbs+=("$db")
done <"$OUT/manifest.tsv"
echo "corpus: ${#ids[@]} queries, $ITER iterations per server, ports $A_PORT vs $B_PORT"

for port in "$A_PORT" "$B_PORT"; do
  "$CH" client --port "$port" -q "SYSTEM STOP MERGES" # freshly restored data must not shift under the benchmark
done

run_one() { # port, query-id, database; returns non-zero on query failure
  local port=$1 qid=$2 db=$3
  "$CH" client --port "$port" --database "$db" --max_execution_time "$MAX_S" \
    --log_comment "bench:$qid" --format Null \
    <"$QDIR/$qid.sql" 2>>"$OUT/errors-$port.log"
}

echo "warm-up pass (unmeasured: fills mark/page caches and the s3-disk cache)"
for idx in "${!ids[@]}"; do
  for port in "$A_PORT" "$B_PORT"; do
    if ! run_one "$port" "${ids[$idx]}" "${dbs[$idx]}"; then
      echo "FAIL ${ids[$idx]} on :$port (see errors-$port.log)"
      touch "$OUT/failed-${ids[$idx]}-$port"
    fi
  done
done

run_start=$(date '+%Y-%m-%d %H:%M:%S') # measured window starts here: warm-up stays out of the stats
echo "$run_start" >"$OUT/run-start.txt"

for i in $(seq 1 "$ITER"); do
  echo "round $i/$ITER"
  if ((i % 2)); then order=("$A_PORT" "$B_PORT"); else order=("$B_PORT" "$A_PORT"); fi
  for idx in "${!ids[@]}"; do
    for port in "${order[@]}"; do
      [[ -e "$OUT/failed-${ids[$idx]}-$port" ]] && continue
      run_one "$port" "${ids[$idx]}" "${dbs[$idx]}" || touch "$OUT/failed-${ids[$idx]}-$port"
    done
  done
done

for port in "$A_PORT" "$B_PORT"; do
  "$CH" client --port "$port" -q "SYSTEM FLUSH LOGS"
  "$CH" client --port "$port" -q "
    SELECT
      substring(log_comment, 7) AS id,
      count() AS n,
      round(median(query_duration_ms)) AS med_ms,
      round(quantile(0.9)(query_duration_ms)) AS p90_ms,
      round(max(memory_usage) / 1048576) AS mem_mb,
      median(result_rows) AS result_rows
    FROM system.query_log
    WHERE type = 'QueryFinish' AND log_comment LIKE 'bench:%' AND event_time >= '$run_start'
    GROUP BY id
    FORMAT JSONEachRow" >"$OUT/stats-$port.jsonl"
done

"$CH" local -q "
  SELECT
    coalesce(a.id, b.id) AS id,
    a.med_ms AS old_ms,
    b.med_ms AS new_ms,
    round(b.med_ms / greatest(a.med_ms, 1), 2) AS ratio,
    a.p90_ms AS old_p90,
    b.p90_ms AS new_p90,
    a.mem_mb AS old_mem_mb,
    b.mem_mb AS new_mem_mb,
    if(a.result_rows != b.result_rows, 'MISMATCH', 'ok') AS rows_check,
    if(a.n != b.n, 'FAILURES', '') AS run_check
  FROM file('$OUT/stats-$A_PORT.jsonl', JSONEachRow) AS a
  FULL OUTER JOIN file('$OUT/stats-$B_PORT.jsonl', JSONEachRow) AS b ON a.id = b.id
  ORDER BY ratio DESC
  FORMAT TSVWithNames" >"$OUT/report.tsv"

column -t "$OUT/report.tsv"
echo
echo "report: $OUT/report.tsv — ratio > 1 means the candidate (:$B_PORT) is slower"
ls "$OUT"/failed-* 2>/dev/null && echo "^ failed query/server pairs — inspect errors-*.log" || true
