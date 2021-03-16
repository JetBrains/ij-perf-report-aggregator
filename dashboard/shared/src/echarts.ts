import { BarSeriesOption, LineSeriesOption, TreemapSeriesOption, SunburstSeriesOption } from "echarts/charts"
import {
  GridComponentOption,
  TooltipComponentOption,
} from "echarts/components"
import { ComposeOption } from "echarts/core"

export type TreeMapChartOptions = ComposeOption<TreemapSeriesOption>
export type SunburstChartOptions = ComposeOption<SunburstSeriesOption>

export type ChartOptions = ComposeOption<
  TooltipComponentOption | BarSeriesOption | LineSeriesOption | GridComponentOption
>

export type BarChartOptions = ComposeOption<
  TooltipComponentOption | BarSeriesOption | GridComponentOption
>