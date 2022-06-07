import * as ecStat from "echarts-stat"
import { LineChart, ScatterChart } from "echarts/charts"
import { DatasetComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { registerTransform, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { debounceTime } from "rxjs"
import { Ref } from "vue"
import { ChartManagerHelper } from "../ChartManagerHelper"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { adaptToolTipFormatter, timeFormat, ToolTipFormatter, ValueUnit } from "../chart"
import { refToObservable } from "../configurators/rxjs"
import { LineChartOptions } from "../echarts"
import { nsToMs, numberFormat } from "../formatter"

const dataZoomConfig = [
  // https://echarts.apache.org/en/option.html#dataZoom-inside
  // type inside means that mouse maybe used to zoom.
  {type: "inside"},
  {},
]

use([DatasetComponent, ToolboxComponent, TooltipComponent, GridComponent, LineChart, ScatterChart, LegendComponent, CanvasRenderer])

// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
// eslint-disable-next-line @typescript-eslint/no-unsafe-argument,@typescript-eslint/no-unsafe-member-access
registerTransform(ecStat["transform"].regression)
export type PopupTrigger = "item" | "axis" | "none"
export class LineChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement,
              private _dataQueryExecutor: DataQueryExecutor,
              dataZoom: Ref<boolean>,
              tooltipFormatter: ToolTipFormatter | null,
              valueUnit: ValueUnit,
              trigger: PopupTrigger = "axis") {
    this.chart = new ChartManagerHelper(container)
    const isMs = valueUnit == "ms"
    this.chart.chart.setOption<LineChartOptions>({
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
        show: true,
        trigger,
        enterable: true,
        // select text in tooltip
        extraCssText: "user-select: text",
        axisPointer: {
          type: "cross",
        },
        formatter: tooltipFormatter == null ? undefined : adaptToolTipFormatter(tooltipFormatter),
        valueFormatter: tooltipFormatter == null ? it => (numberFormat.format(isMs ? it as number : nsToMs(it as number)) + " ms") : undefined
      },
      xAxis: {
        type: "time",
        axisPointer: {
          snap: true,
          label: {
            formatter(data) {
              return timeFormat.format(data["value"] as number)
            },
          },
        },
      },
      yAxis: {
        type: "value",
        axisPointer: {
          snap: true,
        }
      },
      dataZoom: dataZoom.value ? dataZoomConfig : undefined,
    })

    this.chart.enableZoomTool()
    this.subscribe()
    refToObservable(dataZoom)
      .pipe(debounceTime(100))
      .subscribe(value => {
        this.chart.chart.setOption({
          dataZoom: value ? dataZoomConfig : undefined,
        })
      })
  }

  private unsubscribe: () => void = () => {
    return
  }

  private subscribe() {
    this.unsubscribe()
    this.unsubscribe = this.dataQueryExecutor.subscribe((data, configuration) => {
      this.chart.replaceDataSetAndSeries(configuration.chartConfigurator.configureChart(data, configuration))
    })
  }

  get dataQueryExecutor(): DataQueryExecutor {
    return this._dataQueryExecutor
  }

  set dataQueryExecutor(newDataQueryExecutor: DataQueryExecutor) {
    this._dataQueryExecutor = newDataQueryExecutor
    this.subscribe()
  }

  dispose(): void {
    this.unsubscribe()
    this.chart.dispose()
  }
}