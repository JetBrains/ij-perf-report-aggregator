import { BarSeriesOption, LineSeriesOption, TreemapSeriesOption, SunburstSeriesOption, CustomSeriesOption } from "echarts/charts"
import {
  DatasetComponentOption,
  GridComponentOption,
  TooltipComponentOption,
} from "echarts/components"
import { ComposeOption } from "echarts/core"

export type TreeMapChartOptions = ComposeOption<TreemapSeriesOption>
export type SunburstChartOptions = ComposeOption<SunburstSeriesOption>

export type LineChartOptions = ComposeOption<
  TooltipComponentOption | LineSeriesOption | GridComponentOption | DatasetComponentOption
>

export type BarChartOptions = ComposeOption<
  TooltipComponentOption | BarSeriesOption | GridComponentOption
>

export type CustomChartOptions = ComposeOption<
  TooltipComponentOption | CustomSeriesOption | GridComponentOption
>