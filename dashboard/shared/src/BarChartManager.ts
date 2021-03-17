import { BarChart } from "echarts/charts"
import { DatasetComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { ChartManagerHelper } from "./ChartManagerHelper"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { LineChartOptions } from "./echarts"

use([DatasetComponent, ToolboxComponent, TooltipComponent, GridComponent, BarChart, LegendComponent, CanvasRenderer])

export class BarChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement, private readonly dataQueryExecutor: DataQueryExecutor) {
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
      },
      yAxis: {
        type: "category",
      },
    })

    dataQueryExecutor.setListener((data, configuration) => {
      this.chart.replaceDataSetAndSeries(configuration.chartConfigurator.configureChart(data, configuration))
    })
  }

  dispose(): void {
    this.chart.dispose()
    this.dataQueryExecutor.setListener(null)
  }
}