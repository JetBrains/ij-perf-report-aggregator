import { ECBasicOption } from "echarts/types/dist/shared"
import { combineLatest, Observable } from "rxjs"
import { provide, Ref, ref, watch } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { ChartConfigurator } from "../components/common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { timeRangeKey } from "../shared/injectionKeys"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export declare type TimeRange = "1w" | "2w" | "1M" | "3M" | "1y" | "all" | "custom"

export interface TimeRangeItem {
  label: string
  value: TimeRange
}

export class TimeRangeConfigurator implements DataQueryConfigurator, FilterConfigurator, ChartConfigurator {
  readonly value = ref<TimeRange>("1M")
  readonly customRange = ref<string>("")
  public timeRanges: Ref<TimeRangeItem[]> = ref([
    { label: "Last week", value: "1w" },
    { label: "Last two weeks", value: "2w" },
    { label: "Last month", value: "1M" },
    { label: "Last 3 months", value: "3M" },
    { label: "Last year", value: "1y" },
    { label: "All", value: "all" },
    { label: this.customRange, value: "custom" },
  ])

  constructor(persistentStateManager: PersistentStateManager) {
    provide(timeRangeKey, this.value)
    persistentStateManager.add("timeRange", this.value)
    persistentStateManager.add("customRange", this.customRange)
    if (this.customRange.value != "") {
      this.value.value = "custom"
    }
    watch(this.value, (customRangeValue) => {
      if (customRangeValue != "custom") {
        this.customRange.value = ""
      }
    })
  }

  public setCustomRange(start: Date, end: Date) {
    this.value.value = "custom"
    this.customRange.value = `${this.getStringFromDate(start)}:${this.getStringFromDate(end)}`
  }

  private getStringFromDate(date: Date): string {
    return date.getFullYear().toString() + "-" + (date.getMonth() + 1).toString() + "-" + date.getDate().toString()
  }

  createObservable(): Observable<[TimeRange, string]> {
    return combineLatest([refToObservable(this.value), refToObservable(this.customRange)]).pipe()
  }

  configureFilter(query: DataQuery): boolean {
    return this.configureQuery(query, null)
  }

  configureChart(_data: (string | number)[][][], _configuration: DataQueryExecutorConfiguration): Promise<ECBasicOption> {
    const timeRange = this.value.value
    if (this.customRange.value == "") {
      return Promise.resolve({ xAxis: { min: getStartTime(timeRange), max: getEndTime(timeRange) } })
    }
    const between = this.customRange.value.split(":")
    return Promise.resolve({ xAxis: { min: between[0], max: between[1] } })
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration | null): boolean {
    const duration = this.value.value
    configuration?.addChartConfigurator(this)
    if (duration === "all") {
      return true
    }
    if (this.customRange.value == "") {
      const sql = `>${toClickhouseSql(parseDuration(duration))}`
      query.addFilter({ f: "generated_time", q: sql })
    } else {
      const between = this.customRange.value.split(":")
      const startDate = this.getSQLDateTime(between[0])
      const endDate = this.getSQLDateTime(between[1])
      const sql = `BETWEEN toDateTime('${startDate} 00:00:00') AND toDateTime('${endDate} 23:59:59')`
      query.addFilter({ f: "generated_time", q: sql })
    }
    return true
  }

  private getSQLDateTime(dateString: string): string {
    const [year, month, day] = dateString.split("-").map(Number)
    const date = new Date(Date.UTC(year, month - 1, day))
    return date.getFullYear().toString() + "-" + (date.getMonth() + 1).toString().padStart(2, "0") + "-" + date.getDate().toString().padStart(2, "0")
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

export function parseDuration(s: string): DurationParseResult {
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

function getStartTime(range: TimeRange): Date | null {
  const currentDate = new Date()

  switch (range) {
    case "1w": // Subtract 1 week
      currentDate.setDate(currentDate.getDate() - 7)
      break
    case "2w": // Subtract 2 weeks
      currentDate.setDate(currentDate.getDate() - 14)
      break
    case "1M": // Subtract 1 month
      currentDate.setMonth(currentDate.getMonth() - 1)
      break
    case "3M": // Subtract 3 months
      currentDate.setMonth(currentDate.getMonth() - 3)
      break
    case "1y": // Subtract 1 year
      currentDate.setFullYear(currentDate.getFullYear() - 1)
      break
    case "all": // Subtract 10 years
      return null
    case "custom":
      return null
  }

  return currentDate
}

function getEndTime(range: TimeRange): Date | null {
  switch (range) {
    case "1w":
    case "2w":
    case "1M":
    case "3M":
    case "1y":
    case "all":
      return new Date()
    case "custom":
      return null //don't set end time for custom range
  }
}
