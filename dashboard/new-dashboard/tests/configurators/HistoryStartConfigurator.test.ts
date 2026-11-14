import { beforeEach, describe, expect, it } from "vitest"
import { DataQuery, DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"
import { HistoryStartConfigurator } from "../../src/configurators/HistoryStartConfigurator"
import { TimeRangeConfigurator } from "../../src/configurators/TimeRangeConfigurator"
import ConfiguratorTest from "./ConfiguratorTest"

function daysAgo(days: number): string {
  const date = new Date()
  date.setDate(date.getDate() - days)
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")}`
}

describe("History start configurator", () => {
  let timeRange: TimeRangeConfigurator

  beforeEach(() => {
    localStorage.clear()
    timeRange = ConfiguratorTest.setupPreconditions([]).timeRangeConfigurator
  })

  async function axisOf(start: string): Promise<{ min: Date; max: Date | string | null }> {
    const options = await new HistoryStartConfigurator(start, timeRange).configureChart([], new DataQueryExecutorConfiguration())
    return options["xAxis"] as { min: Date; max: Date | string | null }
  }

  it("starts the axis at the history start when the time range starts earlier", async () => {
    const start = daysAgo(3)
    const axis = await axisOf(start)
    expect(axis.min.getTime()).toBe(new Date(`${start}T00:00:00`).getTime())
    expect(axis.max).toBeInstanceOf(Date)
  })

  it("keeps the start of the time range when it is after the history start", async () => {
    timeRange.value.value = "1w"
    const axis = await axisOf(daysAgo(60))
    const weekAgo = new Date()
    weekAgo.setDate(weekAgo.getDate() - 7)
    expect(Math.abs(axis.min.getTime() - weekAgo.getTime())).toBeLessThan(60_000)
  })

  it("starts the axis at the history start for all the time", async () => {
    timeRange.value.value = "all"
    const start = daysAgo(3)
    const axis = await axisOf(start)
    expect(axis.min.getTime()).toBe(new Date(`${start}T00:00:00`).getTime())
  })

  it("keeps the start of a custom range when it is after the history start", async () => {
    timeRange.setCustomRange(new Date(2026, 10, 20), new Date(2026, 10, 25))
    const axis = await axisOf("2026-11-12")
    expect(axis.min.getTime()).toBe(new Date("2026-11-20").getTime())
    expect(axis.max).toBe("2026-11-25")
  })

  it("filters out the points before the history start and asks to configure the chart", () => {
    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()
    new HistoryStartConfigurator("2026-11-12", timeRange).configureQuery(query, configuration)
    expect(query.filters).toContainEqual({ f: "generated_time", q: ">= toDateTime('2026-11-12 00:00:00')" })
    expect(configuration.getChartConfigurators()).toHaveLength(1)
  })
})
