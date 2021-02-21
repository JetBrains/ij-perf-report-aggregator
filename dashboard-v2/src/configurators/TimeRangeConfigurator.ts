import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "@/dataQuery"
import { ref } from "vue"
import { PersistentStateManager } from "@/PersistentStateManager"
import { Item } from "@/configurators/DimensionConfigurator"
import { DataQueryExecutor } from "@/DataQueryExecutor"


export class TimeRangeConfigurator implements DataQueryConfigurator {
  static readonly timeRanges: Array<Item> = [
    {label: "All", value: "all"},
    {label: "Last year", value: "1y"},
    {label: "Last 3 months", value: "3M"},
    {label: "Last month", value: "1M"},
  ]
  static readonly defaultTimeRange = TimeRangeConfigurator.timeRanges[2].value

  readonly value = ref<string>(TimeRangeConfigurator.defaultTimeRange)

  constructor(dataQueryExecutor: DataQueryExecutor, persistentStateManager: PersistentStateManager) {
    persistentStateManager.add("timeRange", this.value)

    dataQueryExecutor.addConfigurator(this)
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  configure(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    const duration = this.value.value ?? TimeRangeConfigurator.defaultTimeRange
    if (duration !== "all") {
      const timeRange = parseDuration(duration)
      query.addFilter({field: "generated_time", sql: `> ${toClickhouseSql(timeRange)}`})
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
    if (result.length === 0) {
      expression += "now()"
    }
    else {
      expression += result
    }
    expression += `, ${value})`
    result = expression
  }
  return result
}

interface DurationParseResult {
  days?: number
  weeks?: number
  months?: number
  years?: number
}

const duration = /(-?\d*\.?\d+(?:e[-+]?\d+)?)\s*([a-zμ]*)/ig

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

const unitToDescriptor: { [key: string]: UnitDescriptor; } = {
  "d": units[0],
  "w": units[1],
  "M": units[2],
  "month": units[2],
  "y": units[3],
}

function parseDuration(s: string): DurationParseResult {
  const result: DurationParseResult = {}
  // ignore commas
  s = s.replace(/(\d),(\d)/g, "$1$2")
  // eslint-disable-next-line
  s.replace(duration, (_, ...args: any[]) => {
    const n = args[0]
    const unit = args[1]
    const unitDescriptor = unitToDescriptor[unit]
      || unitToDescriptor[unit.replace(/s$/, "")]
      || null
    if (unitDescriptor == null) {
      console.error(`unknown unit: ${unit}`)
    }
    else {
      unitDescriptor.apply(parseInt(n, 10), result)
    }
    return ""
  })
  return result
}