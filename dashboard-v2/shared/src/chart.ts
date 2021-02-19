import {
  BarChart, BarSeriesOption,
  LineChart, LineSeriesOption,
} from "echarts/charts"
import {
  BrushComponent,
  DatasetComponent,
  DataZoomComponent,
  GridComponent, GridComponentOption,
  LegendComponent,
  TitleComponent, TitleComponentOption,
  ToolboxComponent,
  TooltipComponent,
} from "echarts/components"
import * as echarts from "echarts/core"
import {
  CanvasRenderer,
} from "echarts/renderers"
import { DataQueryExecutorConfiguration } from "./dataQuery"

initEcharts()

export interface ChartConfigurator {
  configureChart(data: Array<Array<Array<string | number>>>, configuration: DataQueryExecutorConfiguration): ChartOptions
}

export type ChartOptions = echarts.ComposeOption<
  BarSeriesOption | LineSeriesOption | TitleComponentOption | GridComponentOption
>

export const timeFormat = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
  second: "numeric",
})

// https://github.com/apache/echarts/issues/8294
export const numberFormat = new Intl.NumberFormat(undefined, {
  maximumFractionDigits: 0,
})

// register the required components
function initEcharts(): void {
  echarts.use([
    ToolboxComponent, BrushComponent, DataZoomComponent, DatasetComponent,
    TitleComponent, LegendComponent,
    TooltipComponent, GridComponent,
    BarChart, LineChart,
    CanvasRenderer,
  ])
}