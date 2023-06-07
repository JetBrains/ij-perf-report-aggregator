import { LineChart, ScatterChart } from "echarts/charts"
import { DatasetComponent, DataZoomInsideComponent, DataZoomSliderComponent, GridComponent, LegendComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { EChartsType, throttle, use, init as initChart } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { ECBasicOption } from "echarts/types/dist/shared"

use([
  DatasetComponent,
  ToolboxComponent,
  TooltipComponent,
  GridComponent,
  LineChart,
  ScatterChart,
  LegendComponent,
  CanvasRenderer,
  DataZoomInsideComponent,
  DataZoomSliderComponent,
])

export class ChartManager {
  readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver

  constructor(chartContainer: HTMLElement, private resizeContainer: HTMLElement = document.body) {
    this.chart = initChart(chartContainer)

    this.resizeObserver = new ResizeObserver(
      throttle(() => {
        this.chart.resize()
      }, 300)
    )
    this.resizeObserver.observe(resizeContainer)
  }

  updateChart(options: ECBasicOption): void {
    this.chart.setOption(options, {
      replaceMerge: ["dataset", "series"],
    })
  }

  dispose(): void {
    this.resizeObserver.unobserve(this.resizeContainer)
    this.chart.dispose()
  }
}
