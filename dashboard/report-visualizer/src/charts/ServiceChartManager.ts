import { BarChart } from "echarts/charts"
import { DataZoomInsideComponent, GridComponent, MarkLineComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { XAXisOption } from "echarts/types/dist/shared"
import { MarkLine1DDataItemOption } from "echarts/types/src/component/marker/MarkLineModel"
import { TplFormatterParam } from "echarts/types/src/util/format"
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { adaptToolTipFormatter } from "shared/src/chart"
import { BarChartOptions, ChartOptions } from "shared/src/echarts"
import { durationAxisPointerFormatter, numberFormat } from "shared/src/formatter"
import { DataManager, markerNameToRangeTitle } from "../DataManager"
import { ItemV20 } from "../data"
import { ChartManager } from "./ChartManager"
import { buildTooltip, TooltipLineDescriptor } from "./tooltip"

use([ToolboxComponent, TooltipComponent, DataZoomInsideComponent, GridComponent, BarChart, MarkLineComponent, CanvasRenderer])

// 10 ms
const threshold = 10 * 1000

type ChartDataItem = [string, number, ItemV20]

function formatDurationInMicroSeconds(microseconds: number): string {
  return numberFormat.format(microseconds / 1000)
}

export class ServiceChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper

  private lastData!: Array<ChartDataItem>

  constructor(container: HTMLElement) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<BarChartOptions>({
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
      dataZoom: [{type: "inside"}],
      tooltip: {
        trigger: "axis",
        axisPointer: {
          type: "cross",
        },
        enterable: true,
        formatter: adaptToolTipFormatter(params => {
          const info = params[0]
          const chartItem = info.data as ChartDataItem
          const item = chartItem[2]
          const lines: Array<TooltipLineDescriptor> = [
            {name: chartItem[0], main: true, value: durationAxisPointerFormatter(chartItem[1])},
            {name: "range", value: `${(formatDurationInMicroSeconds(item.s))}&ndash;${formatDurationInMicroSeconds(item.s + item.d)}`},
            {name: "thread", selectable: true, value: item.t, extraStyle: item.t === "edt" ? "color: orange" : ""},
            {name: "plugin", selectable: true, value: item.p},
          ]
          if (item.od !== undefined) {
            lines.push({name: "total duration", value: durationAxisPointerFormatter(item.d / 1000)})
          }
          return `${info.marker as string} ${buildTooltip(lines)}`
        })
      },
      xAxis: {
        type: "category",
        axisLabel: {
          interval: 0,
          rotate: 30,
          fontSize: 10,
          color: (_value, index): string => {
            const item = this.lastData[index as number]
            const currentOption: BarChartOptions = this.chart.chart.getOption() as BarChartOptions
            // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
            return item[2].t === "edt" ? "orange" : (currentOption.xAxis as Array<XAXisOption>)[0].axisLine!.lineStyle!.color as string
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
            formatter(data: TplFormatterParam) {
              return numberFormat.format(data["value"])
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
    const list = dataManager.getServiceItems()
    const data: Array<ChartDataItem> = []
    this.lastData = data
    for (const item of list) {
      if (item.d < threshold) {
        continue
      }

      data.push([
        getShortName(item.n),
        (item.od ?? item.d) / 1000,
        item,
      ])
    }
    // sort by start
    data.sort((a, b) => a[2].s - b[2].s)
    const markLineData: Array<MarkLine1DDataItemOption> = []
    for (const markerItem of dataManager.markerItems) {
      if (markerItem != null) {
        const item = data.find(it => (it[2].s / 1000) >= markerItem.end)
        if (item != null) {
          // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
          markLineData.push({
            xAxis: item[0] as never,
            label: {
              show: true,
              formatter: markerNameToRangeTitle.get(markerItem.name),
            },
          })
        }
      }
    }

    this.chart.chart.setOption<ChartOptions>({
      series: [
        {
          type: "bar",
          dimensions: [
            {name: "name", type: "ordinal"},
            {name: "duration", type: "number"},
            {name: "start", type: "number"},
            {name: "thread", type: "ordinal"},
          ],
          // well, ItemV20 is not expected as dimension value, but it is actually allowed and supported
          data: data as never,
          markLine: {
            symbol: "none",
            silent: true,
            lineStyle: {
              type: "dashed",
            },
            data: markLineData,
          },
        },
      ],
    })
  }
}

function getShortName(qualifiedName: string): string {
  const lastDotIndex = qualifiedName.lastIndexOf(".")
  return lastDotIndex < 0 ? qualifiedName : qualifiedName.substring(lastDotIndex + 1)
}