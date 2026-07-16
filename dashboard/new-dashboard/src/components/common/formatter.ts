import humanizeDuration, { HumanizerOptions } from "humanize-duration"
import type { ValueUnit } from "./chart"
import { getMeasureUnit } from "../../shared/metricsDescription"

export function nsToMs(v: number) {
  return v / 1_000_000
}

export function formatPercentage(difference: number): string {
  return difference.toLocaleString(undefined, { style: "percent", minimumFractionDigits: 2 })
}

// https://github.com/apache/echarts/issues/8294
export const numberFormat = new Intl.NumberFormat(undefined, {
  maximumFractionDigits: 0,
})

const durationFormatOptions: HumanizerOptions = {
  language: "shortEn",
  round: true,
  units: ["y", "mo", "w", "d", "h", "m", "s", "ms"],
  languages: {
    shortEn: {
      y: () => "y",
      mo: () => "mo",
      w: () => "w",
      d: () => "d",
      h: () => "h",
      m: () => "min",
      s: () => "s",
      ms: () => "ms",
    },
  },
}

// Humanizes a duration given in milliseconds; sub-millisecond values fall back to µs/ns.
const durationAxisPointerFormatter = (valueInMs: number): string => {
  //humanizer doesn't handle values less than 1ms properly and just round them to either 0 or 1ms
  if (valueInMs == 0) {
    return "0"
  }
  if (valueInMs < 1) {
    return valueInMs * 1000 < 1 ? (valueInMs * 1000000).toLocaleString() + " ns" : (valueInMs * 1000).toLocaleString() + " μs"
  }
  const humanizer = humanizeDuration.humanizer(durationFormatOptions)
  return humanizer(valueInMs)
}

export const timeFormatWithoutSeconds = new Intl.DateTimeFormat("en-US", {
  year: "numeric",
  month: "short",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
})

export const durationFormatterInOneWord: (valueInMs: number) => string = humanizeDuration.humanizer({
  ...durationFormatOptions,
  largest: 1,
})

// Binary (IEC) units for memory sizes, as the JVM and GCViewer report them (base 1024).
const binaryUnits = ["B", "KiB", "MiB", "GiB", "TiB", "PiB"]
// Decimal (SI) units for file sizes and throughput, matching IntelliJ's formatFileSize (base 1000).
const decimalUnits = ["B", "kB", "MB", "GB", "TB", "PB"]

// Whether a memory size is stored in kilobytes (otherwise megabytes), read from the metric name.
function isKilobyteMeasure(measureName: string): boolean {
  return measureName.endsWith("Kb") || measureName.includes("(KB)")
}

// Memory sizes detected by metric name. Used only as a fallback for names without a declared unit:
// the stored type ("c"/"d") cannot tell a size from a plain count. The unit appears in several
// forms: a "Mb"/"Kb" suffix, a "Megabytes" suffix, or a "(MB)"/"(KB)" tag.
function isMemoryMeasure(measureName: string): boolean {
  return (
    isKilobyteMeasure(measureName) ||
    measureName.endsWith("Mb") ||
    measureName.endsWith("Megabytes") ||
    measureName.includes("(MB)") ||
    measureName.includes("totalHeapUsedMax") ||
    measureName.includes("freedMemoryByGC")
  )
}

// Auto-scales `value` (expressed in the unit at `startIndex` of `units`) by `base` to the largest
// unit that keeps the mantissa in [1, base), then appends `suffix`.
function scaleBy(value: number, base: number, units: string[], startIndex: number, suffix = ""): string {
  let unitIndex = startIndex
  let scaled = value
  while (Math.abs(scaled) >= base && unitIndex < units.length - 1) {
    scaled /= base
    unitIndex++
  }
  while (scaled !== 0 && Math.abs(scaled) < 1 && unitIndex > 0) {
    scaled *= base
    unitIndex--
  }
  return `${scaled.toLocaleString(undefined, { maximumFractionDigits: 2 })} ${units[unitIndex]}${suffix}`
}

// A duration metric by name, when neither a declared unit nor the stored type decides it. Fallback only.
function isDurationFormatterApplicable(measureName: string): boolean {
  return !(
    isMemoryMeasure(measureName) ||
    measureName.includes("number") ||
    measureName.includes("Number") ||
    measureName.includes("count") ||
    measureName.includes("Count") ||
    measureName.endsWith("_sources")
  )
}

