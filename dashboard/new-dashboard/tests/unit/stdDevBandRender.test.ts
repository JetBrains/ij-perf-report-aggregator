import { createPinia, setActivePinia } from "pinia"
import { afterEach, beforeEach, describe, expect, it } from "vitest"
import { shallowRef } from "vue"
import { EChartsType, init, use } from "echarts/core"
import { SVGRenderer } from "echarts/renderers"
import { ECBasicOption } from "echarts/types/dist/shared"
// for the side effect: ChartManager is where the app registers its ECharts components, and a type-only import
// of it would be elided along with them
import "../../src/components/charts/ChartManager"
import type { ChartManager } from "../../src/components/charts/ChartManager"
import { HoverFadeController } from "../../src/components/charts/hoverFade"
import { generateQueries, DataQueryResult } from "../../src/components/common/DataQueryExecutor"
import { DataQuery, DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"
import { PredefinedMeasureConfigurator } from "../../src/configurators/MeasureConfigurator"
import { useSettingsStore } from "../../src/components/settings/settingsStore"
import { dbTypeStore } from "../../src/shared/dbTypes"

// happy-dom has no canvas 2D context, so the app's own renderer cannot draw here. Everything else - the
// component set, the option built by the configurator, the style patch the hover fade applies - is the app's.
use([SVGRenderer])

const PROJECT = "caddy/completion/variable"
const MEAN = "completion#mean_value"
const MEAN_SERIES_ID = `${MEAN}@${PROJECT}`
const MEANS = [100, 200]
const STD_DEVS = [10, 20]
const TIMESTAMPS = [1_000, 1_001]

// A filled, closed polygon: the band's area. Symbol circles (arc commands) and the clip rectangles
// (relative `l` commands) are shaped differently and do not match.
const AREA_PATH = /^M[\d.-]+ [\d.-]+(?:L[\d.-]+ [\d.-]+)+Z$/

function queryResult(values: number[], measureName: string): (string | number)[][] {
  return [
    TIMESTAMPS,
    values,
    values.map(() => measureName),
    values.map(() => "d"),
    values.map(() => "Linux EC2"),
    values.map((_, index) => 5_000 + index),
    values.map(() => PROJECT),
    values.map(() => "master"),
  ]
}

function chartOptions(): Promise<ECBasicOption> {
  const configuration = new DataQueryExecutorConfiguration()
  const query = new DataQuery()
  query.db = "perfintDev"
  query.table = "goland"

  const configurator = new PredefinedMeasureConfigurator(shallowRef([MEAN]), shallowRef(true), "line", "auto", {}, null, "item", "lower", true)
  if (!configurator.configureQuery(query, configuration)) {
    throw new Error("the measures were not configured into the query")
  }
  generateQueries(query, configuration)

  const data: DataQueryResult = [queryResult(MEANS, MEAN), queryResult(STD_DEVS, "completion#standard_deviation")]
  return configurator.configureChart(data, configuration)
}

function yAxisExtent(chart: EChartsType): number[] {
  return (chart as unknown as { getModel: () => { getComponent: (type: string, index: number) => { axis: { scale: { getExtent: () => number[] } } } } })
    .getModel()
    .getComponent("yAxis", 0)
    .axis.scale.getExtent()
}

function areaPaths(container: HTMLElement): SVGPathElement[] {
  return [...container.querySelectorAll("path")].filter((path) => (path.getAttribute("fill") ?? "none") !== "none" && AREA_PATH.test(path.getAttribute("d") ?? ""))
}

/** The band's expected outline: along its upper edge, then back along its lower one. */
function expectedBandPath(chart: EChartsType): string {
  const target = { seriesId: MEAN_SERIES_ID }
  const upper = TIMESTAMPS.map((t, i) => chart.convertToPixel(target, [t, MEANS[i] + STD_DEVS[i]]))
  const lower = TIMESTAMPS.map((t, i) => chart.convertToPixel(target, [t, MEANS[i] - STD_DEVS[i]]))
  return `M${point(upper[0])}` + [...upper.slice(1), ...lower.toReversed()].map((it) => `L${point(it)}`).join("") + "Z"
}

// The SVG renderer writes coordinates rounded to four decimals, with trailing zeros dropped.
function point([x, y]: number[]): string {
  return `${round(x)} ${round(y)}`
}

function round(value: number): number {
  return Number(value.toFixed(4))
}

describe("standard deviation band as ECharts renders it", () => {
  let container: HTMLElement
  let chart: EChartsType

  beforeEach(async () => {
    setActivePinia(createPinia())
    dbTypeStore().setDbType("perfintDev", "goland")
    useSettingsStore().stdDevInterval = true

    container = document.createElement("div")
    document.body.append(container)
    chart = init(container, undefined, { renderer: "svg", width: 600, height: 400 })
    // the order the app uses: LineChartVM sets up the axes, then ChartManager.updateChart replaces the data
    chart.setOption({ animation: false, xAxis: { type: "time" }, yAxis: { type: "value" } })
    chart.setOption(await chartOptions(), { replaceMerge: ["dataset", "series"] })
  })

  afterEach(() => {
    chart.dispose()
    container.remove()
  })

  it("stacks the two halves into one band from mean - stdDev to mean + stdDev", () => {
    const paths = areaPaths(container)
    expect(paths).toHaveLength(1)
    // The lower half contributes no fill, so the single area runs between the two edges rather than up from
    // the axis - which is the whole reason the band is stacked.
    expect(paths[0].getAttribute("d")).toBe(expectedBandPath(chart))
  })

  // The band gets no explicit color: it shares the `name` of its line, and ECharts hands out palette colors
  // per series name. Nothing else keeps the two in step.
  it("takes the color of the line it wraps", () => {
    expect(areaPaths(container)[0].getAttribute("fill")).toBe(chart.getVisual({ seriesId: MEAN_SERIES_ID }, "color"))
  })

  // The band's upper half holds the band's *width*, not where it is drawn, so anything that judges a point by
  // the number its series holds gets the band wrong. The y-axis slider used to filter on exactly that: a band
  // 20 wide around a mean of 100 was dropped by any window starting above 20, while both means stayed visible.
  describe("with the y axis narrowed", () => {
    it("survives flexible zero", () => {
      useSettingsStore().flexibleYZero = true
      chart.setOption({ yAxis: { min: (value: { min: number }) => value.min * 0.9 } })

      expect(yAxisExtent(chart)[0]).toBeGreaterThan(Math.min(...STD_DEVS))
      const paths = areaPaths(container)
      expect(paths).toHaveLength(1)
      expect(paths[0].getAttribute("d")).toBe(expectedBandPath(chart))
    })

    it.each([
      [30, 100],
      [50, 100],
      [0, 60],
    ])("survives the slider dragged to %i-%i percent", (start, end) => {
      useSettingsStore().flexibleYZero = true
      chart.setOption({ yAxis: { min: (value: { min: number }) => value.min * 0.9 } })
      chart.dispatchAction({ type: "dataZoom", dataZoomIndex: 0, start, end })

      expect(areaPaths(container)).toHaveLength(1)
    })
  })

  // The band carries none of the columns the popup reads (build id, branch, delta), so it must never reach the
  // formatter. It declares `tooltip: { show: false }` for that, which is what keeps the formatter itself free
  // of band-shaped special cases - if ECharts ever stops honouring it, this is where it shows up.
  it("is left out of the axis tooltip", () => {
    let reported: { seriesId?: string }[] = []
    chart.setOption({
      tooltip: {
        trigger: "axis",
        formatter: (params: { seriesId?: string }[]) => {
          reported = params
          return ""
        },
      },
    })
    chart.dispatchAction({ type: "showTip", seriesIndex: 0, dataIndex: 0 })

    expect(reported.map((it) => it.seriesId)).toStrictEqual([MEAN_SERIES_ID])
  })

  it("fades its fill without the invisible lower edge growing one", () => {
    useSettingsStore().fadeOnHover = true
    const fade = new HoverFadeController({ chart, chartContainer: container } as unknown as ChartManager)
    try {
      // hovering a series that is not this one - the band has to fade with its line
      ;(chart as unknown as { trigger: (name: string, params: unknown) => void }).trigger("mouseover", { seriesId: "another-series" })

      const paths = areaPaths(container)
      expect(paths).toHaveLength(1)
      expect(Number(paths[0].getAttribute("fill-opacity"))).toBeLessThan(0.15)
    } finally {
      fade.dispose()
    }
  })
})
