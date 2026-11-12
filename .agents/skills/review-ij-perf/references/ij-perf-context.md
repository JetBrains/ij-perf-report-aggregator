# IJ Perf context and risk map

Verified on 2026-11-12 against `JetBrains/ij-perf-report-aggregator` master at `107f63c465c7c25aabcfdf6cecae1ff4669873d5`. Refresh from the reviewed revision. These are review questions, not allegations of existing defects.

Contents: navigation; measurement meaning; collection/recovery; detector semantics; SQL/contracts/caches; analysis state/evidence; dashboard behavior; schema/operations; targeted validation.

## Navigation

| Area | Starting paths | Relationship to inspect |
| --- | --- | --- |
| Collection | `cmd/tc-collector/`, `pkg/analyzer/`, `pkg/sql-util/` | TeamCity artifacts → reports → inserts → checkpoint |
| Storage/query | `db-schema/`, `pkg/data-query/`, `pkg/server/clickhouse.go`, `pkg/server/comparison.go` | Schema/query DSL → column-oriented results → dashboard/detector |
| Detection | `pkg/degradation-detector/`, its `setting/` and `statistic/`, `cmd/degradation-analyzer/` | Ordered observations → analysis → accidents/notifications |
| Metadata/actions | `pkg/server/meta/`, `db-schema/meta/`, `openapi/meta.yaml` | PostgreSQL → analysis runs, ownership, TeamCity, YouTrack |
| API/cache/MCP | `pkg/server/server.go`, `responseCache.go`, `compression.go`, `pkg/server/mcp/` | HTTP/auth → SQL → cached/compressed response |
| Dashboard | `dashboard/new-dashboard/src/components/`, `src/configurators/`, `dashboard/app/` | Configurator state → queries → series → point actions |
| Restore/rollout | `BACKUP.md`, `cmd/clickhouse/`, `pkg/clickhouse-backup/`, `deployment/`, `.ko.yaml` | Images → restore-on-start → backups/readiness |

Read README, CLAUDE.md, current manifests, and relevant tests. Some deployed settings, including collected TeamCity configurations, live outside this repository; the README links the deployment repository. Do not infer production policy solely from example Helm values.

## Measurement meaning

- Trace database/table, project/test, metric, branch, machine/group, mode, build configuration, and personal/bisect status. Check that filters, groups, joins, caches, and links preserve the dimensions relevant to the feature.
- Do not impose all dimensions on every query: aggregation and cross-chart matching can deliberately combine them. Establish which populations the feature promises to compare.
- Distinguish TeamCity test build ID, installer build ID, IDE build number, and VCS revision. `PerformanceSettings.query()` uses IDE build numbers for `perfint` and TeamCity IDs elsewhere; Fleet also has its own build encoding. Numeric-looking strings are not interchangeable IDs.
- Preserve aligned metric names/values/types, timestamps/build IDs, and column-oriented response order. Inspect `pkg/data-query/resultJsonEncoder.go`, `fetchMetricValues.go`, and dashboard consumers for shape changes.
- Verify seconds/milliseconds, bytes/scaled units, integers/fractions, direction, overflow, and missing/null/zero distinctions. Missing data must not quietly become zero or "no regression."
- For comparisons, inspect weighting, repeated runs, grouping, percent-change denominator, zero baseline, small samples, variation displays, and whether higher/lower means worse. Means, medians, and pooled statistics are not interchangeable.
- Check chronology and ties; distinguish generated time, build time, run date, and analysis creation time. A rolling N-day window can span N+1 calendar dates.

## Collection and recovery

At the verified revision, collection reads paginated TeamCity builds, processes older chunks before newer ones, waits for analysis/inserts, then persists `collector_state`.

- On download, decode, insert, or flush failure, can a checkpoint pass unpersisted data? Follow `tcCollector.go`, `reportLoader.go`, `ReportAnalyzer.go`, and batch managers together.
- On retry, restart, restore, or overlapping collection, are existing reports recognized without skipping new artifacts or double-counting? Establish actual report identity and table engine; do not assume exactly-once delivery.
- Inspect timestamp boundaries, simultaneous finishes, pagination, partial artifacts, compressed-data limits, and cancellation.
- Follow batch dependencies, buffer ownership, final flush, first-error propagation, channel closure, and worker completion. A successful enqueue is not a durable insert.
- Preserve build/configuration/artifact context needed to diagnose an actual failure.

## Detector semantics

- Trace settings through query selection, extraction, detection, inference, merging, and posting. Compare chart data with the detector's intended population.
- Inspect minimum segments, short/constant/noisy series, missing values, zero baseline, effect-size/percentage/absolute gates, and threshold equality.
- Change analysis and explicit threshold analysis are separate branches at the verified revision. Verify direction and defaults from code; do not hard-code numeric thresholds into review policy.
- Confirm the reported build/timestamp identifies the right boundary. Keep values/builds/timestamps aligned.
- Preserve ignore-before rules, inferred/manual kinds, routing, and deduplication. Consider missed regressions, false regressions, and repeated alerts separately.
- Check changed mathematics with small explicit series whose expected result can be independently established.

## SQL, contracts, and caches