// The resolved rendering unit for a single value. Distinct from ValueUnit, the chart-level *request*.
// Memory is binary (IEC): bytes/kibibytes/mebibytes render as B/KiB/MiB/GiB. File size and throughput
// are decimal (SI): kilobytes/megabytes render as kB/MB/GB and kilobytesPerSecond as kB/s, MB/s.
export type MeasureUnit = "nanoseconds" | "milliseconds" | "counter" | "bytes" | "kibibytes" | "mebibytes" | "kilobytes" | "megabytes" | "kilobytesPerSecond"

// Units carrying a physical quantity (a size or a rate). While scaling these become a baseline ratio.
const physicalUnits: ReadonlySet<MeasureUnit> = new Set<MeasureUnit>(["bytes", "kibibytes", "mebibytes", "kilobytes", "megabytes", "kilobytesPerSecond"])

// Resolves how a value should be rendered. A unit declared in metricsDescription is authoritative;
// otherwise an explicit value-unit wins. Then the metric name is authoritative over the stored type:
// the perf pipeline stores sizes and plain counts as "d"/"c" indiscriminately, so a name that is
// clearly a size or a count ("...Count", "...number", memory suffixes) must win over a mis-typed
// stored "d" duration. Only names that are not obviously a size/count defer to the stored type,
// falling back to a duration by default.
// While scaling, a physical unit is a baseline ratio and renders as a counter.
export function resolveMeasureUnit(measureName: string, opts: { storedType?: string; valueUnit?: ValueUnit; scaling?: boolean } = {}): MeasureUnit {
  const { storedType, valueUnit = "auto", scaling = false } = opts
  const declared = getMeasureUnit(measureName)
  if (scaling && (declared === undefined ? isMemoryMeasure(measureName) : physicalUnits.has(declared))) {
    return "counter"
  }
  if (declared !== undefined) return declared
  if (valueUnit === "counter") return "counter"
  if (valueUnit === "ns") return "nanoseconds"
  if (valueUnit === "ms") return "milliseconds"
  if (isMemoryMeasure(measureName)) return isKilobyteMeasure(measureName) ? "kibibytes" : "mebibytes"
  if (!isDurationFormatterApplicable(measureName)) return "counter"
  if (storedType === "d" || storedType === "duration") return "milliseconds"
  if (storedType === "c" || storedType === "counter") return "counter"
  return "milliseconds"
}

// Formats a single value that is already known to be in `unit`. Nanosecond values are converted to
// milliseconds before humanizing; sizes and throughput auto-scale to a readable unit.
export function formatMeasureValue(value: number, unit: MeasureUnit): string {
  switch (unit) {
    case "nanoseconds":
      return durationAxisPointerFormatter(value / 1_000_000)
    case "milliseconds":
      return durationAxisPointerFormatter(value)
    case "bytes":
      return scaleBy(value, 1024, binaryUnits, 0)
    case "kibibytes":
      return scaleBy(value, 1024, binaryUnits, 1)
    case "mebibytes":
      return scaleBy(value, 1024, binaryUnits, 2)
    case "kilobytes":
      return scaleBy(value, 1000, decimalUnits, 1)
    case "megabytes":
      return scaleBy(value, 1000, decimalUnits, 2)
    case "kilobytesPerSecond":
      return scaleBy(value, 1000, decimalUnits, 1, "/s")
    case "counter":
      return value.toLocaleString()
    default:
      throw new Error(`Unsupported measure unit: ${unit as string}`)
  }
}

// Reduces the per-series units of an aggregate chart to one axis unit: when every series is a
// duration the duration unit is kept; otherwise the first physical family present wins (a rate,
// then a binary memory size, then a decimal size), falling back to a plain counter.
export function reduceToAxisUnit(units: MeasureUnit[]): MeasureUnit {
  if (units.every((unit) => unit === "milliseconds" || unit === "nanoseconds")) {
    return units.includes("nanoseconds") ? "nanoseconds" : "milliseconds"
  }
  const firstOf = (...candidates: MeasureUnit[]): MeasureUnit | undefined => units.find((unit) => candidates.includes(unit))
  return firstOf("kilobytesPerSecond") ?? firstOf("mebibytes", "kibibytes", "bytes") ?? firstOf("megabytes", "kilobytes") ?? "counter"
}
