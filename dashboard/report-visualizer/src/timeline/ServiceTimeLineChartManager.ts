import { Legend, DurationAxisDataItem } from "@amcharts/amcharts4/charts"
import { ColorSet, Color, color as chartColor } from "@amcharts/amcharts4/core"
import { getShortName } from "../ActivityChartDescriptor"
import { DataManager, SERVICE_WAITING } from "../DataManager"
import { ClassItem, ClassItemChartConfig } from "../charts/ActivityChartManager"
import { LegendItem } from "../charts/ChartManager"
import { CompleteTraceEvent } from "../data"
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

function transformTraceEventToClassItem(items: Array<CompleteTraceEvent>, chartConfig: ClassItemChartConfig | null, resultList: Array<ClassItem>,
                                               durationAsOwn: boolean): void {
  for (const item of items) {
    const isServiceWaiting = item.cat === SERVICE_WAITING
    if (durationAsOwn && isServiceWaiting) {
      continue
    }

    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const ownDur = item.args!.ownDur
    const reportedDurationInMicroSeconds = durationAsOwn ? ownDur : item.dur

    if (reportedDurationInMicroSeconds < (10 * 1000)) {
      continue
    }

    const shortName = getShortName(item)
    const result: ClassItem & CompleteTraceEvent = {
      ...item,
      name: isServiceWaiting ? `wait for ${item.name}` : item.name,
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      sourceName: item.cat!,
      shortName: isServiceWaiting ? `wait for ${shortName}` : shortName,
      chartConfig,
      thread: item.tid,
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      plugin: item.args!.plugin,
      start: Math.round(item.ts / 1000),
      end: Math.round((item.ts + item.dur) / 1000),
      duration: Math.round(reportedDurationInMicroSeconds / 1000),
      totalDuration: Math.round(item.dur / 1000),
    }

    if (!durationAsOwn) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any,@typescript-eslint/no-unsafe-member-access
      (result as any).ownDuration = Math.round(ownDur / 1000)
    }
    resultList.push(result)
  }
}