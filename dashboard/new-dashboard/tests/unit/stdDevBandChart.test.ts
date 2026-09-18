import { createPinia, setActivePinia } from "pinia"
import { afterEach, beforeEach, describe, expect, it } from "vitest"
import { nextTick, shallowRef } from "vue"
import { LineSeriesOption } from "echarts/charts"
import { DatasetOption } from "echarts/types/dist/shared"
import { generateQueries, mergeQueries, DataQueryResult } from "../../src/components/common/DataQueryExecutor"
import { DataQuery, DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"
import { PredefinedMeasureConfigurator } from "../../src/configurators/MeasureConfigurator"
import { useSettingsStore } from "../../src/components/settings/settingsStore"
import { dbTypeStore } from "../../src/shared/dbTypes"
import { parseSeriesId } from "../../src/components/charts/seriesId"
import { configureQueryProducer } from "../../src/configurators/DimensionConfigurator"
import { BranchConfigurator } from "../../src/configurators/BranchConfigurator"
import { indexSeries, indexSeriesRuns, seriesKey } from "../../src/components/charts/compareQuery"

const PROJECT = "caddy/completion/variable"

interface ChartOptions {
  dataset: DatasetOption[]
  series: LineSeriesOption[]
}

/** One measure of one project, as the backend returns it for perfintDev. */
function queryResult(
  values: number[],
  measureName: string,
  project: string = PROJECT,
  timestamps: number[] = values.map((_, index) => 1_000 + index),
  buildIds: number[] = values.map((_, index) => 5_000 + index)
): (string | number)[][] {
  return [timestamps, values, values.map(() => measureName), values.map(() => "d"), values.map(() => "Linux EC2"), buildIds, values.map(() => project), values.map(() => "master")]
}

function configure(measures: string[], projects: string[] = [], withStdDevInterval = true) {
  const configuration = new DataQueryExecutorConfiguration()
  const query = new DataQuery()
  query.db = "perfintDev"
  query.table = "goland"
  const configurator = new PredefinedMeasureConfigurator(shallowRef(measures), shallowRef(true), "line", "auto", {}, null, "item", "lower", withStdDevInterval)
  configurator.configureQuery(query, configuration)
  if (projects.length > 0) configureQueryProducer(configuration, { f: "project" }, projects)
  const queries = generateQueries(query, configuration)
  return { configurator, configuration, queries }
}

async function chart(measures: string[], data: DataQueryResult, projects: string[] = []) {
  const { configurator, configuration, queries } = configure(measures, projects)
  return {
    options: (await configurator.configureChart(data, configuration)) as unknown as ChartOptions,
    queriedMeasures: [...configuration.measureNames],
    queries,
  }
}

function bandSeries(options: ChartOptions): LineSeriesOption[] {
  return options.series.filter((it) => parseSeriesId(it.id as string).role?.startsWith("stdDevBand") === true)
}

/** Whether each query skips the rows whose value is zero. */
function skipsZeroValues(queries: DataQuery[]): boolean[] {
  return queries.map((it) => (it.filters ?? []).some((filter) => filter.o === "!=" && filter.v === 0))
}

function sourceOf(options: ChartOptions, series: LineSeriesOption): (string | number)[][] {
  if (series.data != undefined) {
    const points = series.data as number[][]
    return [points.map((point) => point[0]), points.map((point) => point[1])]
  }
  return (options.dataset[series.datasetIndex as number].source ?? []) as (string | number)[][]
}

describe("standard deviation band on a chart", () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
    dbTypeStore().setDbType("perfintDev", "goland")
    useSettingsStore().stdDevInterval = true
  })
  afterEach(() => {
    localStorage.clear()
  })

  const mean = "completion#mean_value"
  const deviation = "completion#standard_deviation"
  const meanWithDeviation = () => chart([mean], [queryResult([100, 200], mean), queryResult([10, 20], deviation)])

  it("does not fetch deviations for comparison consumers", () => {
    expect(configure([mean], [], false).configuration.measureNames).toStrictEqual([mean])
  })

  it.each([
    { withStdDevInterval: false, expectedEmissions: 1 },
    { withStdDevInterval: true, expectedEmissions: 2 },
  ])("reloads on interval changes only for chart consumers (chart=$withStdDevInterval)", async ({ withStdDevInterval, expectedEmissions }) => {
    const { configurator } = configure([mean], [], withStdDevInterval)
    const emissions: unknown[] = []
    const controller = new AbortController()
    configurator.createObservable().subscribe(
      (value) => {
        emissions.push(value)
      },
      { signal: controller.signal }
    )
    try {
      useSettingsStore().stdDevInterval = false
      await nextTick()
      expect(emissions).toHaveLength(expectedEmissions)
    } finally {
      controller.abort()
    }
  })

  it("excludes implicit deviations from compared values and runs", () => {
    const { configuration } = configure([mean])
    const data = [queryResult([100, 200], mean), queryResult([10, 20], deviation)]
    const key = seriesKey("master", PROJECT, mean)
    expect(indexSeries(data, configuration, ["master"], [PROJECT], [mean]).byKey.get(key)).toStrictEqual([100, 200])
    expect(
      indexSeriesRuns(data, configuration, ["master"], [PROJECT], [mean])
        .get(key)
        ?.map((it) => it.v)
    ).toStrictEqual([100, 200])
  })

  it("compares explicitly selected deviations as an ordinary metric", () => {
    const { configuration } = configure([mean, deviation], [], false)
    const data = [queryResult([100, 200], mean), queryResult([10, 20], deviation)]
    expect(indexSeries(data, configuration, ["master"], [PROJECT], [mean, deviation]).byKey.get(seriesKey("master", PROJECT, deviation))).toStrictEqual([10, 20])
  })

  it("scales the mean and its band together", async () => {
    useSettingsStore().scaling = true
    const { options } = await meanWithDeviation()
    const [lower, fill] = bandSeries(options)
    expect(sourceOf(options, lower)[1][0]).toBeCloseTo(30)
    expect(sourceOf(options, lower)[1][1]).toBeCloseTo(60)
    expect(sourceOf(options, fill)[1][0]).toBeCloseTo(20 / 3)
    expect(sourceOf(options, fill)[1][1]).toBeCloseTo(40 / 3)
  })

  it("fetches the deviation alongside the mean", async () => {
    const { queriedMeasures } = await meanWithDeviation()
    expect(queriedMeasures).toStrictEqual(["completion#mean_value", "completion#standard_deviation"])
  })

  it("draws the deviation as a band around the mean rather than as a line of its own", async () => {
    const { options } = await meanWithDeviation()

    const lines = options.series.filter((it) => parseSeriesId(it.id as string).role?.startsWith("stdDevBand") !== true)
    expect(lines.map((it) => it.id)).toStrictEqual([`completion#mean_value@${PROJECT}`])

    const [lowerEdge, band] = bandSeries(options)
    expect(sourceOf(options, lowerEdge)[1]).toStrictEqual([90, 180])
    expect(sourceOf(options, band)[1]).toStrictEqual([20, 40])
    expect([lowerEdge.stack, band.name]).toStrictEqual([band.stack, lines[0].name])
  })

  it("draws nothing extra while the setting is off", async () => {
    useSettingsStore().stdDevInterval = false
    const { options, queriedMeasures } = await chart(["completion#mean_value"], [queryResult([100, 200], "completion#mean_value")])

    expect(queriedMeasures).toStrictEqual(["completion#mean_value"])
    expect(bandSeries(options)).toStrictEqual([])
  })

  it("keeps each branch's mean and band aligned after query merging", async () => {
    const mean = "completion#mean_value"
    const deviation = "completion#standard_deviation"
    const configuration = new DataQueryExecutorConfiguration()
    const query = new DataQuery()
    query.db = "perfintDev"
    query.table = "goland"

    // Branches are configured before measures in LineChart. Keeping zero deviations prevents merging
    // a mean with its companion, so the executor merges across branches instead.
    const branches = new BranchConfigurator()
    branches.selected.value = ["master", "release"]
    branches.configureQuery(query, configuration)
    const configurator = new PredefinedMeasureConfigurator(shallowRef([mean]), shallowRef(true), "line", "auto", {}, null, "item", "lower", true)
    configurator.configureQuery(query, configuration)
    for (const field of ["machine", "tc_build_id", "project", "branch"]) query.addField(field)

    const queries = mergeQueries(generateQueries(query, configuration), configuration)
    expect(queries.map((it) => it.filters?.find((filter) => filter.f === "measures.name")?.v)).toStrictEqual([mean, deviation])
    expect(queries.map((it) => it.filters?.find((filter) => filter.f === "branch"))).toStrictEqual([
      { f: "branch", v: ["master", "release"], s: true },
      { f: "branch", v: ["master", "release"], s: true },
    ])

    // The backend returns both branches' means, then both deviations, in the split filters' order.
    const releaseMean = queryResult([1000, 2000], mean)
    const releaseDeviation = queryResult([100, 200], deviation)
    releaseMean[7].fill("release")
    releaseDeviation[7].fill("release")
    const data = [queryResult([100, 200], mean), releaseMean, queryResult([10, 20], deviation), releaseDeviation]
    const options = (await configurator.configureChart(data, configuration)) as unknown as ChartOptions

    const lines = options.series.filter((it) => parseSeriesId(it.id as string).role == null)
    expect(lines.map((it) => ({ name: it.name, values: sourceOf(options, it)[1] }))).toStrictEqual([
      { name: "master", values: [100, 200] },
      { name: "release", values: [1000, 2000] },
    ])
    expect(bandSeries(options).map((it) => ({ name: it.name, values: sourceOf(options, it)[1] }))).toStrictEqual([
      { name: "master", values: [90, 180] },
      { name: "master", values: [20, 40] },
      { name: "release", values: [900, 1800] },
      { name: "release", values: [200, 400] },
    ])
  })

  it("leaves a deviation the user selected themselves as a line", async () => {
    const { options, queriedMeasures } = await chart(
      ["completion#mean_value", "completion#standard_deviation"],
      [queryResult([100, 200], "completion#mean_value"), queryResult([10, 20], "completion#standard_deviation")]
    )

    expect(queriedMeasures).toStrictEqual(["completion#mean_value", "completion#standard_deviation"])
    expect(bandSeries(options)).toStrictEqual([])
    expect(options.series).toHaveLength(2)
  })

  it("gives two builds that reported in the same second their own deviation", async () => {
    const sameSecond = [1_000, 1_000]
    const twoBuilds = [5_000, 5_001]
    const { options } = await chart(
      ["completion#mean_value"],
      [queryResult([100, 200], "completion#mean_value", PROJECT, sameSecond, twoBuilds), queryResult([10, 20], "completion#standard_deviation", PROJECT, sameSecond, twoBuilds)]
    )

    const [lowerEdge, band] = bandSeries(options)
    expect(sourceOf(options, lowerEdge)[1]).toStrictEqual([90, 180])
    expect(sourceOf(options, band)[1]).toStrictEqual([20, 40])
  })

  it("pairs by build rather than by arrival order within a second", async () => {
    const sameSecond = [1_000, 1_000]
    const { options } = await chart(
      ["completion#mean_value"],
      [
        queryResult([100, 200], "completion#mean_value", PROJECT, sameSecond, [5_000, 5_001]),
        // the same pairing - build 5000 with 10, build 5001 with 20 - returned the other way round
        queryResult([20, 10], "completion#standard_deviation", PROJECT, sameSecond, [5_001, 5_000]),
      ]
    )

    const [lowerEdge, band] = bandSeries(options)
    expect(sourceOf(options, lowerEdge)[1]).toStrictEqual([90, 180])
    expect(sourceOf(options, band)[1]).toStrictEqual([20, 40])
  })

  it("pairs the repeats of a build that reported the measure twice in order", async () => {
    const sameSecond = [1_000, 1_000]
    const oneBuild = [5_000, 5_000]
    const { options } = await chart(
      ["completion#mean_value"],
      [queryResult([100, 200], "completion#mean_value", PROJECT, sameSecond, oneBuild), queryResult([10, 20], "completion#standard_deviation", PROJECT, sameSecond, oneBuild)]
    )

    expect(sourceOf(options, bandSeries(options)[1])[1]).toStrictEqual([20, 40])
  })

  it("keeps the pairing when an outlier is removed from between two points of one build", async () => {
    useSettingsStore().removeOutliers = true
    const means = [100, 101, 1000, 102, 103, 104, 105]
    // the outlier and the point after it come from the same build, in the same second
    const times = [1_000, 1_001, 1_002, 1_002, 1_003, 1_004, 1_005]
    const builds = [5_000, 5_001, 5_002, 5_002, 5_003, 5_004, 5_005]
    const { options } = await chart(
      ["completion#mean_value"],
      [queryResult(means, "completion#mean_value", PROJECT, times, builds), queryResult([1, 1, 1000, 3, 1, 1, 1], "completion#standard_deviation", PROJECT, times, builds)]
    )

    const [lowerEdge, band] = bandSeries(options)
    expect(sourceOf(options, lowerEdge)[0]).toHaveLength(means.length - 1)
    // the surviving 102 keeps its own ±3, and every other point its ±1
    expect(sourceOf(options, band)[1]).toStrictEqual([2, 2, 6, 2, 2, 2])
    expect(sourceOf(options, lowerEdge)[1]).toStrictEqual([99, 100, 99, 102, 103, 104])
  })

  it("keeps the zero deviations of a build that reported the measure twice", async () => {
    const sameSecond = [1_000, 1_000]
    const oneBuild = [5_000, 5_000]
    const { options, queries } = await chart(
      ["completion#mean_value"],
      [queryResult([100, 200], "completion#mean_value", PROJECT, sameSecond, oneBuild), queryResult([0, 20], "completion#standard_deviation", PROJECT, sameSecond, oneBuild)]
    )

    // the mean skips its own zeros, the deviation is asked for all of its values
    expect(skipsZeroValues(queries)).toStrictEqual([true, false])

    const [lowerEdge, band] = bandSeries(options)
    expect(sourceOf(options, lowerEdge)[1]).toStrictEqual([100, 180])
    expect(sourceOf(options, band)[1]).toStrictEqual([0, 40])
  })

  it("draws no band where the deviations of a point cannot be told apart", async () => {
    const { options } = await chart(
      ["completion#mean_value"],
      [
        queryResult([100, 300], "completion#mean_value", PROJECT, [1_000, 1_001], [5_000, 5_001]),
        // build 5000 reported the measure twice, but only the second mean survived the zero filter
        queryResult([10, 20, 30], "completion#standard_deviation", PROJECT, [1_000, 1_000, 1_001], [5_000, 5_000, 5_001]),
      ]
    )

    const [lowerEdge, band] = bandSeries(options)
    // the ambiguous point keeps its mean and loses its band; the unambiguous one is wrapped as usual
    expect(sourceOf(options, lowerEdge)[1]).toStrictEqual([100, 270])
    expect(sourceOf(options, band)[1]).toStrictEqual([0, 60])
  })

  it("wraps one measure drawn across several projects in a band per project", async () => {
    const projects = [PROJECT, "permify/completion/method"]
    const { options } = await chart(
      ["completion#mean_value"],
      [
        queryResult([100, 200], "completion#mean_value", projects[0]),
        queryResult([1000, 2000], "completion#mean_value", projects[1]),
        queryResult([10, 20], "completion#standard_deviation", projects[0]),
        queryResult([100, 200], "completion#standard_deviation", projects[1]),
      ],
      projects
    )

    const lines = options.series.filter((it) => parseSeriesId(it.id as string).role?.startsWith("stdDevBand") !== true)
    expect(lines.map((it) => it.name)).toStrictEqual(projects)

    const bands = bandSeries(options)
    expect(bands.map((it) => it.name)).toStrictEqual([projects[0], projects[0], projects[1], projects[1]])
    expect(sourceOf(options, bands[0])[1]).toStrictEqual([90, 180])
    expect(sourceOf(options, bands[2])[1]).toStrictEqual([900, 1800])
  })

  it("wraps each measure of a chart in its own band", async () => {
    const { options, queriedMeasures } = await chart(
      ["completion#mean_value", "localInspections#mean_value"],
      [
        queryResult([100, 200], "completion#mean_value"),
        queryResult([1000, 2000], "localInspections#mean_value"),
        queryResult([10, 20], "completion#standard_deviation"),
        queryResult([100, 200], "localInspections#standard_deviation"),
      ]
    )

    expect(queriedMeasures).toStrictEqual(["completion#mean_value", "localInspections#mean_value", "completion#standard_deviation", "localInspections#standard_deviation"])

    const bands = bandSeries(options)
    expect(bands.map((it) => it.id)).toStrictEqual([
      "completion#mean_valuestdDevBandLowerEdge",
      "completion#mean_valuestdDevBandFill",
      "localInspections#mean_valuestdDevBandLowerEdge",
      "localInspections#mean_valuestdDevBandFill",
    ])
    expect(sourceOf(options, bands[0])[1]).toStrictEqual([90, 180])
    expect(sourceOf(options, bands[2])[1]).toStrictEqual([900, 1800])
  })
})
