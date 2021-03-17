import { LineChart } from "echarts/charts"
import { DatasetComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { TplFormatterParam } from "echarts/types/src/util/format"
import { watch , Ref } from "vue"
import { ChartManagerHelper } from "./ChartManagerHelper"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { adaptToolTipFormatter, timeFormat, ToolTipFormatter } from "./chart"
import { LineChartOptions } from "./echarts"
import { debounceSync } from "./util/debounce"

const dataZoomConfig = [
  // https://echarts.apache.org/en/option.html#dataZoom-inside
  // type inside means that mouse maybe used to zoom.
  {type: "inside"},
  {},
]

use([DatasetComponent, ToolboxComponent, TooltipComponent, GridComponent, LineChart, LegendComponent, CanvasRenderer])

export class LineChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement,
              private _dataQueryExecutor: DataQueryExecutor,
              dataZoom: Ref<boolean>,
              tooltipFormatter: ToolTipFormatter) {
    this.chart = new ChartManagerHelper(container)

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