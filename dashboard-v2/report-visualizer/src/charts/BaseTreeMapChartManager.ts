import { TreeMap, LabelBullet } from "@amcharts/amcharts4/charts"
import { create, Scrollbar, color } from "@amcharts/amcharts4/core"
import { BaseChartManager, configureCursor } from "./ChartManager"

export abstract class BaseTreeMapChartManager extends BaseChartManager<TreeMap> {
  protected constructor(container: HTMLElement) {
    super(create(container, TreeMap))

    configureCursor(this.chart)

    // cursor tooltip is distracting (cannot be in BaseChartManager because only TreeMap creates axis as part of chart creation, for other charts axis is created customly)
    this.chart.xAxis.cursorTooltipEnabled = false
    this.chart.yAxis.cursorTooltipEnabled = false
  }

  protected enableZoom(): void {
    const chart = this.chart
    chart.mouseWheelBehavior = "zoomX"
    chart.scrollbarX = new Scrollbar()
    chart.mouseWheelBehavior = "zoomXY"
  }

  protected configureLabelBullet(bullet: LabelBullet): void {
    bullet.locationY = 0.5
    bullet.locationX = 0.5
    bullet.label.fill = color("#fff")
  }
}