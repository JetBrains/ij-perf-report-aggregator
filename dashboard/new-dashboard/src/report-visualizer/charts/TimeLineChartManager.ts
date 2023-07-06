import { CustomChart, CustomSeriesOption } from "echarts/charts"
import { DataZoomInsideComponent, GridComponent, LegendComponent, MarkAreaComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { graphic, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { CustomSeriesRenderItemAPI, CustomSeriesRenderItemParams, XAXisOption } from "echarts/types/dist/shared"
import { MarkArea2DDataItemOption } from "echarts/types/src/component/marker/MarkAreaModel"
import { MarkLine1DDataItemOption } from "echarts/types/src/component/marker/MarkLineModel"
import { SeriesLabelOption } from "echarts/types/src/util/types"
import { ChartManagerHelper } from "../../components/common/ChartManagerHelper"
import { adaptToolTipFormatter, collator } from "../../components/common/chart"
import { BarChartOptions, CustomChartOptions } from "../../components/common/echarts"
import { durationAxisPointerFormatter, numberFormat } from "../../components/common/formatter"
import { DataDescriptor, DataManager, formatDuration, getShortName, GroupedItems } from "../DataManager"
import { ItemV20 } from "../data"
import { ChartManager } from "./ChartComponent"
import { buildTooltip, TooltipLineDescriptor } from "./tooltip"

use([ToolboxComponent, TooltipComponent, DataZoomInsideComponent, GridComponent, CustomChart, MarkAreaComponent, CanvasRenderer, LegendComponent])

type ChartDataItem = [string, number, number, number, ItemV20, string, string]

function getDuration(chartItem: ChartDataItem) {
  return chartItem[3]
}

const LABEL_THRESHOLD = 20

export class TimeLineChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper

  constructor(
    container: HTMLElement,
    private readonly dataProvider: (dataManager: DataManager) => GroupedItems,
    private readonly dataDescriptor: DataDescriptor
  ) {
    this.chart = new ChartManagerHelper(container)
  }

  private setInitialOption() {
    this.chart.chart.setOption<CustomChartOptions>(
      {
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: false,
            },
            saveAsImage: {},
          },
        },
        grid: {
          top: 20,
          left: 40,
          // place for item label (as we set max for xAxis)
          right: 80,
          bottom: 20,
          containLabel: true,
        },
        tooltip: {
          enterable: true,
          formatter: adaptToolTipFormatter((params) => {
            const info = params[0]
            const chartItem = info.data as ChartDataItem
            const item = chartItem[4]
            const lines: TooltipLineDescriptor[] = [
              { name: chartItem[5], main: true, value: durationAxisPointerFormatter(getDuration(chartItem)) },
              { name: "range", value: `${formatDuration(item.s, this.dataDescriptor)}&ndash;${formatDuration(item.s + item.d, this.dataDescriptor)}` },
              { name: "thread", selectable: true, value: item.t, extraStyle: item.t === "edt" ? "color: orange" : "" },
            ]
            if (item.p != undefined) {
              lines.push({ name: "plugin", selectable: true, value: item.p })
            }
            if ("od" in item) {
              lines.push({ name: "total duration", value: durationAxisPointerFormatter(this.dataDescriptor.unitConverter.convert(item.d)) })
            }
            return `${info.marker as string} ${buildTooltip(lines)}`
          }),
        },
        // chart is too tall - not convenient to use mouse for zoom, as mouse is better to use for scrolling
        dataZoom: [{ type: "inside", filterMode: "weakFilter", disabled: true }],
        xAxis: {
          scale: true,
          maxInterval: 1_000,
          axisLabel: {
            formatter(value: number) {
              // return Math.max(0, value - startTime) + " ms"
              return numberFormat.format(value)
            },
          },
        },
        yAxis: {
          type: "category",
          splitLine: {
            show: true,
            interval(_index: number, value: string | number) {
              return !(value as string).includes("__")
            },
          },
          axisLabel: {
            formatter(value: string | number, _index: number) {
              return (value as string).includes("__") ? "" : (value as string)
            },
          },
          axisTick: {
            show: false,
          },
        },
      },
      { notMerge: true }
    )
  }

  dispose(): void {
    this.chart.dispose()
  }

  render(dataManager: DataManager): void {
    const data = new Map<string, ChartDataItem[]>()
    const threshold = (this.dataDescriptor.threshold ?? 10) * this.dataDescriptor.unitConverter.factor
    for (const group of this.dataProvider(dataManager)) {
      const namePrefix = group.category === "service waiting" ? "wait for " : ""
      // eslint-disable-next-line @typescript-eslint/no-unnecessary-condition
      if (group.items == null) {
        console.error("No `items` for group", group)
        continue
      }

      for (const item of group.items) {
        if (item.d < threshold) {
          continue
        }

        let list = data.get(item.t)
        if (list == null) {
          list = []
          data.set(item.t, list)
        }

        const chartItem: ChartDataItem = [
          "",
          this.dataDescriptor.unitConverter.convert(item.s),
          this.dataDescriptor.unitConverter.convert(item.s + item.d),
          this.dataDescriptor.unitConverter.convert(item.d),
          item,
          this.dataDescriptor.shortenName === false ? item.n : `${namePrefix}${getShortName(item.n)}`,
          group.category,
        ]
        list.push(chartItem)
      }
    }

    const threadNames = [...data.keys()]
    threadNames.sort((a, b) => {
      const aW = getThreadOrderWeight(a)
      const bW = getThreadOrderWeight(b)
      const wR = aW - bW
      return wR === 0 ? collator.compare(a, b) : wR
    })

    // compute categories - for each thread maybe several categories as one service can include another one and we render it above each other
    const rowToItems = new Map<string, CustomSeriesOption & { label: SeriesLabelOption; itemStyle: unknown }>()
    // const labelThreshold = 20 * this.dataDescriptor.unitConverter.factor
    let minStart = Number.MAX_VALUE
    let maxEnd = 0
    const rowIndexThreshold = (this.dataDescriptor.rowIndexThreshold ?? 300) * this.dataDescriptor.unitConverter.factor
    const rowToEnd = new Map<number, RowInfo>()
    for (const threadName of threadNames) {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      const list = data.get(threadName)!
      list.sort((a, b) => a[4].s - b[4].s)
      if (minStart > list[0][4].s) {
        minStart = list[0][4].s
      }
      const last = list.at(-1) as ChartDataItem
      if (maxEnd < last[4].s + last[4].d) {
        maxEnd = last[4].s + last[4].d
      }

      let rowIndex = 0
      rowToEnd.clear()
      for (const chartItem of list) {
        const item = chartItem[4]
        if (rowToEnd.size > 0) {
          const newRowIndex = findRowIndex(rowIndex, rowToEnd, item, rowIndexThreshold)
          if (newRowIndex === -1) {
            // no place
            rowIndex++
          } else {
            rowIndex = newRowIndex
          }
        }

        const rowName = rowIndex === 0 ? threadName : `${threadName}__${rowIndex}`
        chartItem[0] = rowName
        let series = rowToItems.get(rowName)
        if (series == null) {
          series = {
            // same name to ensure that color will be the same
            name: threadName,
            id: rowName,
            type: "custom",
            renderItem: renderItem as never,
            encode: {
              x: [1, 2, 3],
              y: 0,
            },
            itemStyle: {
              // eslint-disable-next-line @typescript-eslint/ban-ts-comment
              //@ts-expect-error: https://github.com/apache/echarts/issues/16775
              color(value: { data: ChartDataItem; color: string }): string {
                const chartItem = value.data
                return chartItem[6] === "service waiting" ? "#FF0000" : value.color
              },
            },
            label: {
              show: true,
              position: "insideLeft",
              distance: 1,
              fontFamily: "monospace",
              formatter: adaptToolTipFormatter((params) => {
                const info = params[0]
                const chartItem = info.data as ChartDataItem
                return getDuration(chartItem) < LABEL_THRESHOLD ? "" : chartItem[5]
              }),
            },
            data: [],
          }
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          // @ts-expect-error
          rowToItems.set(rowName, series)
        }

        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        ;(series!.data as ChartDataItem[]).push(chartItem)

        const newEnd = item.s + item.d
        const info = rowToEnd.get(rowIndex)
        if (info === undefined) {
          rowToEnd.set(rowIndex, { end: newEnd, item })
        } else if (info.end < newEnd) {
          info.end = newEnd
          info.item = item
        }
      }
    }

    const series = [...rowToItems.values()]

    this.chart.chart.getDom().style.height = `${rowToItems.size * 24}px`
    // for unknown reasons `replaceMerge: ["series"]` doesn't work and data from previous report can be still rendered,
    // as workaround, create chart from scratch
    this.chart.chart.clear()
    this.setInitialOption()

    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const axisLineColor = ((this.chart.chart.getOption() as BarChartOptions).xAxis as XAXisOption[])[0].axisLine!.lineStyle!.color as string
    configureMarkAreas(dataManager, series, axisLineColor)

    this.chart.chart.setOption<CustomChartOptions>(
      {
        xAxis: {
          min: this.dataDescriptor.unitConverter.convert(minStart),
          max: this.dataDescriptor.unitConverter.convert(maxEnd),
        },
        series,
      },
      {
        replaceMerge: ["series"],
      }
    )
    this.chart.enableZoomTool()
    this.chart.chart.resize()
  }
}

