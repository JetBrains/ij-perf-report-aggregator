import { BarChart, BarSeriesOption } from "echarts/charts"
import { DataZoomInsideComponent, GridComponent, LegendComponent, MarkLineComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { XAXisOption } from "echarts/types/dist/shared"
import { MarkLine1DDataItemOption } from "echarts/types/src/component/marker/MarkLineModel"
import { ChartManagerHelper } from "../../components/common/ChartManagerHelper"
import { adaptToolTipFormatter } from "../../components/common/chart"
import { BarChartOptions } from "../../components/common/echarts"
import { durationAxisPointerFormatter, numberFormat } from "../../components/common/formatter"
import { DataDescriptor, DataManager, GroupedItems, markerNameToRangeTitle, getShortName, formatDuration } from "../DataManager"
import { ItemV20, UnitConverter } from "../data"
import { ChartManager } from "./ChartComponent"
import { buildTooltip, TooltipLineDescriptor } from "./tooltip"

use([ToolboxComponent, TooltipComponent, DataZoomInsideComponent, GridComponent, BarChart, MarkLineComponent, CanvasRenderer, LegendComponent])

type ChartDataItem = [string, number, ItemV20, string]

export class ActivityBarChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper
  private lastData!: ChartDataItem[]

  constructor(
    container: HTMLElement,
    private readonly dataProvider: (dataManager: DataManager) => GroupedItems,
    private readonly dataDescriptor: DataDescriptor
  ) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<BarChartOptions>({
      legend: {},
      grid: {
        left: 50,
        right: 5,
        containLabel: true,
      },
      toolbox: {
        feature: {
          dataZoom: {
            yAxisIndex: false,
          },
          saveAsImage: {},
        },
      },
      dataZoom: [{ type: "inside" }],
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "cross",
        },
        enterable: true,
        formatter: adaptToolTipFormatter((params) => {
          const info = params[0]
          const chartItem = info.data as ChartDataItem
          const item = chartItem[2]
          const lines: TooltipLineDescriptor[] = [
            { name: chartItem[0], main: true, value: durationAxisPointerFormatter(chartItem[1]) },
            { name: "range", value: `${formatDuration(item.s, this.dataDescriptor)}&ndash;${formatDuration(item.s + item.d, this.dataDescriptor)}` },
            { name: "thread", selectable: true, value: item.t, extraStyle: item.t === "edt" ? "color: orange" : "" },
          ]
          if (item.p !== undefined) {
            lines.push({ name: "plugin", selectable: true, value: item.p })
          }
          if ("od" in item) {
            lines.push({ name: "total duration", value: durationAxisPointerFormatter(this.dataDescriptor.unitConverter.convert(item.d)) })
          }
          return `${info.marker as string} ${buildTooltip(lines)}`
        }),
      },
      xAxis: {
        type: "category",
        axisLabel: {
          interval: 0,
          rotate: 30,
          fontSize: 10,
          color: (_value, index): string => {
            const item = this.lastData[index as number]
            if (item[2].t === "edt") {
              return "orange"
            } else {
              const currentOption = this.chart.chart.getOption() as BarChartOptions
              // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
              return (currentOption.xAxis as XAXisOption[])[0].axisLine!.lineStyle!.color as string
            }
          },
        },
        axisTick: {
          alignWithLabel: true,
        },
      },
      yAxis: {
        type: "value",
        axisPointer: {
          label: {
            formatter(data): string {
              return typeof data.value == "number" ? numberFormat.format(data.value) : ""
            },
          },
        },
      },
    })
    this.chart.enableZoomTool()
  }

  dispose(): void {
    this.chart.dispose()
  }

  render(dataManager: DataManager): void {
    const data: ChartDataItem[] = []
    const seriesNames: string[] = []
    // 10 ms
    const threshold = (this.dataDescriptor.threshold ?? 10) * this.dataDescriptor.unitConverter.factor
    for (const group of this.dataProvider(dataManager)) {
      let categoryAdded = false
      for (const item of group.items) {
        if (item.d < threshold) {
          continue
        }

        if (!categoryAdded) {
          categoryAdded = true
          seriesNames.push(group.category)
        }

        const chartItem: ChartDataItem = [
          this.dataDescriptor.shortenName === false ? item.n : getShortName(item.n),
          this.dataDescriptor.unitConverter.convert(item.od ?? item.d),
          item,
          group.category,
        ]
        data.push(chartItem)
      }
    }
    // sort by start
    data.sort((a, b) => a[2].s - b[2].s)

    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const axisLineColor = ((this.chart.chart.getOption() as BarChartOptions).xAxis as XAXisOption[])[0].axisLine!.lineStyle!.color

    const series = new Array<BarSeriesOption>(seriesNames.length)
    for (const [seriesIndex, seriesName] of seriesNames.entries()) {
      const seriesData: (ChartDataItem | null)[] = [...data]
      for (let i = 0; i < seriesData.length; i++) {
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        if (seriesData[i]![3] !== seriesName) {
          seriesData[i] = null
        }
      }
      series[seriesIndex] = {
        name: seriesName,
        type: "bar",
        barGap: "-100%",
        // emphasis: {focus: "series"},
        // well, ItemV20 is not expected as dimension value, but it is actually allowed and supported
        data: seriesData as never,
      }

      if (seriesIndex === 0) {
        series[seriesIndex].markLine = {
          symbol: "none",
          silent: true,
          lineStyle: {
            type: "dashed",
            color: axisLineColor,
          },
          data: createMarkLineData(dataManager, data, this.dataDescriptor.unitConverter),
        }
      }
    }

    this.lastData = data
    this.chart.chart.setOption<BarChartOptions>({
      xAxis: {
        data: data.map((it) => it[0]),
      },
      series,
    })
  }
}

function createMarkLineData(dataManager: DataManager, data: ChartDataItem[], unitConverter: UnitConverter): MarkLine1DDataItemOption[] {
  const markLineData: MarkLine1DDataItemOption[] = []
  for (const markerItem of dataManager.markerItems) {
    if (markerItem != null) {
      const item = data.find((it) => unitConverter.convert(it[2].s) >= markerItem.e)
      if (item == null) {
        continue
      }

      markLineData.push({
        xAxis: item[0] as never,
        label: {
          show: true,
          formatter: markerNameToRangeTitle.get(markerItem.n),
        },
      })
    }
  }
  return markLineData
}
