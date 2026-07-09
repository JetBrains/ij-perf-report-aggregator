import type { DefaultLabelFormatterCallbackParams as CallbackDataParams } from "echarts"
import type { EChartsType } from "echarts/core"
import type { OptionDataValue } from "../../shared/echarts-types"
import { SMOOTHED_SERIES_SUFFIX } from "../../configurators/MeasureConfigurator"
import { measureNameToLabel } from "../../shared/metricsMapping"
import { ValueUnit } from "../common/chart"
import { getBasicInfo } from "../common/sideBar/InfoSidebarPerformance"
import { useSettingsStore } from "../settings/settingsStore"

interface MetricPoint {
  buildTime: string
  buildNumber: string | null
  value: number | null
}

interface ChartOption {
  series?: { name?: string; id?: string; datasetIndex?: number }[]
  dataset?: { source?: unknown[] }[]
  legend?: { selected?: Record<string, boolean> }[]
}

/**
 * Exports the currently visible metrics of a chart as a structured YAML file.
 * The file is named after the chart title and contains the time period, branch, mode and agent
 * followed by a map of each line (series) to the visible points (build time, build number, value).
 */
export function exportChartMetricsAsYaml(chart: EChartsType, chartTitle: string, valueUnit: ValueUnit): void {
  const option = chart.getOption() as ChartOption
  const datasets = option.dataset ?? []
  const legendSelected = option.legend?.[0]?.selected ?? {}
  const [xMin, xMax] = getVisibleXRange(chart)
  const scaling = useSettingsStore().scaling

  const branches = new Set<string>()
  const modes = new Set<string>()
  const agents = new Set<string>()
  let minTime = Number.POSITIVE_INFINITY
  let maxTime = Number.NEGATIVE_INFINITY

  const metrics: Record<string, MetricPoint[]> = {}

  for (const series of option.series ?? []) {
    // skip the helper series that renders the smoothed line - it duplicates an existing series
    if (typeof series.id === "string" && series.id.endsWith(SMOOTHED_SERIES_SUFFIX)) {
      continue
    }
    const seriesName = series.name ?? ""

    // series hidden via the legend are not "currently visible"
    if (legendSelected[seriesName] != null && !legendSelected[seriesName]) {
      continue
    }

    const source = datasets[series.datasetIndex ?? 0]?.source
    if (!Array.isArray(source) || !Array.isArray(source[0])) {
      continue
    }
    const timestamps = source[0] as number[]
    const values = (scaling ? source.at(-1) : source[1]) as number[] | undefined

    const points: MetricPoint[] = []
    let key = seriesName
    for (let i = 0; i < timestamps.length; i++) {
      const time = timestamps[i]
      if (time < xMin || time > xMax) {
        continue
      }
      const row = (source as unknown[][]).map((column) => (Array.isArray(column) ? column[i] : undefined)) as OptionDataValue[]
      const info = getBasicInfo({ value: row, seriesName } as unknown as CallbackDataParams, valueUnit)
      const value = values?.[i]
      points.push({
        buildTime: new Date(time).toISOString(),
        buildNumber: info.build ?? null,
        value: typeof value === "number" && Number.isFinite(value) ? value : null,
      })
      if (info.branch) branches.add(info.branch)
      if (info.mode) modes.add(info.mode)
      if (info.machineName) agents.add(info.machineName)
      if (key === "" && info.metricName) key = measureNameToLabel(info.metricName)
      minTime = Math.min(minTime, time)
      maxTime = Math.max(maxTime, time)
    }

    if (points.length === 0) {
      continue
    }
    if (key === "") {
      key = chartTitle
    }
    // if several series share the same label keep them separated
    let uniqueKey = key
    let suffix = 2
    while (Object.hasOwn(metrics, uniqueKey)) {
      uniqueKey = `${key} (${suffix++})`
    }
    metrics[uniqueKey] = points
  }

  const doc: Record<string, unknown> = {
    chart: chartTitle,
    timePeriod:
      minTime <= maxTime
        ? {
            from: new Date(minTime).toISOString(),
            to: new Date(maxTime).toISOString(),
          }
        : null,
    branch: collapse(branches),
    mode: collapse(modes),
    agent: collapse(agents),
    metrics,
  }

  downloadYaml(sanitizeFileName(chartTitle) + ".yaml", toYaml(doc))
}

