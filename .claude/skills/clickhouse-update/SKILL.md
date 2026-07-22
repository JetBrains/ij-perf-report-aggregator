---
name: clickhouse-update
description: Evaluate and ship a ClickHouse version update — restore-compat check, dashboard-query benchmark against prod version, evaluation record, prod rollout
argument-hint: [TargetVersion e.g. v26.3.4.10-lts]
---

Evaluate ClickHouse **$0** against the current prod version and, if it passes, ship it.
The evaluation proves two things on real data: the prod **backup restores** under the new
version (that is how prod upgrades — the data volume is ephemeral, see BACKUP.md), and the
**dashboard query workload** does not regress.

## 1. Resolve versions and binaries

- Current prod version: the release pinned in `deployment/clickhouse/Dockerfile`.
- Download missing binaries (macOS arm64) to `~/clickhouse-<version>` and verify:
  `curl -fL -o ~/clickhouse-<ver> https://github.com/ClickHouse/ClickHouse/releases/download/<tag>/clickhouse-macos-aarch64 && chmod +x ~/clickhouse-<ver> && ~/clickhouse-<ver> --version`
- On a major-version jump, check the pinned clickhouse-backup (same Dockerfile) against
  Altinity's release notes/test matrix — it may need a bump too.

## 2. Restore the same backup into both instances

Preconditions: MinIO running (BACKUP.md "Restore locally"), doppler available, ports 9000/9010 free.

1. Baseline (exact prod version, default paths, port 9000):
   `RESTORE_DB=true CLICKHOUSE_BIN=~/clickhouse-<prod-ver> go run ./cmd/clickhouse`
   Note the backup name from the `DB is restored (backup=...)` log line.
   Skip if the baseline dir already holds a restore of a recent backup made with the prod version —
   but then pin the candidate to that same backup.
2. Candidate (target version, its own config/data/port):
   `RESTORE_DB=true CLICKHOUSE_BIN=~/clickhouse-<target-ver> CLICKHOUSE_CONFIG=$PWD/deployment/ch-local-candidate/config.xml CLICKHOUSE_DATA_PATH=/Volumes/data/ij-perf-db/clickhouse-candidate CLICKHOUSE_PORT=9010 RESTORE_BACKUP_NAME=<same-backup> go run ./cmd/clickhouse`

**This candidate restore IS the backup-compat check**: a prod-version-written backup must restore
under the target server binary. If it fails, the upgrade is blocked regardless of performance.
Gotchas: concurrent restores of the same backup need distinct `TMPDIR`s (pid file is keyed by
backup name); a failed restore leaves the temporary server running — kill it before retrying.

## 3. Start and verify both servers

```sh
~/clickhouse-<prod-ver> server --config-file=deployment/ch-local/config.xml           # port 9000
~/clickhouse-<target-ver> server --config-file=deployment/ch-local-candidate/config.xml  # port 9010
```

On each port: `SELECT version()`, `SELECT sum(rows) FROM system.parts WHERE active` (must match
across servers within ~0.01%), and the s3-tier probe from BACKUP.md
(`SELECT count() FROM ij.report WHERE generated_time < '2021-01-01'` — non-zero and equal).

## 4. Benchmark

Corpus: reuse `benchmark/corpus.jsonl`. Refresh it only if dashboards changed materially:
browse the heavy dashboards against a local backend (defaults to 127.0.0.1:9000), then run
`benchmark/extract-queries.sql` per its header.

```sh
benchmark/run.sh benchmark/corpus.jsonl benchmark/results 16
```

Read `benchmark/results/report.tsv` (`ratio` > 1 = candidate slower):

- Any `rows_check MISMATCH` or failed query → **correctness blocker**, stop here.
- Weighted verdict (frequency-weighted, what users feel):
  `clickhouse local -q "SELECT round(sum(toInt64(c.runs)*r.old_ms)/1000,2) old_s, round(sum(toInt64(c.runs)*r.new_ms)/1000,2) new_s, round(new_s/old_s,3) ratio FROM file('benchmark/results/report.tsv', TSVWithNames) r JOIN file('benchmark/corpus.jsonl', JSONEachRow) c ON r.id=c.id"`
- Investigate queries ≥20 ms with ratio ≥ 1.2: replay on the candidate with
  `--compatibility '<prod-major>'` and `--enable_analyzer 0` (settings-default vs analyzer cause),
  and diff `read_rows`/`ProfileEvents['SelectedMarks']` between servers' query_logs (planner vs
  executor cause). Compare memory columns too — regressions can hide there (prod limit is 32Gi).
- Ignore ratios on sub-10 ms queries: 1 ms rounding dominates.

Pass bar: no correctness issues, weighted ratio ≤ ~1.05, no user-facing query class with a
sustained >20% median regression, no memory blowups.

## 5. Record the evaluation

Write `benchmark/evaluation/<version>.md` (verdict, environment incl. backup name, results,
diagnosis, mitigations tested) and copy `report.tsv` to
`benchmark/evaluation/<version>-report.tsv`. See `26.6.2.81.md` for the shape.

## 6. Ship (only if the evaluation passed)

Before the first push, note the latest prod backup name (`list remote`) — that is the rollback pin.

1. Update the release URL in `deployment/clickhouse/Dockerfile` (and clickhouse-backup if bumped);
   commit and push. CI builds the new **base** image (entrypoint `/usr/bin/entrypoint` — a
   placeholder; this image is not deployable on its own).
2. Find the new base digest in the ghcr `clickhouse` package (upload time after your push;
   `org.opencontainers.image.revision` label = your commit SHA). Update **both**
   `baseImageOverrides` entries in `.ko.yaml` to it; commit and push. CI now ko-builds the pod
   images on the 25.x base.
3. Find the fresh ko-built digests (upload time after the second push; entrypoint
   `/ko-app/clickhouse` in the `clickhouse` package, `/ko-app/clickhouse-backup` in the
   `clickhouse-backup` package — ko images inherit the base's revision label, so match on upload
   time + entrypoint, not revision). Pin both in `images:` in `jb/values.yaml`, commit, push the
   jb repo — this triggers the real rollout. Never pin the raw base image.
4. Watch the rollout: `kubectl logs -n idea <new-pod> -c restore-clickhouse-backup -f` must end
   with `DB is restored`; the old pod keeps serving until the new one answers `/ping`. Verify
   `SELECT version()` afterwards. Then confirm the next sidecar backup succeeds (`create_remote`
   in the clickhouse-backup container logs — the rate-limit clock resets on pod restart, so it
   lands within ~20–40 min of the rollout).
5. Rollback insurance: after the new version writes its first backup, rolling back requires
   `RESTORE_BACKUP_NAME` pinned to the pre-upgrade backup noted above (BACKUP.md "Restore in prod").

## Cleanup

Stop both local servers when done (`SYSTEM STOP MERGES` is not persistent — no state to undo).
The candidate data dir (`/Volumes/data/ij-perf-db/clickhouse-candidate`, ~12 GB) can be deleted
or left for the next evaluation.
