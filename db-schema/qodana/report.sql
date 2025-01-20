create table report
(
  `machine`               LowCardinality(String) CODEC (ZSTD(20)),
  `generated_time`        DateTime CODEC (Delta(4), ZSTD(20)),
  `project`               LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_id`           UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `branch`                LowCardinality(String) CODEC (ZSTD(20)),

  `measures.name`         Array(LowCardinality(String)) CODEC (ZSTD(20)),
  `measures.value`        Array(Int32) CODEC (ZSTD(20)),
  `measures.type`         Array(LowCardinality(String)) CODEC (ZSTD(20)),

  `triggeredBy`           LowCardinality(String) CODEC (ZSTD(20)),
  `mode`                  LowCardinality(String) CODEC (ZSTD(20))
)
  engine = MergeTree
--     partitioning by month gives up to 3x performance boost
    partition by (toYYYYMM(generated_time))
--     we need low cardinality indices first
    order by (project, branch, machine, generated_time)