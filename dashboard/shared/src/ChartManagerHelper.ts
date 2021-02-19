import { EChartsType, throttle, init as initChart } from "echarts/core"
import { ChartOptions } from "./chart"

export class ChartManagerHelper {
  readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(container: HTMLElement) {
    this.chart = initChart(container)

    this.resizeObserver = new ResizeObserver(throttle(() => {
      this.chart.resize()
    }, 300))
    this.resizeObserver.observe(this.chart.getDom())
  }

  enableZoomTool(): void {
    // https://github.com/apache/echarts/issues/10274
    this.chart.dispatchAction({
      type: "takeGlobalCursor",
      key: "dataZoomSelect",
      dataZoomSelectActive: true,
    })
  }

  replaceDataSetAndSeries(options: ChartOptions): void {
    this.chart.setOption(options, {replaceMerge: ["dataset", "series"]})
  }

  dispose(): void {
    this.resizeObserver.unobserve(this.chart.getDom())
    this.chart.dispose()
  }
}