import { TreemapSeriesNodeItemOption } from "echarts/types/src/chart/treemap/TreemapSeries"
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { adaptToolTipFormatter } from "shared/src/chart"
import { ChartOptions, TreeMapChartOptions, useTreeMapChart } from "shared/src/echarts"
import { numberFormat } from "shared/src/formatter"
import { DataManager } from "../state/DataManager"
import { PluginStatItem } from "../state/data"
import { ChartManager } from "./ChartManager"

useTreeMapChart()

interface ItemExtraInfo {
  abbreviatedName: string
  item: PluginStatItem
}

export class PluginClassCountTreeMapChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper
  constructor(container: HTMLElement) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<ChartOptions>({
      tooltip: {
        formatter: adaptToolTipFormatter(params => {
          const info = params[0]
          const value = info.value as number
          let result = `${info.marker as string} ${info.name}<span style="float:right;margin-left:20px;font-weight:900">${numberFormat.format(value)}</span>`
          const item = (info.data as ItemExtraInfo).item
          result += `<br/>  class loading time in EDT<span style="float:right;margin-left:20px">${item.classLoadingEdtTime}</span>`
          result += `<br/>class loading time in background<span style="float:right;margin-left:20px">${item.classLoadingBackgroundTime}</span>`
          return result
        }),
      },
    })
  }

  dispose(): void {
    this.chart.dispose()
  }

  render(data: DataManager): void {
    const items: Array<TreemapSeriesNodeItemOption & ItemExtraInfo> = []

    const loadedClasses = data.data.stats.loadedClasses
    if (loadedClasses == null) {
      // v20+
      for (const item of (data.data.plugins ?? [])) {
        items.push({
          name: item.id,
          abbreviatedName: getAbbreviatedName(item.id),
          value: item.classCount,
          item,
        })
      }
    }
    else {
      for (const name of Object.keys(loadedClasses)) {
        items.push({
          name,
          abbreviatedName: getAbbreviatedName(name),
          value: loadedClasses[name],
          item: {classCount: 0, classLoadingBackgroundTime: 0, classLoadingEdtTime: 0, id: name}
        })
      }
    }

    this.chart.chart.setOption<TreeMapChartOptions>({
      series: [{
        type: "treemap",
        data: items,
        leafDepth: 2,
        roam: "move",
        label: {
          formatter(data) {
            return `${(data.data as ItemExtraInfo).abbreviatedName} (${numberFormat.format(data["value"] as number)})`
          }
        }
      }],
    })
  }
}

function getAbbreviatedName(name: string): string {
  if (!name.includes(".")) {
    return name
  }

  let abbreviatedName = ""
  const names = name.split(".")
  for (let i = 0; i < names.length; i++) {
    const unqualifiedName = names[i]
    if (i == (names.length - 1)) {
      abbreviatedName += unqualifiedName
    } else {
      abbreviatedName += unqualifiedName.substring(0, 1) + "."
    }
  }
  return abbreviatedName
}