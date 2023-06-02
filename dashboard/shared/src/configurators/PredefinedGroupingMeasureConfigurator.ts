import { BarSeriesOption } from "echarts/charts"
import { CallbackDataParams, DimensionDefinition } from "echarts/types/src/util/types"
import { Ref } from "vue"
import { DataQueryResult } from "../DataQueryExecutor"
import { ChartConfigurator, ChartStyle, ValueUnit } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { BarChartOptions } from "../echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, nsToMs, numberFormat } from "../formatter"
import { measureNameToLabel } from "./MeasureConfigurator"
import { TimeRange, TimeRangeConfigurator } from "./TimeRangeConfigurator"

export class PredefinedGroupingMeasureConfigurator implements DataQueryConfigurator, ChartConfigurator {
  constructor(private readonly measures: string[],
              private readonly timeRange: Ref<TimeRange>,
              private readonly chartStyle: ChartStyle) {}
  createObservable() {
    return null
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const timeRange = this.timeRange.value || TimeRangeConfigurator.timeRanges[0].value
    const interval = getClickHouseIntervalByDuration(timeRange)
    query.addDimension({
      n: "t",
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
    const measureNames = [...this.measures].reverse()
    if (query.db === "ij") {
      for (const measureName of measureNames) {
        query.addField(measureName)
      }
    }
    else {
      if (measureNames.length > 1) {
        throw new Error("multiple measures are not supported")
      }

      if (query.table === "measure") {
        query.addField({n: "value"})
        query.addFilter({f: "name", v: measureNames[0]})
      }
      else {
        query.addField({n: "measures", subName: "value"})
        query.addFilter({f: "measures.name", v: measureNames})
        if (measureNames.length > 1) {
          throw new Error("multiple measures are not supported")
        }
      }
    }

    query.order = "t"

    configuration.measures = measureNames
    configuration.chartConfigurator = this
    return true
  }

  configureChart(dataList: DataQueryResult, configuration: DataQueryExecutorConfiguration): BarChartOptions {
    return configureWithQueryProducers(dataList, configuration, this.chartStyle)
  }
}

function configureWithQueryProducers(dataList: (string | number)[][][], configuration: DataQueryExecutorConfiguration, chartStyle: ChartStyle): BarChartOptions {
  let useDurationFormatter = true

  const dimensionNameSet = new Set<string>()
  const source: { [key: string]: string | number }[] = []

  // outer cycle over measures to group it together (e.g. if several machines are selected)
  for (let measureIndex = 0; measureIndex < configuration.measures.length; measureIndex++) {
    for (let dataIndex = 0, n = dataList.length; dataIndex < n; dataIndex++) {
      if (useDurationFormatter && !isDurationFormatterApplicable(configuration.measureNames[dataIndex])) {
        useDurationFormatter = false
      }

      const measure = measureNameToLabel(configuration.measures[measureIndex])
      const classifier = configuration.seriesNames[dataIndex]
      let dimension = measure
      if (classifier.length > 0) {
        dimension += ` | ${classifier}`
      }
      const column: { [key: string]: string | number } = {dimension}
      source.push(column)
      const result = dataList[dataIndex]
      for (let i = 0; i < result[0].length; i++) {
        const valueKey = result[0][i] as string
        dimensionNameSet.add(valueKey)
        column[valueKey] = result[measureIndex + 1][i]
      }
    }
  }

  // https://echarts.apache.org/examples/en/editor.html?c=dataset-simple1
  const dimensions: DimensionDefinition[] = []
  dimensions.push({name: "dimension", type: "ordinal"})
  for (const name of dimensionNameSet) {
    dimensions.push({name, type: "number"})
  }

  const series = new Array<BarSeriesOption>(dimensions.length - 1)
  const formatter = getSeriesLabelFormatter(useDurationFormatter, chartStyle.valueUnit)
  for (let i = 0; i < series.length; i++) {
    series[i] = {
      type: "bar",
      label: {
        show: true,
        formatter,
        position: chartStyle.barSeriesLabelPosition,
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

function getSeriesLabelFormatter(useDurationFormatter: boolean, valueUnit: ValueUnit): (p: CallbackDataParams) => string {
  const converter: (it: number) => number = valueUnit === "ns" ? nsToMs : it => it
  return function (data: CallbackDataParams): string {
    const value = converter((data.value as { [key: string]: string | number })[data.seriesName as string] as number)
    return useDurationFormatter && value > 10_000 ? durationAxisPointerFormatter(value) : numberFormat.format(value)
  }
}

function getClickHouseIntervalByDuration(range: TimeRange) {
  switch (range) {
    case "1w":
      return "7 day"
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