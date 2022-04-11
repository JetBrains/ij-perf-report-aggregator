import { BarSeriesOption } from "echarts/charts"
import { CallbackDataParams, DimensionDefinition } from "echarts/types/src/util/types"
import { Ref } from "vue"
import { DataQueryResult } from "../DataQueryExecutor"
import { ChartConfigurator, ChartStyle } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { BarChartOptions } from "../echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, numberFormat } from "../formatter"
import { measureNameToLabel } from "./MeasureConfigurator"
import { TimeRange, TimeRangeConfigurator } from "./TimeRangeConfigurator"

export class PredefinedGroupingMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  constructor(private readonly measures: Array<string>, private readonly timeRange: Ref<TimeRange>, private readonly chartStyle: ChartStyle) {
  }

  createObservable() {
    return null
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const timeRange = this.timeRange.value || TimeRangeConfigurator.timeRanges[0].value
    const interval = getClickHouseIntervalByDuration(timeRange)
    query.addDimension({
      name: "t",
      sql: `toStartOfInterval(generated_time, interval ${interval}, '${Intl.DateTimeFormat().resolvedOptions().timeZone}')`,
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
    const measureNames = this.measures.slice().reverse()
    if (query.db === "ij") {
      for (let i = 0; i < measureNames.length; i++) {
        query.addField(measureNames[i])
      }
    }
    else {
      query.addField({name: "measures", subName: "value"})
      query.addFilter({field: "measures.name", value: measureNames})
      if (measureNames.length > 1) {
        throw new Error("multiple measures are not supported")
      }
    }

    query.order = ["t"]

    configuration.measures = measureNames
    configuration.chartConfigurator = this
    return true
  }

  configureChart(dataList: DataQueryResult, configuration: DataQueryExecutorConfiguration): BarChartOptions {
    const queryProducers = configuration.extraQueryProducers
    if (queryProducers.length !== 0) {
      let useDurationFormatter = true

      // https://echarts.apache.org/examples/en/editor.html?c=dataset-simple1
      const dimensions: Array<DimensionDefinition> = []
      dimensions.push({name: "dimension", type: "ordinal"})
      const dimensionNameSet = new Set<string>()
      const source: Array<{ [key: string]: string | number }> = []

      for (let dataIndex = 0, n = dataList.length; dataIndex < n; dataIndex++) {
        if (useDurationFormatter && !isDurationFormatterApplicable(configuration.measureNames[dataIndex])) {
          useDurationFormatter = false
        }

        const column: { [key: string]: string | number } = {dimension: configuration.seriesNames[dataIndex]}
        source.push(column)
        for (const data of dataList[dataIndex]) {
          const valueKey = data[0] as string
          dimensionNameSet.add(valueKey)
          column[valueKey] = data[1]
        }
      }

      for (const name of dimensionNameSet) {
        dimensions.push({name, type: "number"})
      }

      const series = new Array<BarSeriesOption>(dimensions.length - 1)
      for (let i = 0; i < series.length; i++) {
        series[i] = {
          type: "bar",
          label: {
            show: true,
            formatter: useDurationFormatter ? formatterForFieldData : formatterForFieldNumericData,
            position: this.chartStyle.barSeriesLabelPosition,
          },
        }
      }
      return {
        dataset: {
          dimensions,
          source,
        },
        series,
      }
    }

    const data = dataList[0]
    const series = new Array<BarSeriesOption>(data.length)
    for (let i = 0; i < data.length; i++) {
      series[i] = {
        type: "bar",
        seriesLayoutBy: "row",
        label: {
          show: true,
          formatter,
          position: this.chartStyle.barSeriesLabelPosition,
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
      series,
    }
  }
}

const formatterForFieldData = function (data: CallbackDataParams) {
  const value = (data.value as { [key: string]: string | number })[data.seriesName as string] as number
  if (value > 10_000) {
    return durationAxisPointerFormatter(value)
  }
  else {
    return numberFormat.format(value)
  }
}

const formatterForFieldNumericData = function (data: CallbackDataParams) {
  const value = (data.value as { [key: string]: string | number })[data.seriesName as string] as number
  return numberFormat.format(value)
}

const formatter = function (data: CallbackDataParams): string {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  return numberFormat.format((data.value as Array<number>)[data.seriesIndex! + 1])
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
      throw new Error(`Unsupported time range: ${range as string}`)
  }
}