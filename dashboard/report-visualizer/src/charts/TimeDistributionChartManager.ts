import { TreemapChart } from "echarts/charts"
import { TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { TreemapSeriesNodeItemOption } from "echarts/types/src/chart/treemap/TreemapSeries"
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { adaptToolTipFormatter } from "shared/src/chart"
import { LineChartOptions, TreeMapChartOptions } from "shared/src/echarts"
import { numberFormat } from "shared/src/formatter"
import { getShortName } from "../ActivityChartDescriptor"
import { DataManager } from "../DataManager"
import { InputDataV20, ItemV20 } from "../data"
import { ChartManager } from "./ChartComponent"

use([TooltipComponent, CanvasRenderer, TreemapChart])

interface ItemExtraInfo {
  count?: string
}

export class TimeDistributionChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<LineChartOptions>({
      toolbox: {
        feature: {
          saveAsImage: {},
        },
      },
      tooltip: {
        formatter: adaptToolTipFormatter(params => {
          const info = params[0]
          let result = `${info.marker as string} ${info.name}`
          const count = (info.data as ItemExtraInfo).count
          if (count !== undefined) {
            result += ` (count=${count})`
          }
          return result + `<span class="tooltipMainValue">${numberFormat.format(info.value as number)}</span>`
        }),
      },
    })
  }

  dispose(): void {
    this.chart.dispose()
  }

  render(data: DataManager): void {
    const items: Array<TreemapSeriesNodeItemOption & ItemExtraInfo> = []

    addServicesOrComponents(data, items, "component", "appComponents", "projectComponents")
    addServicesOrComponents(data, items, "service", "appServices", "projectServices")
    addIcons(data, items)

    this.chart.chart.setOption<TreeMapChartOptions>({
      series: [{
        type: "treemap",
        data: items,
        levels: [
          {},
          {},
          {
            colorSaturation: [0.35, 0.5],
            itemStyle: {
              borderWidth: 5,
              gapWidth: 1,
              borderColorSaturation: 0.6,
            },
            upperLabel: {show: true},
          },
        ],
        leafDepth: 3,
        label: {
          formatter(data) {
            return `${data.name} (${numberFormat.format(data["value"] as number)})`
          },
        },
      }],
    })
  }
}

function addServicesOrComponents(dataManager: DataManager,
                                 items: Array<TreemapSeriesNodeItemOption & ItemExtraInfo>,
                                 statName: "component" | "service",
                                 appFieldName: "appServices" | "appComponents",
                                 projectFieldName: "projectComponents" | "projectServices"): void {
  const children: Array<TreemapSeriesNodeItemOption & ItemExtraInfo> = []
  const data = dataManager.data as InputDataV20
  const stats = data.stats[statName]
  children.push({
    name: `app ${statName}s`,
    children: toTreeMapItem(data[appFieldName]),
    count: numberFormat.format(stats.app),
  }, {
    name: `project ${statName}s`,
    children: toTreeMapItem(data[projectFieldName]),
    count: numberFormat.format(stats.project),
  })

  items.push({
    name: `${statName}s`,
    children,
    count: numberFormat.format(stats.app + stats.project + stats.module),
  })
}

function addIcons(data: DataManager, items: Array<TreemapSeriesNodeItemOption>) {
  const icons = data.data.icons
  if (icons == null) {
    return
  }

  const iconList: Array<TreemapSeriesNodeItemOption> = []
  let duration = 0
  for (const item of icons) {
    duration += item.loading
    iconList.push({
      value: item.loading,
      ...item,
      children: [
        {name: "searching", value: item.loading - item.decoding},
        {name: "decoding", value: item.decoding},
      ],
    })
  }
  items.push({
    name: "icons",
    children: iconList,
    value: duration,
  })
}

function toTreeMapItem(items: Array<ItemV20> | null | undefined): Array<TreemapSeriesNodeItemOption> {
  return items == null ? [] : items.map(it => {
    return {name: getShortName(it), value: it.d / 1000}
  })
}