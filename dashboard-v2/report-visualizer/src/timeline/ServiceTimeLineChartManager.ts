import { Legend, DurationAxisDataItem } from "@amcharts/amcharts4/charts"
import { ColorSet, Color, color as chartColor } from "@amcharts/amcharts4/core"
import { ClassItem } from "../charts/ActivityChartManager"
import { LegendItem } from "../charts/ChartManager"
import { transformTraceEventToClassItem } from "../charts/ServiceChartManager"
import { DataManager } from "../state/DataManager"
import { CompleteTraceEvent } from "../state/data"
import { BaseTimeLineChartManager } from "./BaseTimeLineChartManager"
import { TimeLineGuide } from "./timeLineChartHelper"

export class ServiceTimeLineChartManager extends BaseTimeLineChartManager {
  private readonly rangeLegend: Legend

  constructor(container: HTMLElement) {
    super(container)

    this.configureDurationAxis()
    this.configureSeries("{shortName}")

    const rangeLegend = new Legend()
    this.rangeLegend = rangeLegend
    rangeLegend.parent = this.chart.chartContainer
    rangeLegend.tooltipText = "Range Legend"
    rangeLegend.interactionsEnabled = false
    // make clear that it is not item, but range legend
    rangeLegend.markers.template.width = 4
  }

  protected getToolTipText(): string {
    return "{name}: {ownDuration} ms\nrange: {start}-{end}\nthread: {thread}" + "\nplugin: {plugin}" + "\ntotal duration: {totalDuration} ms"
  }

  render(dataManager: DataManager): void {
    this.guides.length = 0

    this.chart.data = this.transformIjData(dataManager)

    this.computeRangeMarkers(dataManager)
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private transformIjData(dataManager: DataManager): Array<any> {
    const colorSet = new ColorSet()

    const items: Array<ClassItem & CompleteTraceEvent> = []
    transformTraceEventToClassItem(dataManager.serviceEvents, null, items, false)
    this.maxRowIndex = 0

    return this.transformParallelToTimeLineItems(items, colorSet, 0)
  }

  protected collectGuides(dataManager: DataManager, guides: Array<TimeLineGuide>): void {
    const colorSet = new ColorSet()
    colorSet.step = 2
    for (const item of dataManager.data.prepareAppInitActivities) {
      if (item.name.endsWith(" async preloading") || item.name.endsWith(" sync preloading")) {
        const color = colorSet.next()
        guides.push({label: item.name, value: item.start, endValue: item.end, color})
      }
    }

    this.rangeLegend.data = this.guides.map(it => {
      const result: LegendItem = {name: it.label, fill: it.color}
      return result
    })
  }

  protected configureRangeMarker(range: DurationAxisDataItem, _label: string, color: Color, _yOffset = 0): void {
    const axisFill = range.axisFill
    axisFill.stroke = color
    axisFill.strokeDasharray = "2,2"
    axisFill.strokeOpacity = 1

    // HTML WhiteSmoke
    axisFill.fill = chartColor("#F5F5F5")
    axisFill.fillOpacity = 0.3
  }
}