import { TreemapChart } from "echarts/charts"
import { ToolboxComponent, TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { TreemapSeriesNodeItemOption } from "echarts/types/src/chart/treemap/TreemapSeries"
import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { adaptToolTipFormatter } from "shared/src/chart"
import { LineChartOptions, TreeMapChartOptions } from "shared/src/echarts"
import { numberFormat } from "shared/src/formatter"
import { DataManager } from "../DataManager"
import { PluginStatItem } from "../data"
import { ChartManager } from "./ChartComponent"
import { buildTooltip } from "./tooltip"

use([TooltipComponent, CanvasRenderer, TreemapChart, ToolboxComponent])

interface ItemExtraInfo {
  abbreviatedName: string
  item: PluginStatItem
}

export class PluginClassCountTreeMapChartManager implements ChartManager {
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
          const item = (info.data as ItemExtraInfo).item
          return `${info.marker as string} ` + buildTooltip([
            {name: info.name, main: true, value: numberFormat.format(info.value as number)},
            {name: "class loading time in EDT", value: numberFormat.format(item.classLoadingEdtTime)},
            {name: "class loading time in background", value: numberFormat.format(item.classLoadingBackgroundTime)},
          ])
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