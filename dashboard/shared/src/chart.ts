import { BarChart, BarSeriesOption, LineChart, LineSeriesOption } from "echarts/charts"
import {
  BrushComponent,
  DatasetComponent,
  DataZoomComponent,
  GridComponent, GridComponentOption,
  LegendComponent,
  TitleComponent, TitleComponentOption,
  ToolboxComponent, TooltipComponentOption,
} from "echarts/components"
import { use, ComposeOption } from "echarts/core"
import { CallbackDataParams } from "echarts/types/src/util/types"
import { DataQueryExecutorConfiguration } from "./dataQuery"
import { useCanvasRenderer } from "./echarts"

export const DEFAULT_LINE_CHART_HEIGHT = 340

export declare type ToolTipFormatter = (params: CallbackDataParams[]) => string | null

export function adaptToolTipFormatter(formatter: ToolTipFormatter): (params: CallbackDataParams | CallbackDataParams[], _ticket: string) => string {
  return function (params: CallbackDataParams | CallbackDataParams[], _ticket: string): string {
    // function return type doesn't allow null, but actually it can be returned
    return formatter(Array.isArray(params) ? params : [params]) as never
  }
}

export interface ChartConfigurator {
  configureChart(data: Array<Array<Array<string | number>>>, configuration: DataQueryExecutorConfiguration): ChartOptions
}

export type ChartOptions = ComposeOption<
  TooltipComponentOption | BarSeriesOption | LineSeriesOption | TitleComponentOption | GridComponentOption
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
export function useLineAndBarCharts(): void {
  useCanvasRenderer()

  use([
    ToolboxComponent, BrushComponent, DataZoomComponent, DatasetComponent,
    TitleComponent, LegendComponent,
    GridComponent,
    BarChart, LineChart,
  ])
}