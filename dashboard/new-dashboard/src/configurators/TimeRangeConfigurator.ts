import { Observable } from "rxjs"
import { provide, ref } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { timeRangeKey } from "../shared/injectionKeys"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export declare type TimeRange = "1w" | "1M" | "3M" | "1y" | "all" | "custom"

export interface TimeRangeItem {
  label: string
  value: TimeRange
  customRange: string
}

export class TimeRangeConfigurator implements DataQueryConfigurator, FilterConfigurator {
  static readonly timeRanges: TimeRangeItem[] = [
    { label: "Last week", value: "1w" , customRange : ""},
    { label: "Last month", value: "1M" , customRange : ""},
    { label: "Last 3 months", value: "3M" , customRange : ""},
    { label: "Last year", value: "1y" , customRange : ""},
    { label: "All", value: "all" , customRange : ""},
    { label: "", value: "custom" , customRange : ""},
  ]

  static readonly timeRangeValueToItem = new Map<string, TimeRangeItem>(TimeRangeConfigurator.timeRanges.map((it) => [it.value, it]))

  readonly value = ref<TimeRange>(TimeRangeConfigurator.timeRanges[0].value)
  readonly customRange = ref<string>(TimeRangeConfigurator.timeRanges[0].customRange)

  constructor(persistentStateManager: PersistentStateManager) {
    provide(timeRangeKey, this.value)
    persistentStateManager.add("timeRange", this.value)
    persistentStateManager.add("customRange", this.customRange)
  }

  createObservable(): Observable<TimeRange> {
    return refToObservable(this.value)
  }

  configureFilter(query: DataQuery): boolean {
    return this.configureQuery(query, null)
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration | null): boolean {
    const duration = this.value.value
    if (duration === "all") {
      return true
    }
    if (this.value.value == "custom" && this.customRange.value != "") {
      const between = this.customRange.value.split(":")
      // const ago = getDateAgoByDuration(this.value.value)
      const sql = `BETWEEN toDate('${between[0]}') AND toDate('${between[1]}')`
      // const sql = `BETWEEN toDate('${ago.getFullYear()}-${ago.getMonth() + 1}-${ago.getDate()}') AND toDate('${now.getFullYear()}-${now.getMonth() + 1}-${now.getDate()}')`
      query.addFilter({ f: "generated_time", q: sql })

    } else {
      const sql = `>${toClickhouseSql(parseDuration(duration))}`
      query.addFilter({ f: "generated_time", q: sql })
    }
    return true
  }
}

function toClickhouseSql(duration: DurationParseResult): string {
  let result = ""
  for (const unit of units) {
    const value = unit.getValue(duration)
    if (value == undefined) {
      continue
    }

    let expression = `${unit.subtractFunction}(`
    expression += result.length === 0 ? "now()" : result
    expression += `,${value})`
    result = expression
  }
  return result
}

export interface DurationParseResult {
  days?: number
  weeks?: number
  months?: number
  years?: number
}

const duration = /(-?\d*\.?\d+(?:e[+-]?\d+)?)\s*([a-zμ]*)/gi

interface UnitDescriptor {
  readonly subtractFunction: string

  apply(value: number, result: DurationParseResult): void

  getValue(result: DurationParseResult): number | undefined
}

const units: UnitDescriptor[] = [
  {
    subtractFunction: "subtractDays",
    apply(value: number, result: DurationParseResult): void {
      result.days = value
    },
    getValue: (result) => result.days,
  },
  {
    subtractFunction: "subtractWeeks",
    apply(value: number, result: DurationParseResult): void {
      result.weeks = value
    },
    getValue: (result) => result.weeks,
  },
  {
    subtractFunction: "subtractMonths",
    apply(value: number, result: DurationParseResult): void {
      result.months = value
    },
    getValue: (result) => result.months,
  },
  {
    subtractFunction: "subtractYears",
    apply(value: number, result: DurationParseResult): void {
      result.years = value
    },
    getValue: (result) => result.years,
  },
]

const unitToDescriptor = new Map<string, UnitDescriptor>([
  ["d", units[0]],
  ["w", units[1]],
  ["M", units[2]],
  ["month", units[2]],
  ["y", units[3]],
])

function parseDuration(s: string): DurationParseResult {
  const result: DurationParseResult = {}
  // ignore commas
  s = s.replaceAll(/(\d),(\d)/g, "$1$2")
  s.replaceAll(duration, (_, ...args: string[]) => {
    const [n, unit] = args
    const unitDescriptor = unitToDescriptor.get(unit) ?? unitToDescriptor.get(unit.replace(/s$/, ""))
    if (unitDescriptor == null) {
      console.error(`unknown unit: ${unit}`)
    } else {
      unitDescriptor.apply(Number.parseInt(n, 10), result)
    }
    return ""
  })
  return result
}
