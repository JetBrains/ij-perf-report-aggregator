import { ChartManagerHelper } from "./ChartManagerHelper"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { ChartOptions, useLineAndBarCharts } from "./chart"

useLineAndBarCharts()

export class BarChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement, private readonly dataQueryExecutor: DataQueryExecutor) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<ChartOptions>({
      animation: false,
      legend: {},
      grid: {
        top: 30,
        left: 5,
        // place for bar label
        right: 25,
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