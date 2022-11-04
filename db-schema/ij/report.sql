create table report2
(
  `product`                    LowCardinality(String) CODEC(ZSTD(20)),
  `machine`                    LowCardinality(String) CODEC(ZSTD(20)),
  `build_time`                 DateTime CODEC(Gorilla, ZSTD(20)),
  `generated_time`             DateTime CODEC(Gorilla, ZSTD(20)),
  `project`                    LowCardinality(String) CODEC(ZSTD(20)),
  `tc_build_id`                UInt32 CODEC(Gorilla, ZSTD(20)),
  `tc_installer_build_id`      UInt32 CODEC(Gorilla, ZSTD(20)),
  `branch`                     LowCardinality(String) CODEC(ZSTD(20)),
  `raw_report`                 String CODEC(ZSTD(22)),

  `build_c1`                   UInt8 CODEC(Gorilla, ZSTD(20)),
  `build_c2`                   UInt16 CODEC(Gorilla, ZSTD(20)),
  `build_c3`                   UInt16 CODEC(Gorilla, ZSTD(20)),

  `triggeredBy`           LowCardinality(String) CODEC (ZSTD(20)),

  `bootstrap_d`          UInt32 CODEC (Gorilla, ZSTD(20)),
  `appInitPreparation_d` Int32 CODEC (Gorilla, ZSTD(20)),
  `appInit_d`            Int32 CODEC (Gorilla, ZSTD(20)),

  `pluginDescriptorLoading_d`    UInt16 CODEC(Gorilla, ZSTD(20)),
  `pluginDescriptorInitV18_d`    UInt16 CODEC(Gorilla, ZSTD(20)),

  `euaShowing_d` UInt16 CODEC(Gorilla, ZSTD(20)),
  `appStarter_d` UInt16 CODEC(Gorilla, ZSTD(20)),

  `appComponentCreation_d`     UInt16 CODEC(Gorilla, ZSTD(20)),
  `projectComponentCreation_d` UInt16 CODEC(Gorilla, ZSTD(20)),

  `serviceSyncPreloading_d`         UInt32 CODEC(Gorilla, ZSTD(20)),
  `serviceAsyncPreloading_d`        UInt32 CODEC(Gorilla, ZSTD(20)),
  `projectServiceSyncPreloading_d`  UInt32 CODEC(Gorilla, ZSTD(20)),
  `projectServiceAsyncPreloading_d` UInt32 CODEC(Gorilla, ZSTD(20)),

  `projectFrameInit_d`         UInt16 CODEC(Gorilla, ZSTD(20)),
  `projectProfileLoading_d`    UInt16 CODEC(Gorilla, ZSTD(20)),

  `projectDumbAware`         Int32 CODEC(Gorilla, ZSTD(20)),

  `editorRestoring`          Int32 CODEC(Gorilla, ZSTD(20)),
  `editorRestoringTillPaint` Int32 CODEC(Gorilla, ZSTD(20)),

  `splash_i`                   Int32 CODEC(Gorilla, ZSTD(20)),
  `startUpCompleted`         Int32 CODEC(Gorilla, ZSTD(20)),

  `classLoadingTime`       Int32 CODEC (Gorilla, ZSTD(20)),
  `classLoadingSearchTime` Int32 CODEC (Gorilla, ZSTD(20)),
  `classLoadingDefineTime` Int32 CODEC (Gorilla, ZSTD(20)),
  `classLoadingCount`      Int32 CODEC (Gorilla, ZSTD(20)),
  `classLoadingPreparedCount` Int32 CODEC (Gorilla, ZSTD(20)),
  `classLoadingLoadedCount`   Int32 CODEC (Gorilla, ZSTD(20)),

  `resourceLoadingTime`       Int32 CODEC (Gorilla, ZSTD(20)),
  `resourceLoadingCount`      Int32 CODEC (Gorilla, ZSTD(20)),

  `service.name` Array(LowCardinality(String)) CODEC (ZSTD(20)),
  `service.start` Array(UInt32) CODEC (Gorilla, ZSTD(20)),
  `service.duration` Array(UInt32) CODEC (ZSTD(20)),
  `service.thread` Array(LowCardinality(String)) CODEC (ZSTD(20)),
  `service.plugin` Array(LowCardinality(String)) CODEC (ZSTD(20)),

  `measure.name` Array(LowCardinality(String)) CODEC(ZSTD(20)),
  `measure.start` Array(UInt32) CODEC(Gorilla, ZSTD(20)),
  `measure.duration` Array(UInt32) CODEC(ZSTD(20)),
  `measure.thread` Array(LowCardinality(String)) CODEC(ZSTD(20))
)
  engine = MergeTree
    partition by (product, toYYYYMM(generated_time))
    order by (product, machine, branch, project, build_c1, build_c2, build_c3, build_time, generated_time);