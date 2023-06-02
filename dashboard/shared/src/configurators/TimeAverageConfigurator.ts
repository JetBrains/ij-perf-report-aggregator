import { LineSeriesOption, ScatterSeriesOption } from "echarts/charts"
import { DatasetOption } from "echarts/types/dist/shared"
import { ChartConfigurator } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { LineChartOptions } from "../echarts"
import { durationAxisPointerFormatter, isDurationFormatterApplicable, numberAxisLabelFormatter } from "../formatter"

export class TimeAverageConfigurator implements DataQueryConfigurator, ChartConfigurator {

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    query.addDimension({
      n: "t",
      //2018-04-10T00:00:00
      sql: "toYYYYMMDD(generated_time)",
    })
    query.addField({n: "measures", subName: "value"})
    query.aggregator = "avg"
    query.order = "t"

    configuration.chartConfigurator = this

    return true
  }

  createObservable() {
    return null
  }

  configureChart(dataList: (string | number)[][][], configuration: DataQueryExecutorConfiguration): LineChartOptions {
    const series = new Array<LineSeriesOption | ScatterSeriesOption>()
    let useDurationFormatter = true

    const dataset: DatasetOption[] = []

    for (let dataIndex = 0, n = dataList.length; dataIndex < n; dataIndex++) {
      const measureName = configuration.measureNames[dataIndex]
      const seriesName = "aggregate"
      const seriesData = dataList[dataIndex]
      seriesData[0] = seriesData[0].map((value, _n, _a) => {
        if (typeof value == "number") {
          value = value.toString()
          return value.slice(0, 4) + "-" + value.slice(4, 6) + "-" + value.slice(6)
        }
        return ""
      })

      if (seriesData.length > 2) {
        // we take only the last type of the metric since it's not clear how to show different types and last type helps to change the type if necessary
        const type = seriesData[2].at(-1)
        if (type === "c") {
          useDurationFormatter = false
        }
        else if (type === "d") {
          useDurationFormatter = true
        }
      }

      series.push({
        // formatter is detected by measure name - that's why series id is specified (see usages of seriesId)
        id: measureName === seriesName ? seriesName : `${measureName}@${seriesName}`,
        name: seriesName,
        type: "line",
        showSymbol: seriesData[0].length < 100,
        // 10 is a default value for scatter (undefined doesn't work to unset)
        symbolSize: Math.min(800 / seriesData[0].length, 9),
        symbol: "circle",
        legendHoverLink: true,
        // applicable only for line chart
        sampling: "lttb",
        seriesLayoutBy: "row",
        datasetIndex: dataIndex,
        dimensions: [{name: "date", type: "time"}, {name: seriesName, type: "int"}],
      })


      if (useDurationFormatter && !isDurationFormatterApplicable(measureName)) {
        useDurationFormatter = false
      }

      dataset.push({
        source: seriesData,
        sourceHeader: false,
      })
    }

    const formatter: (valueInMs: number) => string = useDurationFormatter ? durationAxisPointerFormatter : numberAxisLabelFormatter
    return {
      dataset,
      yAxis: {
        axisLabel: {
          formatter,
        },
        axisPointer: {
          label: {
            formatter(data): string {
              return formatter(data.value as number)
            },
          },
        },
      },
      series: series as LineSeriesOption,
    }
  }
}