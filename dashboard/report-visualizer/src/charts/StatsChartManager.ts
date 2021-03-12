import { ChartManagerHelper } from "shared/src/ChartManagerHelper"
import { SunburstChartOptions, useSunburstChart } from "shared/src/echarts"
import { numberFormat } from "shared/src/formatter"
import { DataManager } from "../state/DataManager"
import { ChartManager } from "./ChartManager"

useSunburstChart()

export class StatsChartManager implements ChartManager {
  private readonly chart: ChartManagerHelper

  constructor(container: HTMLElement) {
    this.chart = new ChartManagerHelper(container)
    this.chart.chart.setOption<SunburstChartOptions>({
      tooltip: {},
    })
  }

  dispose(): void {
    this.chart.dispose()
  }

  render(data: DataManager): void {
    const stats = data.data.stats

    this.chart.chart.setOption<SunburstChartOptions>({
      series: [{
        type: "sunburst",
        "emphasis": {
          focus: "none",
        },
        data: [
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
        ],
        label: {
          formatter(data) {
            return `${data.name} (${numberFormat.format(data["value"] as number)})`
          }
        }
      }],
    })
  }
}
