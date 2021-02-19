import { Legend } from "@amcharts/amcharts4/charts"
import { create, percent } from "@amcharts/amcharts4/core"
import { SunburstSeries, Sunburst } from "@amcharts/amcharts4/plugins/sunburst"
import { DataManager } from "../state/DataManager"
import { BaseChartManager } from "./ChartManager"

export class StatsChartManager extends BaseChartManager<Sunburst> {
  constructor(container: HTMLElement) {
    super(create(container, Sunburst))

    const chart = this.chart
    chart.dataFields.value = "value"
    chart.dataFields.name = "name"
    chart.dataFields.children = "children"
    chart.legend = new Legend()
    chart.fontSize = 11

    // to allocate more space for central labels - still overlapped, but in any case more readable
    chart.innerRadius = percent(10)

    const level0SeriesTemplate = new SunburstSeries()
    chart.seriesTemplates.setKey("0", level0SeriesTemplate)
    chart.seriesTemplates.setKey("0", level0SeriesTemplate)
    level0SeriesTemplate.labels.template.text = "{name}: {value}"

    const level1SeriesTemplate = new SunburstSeries()
    chart.seriesTemplates.setKey("1", level1SeriesTemplate)
    level1SeriesTemplate.fillOpacity = 0.75
    level1SeriesTemplate.hiddenInLegend = true
    level1SeriesTemplate.labels.template.text = "{name}: {value}"

    chart.legend.valueLabels.template.text = "{value.value}"
  }

  render(data: DataManager): void {
    const stats = data.data.stats
    this.chart.data = [
      {
        name: "Components",
        children: [
          {name: "app", value: stats.component.app},
          {name: "project", value: stats.component.project},
          {name: "module", value: stats.component.module},
        ],
      },
      {
        name: "Services",
        children: [
          {name: "app", value: stats.service.app},
          {name: "project", value: stats.service.project},
          {name: "module", value: stats.service.module},
        ],
      },
    ]
  }
}