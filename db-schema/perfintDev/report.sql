create table perfintDev.phpstorm
(
  `machine`        LowCardinality(String) CODEC (ZSTD(20)),
  `build_time`     DateTime CODEC (Delta(4), ZSTD(20)),
  `generated_time` DateTime CODEC (Delta(4), ZSTD(20)),
  `project`        LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_id`    UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `branch`         LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_type`  LowCardinality(String) CODEC (ZSTD(20)),

  `measures.name`  Array(LowCardinality(String)) CODEC (ZSTD(20)),
  `measures.value` Array(Int32) CODEC (Gorilla, ZSTD(20)),
  `measures.type`  Array(LowCardinality(String)) CODEC (ZSTD(20)),

  `triggeredBy`    LowCardinality(String) CODEC (ZSTD(20)),
  `mode`           LowCardinality(String) CODEC (ZSTD(20))
)
  engine = MergeTree
    partition by (toYYYYMM(generated_time))
    order by (branch, mode, project, machine, build_time, generated_time)