import { ECBasicOption } from "echarts/types/dist/shared"
import { ChartConfigurator } from "../components/common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { FilterConfigurator } from "./filter"
import { TimeRangeConfigurator } from "./TimeRangeConfigurator"

/**
 * Leaves out the points generated before [start], a date in the `YYYY-MM-DD` form. A dashboard whose tests changed how
 * they measure shows only the points that the current method produced, so a chart and its alerts compare like with like.
 *
 * The x axis of a chart starts at [start] when the time range of [timeRange] starts earlier, so the chart has no empty
 * part before its first point.
 */
export class HistoryStartConfigurator implements DataQueryConfigurator, FilterConfigurator, ChartConfigurator {
  private readonly startTime: number

  constructor(
    private readonly start: string,
    private readonly timeRange: TimeRangeConfigurator
  ) {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(start)) {
      throw new Error(`The history start must be a YYYY-MM-DD date, not ${start}`)
    }
    this.startTime = new Date(`${start}T00:00:00`).getTime()
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    // The time range adds its own axis bounds before this one, and the bound of this one comes after them.
    configuration.addChartConfigurator(this)
    return this.configureFilter(query)
  }

  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "generated_time", q: `>= toDateTime('${this.start} 00:00:00')` })
    return true
  }

  async configureChart(data: (string | number)[][][], configuration: DataQueryExecutorConfiguration): Promise<ECBasicOption> {
    const rangeOptions = await this.timeRange.configureChart(data, configuration)
    const rangeAxis = rangeOptions["xAxis"] as { min: Date | string | null; max: Date | string | null }
    const rangeStart = rangeAxis.min == null ? Number.NEGATIVE_INFINITY : new Date(rangeAxis.min).getTime()
    return { xAxis: { min: new Date(Math.max(rangeStart, this.startTime)), max: rangeAxis.max } }
  }

  createObservable(): null {
    return null
  }
}
