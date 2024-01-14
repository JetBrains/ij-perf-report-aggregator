export const GRADLE_METRICS = [
  "gradle.sync.duration",
  "GRADLE_CALL",
  "DATA_SERVICES",
  "PROJECT_RESOLVERS",
  "WORKSPACE_MODEL_APPLY",
  "fus_gradle.sync",

  "AWTEventQueue.dispatchTimeTotal",
  "CPU | Load |Total % 95th pctl",
  "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
  "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
  "gcPause",
  "gcPauseCount",
  "fullGCPause",
  "freedMemoryByGC",
  "totalHeapUsedMax",
  "JVM.GC.collectionTimesMs",
  "JVM.GC.collections",
  "JVM.maxHeapMegabytes",
  "JVM.threadCount",
  "JVM.totalCpuTimeMs",
  "JVM.totalMegabytesAllocated",
  "JVM.usedHeapMegabytes",
  "JVM.usedNativeMegabytes",
]

export const GRADLE_METRICS_NEW_DASHBOARD = [
  // total sync time
  "ExternalSystemSyncProjectTask",
  // full time of the sink operation, with all our overhead for preparation
  "GradleExecution",
  // work inside Gradle connection, operations that are performed inside connection
  "GradleConnection",
  // resolving models from daemon
  "GradleCall",
  // processing the data we received from Gradle
  "ExternalSystemSyncResultProcessing",
  // work of data services
  "ProjectDataServices",
  // project resolve
  "GradleProjectResolvers",
  // apply ws model
  "WorkspaceModelApply",
]
