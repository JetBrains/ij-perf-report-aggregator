import { TreemapChart, TreemapSeriesOption } from "echarts/charts"
import { TooltipComponent } from "echarts/components"
import { ComposeOption, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"

export function useCanvasRenderer(): void {
  use([
    TooltipComponent, CanvasRenderer,
  ])
}

export type TreeMapChartOptions = ComposeOption<TreemapSeriesOption>

export function useTreeMapChart(): void {
  useCanvasRenderer()
  use([TreemapChart])
}