import { BarSeriesOption } from "echarts/charts"
import { CallbackDataParams, DimensionDefinition } from "echarts/types/src/util/types"
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
    if (query.db !== "ij" && query.db !== "fleet") {
      query.addField({name: "measures", subName: "value"})
      query.addFilter({field: "measures.name", value: measureNames})
      if (measureNames.length > 1) {
        throw new Error("multiple measures are not supported")
      }
    }
    else {
      for (let i = 0; i < measureNames.length; i++) {
        query.addField(measureNames[i])
      }
    }

    query.order = ["t"]

    configuration.measures = measureNames
    configuration.chartConfigurator = this
    return true
  }

  configureChart(dataSets: DataQueryResult, configuration: DataQueryExecutorConfiguration): ChartOptions {
    const extraQueryProducer = configuration.extraQueryProducer
    if (extraQueryProducer != null) {
      // https://echarts.apache.org/examples/en/editor.html?c=dataset-simple1
      const dimensions: Array<DimensionDefinition> = []
      dimensions.push({name: "dimension", type: "ordinal"})
      const dimensionNameSet = new Set<string>()
      const source: Array<{ [key: string]: string | number }> = []
      for (let dataSetIndex = 0; dataSetIndex < dataSets.length; dataSetIndex++){
        const dataSet = dataSets[dataSetIndex]
        const dataSetLabel = extraQueryProducer.getDataSetLabel(dataSetIndex)
        const column: { [key: string]: string | number } = {dimension: dataSetLabel}
        source.push(column)
        for (const data of dataSet) {
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
            formatter(data: CallbackDataParams) {
              return numberFormat.format((data.value as { [key: string]: string | number })[data.seriesName as string] as number)
            },
            position: "right",
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

    const data = dataSets[0]
    const series = new Array<BarSeriesOption>(data.length)
    for (let i = 0; i < data.length; i++) {
      series[i] = {
        type: "bar",
        seriesLayoutBy: "row",
        label: {
          show: true,
          formatter(data: CallbackDataParams) {
            return numberFormat.format(Math.round((data.value as Array<number>)[i + 1]))
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
      series,
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
      throw new Error(`Unsupported time range: ${range as string}`)
  }
}