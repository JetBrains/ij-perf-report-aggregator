create table measure
(
  `machine`               LowCardinality(String) CODEC (ZSTD(20)),
  `generated_time`        DateTime CODEC (Gorilla, ZSTD(20)),
  `project`               LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_id`           UInt32 CODEC (Gorilla, ZSTD(20)),
  `branch`                LowCardinality(String) CODEC (ZSTD(20)),
  `triggeredBy`           LowCardinality(String) CODEC (ZSTD(20)),

  `name` LowCardinality(String) CODEC(ZSTD(20)),
  `value` UInt64 CODEC(Gorilla, ZSTD(20))
)
  engine = MergeTree
    partition by toYYYYMM(generated_time)
    order by (branch, machine, project, generated_time)
    settings old_parts_lifetime = 10;