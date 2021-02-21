import * as echarts from "echarts/core"

import {
  BarChart,
  LineChart,
} from "echarts/charts"

import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  DataZoomComponent,
  LegendComponent,
  ToolboxComponent,
  BrushComponent
} from "echarts/components"

import {
  CanvasRenderer,
} from "echarts/renderers"

// Register the required components
export function initEcharts(): void {
  echarts.use([
    ToolboxComponent, BrushComponent, DataZoomComponent, DatasetComponent,
    TitleComponent, LegendComponent,
    TooltipComponent, GridComponent,
    BarChart, LineChart,
    CanvasRenderer,
  ])
}