import { ECBasicOption } from "echarts/types/dist/shared"
import { CallbackDataParams } from "echarts/types/src/util/types"
import { DataQueryExecutorConfiguration } from "./dataQuery"

export const DEFAULT_LINE_CHART_HEIGHT = 340

export declare type ToolTipFormatter = (params: CallbackDataParams[]) => string | null

// natural sort of alphanumerical strings
export const collator = new Intl.Collator(undefined, { numeric: true, sensitivity: "base" })

export function adaptToolTipFormatter(formatter: ToolTipFormatter): (params: CallbackDataParams | CallbackDataParams[]) => string {
  return function (params: CallbackDataParams | CallbackDataParams[]): string {
    // function return type doesn't allow null, but actually it can be returned
    return formatter(Array.isArray(params) ? params : [params]) as never
  }
}

export type ChartType = "line" | "scatter"

export interface SymbolOptions {
  symbol?: ChartSymbolType
  symbolSize?: number
  showSymbol?: boolean
}

export interface ChartConfigurator {
  configureChart(data: (string | number)[][][], configuration: DataQueryExecutorConfiguration): ECBasicOption
}

export const timeFormat = new Intl.DateTimeFormat(undefined, {
  year: "numeric",
  month: "short",
  day: "numeric",
  hour: "numeric",
  minute: "numeric",
  second: "numeric",
})

export const chartDefaultStyle: ChartStyle = {
  barSeriesLabelPosition: "insideRight",
  valueUnit: "ms",
}

export type ValueUnit = "ms" | "ns" | "counter"

export type ChartSymbolType = "circle" | "rect" | "roundRect" | "triangle" | "diamond" | "pin" | "arrow" | "none"

export interface ChartStyle {
  barSeriesLabelPosition:
    | "left"
    | "right"
    | "top"
    | "bottom"
    | "inside"
    | "insideLeft"
    | "insideRight"
    | "insideTop"
    | "insideBottom"
    | "insideTopLeft"
    | "insideTopRight"
    | "insideBottomLeft"
    | "insideBottomRight"

  valueUnit: ValueUnit
}