- Read actual engines, sorting/partition keys, indexes, and execution paths. PostgreSQL and ClickHouse differ in update and consistency semantics.
- For `ReplacingMergeTree`, background merges provide eventual deduplication. Determine if this query needs query-time deduplication; do not prescribe `FINAL` everywhere.
- Check bounded time ranges, pruning, join cardinality, repeated queries, memory, truncation, and ordering. Ground performance concerns in concrete queries, cardinality, benchmarks, or measured plans.
- Review values, identifiers, and deliberate SQL expressions separately. The query DSL exposes SQL fragments; string construction alone does not prove introduced injection. Trace validation, reachability, database privileges, accessible tables, and what the PR changes. Read-only access does not prevent every disclosure or resource-exhaustion problem.
- Trace Go structs, JSON, TypeScript consumers, OpenAPI, null/empty behavior, HTTP errors, and old records across contract changes.
- Inspect `responseCache.go`, compression, and NATS `server.clearCache` together. Verify relevant query/format/auth dimensions, invalidation, ETags, compressed/uncompressed representations, and pooled-buffer lifetime.
- For MCP tools, verify limits and reported coverage. Truncated results must not appear complete.

## Analysis state and evidence

At the verified revision, `analyses` is PostgreSQL metadata. Starting analysis inserts a row, starts a TeamCity build, then stores its ID. States include `in_progress`, `success`, and `failed`.

- Trace database write → external build/issue creation → persistence of external ID → response. Check retries, concurrency, orphaned runs, duplicate side effects, terminal state, and recovery.
- Inspect PATCH semantics: absent/null/empty/zero; callback authentication/ownership; repeated/out-of-order updates; old rows after schema changes.
- Preserve project/metric/build/commit-range and attachment provenance. Distinguish run date from analysis creation date.
- For similar degradations, inspect `llmAnalysisMatches.go`, `pkg/server/buildCommits.go`, and UI. Full SHAs, installer ranges, successful analyses, self-match exclusion, related-metric evidence, bounded windows, and machine/branch resolution affect meaning.
- A shared commit suggests a possible explanation; it does not prove causality. Preserve uncertainty and range size when communicating evidence strength.
- Preserve usable partial results where intended. A failed lookup must not masquerade as a confirmed empty result.
- Trace YouTrack/attachment failures through `llmAnalysisYoutrack.go` and shared helpers. An issue already created externally must remain discoverable if its local link write fails.

## Dashboard behavior

- Follow configurators → `DataQueryExecutor.ts` → charts/sidebar and URL state. Check merging/splitting, series identity, and shared mutable state.
- On rapid filter/point/server changes, prove old responses cannot overwrite the current selection. Accept cancellation, generation guards, or equivalent mechanisms; a missing AbortController alone is not a bug.
- Check stale rejection/finally paths, unmount cleanup, subscriptions, timers, and chart disposal. Include A→B→A selection when request identity matters.
- Verify links retain required database, table, project, metric, branch, machine, time range, build, and analysis ID. Check cross-table results and unresolved machines before creating exact-point links.
- Read raw values and formatters before alleging wrong labels, percentages, uncertainty bands, or missing-data displays. Rendering issues matter when they mislead investigation or break a supported action.
- Trace trust boundaries for rendered logs/LLM text, URLs, attachments, requests, and identity headers. Inspect ingress/auth topology before alleging a bypass; an absent handler check alone is insufficient evidence.

## Schema and operations

- Separate bootstrap/example DDL from executed migrations. A `DROP TABLE` in a static schema file is not proof of production data loss. Find the caller and deployment plan.
- For actual migrations, inspect old rows, backfill, uniqueness, concurrent versions, rollback, and locks. Updating creation DDL does not establish how an existing database gets a new field.
- The verified backup design restores ClickHouse on every pod start using ephemeral local data and self-contained remote backups. Collector state restores with data. Evaluate compatibility/recovery, not just startup.
- For ClickHouse upgrades, read BACKUP.md, the repository's ClickHouse-update reference, and benchmark evaluations. Check a shared backup baseline, result correctness, representative query performance, memory, image layers/digests, readiness, and compatible rollback backup.
- Check retention/storage prefixes against live data and backup data. Do not run restore, deletion, rollout, or backup commands during review.
- Report new queues, goroutines, retries, metric labels, or pools only with a reachable leak, exhaustion, deadlock, lost-work, or unbounded-cardinality scenario.

## Targeted validation

Read current `.github/workflows/lint-and-test.yml`, `package.json`, `go.mod`, and `CLAUDE.md` before choosing commands/runtime versions.

| Change | Useful checks when safe and available |
| --- | --- |
| Go logic | `go test ./pkg/<affected-package>` or `./cmd/<affected-command>`; `-race` for changed concurrency |
| Backend cross-package | `go test ./cmd/... ./pkg/...` |
| Dashboard logic | `pnpm exec vitest run <relevant-test-path>` |
| Vue/type contract | `pnpm vue-tsc --noEmit -p dashboard/new-dashboard`; `pnpm build` when relevant |
| SQL/schema | Representative queries/rows in a disposable database at the applicable version |
| Detector math | Boundary/representative series, unchanged controls, expected build attribution |
| Cache/concurrency | Vary inputs/schedule order and assert response identity, not implementation details |

CI creates `cmd/frontend/resources/make-typecheck-happy.txt` for embedded resources when needed. Inspect that setup before reproducing it. Current formatting/linting uses oxfmt/oxlint; do not introduce Prettier or reformat during review. Tests live in `pkg/`, `cmd/`, and `dashboard/new-dashboard/tests/`. An environment setup failure is not proof the PR is broken.
