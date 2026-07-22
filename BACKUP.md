# ClickHouse Backup & Restore

Backups are made by [Altinity clickhouse-backup](https://github.com/Altinity/clickhouse-backup) **v2.8.0**,
invoked as a binary (`/usr/bin/clickhouse-backup` in images, `~/clickhouse-backup` locally) by thin Go wrappers.
Backups are **self-contained**: data on the s3-tiered disk is server-side copied into the backup, so a backup
does not depend on the live database or its live S3 objects.

## How it works

```
tc-collector (cron, every 20m)
  └─ publishes NATS `db.backup` after each collection
       └─ clickhouse-backup sidecar (same pod as clickhouse, shares /var/lib/clickhouse)
            ├─ rate-limits to one backup per 24h (in-memory clock, resets on pod restart)
            └─ exec `clickhouse-backup create_remote --delete-source`
                 ├─ FREEZE tables → hardlink shadow copy → tar → upload
                 ├─ s3-disk parts: server-side S3 copy into object_disks/<name>/
                 └─ retention: keeps last 30 remote backups (BackupsToKeepRemote)
```

Restore happens **on every pod start**: the ClickHouse data volume is ephemeral by design. An init container
runs the same image with `RESTORE_DB=true` → wipes the data dir, starts a temporary server, picks the latest
non-broken remote backup, `restore_remote`, deletes the local download, exits. The `collector_state` table is
restored with the data, so tc-collector automatically re-collects anything newer than the backup.

Configuration lives in `pkg/clickhouse-backup/clickhouse-backup.go` (`backupEnv()`), passed to the binary as
env vars. S3 credentials: doppler project `s3`/config `prd` locally, `/etc/s3/*` files (secret
`ij-perf-data-s3-rw`) in k8s.

## Where the backups are

Bucket `eks-eu-west-1-idea-ij-perf-data-zznrqycixv` (eu-west-1):

| Prefix                          | Content                                                                                 |
| ------------------------------- | --------------------------------------------------------------------------------------- |
| `backup/<YYYY-MM-DDTHH-MM-SS>/` | backups (metadata.json + metadata/ + shadow/ tars), ~12 GB each                         |
| `object_disks/<backup-name>/`   | server-side copies of s3-disk objects for that backup, ~0.3 GB (grows with the s3 tier) |
| `data/`                         | the **live** s3-tiered disk of prod ClickHouse — not a backup                           |
| `<timestamps at bucket root>`   | legacy v1 backup (pre-2026-07-21, pointer-based, no longer purged) — cleanup pending    |

List backups (any command needs a reachable ClickHouse on 127.0.0.1:9000 — the binary connects to it even
for remote-only operations):

```sh
doppler run --project s3 --config prd -- env REMOTE_STORAGE=s3 S3_PATH=backup S3_OBJECT_DISK_PATH=object_disks ~/clickhouse-backup list remote
```

## Restore locally (drill)

Run this periodically — a backup you have not restored is not a backup.

1. Get the binary (darwin-arm64) if missing:
   `curl -sL https://github.com/Altinity/clickhouse-backup/releases/download/v2.8.0/clickhouse-backup-darwin-arm64.tar.gz | tar xz --strip-components=3 -C ~ build/darwin/arm64/clickhouse-backup`
2. Start MinIO (backs the local `s3` disk defined in `deployment/ch-local/config.xml`):
   `MINIO_ROOT_USER=minio MINIO_ROOT_PASSWORD=minio123 minio server --console-address ":9001" --address "127.0.0.1:9002" ~/ij-perf-db/s3`
3. Make sure nothing listens on 127.0.0.1:9000 (the drill refuses to wipe under a running server).
4. Run the drill — **wipes `/Volumes/data/ij-perf-db/clickhouse`**, downloads ~12 GB:
   `RESTORE_DB=true go run ./cmd/clickhouse`

   Local knobs (env): `CLICKHOUSE_BIN` (server binary, default `~/clickhouse`), `CLICKHOUSE_CONFIG`
   (must be named `config.xml`, default `deployment/ch-local/config.xml`), `CLICKHOUSE_DATA_PATH`
   (the dir that gets wiped and restored into), `CLICKHOUSE_PORT` (numeric only — non-numeric values
   fall back to 9000), `RESTORE_BACKUP_NAME` (pin a specific backup instead of latest). Together these
   allow a second instance side by side (e.g. `deployment/ch-local-candidate` on ports 9010/8124 for
   version comparison — see the `clickhouse-update` skill). Two *concurrent* restores of the same backup also need distinct `TMPDIR`s:
   clickhouse-backup's pid file is keyed by backup name. If a restore fails, the temporary server may
   be left running — kill it before retrying.
5. Start ClickHouse from the restored data and verify:
   ```sh
   ~/clickhouse server --config-file=deployment/ch-local/config.xml
   ~/clickhouse client -q "SELECT sum(rows) FROM system.parts WHERE active"
   ~/clickhouse client -q "SELECT count() FROM ij.report WHERE generated_time < '2021-01-01'"  # exercises the s3 tier
   ```

The second query must return non-zero: old `ij.report` partitions live on the s3 disk, so it proves the
object-disk restore path, not just the plain data path.

## Restore in prod

Roll the deployment — restore-on-boot does the rest:

```sh
kubectl rollout restart -n idea deployment/clickhouse
kubectl get pods -n idea -l app=clickhouse                       # find the new (surge) pod
kubectl logs -n idea <new-pod> -c restore-clickhouse-backup -f   # expect: DB is restored (backup=...)
```

The deployment uses the default RollingUpdate strategy (with 1 replica: maxSurge=1, maxUnavailable=0) and a
per-pod ephemeral data volume, so the new pod downloads and restores the backup while the old pod keeps
serving; the old one is killed only when the new one is Ready. Do **not** `kubectl delete pod`: the pod
drops out of the Service the moment it starts terminating, and the replacement cannot serve until the whole
~12 GB restore finishes — full downtime for the duration. If the surge pod stays Pending (no node capacity
for two ClickHouse pods at once), nothing breaks — the old pod keeps serving and the restore just waits.

Up to 24h of data since the last backup is re-collected automatically by tc-collector after restore.
To restore a _specific_ (non-latest) backup, set `RESTORE_BACKUP_NAME=<backup-name>` on the
`restore-clickhouse-backup` init container and roll the deployment (`kubectl edit deployment` or a chart
change — `kubectl set env` does not reach init containers). This is also the rollback path after a
ClickHouse **version upgrade**: once the new version has written its first backup, plain restore-latest
on the old version would try to restore a newer-version backup — pin a pre-upgrade backup instead.

## Deploying changes to backup code (the pin dance)

The helm deploy (TeamCity `Build & Push & Deploy`, auto-triggered by pushes to either repo) resolves images
with ko but then applies `jb/values.yaml` (separate repo, cloned into `jb/`) **last** — its hardcoded
`images:` digest pins always win:

1. Push the main repo. The TeamCity build publishes fresh images to ghcr (`Published ghcr.io/...@sha256:...`
   lines in the build log) — but deploys the **old** pins. This is expected.
2. Put the new digests into `images:` in `jb/values.yaml`, commit, push the `jb/` repo.
3. The build triggered by that push performs the real rollout.

## Gotchas (all bitten in production)

- **Service links**: k8s injects `CLICKHOUSE_PORT=tcp://<ip>:9000` per Service; the binary's envconfig needs
  a numeric port. Fixed by `enableServiceLinks: false` in the pod spec and a pinned `CLICKHOUSE_PORT` in
  `backupEnv()` — keep both.
- **Memory**: the binary derives upload concurrency from the _node_ CPU count and OOMs small containers.
  Fixed by `UPLOAD_CONCURRENCY=2`, a 2Gi limit, and `GOMEMLIMIT` — keep all three.
- **`s3->path` must be non-empty** when object disks are used, and disjoint from `object_disk_path`. Legacy
  root-path backups can be listed with explicit `S3_PATH=""`, but v2 **cannot restore** v1 pointer-based
  backups unless `object_disks/<name>/s3/` is pre-seeded with a copy of `data/` (see git history of the
  2026-07-21 migration).
- **Broken backups** (interrupted uploads, no `metadata.json`) are ignored by restore and by `list ... latest`
  (zero creation date). Clean them with `clean_remote_broken`.
- ⚠️ **Never run `clean_remote_broken` or `delete remote` with `S3_PATH=""`** (root scope): the binary then
  sees every root prefix — including `data/`, the live s3-tiered disk — as a "broken backup" and will happily
  delete it. Clean up root-level leftovers only with targeted `rclone purge`/`aws s3 rm` on explicit prefixes.
- **Cross-endpoint restores need streaming**: restoring a prod backup locally copies object-disk data from
  AWS into MinIO — different endpoints, so v2's server-side `CopyObject` fails with `NoSuchBucket`.
  `backupEnv()` defaults `ALLOW_OBJECT_DISK_STREAMING=true` for local runs (prod keeps the faster
  server-side copy). When invoking the binary manually against prod from a laptop, set it yourself.
- The binary **hangs without a reachable ClickHouse server**, even for remote-only commands like
  `list remote` or `delete remote`.
- Local drill config must be named `config.xml` (`deployment/ch-local/config.xml`): the binary reads
  object-disk credentials from the preprocessed config, whose name follows the main config file's.
