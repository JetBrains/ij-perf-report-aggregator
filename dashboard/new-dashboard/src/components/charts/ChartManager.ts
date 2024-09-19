import { LineChart, ScatterChart } from "echarts/charts"
import {
  DatasetComponent,
  DataZoomInsideComponent,
  DataZoomSliderComponent,
  GridComponent,
  LegendComponent,
  TitleComponent,
  ToolboxComponent,
  TooltipComponent,
} from "echarts/components"
import { EChartsType, init as initChart, throttle, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { ECBasicOption } from "echarts/types/dist/shared"
import { useDarkModeStore } from "../../shared/useDarkModeStore"

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
  TitleComponent,
])

export class ChartManager {
  readonly chart: EChartsType
  private readonly resizeObserver: ResizeObserver
  private readonly resizeContainer: HTMLElement

  constructor(
    public chartContainer: HTMLElement,
    resizeContainer: HTMLElement | null = document.body
  ) {
    this.chart = initChart(chartContainer, useDarkModeStore().darkMode ? "chalk" : "")

    this.resizeObserver = new ResizeObserver(
      throttle(() => {
        this.chart.resize()
      }, 300)
    )
    this.resizeContainer = resizeContainer ?? document.body
    this.resizeObserver.observe(this.resizeContainer)
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
