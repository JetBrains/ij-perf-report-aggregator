/* eslint-disable unicorn/prefer-ternary */
import { CustomChart, CustomSeriesOption } from "echarts/charts"
import { DataZoomInsideComponent, GridComponent, LegendComponent, MarkAreaComponent, ToolboxComponent, TooltipComponent } from "echarts/components"
import { graphic, use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { XAXisOption, CustomSeriesRenderItemAPI, CustomSeriesRenderItemParams, CustomSeriesRenderItemReturn } from "echarts/types/dist/shared"
import { MarkArea2DDataItemOption } from "echarts/types/src/component/marker/MarkAreaModel"
import { MarkLine1DDataItemOption } from "echarts/types/src/component/marker/MarkLineModel"
import { SeriesLabelOption } from "echarts/types/src/util/types"
import { ChartManagerHelper } from "../../components/common/ChartManagerHelper"
import { adaptToolTipFormatter } from "../../components/common/chart"
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

export interface TimeLineChartConfiguration {
  readonly hasParent: boolean
  readonly dataProvider: (dataManager: DataManager) => GroupedItems
  readonly groupData: (item: ItemV20, category: string) => string
  readonly sortGroups: (names: string[]) => void
  readonly dataDescriptor: DataDescriptor
}

export class TimeLineChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper

  constructor(
    container: HTMLElement,
    private readonly config: TimeLineChartConfiguration
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
            const dataDescriptor = this.config.dataDescriptor
            const lines: TooltipLineDescriptor[] = [
              { name: chartItem[5], main: true, value: durationAxisPointerFormatter(getDuration(chartItem)) },
              { name: "range", value: `${formatDuration(item.s, dataDescriptor)}&ndash;${formatDuration(item.s + item.d, dataDescriptor)}` },
              { name: "thread", selectable: true, value: item.t, extraStyle: item.t === "edt" ? "color: orange" : "" },
            ]
            if (item.p != undefined) {
              lines.push({ name: "plugin", selectable: true, value: item.p })
            }
            if ("od" in item) {
              lines.push({ name: "total duration", value: durationAxisPointerFormatter(dataDescriptor.unitConverter.convert(item.d)) })
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
              // return value as string
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

  // prettier-ignore
  render(dataManager: DataManager): void {
    const data = new Map<string, ChartDataItem[]>()
    const config = this.config
    const dataDescriptor = config.dataDescriptor
    const threshold = (dataDescriptor.threshold ?? 10) * dataDescriptor.unitConverter.factor

    const activityList = config.dataProvider(dataManager)

    for (const group of activityList) {
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

        const category = config.groupData(item, group.category)
        let list = data.get(category)
        if (list == null) {
          list = []
          data.set(category, list)
        }

        const chartItem: ChartDataItem = [
          "",
          dataDescriptor.unitConverter.convert(item.s),
          dataDescriptor.unitConverter.convert(item.s + item.d),
          dataDescriptor.unitConverter.convert(item.d),
          item,
          dataDescriptor.shortenName === false ? item.n : `${namePrefix}${getShortName(item.n)}`,
          group.category
        ]
        list.push(chartItem)
      }
    }

    // compute categories - for each thread maybe several categories as one service can include another one and we render it above each other
    const rowToItems = new Map<string, CustomSeriesOption & { label: SeriesLabelOption; itemStyle: unknown }>()
    const labelThreshold = 20 * config.dataDescriptor.unitConverter.factor
    let minStart = Number.MAX_VALUE
    let maxEnd = 0
    // const rowIndexThreshold = 300 * dataDescriptor.unitConverter.factor
    const rowToInfo = new Map<number, RowInfo>()
    const itemToActualIndex = new Map<ItemV20, number>()

    let groupRowStart = 0

    const groupNames = [...data.keys()]
    config.sortGroups(groupNames)

    for (const groupName of groupNames) {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      const list = data.get(groupName)!
      list.sort((a, b) => a[4].s - b[4].s)
      if (minStart > list[0][4].s) {
        minStart = list[0][4].s
      }
      const last = list.at(-1) as ChartDataItem
      if (maxEnd < last[4].s + last[4].d) {
        maxEnd = last[4].s + last[4].d
      }

      if (rowToInfo.size > 0) {
        groupRowStart = Math.max(...rowToInfo.keys())
        rowToInfo.clear()
      }

      itemToActualIndex.clear()

      for (const chartItem of list) {
        const item = chartItem[4]

        let rowIndex = groupRowStart
        // step one: put by one row above where parent item is
        if (config.hasParent) {
          const items = activityList[0].items
          if (item.pa != null) {
            const parent = items[item.pa]
            // here is the issue - if parent belongs to another group, we do not check
            // (so, there is chance that row index will be incorrect and splitLine will be drawn incorrectly)
            if (config.groupData(parent, "") === config.groupData(item, "")) {
              let parentRowIndex = itemToActualIndex.get(parent)
              if (parentRowIndex == null) {
                console.error("parentRowIndex is null")
                parentRowIndex = 0
              }
              rowIndex = parentRowIndex + 1
            }
          }
        }

        // step two: make sure, that no overlap with other items in the row
        rowIndex = findRowIndex(rowIndex, rowToInfo, item)

        let info = rowToInfo.get(rowIndex)
        if (info === undefined) {
          info = { items: [] }
          rowToInfo.set(rowIndex, info)
        }
        info.items.push(item)
        itemToActualIndex.set(item, rowIndex)

        const rowName = rowIndex === groupRowStart ? groupName : `${groupName}__${rowIndex}`

        chartItem[0] = rowName
        let series = rowToItems.get(rowName)
        if (series == null) {
          if (groupName == "completing") {
            debugger
          }
          // noinspection JSUnusedGlobalSymbols,TypeScriptValidateTypes
          series = {
            label: {},
            colorBy: "data",
            id: rowName,
            type: "custom",
            renderItem(params, api) {
              return renderItem(params, config.hasParent, api, labelThreshold)
            },
            encode: {
              x: [1, 2, 3],
              y: 0
            },
            itemStyle:
              groupName === "service waiting"
                ? ({
                  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                  //@ts-expect-error: https://github.com/apache/echarts/issues/16775
                  color(_): string {
                    return "#FF0000"
                  }
                } as never)
                : undefined,
            data: []
          }
          // eslint-disable-next-line @typescript-eslint/ban-ts-comment
          rowToItems.set(rowName, series)
        }

        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        ;(series.data as ChartDataItem[]).push(chartItem)
      }
    }

    const series = [...rowToItems.values()]
    series.sort((a, b) => {
      return (a.id as string).localeCompare(b.id as string, undefined, {numeric: true, sensitivity: "base"})
    })

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
          min: dataDescriptor.unitConverter.convert(minStart),
          max: dataDescriptor.unitConverter.convert(maxEnd)
        },
        series
      },
      {
        replaceMerge: ["series"]
      }
    )
    this.chart.enableZoomTool()
    this.chart.chart.resize()
  }
}

