/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { DurationAxis } from "@amcharts/amcharts4/charts"
import { DataManager, SERVICE_WAITING } from "../state/DataManager"
import { CompleteTraceEvent } from "../state/data"
import { ActivityChartDescriptor, getShortName } from "./ActivityChartDescriptor"
import { ActivityChartManager, ActivityLegendItem, ClassItem, ClassItemChartConfig } from "./ActivityChartManager"

export class ServiceChartManager extends ActivityChartManager {
  constructor(container: HTMLElement, sourceNames: Array<string>, descriptor: ActivityChartDescriptor) {
    super(container, sourceNames, descriptor)
  }

  protected configureDurationAxis(): DurationAxis {
    const durationAxis = super.configureDurationAxis()
    durationAxis.baseUnit = "millisecond"
    return durationAxis
  }

  protected getTooltipText(): string {
    return super.getTooltipText() + "\ntotal duration: {totalDuration} ms"
  }

  render(dataManager: DataManager): void {
    const list = dataManager.serviceEvents
    const categoryToEvents = new Map<string, Array<CompleteTraceEvent>>()
    for (const event of list) {
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
      const category = event.cat!
      let categoryEvents = categoryToEvents.get(category)
      if (categoryEvents == null) {
        categoryEvents = []
        categoryToEvents.set(category, categoryEvents)
      }
      categoryEvents.push(event)
    }

    let colorIndex = 0
    const legendData: Array<ActivityLegendItem> = []
    const applicableSources = new Set<string>()
    const concatenatedData: Array<ClassItem> = []

    for (const sourceName of this.sourceNames) {
      // see XYChart.handleSeriesAdded method - fill and stroke are required to be set, and stroke is set to fill if not set otherwise (in our case fill and stroke are equals)
      // do not use XYChart.handleSeriesAdded because on add series it is already called, so, we need to reuse already generated color for first index
      // noinspection DuplicatedCode
      const color = this.chart.colors.getIndex(colorIndex++)
      const chartConfig: ClassItemChartConfig = {
        fill: color,
        stroke: color,
      }

      // generate color before - even if no data for this type of items, still color should be the same regardless of current data set
      // so, if currently no data for project, but there is data for modules, color for modules should use index 3 and not 2
      const items = categoryToEvents.get(sourceName)
      if (items == null || items.length === 0) {
        continue
      }

      const legendItem: ActivityLegendItem = {
        name: sourceNameToLegendName(sourceName, items.length),
        fill: color,
        sourceName,
      }
      legendData.push(legendItem)
      applicableSources.add(sourceName)

      transformTraceEventToClassItem(items, chartConfig, concatenatedData, true)
    }

    // noinspection DuplicatedCode
    if (this.chart.legend != null) {
      this.chart.legend.data = legendData
      this.legendHitHandler = (legendItem, isActive) => {
        if (isActive) {
          applicableSources.add(legendItem.sourceName)
        }
        else {
          applicableSources.delete(legendItem.sourceName)
        }
        // maybe there is a more effective way to hide data, but this one is reliable and simple
        this.chart.data = concatenatedData.filter(it => applicableSources.has(it.sourceName))
      }
    }

    concatenatedData.sort((a, b) => a.start - b.start)
    this.computeRangeMarkers(dataManager, concatenatedData)
    this.chart.data = concatenatedData
  }
}

// 10 ms
const threshold = 10 * 1000

export function transformTraceEventToClassItem(items: Array<CompleteTraceEvent>, chartConfig: ClassItemChartConfig | null, resultList: Array<ClassItem>,
                                               durationAsOwn: boolean): void {
  for (const item of items) {
    const isServiceWaiting = item.cat === SERVICE_WAITING
    if (durationAsOwn && isServiceWaiting) {
      continue
    }

    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    const ownDur = item.args!.ownDur
    const reportedDurationInMicroSeconds = durationAsOwn ? ownDur : item.dur

    if (reportedDurationInMicroSeconds < threshold) {
      continue
    }

    const shortName = getShortName(item)
    const result: ClassItem & CompleteTraceEvent = {
      ...item,
      name: isServiceWaiting ? `wait for ${item.name}` : item.name,
      sourceName: item.cat!,
      shortName: isServiceWaiting ? `wait for ${shortName}` : shortName,
      chartConfig,
      thread: item.tid,
      plugin: item.args!.plugin,
      start: Math.round(item.ts / 1000),
      end: Math.round((item.ts + item.dur) / 1000),
      duration: Math.round(reportedDurationInMicroSeconds / 1000),
      totalDuration: Math.round(item.dur / 1000),
    }

    if (!durationAsOwn) {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (result as any).ownDuration = Math.round(ownDur / 1000)
    }
    resultList.push(result)
  }
}

function sourceNameToLegendName(sourceName: string, itemCount: number): string {
  let prefix
  if (sourceName.startsWith("app")) {
    prefix = "App"
  }
  else if (sourceName.startsWith("project")) {
    prefix = "Project"
  }
  else if (sourceName.startsWith("module")) {
    prefix = "Module"
  }
  return `${prefix} ${sourceName.includes("Component") ? "component" : "service"} (${itemCount})`
}