import { SunburstChart } from "echarts/charts"
import { TooltipComponent } from "echarts/components"
import { use } from "echarts/core"
import { CanvasRenderer } from "echarts/renderers"
import { ChartManagerHelper } from "../../components/common/ChartManagerHelper"
import { SunburstChartOptions } from "../../components/common/echarts"
import { numberFormat } from "../../components/common/formatter"
import { DataManager } from "../DataManager"
import { ChartManager } from "./ChartComponent"

use([TooltipComponent, CanvasRenderer, SunburstChart])

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
      series: [
        {
          type: "sunburst",
          emphasis: {
            focus: "none",
          },
          data: [
            {
              name: "Components",
              children: [
                { name: "app", value: stats.component.app },
                { name: "project", value: stats.component.project },
                { name: "module", value: stats.component.module },
              ],
            },
            {
              name: "Services",
              children: [
                { name: "app", value: stats.service.app },
                { name: "project", value: stats.service.project },
                { name: "module", value: stats.service.module },
              ],
            },
          ],
          label: {
            formatter(data) {
              return `${data.name} (${numberFormat.format(data.value as number)})`
            },
          },
        },
      ],
    })
  }
}
