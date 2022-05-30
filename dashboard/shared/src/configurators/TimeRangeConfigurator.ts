import { Observable } from "rxjs"
import { provide, ref } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { timeRangeKey } from "../injectionKeys"
import { refToObservable } from "./rxjs"

export declare type TimeRange = "1M" | "3M" | "1y" | "all"

export interface TimeRangeItem {
  label: string
  value: TimeRange
}

export class TimeRangeConfigurator implements DataQueryConfigurator {
  static readonly timeRanges: Array<TimeRangeItem> = [
    {label: "Last month", value: "1M"},
    {label: "Last 3 months", value: "3M"},
    {label: "Last year", value: "1y"},
    {label: "All", value: "all"},
  ]

  readonly value = ref<string>(TimeRangeConfigurator.timeRanges[0].value)

  constructor(persistentStateManager: PersistentStateManager) {
    provide(timeRangeKey, this.value)

    persistentStateManager.add("timeRange", this.value)
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.value)
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    const duration = this.value.value ?? TimeRangeConfigurator.timeRanges[0].value
    if (duration === "all") {
      return true
    }
    if (typeof duration !== "string") {
      return false
    }

    const sql = `>${toClickhouseSql(parseDuration(duration))}`
    query.addFilter({f: "generated_time", q: sql})
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

const units: Array<UnitDescriptor> = [
  {
    subtractFunction: "subtractDays",
    apply(value: number, result: DurationParseResult): void {
      result.days = value
    },
    getValue: result => result.days,
  },
  {
    subtractFunction: "subtractWeeks",
    apply(value: number, result: DurationParseResult): void {
      result.weeks = value
    },
    getValue: result => result.weeks,
  },
  {
    subtractFunction: "subtractMonths",
    apply(value: number, result: DurationParseResult): void {
      result.months = value
    },
    getValue: result => result.months,
  },
  {
    subtractFunction: "subtractYears",
    apply(value: number, result: DurationParseResult): void {
      result.years = value
    },
    getValue: result => result.years,
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
  s = s.replace(/(\d),(\d)/g, "$1$2")
  s.replace(duration, (_, ...args: string[]) => {
    const n = args[0]
    const unit = args[1]
    const unitDescriptor = unitToDescriptor.get(unit) ?? unitToDescriptor.get(unit.replace(/s$/, ""))
    if (unitDescriptor == null) {
      console.error(`unknown unit: ${unit}`)
    }
    else {
      unitDescriptor.apply(Number.parseInt(n, 10), result)
    }
    return ""
  })
  return result
}