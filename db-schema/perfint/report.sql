create table report
(
  `machine`               LowCardinality(String) CODEC (ZSTD(20)),
  `build_time`            DateTime CODEC (Delta(4), ZSTD(20)),
  `generated_time`        DateTime CODEC (Delta(4), ZSTD(20)),
  `project`               LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_id`           UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `tc_installer_build_id` UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `tc_build_properties`   String CODEC (ZSTD(20)),
  `branch`                LowCardinality(String) CODEC (ZSTD(20)),
  `raw_report`            String CODEC (ZSTD(20)),

  measures Nested(
    name LowCardinality(String),
    value Int32
  ) CODEC (ZSTD(20)),

  `build_c1`              UInt8 CODEC (DoubleDelta, ZSTD(20)),
  `build_c2`              UInt16 CODEC (DoubleDelta, ZSTD(20)),
  `build_c3`              UInt16 CODEC (DoubleDelta, ZSTD(20))
)
  engine = MergeTree
    partition by (toYYYYMM(generated_time))
    order by (machine, branch, project, build_c1, build_c2, build_c3, build_time, generated_time)
    settings old_parts_lifetime = 10;