import { BarChart, BarSeriesOption, LineChart, LineSeriesOption, TreemapChart, TreemapSeriesOption } from "echarts/charts"
import {
  BrushComponent,
  DatasetComponent,
  DataZoomComponent,
  GridComponent, GridComponentOption,
  LegendComponent,
  TitleComponent, TitleComponentOption,
  ToolboxComponent,
  TooltipComponent,
  TooltipComponentOption,
} from "echarts/components"
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

// register the required components
export function useLineAndBarCharts(): void {
  useCanvasRenderer()

  use([
    ToolboxComponent, BrushComponent, DataZoomComponent, DatasetComponent,
    TitleComponent, LegendComponent,
    GridComponent,
    BarChart, LineChart,
  ])
}

export type ChartOptions = ComposeOption<
  TooltipComponentOption | BarSeriesOption | LineSeriesOption | TitleComponentOption | GridComponentOption
>