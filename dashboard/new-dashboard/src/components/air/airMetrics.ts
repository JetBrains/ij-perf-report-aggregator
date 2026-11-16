import { MetricInfo } from "../../shared/metricsDescription"

const sizes = ["M", "XL", "XXL"]

function sized(template: string): string[] {
  return sizes.map((size) => template.replace("%s", size))
}

function entries(names: string[], description: string, unit: MetricInfo["unit"]): [string, MetricInfo][] {
  return names.map((name) => [name, { description, unit }])
}

/**
 * The units of the size metrics of the Air performance tests. Without a declared unit, the dashboard reads a metric name
 * that ends in `Bytes` or `MiB` as a time.
 */
export const AIR_METRICS: Map<string, string | MetricInfo> = new Map<string, string | MetricInfo>([
  ...entries(sized("continue.%s.requestBytes"), "The size of the request body that the ACP adapter sends to the model for the continued turn.", "bytes"),
  ...entries(["live.backgroundBytes", "live.second.backgroundBytes"], "The bytes of the requests that the ACP adapter sends besides the steps of the turn.", "bytes"),
  ...entries(
    [...sized("open.%s.agentRssPeakMiB"), "live.agentRssPeakMiB", "live.second.agentRssPeakMiB"],
    "The largest sum of the resident memory of the ACP adapter process tree.",
    "mebibytes"
  ),
  ...entries([...sized("heap.%s.heapLive.MiB"), "manyAgents.heapLive.MiB", "sessionList.heapLive.MiB"], "The live heap of the IDE after a full GC.", "mebibytes"),
])
