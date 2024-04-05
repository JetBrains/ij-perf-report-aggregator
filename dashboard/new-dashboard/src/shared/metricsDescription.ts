import { GRADLE_METRICS_NEW_DASHBOARD } from "../components/intelliJ/build-tools/gradle/gradle-metrics"

/**
 * Map of metric names to their descriptions.
 * The syntax is either:
 * ```
 * ["metricName", "MetricDescription"]
 * ["metricNamePrefix#*", "MetricDescription"]
 * ```
 * where `*` is a wildcard that matches any suffix.
 * For example, `processingTime#*` matches `processingTime#java` and `processingTime#cpp`.
 *
 * If you want to provide URL, use `metricInfo` function:
 * ```
 * ["metricName", metricInfo("MetricDescription", "https://example.com")]
 * ```
 */
export const metricsDescription: Map<string, string | MetricInfo> = new Map<string, string | MetricInfo>([
  //completion
  ["completion", "Total time of each completion invocation in test. Completion invocation time is a time that it takes to load all completion variants."],
  ["completion#mean_value", "Mean value of all completion invocation in test. Completion invocation time is a time that it takes to load all completion variants."],
  //find usages
  ["findUsage_popup", "Time to show the find usages popup"],
  ["findUsages", "Time to show all usages in the popup"],
  ["findUsages#number", "Number of found usages"],
  ["findUsages#mean_value", "Mean time to show all usages in the popup"],
  ["findUsages_firstUsage", "Time to show the first usage in the popup"],
  //analysis
  ["firstCodeAnalysis", "Time it takes to perform code analysis on file opening"],
  ["localInspections", "Sum time of all analysis. From Daemon#restart till DaemonListener#daemonFinished."],
  ["localInspections#mean_value", "Code analysis mean time. From Daemon#restart till DaemonListener#daemonFinished."],
  ["semanticHighlighting#mean_value", "Semantic highlighting mean time. From Daemon#restart till end of GeneralHighlightingPass."],
  ["runDaemon/executionTime", "Time it takes to complete a first daemon run. It might be restarted so it's not a full time."],
  ["codeAnalysisDaemon/fusExecutionTime", "Full highlighting duration since the file was modified and/or dumb mode status changed. It should be equal to the sum of segments."],
  ["globalInspections", "Time of all inspections runned in batch mode (Inspect Project)."],
  //indexing
  ["indexSize", "Index size in (in kb)"],
  ["numberOfIndexedFiles", "Number of indexed files"],
  ["processingSpeedAvg#*", "Speed of indexing file type (in kb/s)"],
  ["processingTime#*", "CPU time spent on processing file type"],
  ["indexingTimeWithoutPauses", "Indexing time without interruptions"],
  ["scanningTimeWithoutPauses", "Scanning time without interruptions"],
  ["pageLoad", "Number of regular Pages' loads."],
  ["pageMiss", "If the needed page has not existed in the main memory (RAM), it is known as PAGE MISS. The metric displays the number of unsuccessful Pages' obtainment."],
  [
    "pageHit",
    "CPU attempts to obtain a needed page from main memory and the page exists in main memory (RAM), it is referred to as a PAGE HIT. This metric displays the number of successful Pages' obtainment.",
  ],
  // startup
  ["reopenProjectPerformance/fusCodeVisibleInEditorDurationMs", metricInfo("New main metric for startup", "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric")],
  ["startup/fusTotalDuration", metricInfo("Old metric (outdated)", "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric")],
  //typing
  ["typing", "Typing executing time (usually equal to number of typed characters times delay between key presses)"],
  ["typing#average_awt_delay", "How long on average it takes to process a single empty AWT event in the queue during typing."],
  ["typing#max_awt_delay", "Max value of how long it takes to process a single empty AWT event in the queue during typing."],
  [
    "typing#latency",
    "Average time in ms of inserting a letter in the Editor (approximation of how long does it take from keyboard press till the appearance of the letter); measured during typing",
  ],
  //GC
  ["freedMemoryByGC", metricInfo("Freed memory by GC (in Mb/s)", "https://github.com/chewiebug/GCViewer#readme")],
  ["fullGCPause", metricInfo("Time that full GC was active (IDE is fully paused)", "https://github.com/chewiebug/GCViewer#readme")],
  ["gcPause", metricInfo("Time spent in GC (including minor collections without pausing)", "https://github.com/chewiebug/GCViewer#readme")],
  ["gcPauseCount", metricInfo("Number of minor GCs pauses", "https://github.com/chewiebug/GCViewer#readme")],
  //others
  ["searchEverywhere_*", "Time to fill all search everywhere results"],
  ["test#average_awt_delay", "The average time it takes to process a single empty AWT event in the queue during the whole test."],
  ["showQuickFixes", "Time to show the quick fixes after calling Alt + Enter."],
  ...GRADLE_METRICS_NEW_DASHBOARD,
])

export interface MetricInfo {
  description: string
  url?: string
}

function metricInfo(description: string, url?: string): MetricInfo {
  return { description, url }
}

function extractMainPrefix(inputString: string): string {
  const regex = /^(\w+#).*/
  const match = inputString.match(regex)
  return match ? match[1] : "non-matching"
}

export function getMetricDescription(metric: string | undefined): MetricInfo | null {
  if (metric == undefined) return null
  const metricDescription = metricsDescription.get(metric) ?? metricsDescription.get(extractMainPrefix(metric) + "*") ?? null
  return typeof metricDescription == "string" ? metricInfo(metricDescription) : metricDescription
}
