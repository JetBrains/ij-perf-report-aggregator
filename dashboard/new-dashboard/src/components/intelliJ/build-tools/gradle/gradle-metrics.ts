import { MetricInfo } from "../../../../shared/metricsDescription"

export const GRADLE_METRICS_NEW_DASHBOARD: Map<string, string | MetricInfo> = new Map<string, string | MetricInfo>([
  ["ExternalSystemSyncProjectTask", "Total gradle sync time"],
  ["GradleExecution", "Full time of the sink operation, with all our overhead for preparation"],
  ["GradleConnection", "Work inside Gradle connection, operations that are performed inside connection"],
  ["GradleCall", "Resolving models from daemon"],
  ["ExternalSystemSyncResultProcessing", "Processing the data we received from Gradle"],
  ["ProjectDataServices", "Work of data services"],
  ["GradleProjectResolverDataProcessing", "Project resolve"],
  ["WorkspaceModelApply", "Apply ws model"],
  ["fus_gradle.sync", "Total gradle sync time (FUS)"],
])