// prettier-ignore
function configureMarkAreas(dataManager: DataManager, series: CustomSeriesOption[], axisLineColor: string): void {
  const areaData: MarkArea2DDataItemOption[] = []
  for (const item of dataManager.items) {
    if (!(item.n.endsWith(" async preloading") || item.n.endsWith(" sync preloading"))) {
      continue
    }

    const isAsync = item.n.includes("async")
    areaData.push([
      {
        name: item.n.replace(" service", "").replace("service", ""),
        label: {
          verticalAlign: isAsync ? "top" : "bottom"
        },
        itemStyle: {
          borderType: isAsync ? "dotted" : "solid",
          borderWidth: 1,
          borderColor: axisLineColor,
          color: "rgba(0, 0, 0, 0)"
        },
        xAxis: item.s
      },
      { xAxis: item.s + item.d }
    ])
  }
  const lastSeries = series.at(-1) as CustomSeriesOption
  if (areaData.length === 0) {
    lastSeries.markArea = undefined
  }
  else {
    lastSeries.markArea = {
      silent: true,
      data: areaData,
    }
  }

  const markLineData: MarkLine1DDataItemOption[] = []
  for (const item of dataManager.data.traceEvents) {
    if (item.name !== "splash shown" && item.name !== "project opened") {
      continue
    }

    markLineData.push({
      label: { formatter: item.name },
      xAxis: item.ts / 1000
    })
  }
  if (markLineData.length === 0) {
    lastSeries.markLine = undefined
  }
  else {
    lastSeries.markLine = {
      symbol: "none",
      silent: true,
      lineStyle: {
        type: "dashed",
        color: axisLineColor
      },
      data: markLineData
    }
  }
}

// prettier-ignore
function renderItem(params: CustomSeriesRenderItemParams, hasParent: boolean, api: CustomSeriesRenderItemAPI, labelThreshold: number): CustomSeriesRenderItemReturn {
  const categoryIndex = api.value(0)
  const start = api.coord([api.value(1), categoryIndex])
  const end = api.coord([api.value(2), categoryIndex])
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const height = (api.size!([0, 1]) as number[])[1] * 0.8

  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  const coordinateSystem = params.coordSys as { x: number; y: number; width: number; height: number }

  const width = end[0] - start[0]
  // clip on zoom
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
    style: {
      fill: api.visual("color"),
    },
    emphasis: {
      style: {
        stroke: "#000",
        lineWidth: 2,
      },
    },
    textConfig: {
      position: "insideLeft",
    },
    textContent: {
      type: "text",
      style: {
        text: api.value(5).toString(),
        fontFamily: "monospace",
        width: width - 2,
        overflow: hasParent || width < labelThreshold ? "truncate" : undefined,
        ellipsis: "…",
        truncateMinChar: 1,
      },
    },
  }
}

// prettier-ignore
function findRowIndex(rowIndex: number, rowToInfo: Map<number, RowInfo>, item: ItemV20): number {
  rowLoop: for (let i = rowIndex; ; i++) {
    const rowInfo = rowToInfo.get(i)
    if (rowInfo === undefined) {
      return i
    }

    const itemStart = item.s
    const itemEnd = item.s + item.d

    for (const sibling of rowInfo.items) {
      const siblingEnd = sibling.s + sibling.d
      if (siblingEnd <= itemStart) {
        continue
      }

      const siblingStart = sibling.s
      if (siblingStart >= itemEnd) {
        continue
      }

      // overlap - check next row
      continue rowLoop
    }
    return i
  }
}

interface RowInfo {
  items: ItemV20[]
}
