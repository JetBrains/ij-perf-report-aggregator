import { EChartsType, throttle, init as initChart } from "echarts/core"
import { ECBasicOption } from "echarts/types/dist/shared"

export class ChartManagerHelper {
  readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(
    chartContainer: HTMLElement,
    private readonly resizeContainer: HTMLElement = document.body
  ) {
    this.chart = initChart(chartContainer)

    this.resizeObserver = new ResizeObserver(
      throttle(() => {
        this.chart.resize()
      }, 300)
    )
    this.resizeObserver.observe(resizeContainer)
  }

  replaceDataSetAndSeries(options: ECBasicOption): void {
    this.chart.setOption(options, { replaceMerge: ["dataset", "series"] })
  }

  dispose(): void {
    this.resizeObserver.unobserve(this.resizeContainer)
    this.chart.dispose()
  }
}
