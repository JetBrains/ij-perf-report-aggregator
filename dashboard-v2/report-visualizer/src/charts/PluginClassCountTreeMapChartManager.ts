import { LabelBullet } from "@amcharts/amcharts4/charts"
import { DataManager } from "../state/DataManager"
import { BaseTreeMapChartManager } from "./BaseTreeMapChartManager"

export class PluginClassCountTreeMapChartManager extends BaseTreeMapChartManager {
  constructor(container: HTMLElement) {
    super(container)

    const chart = this.chart
    chart.dataFields.value = "count"
    chart.dataFields.name = "name"

    this.enableZoom()

    const level1 = chart.seriesTemplates.create("0")
    const level1Bullet = level1.bullets.push(new LabelBullet())
    this.configureLabelBullet(level1Bullet)
    level1Bullet.label.text = "{abbreviatedName} ({count})"
  }

  render(data: DataManager): void {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const items: Array<any> = []

    const loadedClasses = data.data.stats.loadedClasses
    if (loadedClasses != null) {
      for (const name of Object.keys(loadedClasses)) {
        items.push({
          name,
          abbreviatedName: getAbbreviatedName(name),
          count: loadedClasses[name],
        })
      }
    }

    this.chart.data = items
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