function configureMarkAreas(dataManager: DataManager, series: CustomSeriesOption[], axisLineColor: string): void {
  const areaData: MarkArea2DDataItemOption[] = []
  for (const item of dataManager.isUnifiedItems ? dataManager.items : dataManager.data.prepareAppInitActivities) {
    if (!(item.n.endsWith(" async preloading") || item.n.endsWith(" sync preloading"))) {
      continue
    }

    const isAsync = item.n.includes("async")
    areaData.push([
      {
        name: item.n.replace(" service", "").replace("service", ""),
        label: {
          verticalAlign: isAsync ? "top" : "bottom",
        },
        itemStyle: {
          borderType: isAsync ? "dotted" : "solid",
          borderWidth: 1,
          borderColor: axisLineColor,
          color: "rgba(0, 0, 0, 0)",
        },
        xAxis: item.s,
      },
      { xAxis: item.s + item.d },
    ])
  }
  const lastSeries = series.at(-1) as CustomSeriesOption
  lastSeries.markArea =
    areaData.length === 0
      ? undefined
      : {
          silent: true,
          data: areaData,
        }

  const markLineData: MarkLine1DDataItemOption[] = []
  for (const item of dataManager.data.traceEvents) {
    if (item.name !== "splash shown" && item.name !== "project opened") {
      continue
    }

    markLineData.push({
      label: { formatter: item.name },
      xAxis: item.ts / 1000,
    })
  }
  lastSeries.markLine =
    markLineData.length === 0
      ? undefined
      : {
          symbol: "none",
          silent: true,
          lineStyle: {
            type: "dashed",
            color: axisLineColor,
          },
          data: markLineData,
        }
}

