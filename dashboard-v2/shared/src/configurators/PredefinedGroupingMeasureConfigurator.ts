import { BarSeriesOption } from "echarts/charts"
import { TplFormatterParam } from "echarts/types/src/util/format"
import { Ref } from "vue"
import { DataQueryResult } from "../DataQueryExecutor"
import { ChartConfigurator, ChartOptions, numberFormat } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { measureNameToLabel } from "./MeasureConfigurator"
import { TimeRange, TimeRangeConfigurator } from "./TimeRangeConfigurator"

export class PredefinedGroupingMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  constructor(private readonly measures: Array<string>, private readonly timeRange: Ref<TimeRange>) {
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const timeRange = this.timeRange.value || TimeRangeConfigurator.timeRanges[0].value
    const interval = getClickHouseIntervalByDuration(timeRange)
    query.addDimension({
      name: "t",
      sql: `toStartOfInterval(generated_time, interval ${interval}, '${Intl.DateTimeFormat().resolvedOptions().timeZone}')`,
      // sql: `intDiv(toRelativeDayNum(date) - toRelativeDayNum(toDateTime('${new Date().toISOString()}'))`,
    })

    // do not use "Jan 06" because not clear - 06 here it is month or year
    if (timeRange === "1M") {
      query.timeDimensionFormat = "2 Jan"
    }
    else if (timeRange === "3M") {
      query.timeDimensionFormat = "Jan"
    }
    else {
      query.timeDimensionFormat = "2 Jan 2006"
    }

    // do not sort - bar chart shows series exactly in the same order as provided measure name list
    // reverse because echarts layout from bottom to top, but we need to put first measures to top
    const measures = this.measures.slice().reverse()
    for (let i = 0; i < measures.length; i++) {
      query.addField(measures[i])
    }

    query.order = ["t"]

    configuration.measures = measures
    configuration.chartConfigurator = this
    return true
  }

  configureChart(dataSets: DataQueryResult, configuration: DataQueryExecutorConfiguration): ChartOptions {
    const data = dataSets[0]
    const series = new Array<BarSeriesOption>(data.length)
    for (let i = 0; i < data.length; i++) {
      series[i] = {
        type: "bar",
        seriesLayoutBy: "row",
        label: {
          show: true,
          formatter(data: TplFormatterParam) {
            return numberFormat.format(Math.round(data["value"][i + 1]))
          },
          position: "right",
        },
      }
    }

    const measures = configuration.measures
    const header = new Array<string | number>(1 + measures.length)
    header[0] = "date"

    for (let i = 0; i < measures.length; i++) {
      header[i + 1] = measureNameToLabel(measures[i])
    }

    return {
      dataset: {
        sourceHeader: true,
        source: [header, ...data],
      },
      series: series as never,
    }
  }
}

function getClickHouseIntervalByDuration(range: TimeRange) {
  switch (range) {
    case "1M":
      return "7 day"
    case "3M":
      return "30 day"
    case "1y":
      return "90 day"
    case "all":
      return "180 day"
    default:
      throw new Error(`Unsupported time range: ${range}`)
  }
}