import { createPinia, setActivePinia } from "pinia"
import { beforeEach, describe, expect, it } from "vitest"
import { shallowRef } from "vue"
import { LineSeriesOption } from "echarts/charts"
import { DataQueryResult, generateQueries } from "../../src/components/common/DataQueryExecutor"
import { DataQuery, DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"
import { PredefinedMeasureConfigurator } from "../../src/configurators/MeasureConfigurator"
import { dbTypeStore } from "../../src/shared/dbTypes"
import { configureQueryProducer } from "../../src/configurators/DimensionConfigurator"

const PHASE_A = "maven.projects.processor.reading.task"
const PHASE_B = "maven.projects.processor.resolving.task"
const TOTAL = "maven.import.stats.sync.project.task"
const PROJECT = "project-import-maven-5000-modules/fastInstaller"

/** One measure of one project, as the backend returns it; `builds` picks which builds reported it. */
function queryResult(values: number[], measureName: string, project: string = PROJECT, builds: number[] = values.map((_, index) => index)): (string | number)[][] {
  return [
    builds.map((build) => 1_000 + build),
    values,
    values.map(() => measureName),
    values.map(() => "d"),
    values.map(() => "Linux EC2"),
    builds.map((build) => 5_000 + build),
    values.map(() => project),
    values.map(() => "master"),
  ]
}

function configure(measures: string[], stacked: string[], projects: string[]) {
  const configuration = new DataQueryExecutorConfiguration()
  const query = new DataQuery()
  query.db = "perfintDev"
  query.table = "idea"
  const configurator = new PredefinedMeasureConfigurator(shallowRef(measures), shallowRef(true), "line", "auto", {}, null, "item", "lower", false, "sync", stacked)
  configurator.configureQuery(query, configuration)
  // A second producer makes the series measure name composite ("project - metric"), as on a real dashboard.
  const aliases = projects.length === 1 ? new Map([[projects[0], ""]]) : undefined
  configureQueryProducer(configuration, { f: "project" }, projects, aliases)
  generateQueries(query, configuration)
  return { configurator, configuration }
}

async function chart(measures: string[], stacked: string[]) {
  const { configurator, configuration } = configure(measures, stacked, [PROJECT])
  // Query producers run measure by measure; the chart data must come back in the same order.
  const data: DataQueryResult = [...configuration.measureNames].map((measure) => queryResult([10, 20, 30], measure))
  const options = (await configurator.configureChart(data, configuration)) as unknown as { series: LineSeriesOption[] }
  // The series id is "<measureName>@<seriesName>"; the series name is the measure on a one-project chart.
  return new Map(options.series.map((series) => [(series.id as string).split("@")[1], series.stack]))
}

describe("stacked phases chart", () => {
  beforeEach(() => {
    localStorage.clear()
    setActivePinia(createPinia())
    dbTypeStore().setDbType("perfintDev", "idea")
  })

  it("stacks phases but leaves the reported total as a plain line", async () => {
    const stacks = await chart([PHASE_A, PHASE_B, TOTAL], [PHASE_B, PHASE_A])
    expect(stacks.get(PHASE_A)).toBe("sync")
    expect(stacks.get(PHASE_B)).toBe("sync")
    expect(stacks.get(TOTAL)).toBeUndefined()
    // Reference line first, then the phases in the declared (not alphabetical) order.
    expect([...stacks.keys()]).toStrictEqual([TOTAL, PHASE_B, PHASE_A])
  })

  it("keeps one stack per project, so projects are not summed together", async () => {
    const other = "project-import-maven-1000-modules/fastInstaller"
    const { configurator, configuration } = configure([PHASE_A, PHASE_B], [PHASE_A, PHASE_B], [PROJECT, other])
    const data: DataQueryResult = [...configuration.measureNames].map((_, index) => queryResult([10, 20, 30], configuration.measureNames[index], [PROJECT, other][index % 2]))
    const options = (await configurator.configureChart(data, configuration)) as unknown as { series: LineSeriesOption[] }
    const stacks = options.series.map((series) => series.stack)
    expect(new Set(stacks).size).toBe(2)
  })

  it("aligns stacked phases on the build, so a missing point is not added to another build", async () => {
    const { configurator, configuration } = configure([PHASE_A, PHASE_B], [PHASE_A, PHASE_B], [PROJECT])
    // PHASE_B misses the middle build: without alignment its 1000 would land on build 1 and sum to 1003.
    const data: DataQueryResult = [queryResult([1, 2, 3], PHASE_A), queryResult([10, 1000], PHASE_B, PROJECT, [0, 2])]
    const options = (await configurator.configureChart(data, configuration)) as unknown as {
      series: LineSeriesOption[]
      dataset: { source: (number | null)[][] }[]
    }
    const sourceOf = (measure: string) => options.dataset[options.series.find((series) => (series.id as string).includes(measure))?.datasetIndex as number].source
    expect(sourceOf(PHASE_A)[0]).toStrictEqual([1000, 1001, 1002])
    expect(sourceOf(PHASE_B)[0]).toStrictEqual([1000, 1001, 1002])
    expect(sourceOf(PHASE_B)[1]).toStrictEqual([10, null, 1000])
  })
})
