import { LineSeriesOption } from "echarts/charts"
import { DatasetOption, ECBasicOption } from "echarts/types/dist/shared"
import { CallbackDataParams, OptionDataItem, OptionSourceData, ScaleDataValue } from "echarts/types/src/util/types"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ref } from "vue"
import { ChartManager } from "./ChartManager"

// LabelFormatterParams isn't exported from lib
interface LabelFormatterParams {
  value: ScaleDataValue
  axisDimension: string
  axisIndex: number
  seriesData: CallbackDataParams[]
}

const dateFormatter = new Intl.DateTimeFormat(undefined, {
  month: "short",
  day: "numeric"
})

export class AggregationChartVM {
  average = ref(0)

  private chartManager?: ChartManager

  constructor(
    private readonly query: DataQueryExecutor,
    private readonly color: string = "#4B84EE",
  ) {}

  initChart(element: HTMLElement, resizeContainer?: HTMLElement): () => void {
    this.chartManager = new ChartManager(element, resizeContainer)

    this.chartManager.chart.setOption({
      legend: {
        show: false,
      },
      animation: false,
      grid: {
        top: 0,
        bottom: 10,
        left: 5,
        right: 5,
      },
      tooltip: {
        show: false,
        trigger: "axis",
        enterable: true,
      },
      xAxis: {
        type: "time",
        show: false,
        triggerTooltip: false,
        axisPointer: {
          show: true,
          type: "line",
          snap: true,
          label: {
            formatter(params: LabelFormatterParams) {
              const series = params.seriesData[0]
              const [date, durationMs] = series.data as OptionDataItem[]
              const dateLabel = dateFormatter.format(new Date(date as string))
              const durationLabel = `${Math.round(Number(durationMs))} ms`
              return `${dateLabel}, ${durationLabel}`
            },
            color: "#6B7280",
            backgroundColor: "transparent",
            padding: [5, 0, 0, 0]
          }
        }
      },
      yAxis: {
        type: "value",
        show: false,
      }
    })

    const unsubscribe = this.dataSubscribe()

    return () => {
      this.chartManager?.dispose()
      unsubscribe()
    }
  }

  private dataSubscribe(): () => void {
    return this.query.subscribe(
      (data, configuration, isLoading: boolean) => {
        if(isLoading || data == null){
          this.chartManager?.chart.showLoading("default", {showSpinner: false})
          return
        }
        this.chartManager?.chart.hideLoading()
        const options = configuration.chartConfigurator.configureChart(data, configuration)

        this.updateChartData(options)
      })
  }

  private calculateAverage(values: number[]) {
    if (values.length === 0) {
      return 0
    }

    const sum = values.reduce((acc, value) => acc + value, 0)
    const average = sum / values.length

    return average % 1 === 0 ? average : Number(average.toFixed(1))
  }


  private updateChartData(options: ECBasicOption) {
    if (options["series"]) {
      options["series"] = (options["series"] as LineSeriesOption[]).map(item => ({
        ...item,
        showSymbol: false,
        color: this.color,
      }))
    }

    this.chartManager?.chart.setOption(options, { replaceMerge: ["dataset", "series"] })

    const dataset = options["dataset"] as DatasetOption[]
    const [_, values] = dataset[0].source as OptionSourceData[]

    this.average.value = this.calculateAverage(values as number[])
  }
}