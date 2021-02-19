import { EChartsType, init as initChart, throttle } from "echarts/core"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { ChartOptions } from "./chart"

export class BarChartManager {
  private readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(container: HTMLElement, private readonly dataQueryExecutor: DataQueryExecutor) {
    this.chart = initChart(container)

    this.resizeObserver = new ResizeObserver(throttle(() => {
      this.chart.resize()
    }, 300))
    this.resizeObserver.observe(this.chart.getDom())

    this.chart.setOption<ChartOptions>({
      animation: false,
      legend: {},
      // tooltip: {
      //   trigger: "axis",
      //   // position: {bottom: 0, left: 0},
      //   axisPointer: {
      //     type: "shadow",
      //   },
      // },
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

    // let isFirstLoad = true
    dataQueryExecutor.setListener((data, configuration) => {
      const option = configuration.chartConfigurator.configureChart(data, configuration)
      // if (isFirstLoad) {
      //   isFirstLoad = false
      //   nextTick(() => {
      //     this.chart.setOption({animation: true})
      //   })
      // }
      this.chart.setOption<ChartOptions>(option, {replaceMerge: ["dataset", "series"]})
    })
  }

  dispose(): void {
    this.resizeObserver.unobserve(this.chart.getDom())
    this.chart.dispose()
    this.dataQueryExecutor.setListener(null)
  }
}