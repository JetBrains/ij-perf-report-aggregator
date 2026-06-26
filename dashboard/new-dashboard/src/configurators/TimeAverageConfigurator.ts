import { LineSeriesOption, ScatterSeriesOption } from "echarts/charts"
import { DatasetOption } from "echarts/types/dist/shared"
import { ChartConfigurator } from "../components/common/chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { LineChartOptions } from "../components/common/echarts"
import { formatMeasureValue, MeasureUnit, reduceToAxisUnit, resolveMeasureUnit } from "../components/common/formatter"

export class TimeAverageConfigurator implements DataQueryConfigurator, ChartConfigurator {
  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    query.addDimension({
      n: "t",
      //2018-04-10T00:00:00
      sql: "toYYYYMMDD(generated_time)",
    })
    query.addField({ n: "measures", subName: "value" })
    query.aggregator = "avg"
    query.order = "t"

    configuration.addChartConfigurator(this)

    return true
  }

  createObservable() {
    return null
  }

  configureChart(dataList: (string | number)[][][], configuration: DataQueryExecutorConfiguration): Promise<LineChartOptions> {
    const series = new Array<LineSeriesOption | ScatterSeriesOption>()

    const dataset: DatasetOption[] = []

    const measureUnits: MeasureUnit[] = []
    for (let dataIndex = 0, n = dataList.length; dataIndex < n; dataIndex++) {
      const measureName = configuration.measureNames[dataIndex]
      const seriesName = "aggregate"
      const seriesData = dataList[dataIndex]
      if (seriesData.length === 0) {
        continue
      }
      seriesData[0] = seriesData[0].map((value, _n, _a) => {
        if (typeof value == "number") {
          value = value.toString()
          return value.slice(0, 4) + "-" + value.slice(4, 6) + "-" + value.slice(6)
        }
        return ""
      })

      // we take only the last type of the metric since it is not clear how to show different types
      const storedType = seriesData.length > 2 ? (seriesData[2].at(-1) as string) : undefined

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
        dimensions: [
          { name: "date", type: "time" },
          { name: seriesName, type: "int" },
        ],
      })

      measureUnits.push(resolveMeasureUnit(measureName, { storedType }))

      dataset.push({
        source: seriesData,
        sourceHeader: false,
      })
    }

    const axisUnit = reduceToAxisUnit(measureUnits)
    const formatter: (value: number) => string = (value) => formatMeasureValue(value, axisUnit)
    return Promise.resolve({
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
    })
  }
}
