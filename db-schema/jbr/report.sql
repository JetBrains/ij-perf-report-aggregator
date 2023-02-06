drop table jbr.report;
create table jbr.report
(
  `machine`               LowCardinality(String) CODEC (ZSTD(20)),
  `generated_time`        DateTime CODEC (Delta(4), ZSTD(20)),
  `project`               LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_id`           UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `branch`                LowCardinality(String) CODEC (ZSTD(20)),
  `build_number`          LowCardinality(String) CODEC (ZSTD(20)),

  `measures.name`         Array(LowCardinality(String)) CODEC (ZSTD(20)),
  `measures.value`        Array(Float64) CODEC (Gorilla, ZSTD(20)),
  `measures.type`         Array(LowCardinality(String)) CODEC (ZSTD(20)),

  `triggeredBy`           LowCardinality(String) CODEC (ZSTD(20))
)
  engine = MergeTree
    partition by (toYear(generated_time))
    order by (machine, branch, project,  generated_time)