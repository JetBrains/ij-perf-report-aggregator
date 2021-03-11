import { TplFormatterParam } from "echarts/types/src/util/format"
import { OrdinalRawValue } from "echarts/types/src/util/types"
import humanizeDuration, { HumanizerOptions } from "humanize-duration"
import { watch , Ref } from "vue"
import { ChartManagerHelper } from "./ChartManagerHelper"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { adaptToolTipFormatter, timeFormat, ToolTipFormatter } from "./chart"
import { DataQuery } from "./dataQuery"
import { ChartOptions, useLineAndBarCharts } from "./echarts"
import { debounceSync } from "./util/debounce"

export type ChartTooltipLinkProvider = (name: string, query: DataQuery) => string

const durationFormatOptions: HumanizerOptions = {
  language: "shortEn",
  round: true,
  units: ["y", "mo", "w", "d", "h", "m", "s", "ms"],
  languages: {
    shortEn: {
      y: () => "y",
      mo: () => "mo",
      w: () => "w",
      d: () => "d",
      h: () => "h",
      m: () => "m",
      s: () => "s",
      ms: () => "ms",
    },
  },
}
const shortEnglishHumanizer = humanizeDuration.humanizer(durationFormatOptions)
export const axisDurationFormatter = humanizeDuration.humanizer({
  ...durationFormatOptions,
  delimiter: " "
})

export function formatDuration(value: number): string {
  return shortEnglishHumanizer(value)
}

const dataZoomConfig = [
  // https://echarts.apache.org/en/option.html#dataZoom-inside
  // type inside means that mouse maybe used to zoom.
  {type: "inside"},
  {},
]

useLineAndBarCharts()

export class LineChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement,
              private _dataQueryExecutor: DataQueryExecutor,
              dataZoom: Ref<boolean>,
              tooltipFormatter: ToolTipFormatter) {
    this.chart = new ChartManagerHelper(container)

    this.chart.chart.setOption<ChartOptions>({
      legend: {},
      animation: false,
      toolbox: {
        feature: {
          dataZoom: {
            yAxisIndex: false,
          },
          saveAsImage: {},
        },
      },
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "cross",
        },
        formatter: adaptToolTipFormatter(tooltipFormatter),
      },
      xAxis: {
        type: "time",
        axisPointer: {
          label: {
            formatter(data: TplFormatterParam) {
              return timeFormat.format(data["value"])
            },
          },
        },
      },
      yAxis: {
        type: "value",
        axisLabel: {
          formatter(value: OrdinalRawValue, _index: number): string {
            return axisDurationFormatter(value as number)
          },
        },
        axisPointer: {
          label: {
            formatter(data: TplFormatterParam): string {
              return formatDuration(data["value"])
            },
          },
        },
      },
      dataZoom: dataZoom.value ? dataZoomConfig : undefined,
    })

    this.chart.enableZoomTool()
    this.setDataListener()
    watch(dataZoom, debounceSync(() => {
      this.chart.chart.setOption({
        dataZoom: dataZoom.value ? dataZoomConfig : undefined,
      })
    }))
  }

  private setDataListener() {
    this.dataQueryExecutor.setListener((data, configuration) => {
      this.chart.replaceDataSetAndSeries(configuration.chartConfigurator.configureChart(data, configuration))
      // console.log(JSON.stringify(this.chart.getOption(), null, 2))
    })
  }

  get dataQueryExecutor(): DataQueryExecutor {
    return this._dataQueryExecutor
  }

  set dataQueryExecutor(newDataQueryExecutor: DataQueryExecutor) {
    this._dataQueryExecutor.setListener(null)
    this._dataQueryExecutor = newDataQueryExecutor
    this.setDataListener()
  }

  dispose(): void {
    this.chart.dispose()
    this.dataQueryExecutor.setListener(null)
  }
}