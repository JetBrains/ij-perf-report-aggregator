import { BarChart, BarSeriesOption, LineChart, LineSeriesOption, TreemapChart, SunburstChart, TreemapSeriesOption, SunburstSeriesOption } from "echarts/charts"
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
export type SunburstChartOptions = ComposeOption<SunburstSeriesOption>

export function useTreeMapChart(): void {
  useCanvasRenderer()
  use([TreemapChart])
}

export function useSunburstChart(): void {
  useCanvasRenderer()
  use([SunburstChart])
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