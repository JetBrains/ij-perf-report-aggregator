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

function queryResult(values: number[], measureName: string): (string | number)[][] {
  return [
    values.map((_, index) => 1_000 + index),
    values,
    values.map(() => measureName),
    values.map(() => "d"),
    values.map(() => "Linux EC2"),
    values.map((_, index) => 5_000 + index),
    values.map(() => PROJECT),
    values.map(() => "master"),
  ]
}

async function chart(measures: string[], stacked: string[]) {
  const configuration = new DataQueryExecutorConfiguration()
  const query = new DataQuery()
  query.db = "perfintDev"
  query.table = "idea"
  const configurator = new PredefinedMeasureConfigurator(shallowRef(measures), shallowRef(true), "line", "auto", {}, null, "item", "lower", false, "sync", stacked)
  configurator.configureQuery(query, configuration)
  // A second producer makes the series measure name composite ("project - metric"), as on a real dashboard.
  configureQueryProducer(configuration, { f: "project" }, [PROJECT], new Map([[PROJECT, ""]]))
  generateQueries(query, configuration)
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
})
