-- Description: Table to store the measurements of the test runs
create table fleet.measure_new
(
  `machine`               LowCardinality(String) CODEC (ZSTD(20)),
  `build_time`            DateTime CODEC (Delta(4), ZSTD(20)),
  `generated_time`        DateTime CODEC (Delta(4), ZSTD(20)),
  `project`               LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_id`           UInt32 CODEC (DoubleDelta, ZSTD(20)),
  `branch`                LowCardinality(String) CODEC (ZSTD(20)),
  `tc_build_type`         LowCardinality(String) CODEC (ZSTD(20)),

  `measures.name`         Array(LowCardinality(String)) CODEC (ZSTD(20)),
  `measures.value`        Array(Float64) CODEC (Gorilla, ZSTD(20)),
  `measures.type`         Array(LowCardinality(String)) CODEC (ZSTD(20)),

  `triggeredBy`           LowCardinality(String) CODEC (ZSTD(20))
)
  engine = MergeTree
    partition by (toYYYYMM(generated_time))
    order by (branch, project, machine, build_time, generated_time);

-- Copy data to new DB
INSERT INTO fleet.measure_new
SELECT
  machine,
  generated_time AS build_time,
  generated_time,
  project,
  tc_build_id,
  branch,
  '' AS tc_build_type,
  groupArray(name) AS `measures.name`,
  groupArray(CAST(value AS Float64) / 1000000) AS `measures.value`,
  arrayWithConstant(length(groupArray(name)),'d') AS `measures.type`,
  triggeredBy
FROM fleet.measure
GROUP BY
  machine,
  generated_time,
  project,
  tc_build_id,
  branch,
  triggeredBy;
