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
  ["completion", "Total time of each completion invocation in test. Completion invocation time is a time that it takes to load all completion variants."],
  ["findUsage_popup", "Time to show the find usages popup"],
  ["findUsages", "Time to show all usages in the popup"],
  ["findUsages#number", "Number of found usages"],
  ["findUsages_firstUsage", "Time to show the first usage in the popup"],
  ["firstCodeAnalysis", "Time it takes to perform code analysis on file opening"],
  ["freedMemoryByGC", metricInfo("Freed memory by GC (in Mb/s)", "https://github.com/chewiebug/GCViewer#readme")],
  ["indexSize", "Index size in (in kb)"],
  ["numberOfIndexedFiles", "Number of indexed files"],
  ["processingSpeedAvg#*", "Speed of indexing file type (in kb/s)"],
  ["processingTime#*", "CPU time spent on processing file type"],
  ["test#average_awt_delay", "The average time it takes to process a single empty AWT event in the queue during the whole test."],
  ["reopenProjectPerformance/fusCodeVisibleInEditorDurationMs", metricInfo("New main metric for startup", "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric")],
  ["startup/fusTotalDuration", metricInfo("Old metric (outdated)", "https://youtrack.jetbrains.com/articles/IJPL-A-286/Startup-Metric")],
  ["completion#mean_value", "Mean value of all completion invocation in test. Completion invocation time is a time that it takes to load all completion variants."],
  ["localInspections#mean_value", "Code analysis mean time. From Daemon#restart till DaemonListener#daemonFinished."],
  ["semanticHighlighting#mean_value", "Semantic highlighting mean time. From Daemon#restart till end of GeneralHighlightingPass."],
  ["findUsages#mean_value", "Mean time to show all usages in the popup"],
  ["test#average_awt_delay", "Average time in ms of processing an event in Event Dispatch Queue (approximation of UI delays); measured during the whole test"],
  [
    "typing#latency",
    "Average time in ms of inserting a letter in the Editor (approximation of how long does it take from keyboard press till the appearance of the letter); measured during typing",
  ],
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
