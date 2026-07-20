import { ECBasicOption } from "echarts/types/dist/shared"
import { DataQueryExecutorConfiguration } from "./dataQuery"

export const DEFAULT_LINE_CHART_HEIGHT = 340

// natural sort of alphanumerical strings
export const collator = new Intl.Collator(undefined, { numeric: true, sensitivity: "base" })

export type ChartType = "line" | "scatter"

export interface SymbolOptions {
  symbol?: ChartSymbolType
  symbolSize?: number
  showSymbol?: boolean
}

export interface ChartConfigurator {
  configureChart(data: (string | number)[][][], configuration: DataQueryExecutorConfiguration): Promise<ECBasicOption>
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

// "auto" (the default) infers the unit from each series' stored metric type; "ms"/"ns"/"counter"
// are explicit overrides that win over the stored type (some counts/durations are mis-typed).
export type ValueUnit = "ms" | "ns" | "counter" | "auto"

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
