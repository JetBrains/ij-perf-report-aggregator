create table report2
(
  `product`                    Enum8('IU' = 1, 'WS' = 2, 'PS' = 3, 'DB' = 4, 'GO' = 5, 'RM' = 6) CODEC(ZSTD(20)),
  `machine`                    Enum8('intellij-macos-hw-unit-1550' = 1, 'intellij-macos-hw-unit-1551' = 2, 'intellij-windows-hw-unit-499' = 3, 'intellij-windows-hw-unit-498' = 4,
    'intellij-linux-hw-unit-558' = 5, 'intellij-linux-hw-unit-449' = 6, 'intellij-linux-hw-unit-450' = 7, 'intellij-linux-hw-unit-463' = 8,
    'intellij-linux-hw-unit-504' = 9, 'intellij-linux-hw-unit-493' = 10, 'intellij-linux-hw-unit-556' = 11, 'intellij-linux-hw-unit-531' = 12,
    'intellij-linux-hw-unit-484' = 13, 'intellij-linux-hw-unit-534' = 14, 'intellij-macos-hw-unit-1773' = 15, 'intellij-macos-hw-unit-1772' = 16,
    'intellij-windows-hw-unit-449' = 17, 'intellij-windows-hw-unit-463' = 18, 'intellij-windows-hw-unit-504' = 19, 'intellij-windows-hw-unit-493' = 20) CODEC(ZSTD(20)),
  `build_time`                 DateTime CODEC(Delta(4), ZSTD(20)),
  `generated_time`             DateTime CODEC(Delta(4), ZSTD(20)),
  `project`                    Enum8('/q9N7EHxr8F1NHjbNQnpqb0Q0fs' = 0, '1PbxeQ044EEghMOG9hNEFee05kM' = 1, '2j97HQ/UaQtIjdLtXG+0rvRD2Dw' = 2, '2tIppCeSStMfnlpjSo6a4ioxRd4' = 3, '73YWaW9bytiPDGuKvwNIYMK5CKI' = 4,
    'Encc+8LMF96ZwaDTGMC5f2QGj3M' = 5, 'JeNLJFVa04IA+Wasc+Hjj3z64R0' = 6, 'Xd232IXwZwqcumXkug0M3j82nCM' = 8,
    'Z10wDEdKQNOZUJiPgQ3mKGsMrOU' = 9, 'dEQpp0K1M7XEBh1IYVIg0Qf2pkc' = 10, 'j1a8nhKJexyL/zyuOXJ5CFOHYzU' = 11, 'nC4MRRFMVYUSQLNIvPgDt+B3JqA' = 12,
    'rF+W6q6L8SY62UPQ8EQhbENYOD0' = 13, 'uCOz0OEnAx055YJB55nc/baK48g' = 14, '36S9pX28rRGM06EN52YRr5qMIO0' = 15) CODEC(ZSTD(20)),
  `tc_build_id`                UInt32 CODEC(DoubleDelta, ZSTD(20)),
  `tc_installer_build_id`      UInt32 CODEC(DoubleDelta, ZSTD(20)),
  `tc_build_properties`        String CODEC(ZSTD(20)),
  `branch`                     Enum8('master' = 0, '193' = 1) CODEC(ZSTD(20)),
  `raw_report`                 String CODEC(ZSTD(20)),

  `build_c1`                   UInt8 CODEC(DoubleDelta, ZSTD(20)),
  `build_c2`                   UInt16 CODEC(DoubleDelta, ZSTD(20)),
  `build_c3`                   UInt16 CODEC(DoubleDelta, ZSTD(20)),

  `bootstrap_d`          UInt16 CODEC (Gorilla, ZSTD(20)),
  `appInitPreparation_d` UInt16 CODEC (Gorilla, ZSTD(20)),
  `appInit_d`            UInt16 CODEC (Gorilla, ZSTD(20)),

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

  `moduleLoading_d`            UInt16 CODEC(Gorilla, ZSTD(20)),
  `projectDumbAware_d`         UInt16 CODEC(Gorilla, ZSTD(20)),

  `editorRestoring_d`          UInt16 CODEC(Gorilla, ZSTD(20)),
  `editorRestoringTillPaint_d` UInt16 CODEC(Gorilla, ZSTD(20)),

  `splash_i`                   Int32 CODEC(Gorilla, ZSTD(20)),
  `startUpCompleted_i`         Int32 CODEC(Gorilla, ZSTD(20))
)
  engine = MergeTree
    partition by (product, toYYYYMM(generated_time))
    order by (product, machine, branch, project, build_c1, build_c2, build_c3, build_time, generated_time)
    settings old_parts_lifetime = 10;