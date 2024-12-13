drop table if exists diogen.report;

create table diogen.report
(
  `machine`               LowCardinality(String) CODEC (ZSTD(20)),
  `generated_time`        DateTime CODEC (Delta(4), ZSTD(20)),
  `project`               LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_id`          UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `branch`                LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_type`         LowCardinality(String) CODEC (ZSTD(20)),

  `measures.name`         Array(LowCardinality(String)) CODEC (ZSTD(20)),
  `measures.value`        Array(Float64) CODEC (Gorilla, ZSTD(20)),
  `measures.type`         Array(LowCardinality(String)) CODEC (ZSTD(20)),

  `triggeredBy`           LowCardinality(String) CODEC (ZSTD(20)),
  `mode`                  LowCardinality(String) CODEC (ZSTD(20))
)
  engine = MergeTree
    partition by (toYYYYMM(generated_time))
    order by (branch, project, machine, generated_time);
