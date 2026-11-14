import { beforeEach, describe, expect, it } from "vitest"
import { PersistentStateManager } from "../../src/components/common/PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"
import { HistoryStartConfigurator } from "../../src/configurators/HistoryStartConfigurator"
import { TimeRangeConfigurator } from "../../src/configurators/TimeRangeConfigurator"
import { withSetup } from "../withSetup"

function daysAgo(days: number): string {
  const date = new Date()
  date.setDate(date.getDate() - days)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`
}

function timeRangeWithHistoryStart(historyStart: string | null): TimeRangeConfigurator {
  return withSetup(() => {
    const persistence = new PersistentStateManager("test-dashboard", { machine: "machine", project: [], branch: "b1" }, null)
    return new TimeRangeConfigurator(persistence, historyStart)
  })
}

async function axisOf(timeRange: TimeRangeConfigurator): Promise<{ min: Date | string | null; max: Date | string | null }> {
  const options = await timeRange.configureChart([], new DataQueryExecutorConfiguration())
  return options["xAxis"] as { min: Date | string | null; max: Date | string | null }
}

describe("History start", () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it("starts the axis at the history start when the time range starts earlier", async () => {
    const start = daysAgo(3)
    const axis = await axisOf(timeRangeWithHistoryStart(start))
    expect(axis.min).toStrictEqual(new Date(`${start}T00:00:00`))
    expect(axis.max).toBeInstanceOf(Date)
  })

  it("keeps the start of the time range when it is after the history start", async () => {
    const timeRange = timeRangeWithHistoryStart(daysAgo(60))
    timeRange.value.value = "1w"
    const axis = await axisOf(timeRange)
    const weekAgo = new Date()
    weekAgo.setDate(weekAgo.getDate() - 7)
    expect(Math.abs((axis.min as Date).getTime() - weekAgo.getTime())).toBeLessThan(60_000)
  })

  it("starts the axis at the history start for all the time", async () => {
    const start = daysAgo(3)
    const timeRange = timeRangeWithHistoryStart(start)
    timeRange.value.value = "all"
    const axis = await axisOf(timeRange)
    expect(axis.min).toStrictEqual(new Date(`${start}T00:00:00`))
  })

  it("keeps the start of a custom range when it is after the history start", async () => {
    const timeRange = timeRangeWithHistoryStart("2026-11-12")
    timeRange.setCustomRange(new Date(2026, 10, 20), new Date(2026, 10, 25))
    const axis = await axisOf(timeRange)
    expect(axis).toStrictEqual({ min: "2026-11-20", max: "2026-11-25" })
  })

  it("keeps the axis of the time range without a history start", async () => {
    const timeRange = timeRangeWithHistoryStart(null)
    timeRange.value.value = "all"
    const axis = await axisOf(timeRange)
    expect(axis.min).toBeNull()
  })

  // A chart applies each chart configurator with `replaceMerge: ["dataset", "series"]`, so an option of only an axis that
  // comes after the series removes every series of the chart.
  it("filters out the earlier points and adds no chart configurator", () => {
    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()
    new HistoryStartConfigurator("2026-11-12").configureQuery(query, configuration)
    expect(query.filters).toContainEqual({ f: "generated_time", q: ">= toDateTime('2026-11-12 00:00:00')" })
    expect(configuration.getChartConfigurators()).toHaveLength(0)
  })
})
