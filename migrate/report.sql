/*

1. update table sql if needed (see actual sql in clickhouse metadata)
2. backup old database (the whole clickhouse directory).
3. use clickhouse-client because other client read timeout maybe not enough (copy data is a long operation)

  docker run -it --rm --link ij-perf-clickhouse-server:clickhouse-server yandex/clickhouse-client --host clickhouse-server

*/

create table report2
(
  `product`                    Enum8('IU' = 1, 'WS' = 2, 'PS' = 3) CODEC(ZSTD(19)),
  `machine`                    Enum8('intellij-macos-hw-unit-1550' = 1, 'intellij-macos-hw-unit-1551' = 2, 'intellij-windows-hw-unit-499' = 3, 'intellij-windows-hw-unit-498' = 4, 'intellij-linux-hw-unit-558' = 5, 'intellij-linux-hw-unit-449' = 6, 'intellij-linux-hw-unit-450' = 7, 'intellij-linux-hw-unit-463' = 8, 'intellij-linux-hw-unit-504' = 9, 'intellij-linux-hw-unit-493' = 10, 'intellij-linux-hw-unit-556' = 11, 'intellij-linux-hw-unit-531' = 12, 'intellij-linux-hw-unit-484' = 13, 'intellij-linux-hw-unit-534' = 14, 'Dead agent' = 15) CODEC(ZSTD(19)),
  `build_time`                 DateTime CODEC(Delta(4), ZSTD(19)),
  `generated_time`             DateTime CODEC(Delta(4), ZSTD(19)),
  `tc_build_id`                UInt32 CODEC(DoubleDelta, ZSTD(19)),
  `tc_installer_build_id`      UInt32 CODEC(DoubleDelta, ZSTD(19)),
  `tc_build_properties`        String CODEC(ZSTD(19)),
  `branch`                     Enum8('master' = 0, '193' = 1) CODEC(ZSTD(19)),
  `raw_report`                 String CODEC(ZSTD(19)),
  `build_c1`                   UInt8 CODEC(DoubleDelta, ZSTD(19)),
  `build_c2`                   UInt16 CODEC(DoubleDelta, ZSTD(19)),
  `build_c3`                   UInt16 CODEC(DoubleDelta, ZSTD(19)),
  `bootstrap_d`                UInt16 CODEC(Gorilla, ZSTD(19)),
  `appInitPreparation_d`       UInt16 CODEC(Gorilla, ZSTD(19)),
  `appInit_d`                  UInt16 CODEC(Gorilla, ZSTD(19)),
  `pluginDescriptorLoading_d`  UInt16 CODEC(Gorilla, ZSTD(19)),
  `appComponentCreation_d`     UInt16 CODEC(Gorilla, ZSTD(19)),
  `projectComponentCreation_d` UInt16 CODEC(Gorilla, ZSTD(19)),
  `moduleLoading_d`            UInt16 CODEC(Gorilla, ZSTD(19)),
  `splash_i`                   Int32 CODEC(Gorilla, ZSTD(19)),
  `startUpCompleted_i`         Int32 CODEC(Gorilla, ZSTD(19))
)
  ENGINE = MergeTree
    PARTITION BY (product, toYYYYMM(generated_time))
    ORDER BY (product, machine, branch, build_c1, build_c2, build_c3, build_time, generated_time)
    SETTINGS old_parts_lifetime = 10, index_granularity = 8192;

/* use clickhouse-client and not IDEA executor - to get progress and proper read-timeout  */
insert into report2 select * from report;

drop table report;
rename table report2 to report