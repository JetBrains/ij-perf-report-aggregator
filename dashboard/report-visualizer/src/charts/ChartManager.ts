// force order in chunk
import "@amcharts/amcharts4/.internal/core/elements/Modal"
import { AxisDataItem, XYChart, XYCursor } from "@amcharts/amcharts4/charts"
import { color, Color, create, ExportMenu, options as amChartOptions, Scrollbar } from "@amcharts/amcharts4/core"
import { DataManager } from "../DataManager"

// helps during hot-reload
amChartOptions.autoDispose = true

export interface ChartManager {
  render(data: DataManager): void

  dispose(): void
}

function addExportMenu(chart: XYChart): void {
  const exportMenu = new ExportMenu()
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const topItems = exportMenu.items[0].menu!
  for (let i = topItems.length - 1; i >= 0; i--) {
    const chartElement = topItems[i]
    if (chartElement.label == "Data") {
      topItems.splice(i, 1)
    }
    else if (chartElement.label == "Image") {
      // remove PDF
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
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

export abstract class XYChartManager implements ChartManager {
  protected readonly blackColor = color("#000000")
  protected readonly chart: XYChart

  protected constructor(container: HTMLElement) {
    this.chart = create(container, XYChart)
    this.chart.mouseWheelBehavior = "zoomX"
    this.chart.scrollbarX = new Scrollbar()
    addExportMenu(this.chart)
    const cursor = new XYCursor()
    cursor.lineY.disabled = true
    cursor.lineX.disabled = true
    cursor.behavior = "zoomXY"
    this.chart.cursor = cursor
  }

  dispose(): void {
    this.chart.dispose()
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
      // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
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