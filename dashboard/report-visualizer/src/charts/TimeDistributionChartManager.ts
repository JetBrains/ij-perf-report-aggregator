import { TreemapSeriesNodeItemOption } from "echarts/types/src/chart/treemap/TreemapSeries"
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { adaptToolTipFormatter, ChartOptions, numberFormat } from "shared/src/chart"
import { TreeMapChartOptions, useTreeMapChart } from "shared/src/echarts"
import { DataManager } from "../state/DataManager"
import { IconData, InputDataV20, ItemV20 } from "../state/data"
import { getShortName } from "./ActivityChartDescriptor"
import { ChartManager } from "./ChartManager"

useTreeMapChart()

interface ItemExtraInfo {
  count?: string
}

export class TimeDistributionChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper
  constructor(container: HTMLElement) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<ChartOptions>({
      tooltip: {
        formatter: adaptToolTipFormatter(params => {
          const info = params[0]
          let result = `${info.marker} ${info.name}`
          const count = (info.data as ItemExtraInfo).count
          if (count !== undefined) {
            result += ` (count=${count})`
          }
          return result + `<span style="float:right;margin-left:20px;font-weight:900">${numberFormat.format(info.value as number)}</span>`
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
        leafDepth: 2,
        roam: "move",
        label: {
          formatter(data) {
            return `${data.name} (${numberFormat.format(data["value"] as number)})`
          }
        }
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
  })
  children.push({
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
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  let count = 0
  let duration = 0
  for (const [key, value] of Object.entries(icons)) {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const info = value as IconData
    count += info.count
    duration += info.loading
    iconList.push({
      name: key,
      value: info.loading,
      ...info,
      children: [
        {name: "searching", value: info.loading - info.decoding},
        {name: "decoding", value: info.decoding},
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