function collapse(values: Set<string>): string | string[] | null {
  if (values.size === 0) return null
  if (values.size === 1) return [...values][0]
  return [...values]
}

function getVisibleXRange(chart: EChartsType): [number, number] {
  try {
    const model = (
      chart as unknown as {
        getModel: () => { getComponent: (mainType: string, index: number) => { axis?: { scale?: { getExtent?: () => [number, number] } } } | undefined }
      }
    ).getModel()
    const extent = model.getComponent("xAxis", 0)?.axis?.scale?.getExtent?.()
    if (extent != null && Number.isFinite(extent[0]) && Number.isFinite(extent[1])) {
      return extent
    }
  } catch {
    // fall back to the full range below
  }
  return [Number.NEGATIVE_INFINITY, Number.POSITIVE_INFINITY]
}

function sanitizeFileName(name: string): string {
  const sanitized = name.replaceAll(/[^\w.-]+/g, "_").replaceAll(/^_+|_+$/g, "")
  return sanitized.length > 0 ? sanitized : "metrics"
}

function downloadYaml(fileName: string, content: string): void {
  const blob = new Blob([content], { type: "application/yaml;charset=utf-8" })
  const url = URL.createObjectURL(blob)
  const link = document.createElement("a")
  link.href = url
  link.download = fileName
  document.body.append(link)
  link.click()
  link.remove()
  URL.revokeObjectURL(url)
}

// Minimal YAML serializer sufficient for the plain object/array/scalar tree produced above.
function toYaml(value: unknown): string {
  return emit(value, 0).join("\n") + "\n"
}

function emit(value: unknown, indent: number): string[] {
  if (isRecord(value)) return emitObject(value, indent)
  if (Array.isArray(value)) return emitArray(value, indent)
  return [pad(indent) + formatScalar(value)]
}

function emitObject(value: Record<string, unknown>, indent: number): string[] {
  const entries = Object.entries(value)
  if (entries.length === 0) return [pad(indent) + "{}"]
  const lines: string[] = []
  for (const [key, child] of entries) {
    const prefix = pad(indent) + formatKey(key) + ":"
    if (isNonEmptyContainer(child)) {
      lines.push(prefix, ...emit(child, indent + 1))
    } else {
      lines.push(prefix + " " + formatScalar(child))
    }
  }
  return lines
}

function emitArray(value: unknown[], indent: number): string[] {
  if (value.length === 0) return [pad(indent) + "[]"]
  const lines: string[] = []
  for (const item of value) {
    if (isNonEmptyContainer(item)) {
      const [first, ...rest] = emit(item, indent + 1)
      // replace the leading indentation of the first line with the "- " marker
      lines.push(pad(indent) + "- " + first.trimStart(), ...rest)
    } else {
      lines.push(pad(indent) + "- " + formatScalar(item))
    }
  }
  return lines
}

function isNonEmptyContainer(value: unknown): boolean {
  return (isRecord(value) && Object.keys(value).length > 0) || (Array.isArray(value) && value.length > 0)
}

function pad(indent: number): string {
  return "  ".repeat(indent)
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value != null && !Array.isArray(value)
}

function formatKey(key: string): string {
  return /^[\w.-]+$/.test(key) ? key : quote(key)
}

function formatScalar(value: unknown): string {
  if (value == null) return "null"
  if (typeof value === "number") return Number.isFinite(value) ? String(value) : "null"
  if (typeof value === "boolean") return String(value)
  if (typeof value === "string") return quote(value)
  if (Array.isArray(value)) return "[" + value.map((it) => formatScalar(it)).join(", ") + "]"
  return quote(JSON.stringify(value))
}

function quote(value: string): string {
  return `"${value
    .replaceAll("\\", String.raw`\\`)
    .replaceAll('"', String.raw`\"`)
    .replaceAll("\n", String.raw`\n`)}"`
}
