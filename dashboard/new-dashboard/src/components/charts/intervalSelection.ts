import { ChartManager } from "./ChartManager"
import "./intervalSelection.css"

const HOVER_CLASS = "interval-selection-hovering-series"

/**
 * Activates the toolbox "zoom" tool (drag to select an interval) without the user having to click the toolbox icon.
 * The global cursor is dropped whenever the toolbox is rebuilt, so this has to be re-applied after every chart update.
 *
 * https://github.com/apache/echarts/issues/10274
 */
export function activateIntervalSelection(chartManager: ChartManager): void {
  if (chartManager.chart.isDisposed()) {
    return
  }
  chartManager.chart.dispatchAction({
    type: "takeGlobalCursor",
    key: "dataZoomSelect",
    dataZoomSelectActive: true,
  })
}

/**
 * The brush cursor of the active zoom tool hides the "this point is clickable" affordance, so it is restored while the
 * pointer is over a series. Driven by HoverFadeController, which already tracks the hovered series.
 */
export function setSeriesHoverCursor(chartManager: ChartManager, hovering: boolean): void {
  chartManager.chartContainer.classList.toggle(HOVER_CLASS, hovering)
}
