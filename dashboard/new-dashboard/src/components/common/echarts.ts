import { BarSeriesOption, LineSeriesOption, ScatterSeriesOption } from "echarts/charts"
import { DatasetComponentOption, GridComponentOption, TooltipComponentOption } from "echarts/components"
import { ComposeOption } from "echarts/core"

export type LineChartOptions = ComposeOption<TooltipComponentOption | LineSeriesOption | GridComponentOption | DatasetComponentOption>

export type ScatterChartOptions = ComposeOption<TooltipComponentOption | ScatterSeriesOption | GridComponentOption | DatasetComponentOption>

export type BarChartOptions = ComposeOption<TooltipComponentOption | BarSeriesOption | GridComponentOption>
