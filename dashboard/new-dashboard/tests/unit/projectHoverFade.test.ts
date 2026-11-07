import { createPinia, setActivePinia } from "pinia"
import { afterEach, beforeEach, describe, expect, it } from "vitest"
import { shallowRef } from "vue"
import { EChartsType, init, use } from "echarts/core"
import { SVGRenderer } from "echarts/renderers"
// for the side effect: ChartManager is where the app registers its ECharts components
import "../../src/components/charts/ChartManager"
import type { ChartManager } from "../../src/components/charts/ChartManager"
import { HOVER_GROUP_KEY, HoverFadeController } from "../../src/components/charts/hoverFade"
import { STD_DEV_BAND_OPACITY } from "../../src/components/charts/stdDevInterval"
import { generateQueries, DataQueryResult } from "../../src/components/common/DataQueryExecutor"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter } from "../../src/components/common/dataQuery"
import { configureQueryProducer } from "../../src/configurators/DimensionConfigurator"
import { PredefinedMeasureConfigurator } from "../../src/configurators/MeasureConfigurator"
import { useSettingsStore } from "../../src/components/settings/settingsStore"
import { dbTypeStore } from "../../src/shared/dbTypes"

// happy-dom has no canvas 2D context, so the SVG renderer stands in for the app's canvas one.
use([SVGRenderer])

const PROJECTS = ["river/debug", "trufflehog/debug"]
const MEAN = "completion#mean_value"
const P90 = "fus_completion_duration_90p"
const STD_DEV = "completion#standard_deviation"
const TIMESTAMPS = [1_000, 1_001]

interface SeriesOption {
  id: string
  color?: string
  lineStyle?: { type?: string; opacity?: number }
  areaStyle?: { opacity?: number }
  [HOVER_GROUP_KEY]?: string
}

function queryResult(project: string, measure: string): (string | number)[][] {
  const values = measure === STD_DEV ? [10, 20] : [100, 200]
  return [
    TIMESTAMPS,
    values,
    values.map(() => measure),
    values.map(() => "d"),
    values.map(() => "Linux EC2"),
    values.map((_, index) => 5_000 + index),
    values.map(() => project),
    values.map(() => "master"),
  ]
}

/** The option the configurator builds for `measures` of both projects, with one query result per series. */
function chartOptions(measures: string[]) {
  const configuration = new DataQueryExecutorConfiguration()
  const query = new DataQuery()
  query.db = "perfintDev"
  query.table = "goland"

  // the order the app uses: the project dimension produces its series before the measure configurator does
  const filter: DataQueryFilter = { f: "project", v: PROJECTS[0] }
  configureQueryProducer(configuration, filter, PROJECTS)
  query.addFilter(filter)
  const configurator = new PredefinedMeasureConfigurator(shallowRef(measures), shallowRef(true), "line", "auto", {}, null, "item", "lower", true)
  if (!configurator.configureQuery(query, configuration)) {
    throw new Error("the measures were not configured into the query")
  }
  generateQueries(query, configuration)

  const data: DataQueryResult = configuration.seriesNames.map((seriesName, index) => {
    const project = PROJECTS.find((it) => seriesName.startsWith(it)) as string
    return queryResult(project, configuration.measureNames[index].split(" – ").at(-1) as string)
  })
  return configurator.configureChart(data, configuration)
}

function seriesOf(chart: EChartsType): SeriesOption[] {
  return (chart.getOption() as { series: SeriesOption[] }).series
}

function lineOf(chart: EChartsType, project: string, measure: string): SeriesOption {
  // The id leads with the measure part of every producer, and the project producer contributes the first measure.
  const series = seriesOf(chart).find((it) => it.id.endsWith(`${measure}@${project} – ${measure}`))
  if (series == undefined)
    throw new Error(
      `no series of ${measure} for ${project}: ${seriesOf(chart)
        .map((it) => it.id)
        .join(", ")}`
    )
  return series
}

describe("hover fade by project", () => {
  let container: HTMLElement
  let chart: EChartsType
  let fade: HoverFadeController

  async function render(measures: string[]) {
    chart.setOption({ animation: false, xAxis: { type: "time" }, yAxis: { type: "value" } })
    chart.setOption(await chartOptions(measures), { replaceMerge: ["dataset", "series"] })
    fade = new HoverFadeController({ chart, chartContainer: container } as unknown as ChartManager)
  }

  function hover(seriesId: string) {
    ;(chart as unknown as { trigger: (name: string, params: unknown) => void }).trigger("mouseover", { seriesId })
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    dbTypeStore().setDbType("perfintDev", "goland")
    useSettingsStore().stdDevInterval = true
    useSettingsStore().fadeOnHover = true

    container = document.createElement("div")
    document.body.append(container)
    chart = init(container, undefined, { renderer: "svg", width: 600, height: 400 })
  })

  afterEach(() => {
    fade.dispose()
    chart.dispose()
    container.remove()
  })

  it("colors by project and draws each measure with its own line type", async () => {
    await render([MEAN, P90])

    expect(lineOf(chart, PROJECTS[0], MEAN).color).toBe(lineOf(chart, PROJECTS[0], P90).color)
    expect(lineOf(chart, PROJECTS[0], MEAN).color).not.toBe(lineOf(chart, PROJECTS[1], MEAN).color)
    expect([lineOf(chart, PROJECTS[0], MEAN).lineStyle?.type, lineOf(chart, PROJECTS[0], P90).lineStyle?.type]).toStrictEqual(["solid", "dashed"])
  })

  /** The line opacity of the mean and the p90 line of `project`, then the fill opacity of its mean's band. */
  function opacities(project: string): (number | undefined)[] {
    const band = seriesOf(chart).find((it) => it.id.endsWith(`${MEAN}@${project} – ${MEAN}stdDevBandFill`))
    return [lineOf(chart, project, MEAN).lineStyle?.opacity, lineOf(chart, project, P90).lineStyle?.opacity, band?.areaStyle?.opacity]
  }

  it("keeps every line and band of the hovered project", async () => {
    await render([MEAN, P90])
    hover(lineOf(chart, PROJECTS[0], P90).id)

    // the mean, its two band halves and the p90 line
    expect(seriesOf(chart).filter((it) => it[HOVER_GROUP_KEY] === PROJECTS[0])).toHaveLength(4)
    expect(opacities(PROJECTS[0])).toStrictEqual([1, 1, STD_DEV_BAND_OPACITY])
  })

  it("fades every line and band of the other project", async () => {
    await render([MEAN, P90])
    hover(lineOf(chart, PROJECTS[0], P90).id)

    expect(opacities(PROJECTS[1])).toStrictEqual([0.2, 0.2, STD_DEV_BAND_OPACITY * 0.2])
  })

  it("restores every series when the pointer leaves the chart", async () => {
    await render([MEAN, P90])
    hover(lineOf(chart, PROJECTS[0], P90).id)
    ;(chart as unknown as { trigger: (name: string, params: unknown) => void }).trigger("globalout", {})

    expect(lineOf(chart, PROJECTS[1], MEAN).lineStyle?.opacity).toBe(1)
  })

  it("fades a chart of one measure by series, as before", async () => {
    await render([P90])

    const series = seriesOf(chart)
    expect(series.filter((it) => it[HOVER_GROUP_KEY] != undefined)).toHaveLength(0)
    hover(series[0].id)
    expect(seriesOf(chart).map((it) => it.lineStyle?.opacity)).toStrictEqual([1, 0.2])
  })
})