function renderItem(params: CustomSeriesRenderItemParams, api: CustomSeriesRenderItemAPI) {
  const categoryIndex = api.value(0)
  const start = api.coord([api.value(1), categoryIndex])
  const end = api.coord([api.value(2), categoryIndex])
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const height = (api.size!([0, 1]) as number[])[1] * 0.8

  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  const coordinateSystem = params.coordSys as { x: number; y: number; width: number; height: number }
  const rectShape = graphic.clipRectByRect(
    {
      x: start[0],
      y: start[1] - height / 2,
      width: end[0] - start[0],
      height,
    },
    {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      x: coordinateSystem.x,
      y: coordinateSystem.y,
      width: coordinateSystem.width,
      height: coordinateSystem.height,
    }
  )

  return {
    type: "rect",
    transition: ["shape"],
    shape: rectShape,
    style: api.style(),
  }
}

function getThreadOrderWeight(name: string): number {
  switch (name) {
    case "main":
      return 1
    case "idea main":
      return 2
    case "edt":
      return 3
    default:
      return 100
  }
}

function findRowIndex(rowIndex: number, rowToEnd: Map<number, RowInfo>, item: ItemV20, rowIndexThreshold: number): number {
  for (let i = rowIndex; i >= 0; i--) {
    const rowItem = rowToEnd.get(i)
    if (rowItem === undefined) {
      return -1
    }

    // for parallel activities ladder is used only to avoid text overlapping,
    // so two adjacent items are rendered in the same row if next one will not have a label
    // item.d < LABEL_THRESHOLD ||

    if (item.s + item.d < rowItem.end) {
      // sub item should be at higher level
      return -1
    } else if (item.d < LABEL_THRESHOLD || item.s - rowItem.end > rowIndexThreshold) {
      return i
    }
  }
  return -1
}

interface RowInfo {
  end: number
  item: ItemV20
}
