/* eslint-disable @typescript-eslint/no-non-null-assertion */
import { XYChart, XYCursor, Chart, AxisDataItem } from "@amcharts/amcharts4/charts"
import { options as amChartOptions, ExportMenu, Scrollbar, color, create, Color } from "@amcharts/amcharts4/core"
import { DataManager } from "../state/DataManager"
amChartOptions.onlyShowOnViewport = true
// helps during hot-reload
amChartOptions.autoDispose = true

export interface ChartManager {
  render(data: DataManager): void

  dispose(): void
}

export function addExportMenu(chart: XYChart): void {
  const exportMenu = new ExportMenu()
  const topItems = exportMenu.items[0].menu!
  for (let i = topItems.length - 1; i >= 0; i--) {
    const chartElement = topItems[i]
    if (chartElement.label == "Data") {
      topItems.splice(i, 1)
    }
    else if (chartElement.label == "Image") {
      // remove PDF
      const subMenu = chartElement.menu!
      const length = subMenu.length
      if (subMenu[length - 1].label == "PDF") {
        subMenu.length = length - 1
      }
    }
  }

  chart.exporting.menu = exportMenu
  chart.exporting.menu.align = "right"
  chart.exporting.menu.verticalAlign = "bottom"
}

export function configureCommonChartSettings(chart: XYChart): void {
  chart.mouseWheelBehavior = "zoomX"
  chart.scrollbarX = new Scrollbar()
  addExportMenu(chart)
}

export function configureCursor(chart: XYChart): void {
  const cursor = new XYCursor()
  cursor.lineY.disabled = true
  cursor.lineX.disabled = true
  cursor.behavior = "zoomXY"
  chart.cursor = cursor
}

export abstract class BaseChartManager<T extends Chart> implements ChartManager {
  protected constructor(protected readonly chart: T) {
  }

  abstract render(data: DataManager): void

  /** @override */
  dispose(): void {
    this.chart.dispose()
  }
}

export abstract class XYChartManager extends BaseChartManager<XYChart> {
  protected readonly blackColor = color("#000000")

  protected constructor(container: HTMLElement) {
    super(create(container, XYChart))

    configureCommonChartSettings(this.chart)
    configureCursor(this.chart)
  }

  abstract render(data: DataManager): void

  protected configureRangeMarker(range: AxisDataItem, label: string, color: Color, yOffset = 0): void {
    range.label.inside = true
    range.label.horizontalCenter = "middle"
    range.label.fontSize = 12
    range.label.valign = "bottom"
    range.label.text = label
    range.grid.stroke = color
    range.grid.strokeDasharray = "2,2"
    range.grid.strokeOpacity = 1

    range.label.adapter.add("dy", (_y, _target) => {
      return -this.chart.yAxes.getIndex(0)!.pixelHeight + yOffset
    })
    range.label.adapter.add("x", (x, _target) => {
      const rangePoint = range.point
      return rangePoint == null ? x : rangePoint.x
    })
  }
}

export interface LegendItem {
  readonly name: string
  readonly fill: Color
}