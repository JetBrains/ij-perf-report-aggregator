import { BarChart } from "echarts/charts"
import { DatasetComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { ChartManagerHelper } from "./ChartManagerHelper"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { ChartStyle } from "./chart"
import { LineChartOptions } from "./echarts"
import { durationFormatterInOneWord, nsToMs } from "./formatter"

use([DatasetComponent, ToolboxComponent, TooltipComponent, GridComponent, BarChart, LegendComponent, CanvasRenderer])

export class BarChartManager {
  private readonly chart: ChartManagerHelper

  private readonly unsubscribe: () => void

  constructor(container: HTMLElement, dataQueryExecutor: DataQueryExecutor, chartStyle: ChartStyle) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<LineChartOptions>({
      animation: false,
      legend: {},
      grid: {
        top: 30,
        left: 5,
        // place for bar label
        right: 20,
        bottom: 5,
        containLabel: true,
      },
      xAxis: {
        type: "value",
        axisLabel: {
          hideOverlap: true,
          formatter: chartStyle.valueUnit == "ms" ? durationFormatterInOneWord : (it: number) => durationFormatterInOneWord(nsToMs(it)),
        },
      },
      yAxis: {
        type: "category",
      },
    })

    this.unsubscribe = dataQueryExecutor.subscribe((data, configuration, isLoading) => {
      if (isLoading || data == null) {
        this.chart.chart.showLoading("default", { showSpinner: false })
        return
      }
      this.chart.chart.hideLoading()

      for (const it of configuration.getChartConfigurators()) {
        it.configureChart(data, configuration)
          .then((options) => {
            this.chart.replaceDataSetAndSeries(options)
          })
          .catch((error: unknown) => {
            console.error(error)
          })
      }
    })
  }

  dispose(): void {
    this.unsubscribe()
    this.chart.dispose()
  }
}
