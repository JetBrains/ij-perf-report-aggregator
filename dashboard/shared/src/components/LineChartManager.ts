import { LineChart, ScatterChart } from "echarts/charts"
import { DatasetComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { registerTransform, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { CallbackDataParams } from "echarts/types/src/util/types"
import * as ecStat from "echarts-stat"
import { debounceTime } from "rxjs"
import { Ref } from "vue"
import { ChartManagerHelper } from "../ChartManagerHelper"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { adaptToolTipFormatter, timeFormat, ValueUnit } from "../chart"
import { refToObservable } from "../configurators/rxjs"
import { LineChartOptions } from "../echarts"
import { nsToMs, numberFormat } from "../formatter"
import { ChartToolTipManager } from "./ChartToolTipManager"

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
registerTransform(ecStat.transform.regression)
export type PopupTrigger = "item" | "axis" | "none"
export class LineChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement,
              private _dataQueryExecutor: DataQueryExecutor,
              dataZoom: Ref<boolean>,
              chartToolTipManager: ChartToolTipManager | null,
              valueUnit: ValueUnit,
              trigger: PopupTrigger = "axis") {
    this.chart = new ChartManagerHelper(container)
    const isMs = valueUnit == "ms"

    // https://github.com/apache/echarts/issues/2941
    let lastParams: CallbackDataParams[] | null = null
    if (chartToolTipManager != null) {
      this.chart.chart.getZr().on("click", event => {
        chartToolTipManager.showTooltip(lastParams, event.event)
      })
    }

    const isCompoundTooltip = chartToolTipManager == null
    this.chart.chart.setOption<LineChartOptions>({
      legend: {
        top: 0,
        left: 0,
        itemHeight: 3,
        itemWidth: 15,
        icon: "rect",
        type: "scroll",
      },
      animation: false,
      toolbox: {
        top: 20,
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
          snap: true
        },
        formatter: isCompoundTooltip ? undefined : adaptToolTipFormatter(params => {
          lastParams = params
          return null
        }),
        valueFormatter: isCompoundTooltip ? it => (numberFormat.format(isMs ? it as number : nsToMs(it as number)) + " ms") : undefined
      },
      xAxis: {
        type: "time",
        axisPointer: {
          snap: true,
          label: {
            formatter(data) {
              return timeFormat.format(data.value as number)
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
    }, {
      replaceMerge: ["legend"],
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
    this.unsubscribe = this.dataQueryExecutor.subscribe((data, configuration,isLoading) => {
      if(isLoading || data == null){
        this.chart.chart.showLoading("default", {showSpinner: false})
        return
      }
      this.chart.chart.hideLoading()
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