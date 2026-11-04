import { describe, expect, it } from "vitest"
import { columnName, indexSeries, seriesKey } from "../../src/components/charts/compareQuery"
import { DataQueryResult, SERIES_NAME_SEPARATOR } from "../../src/components/common/DataQueryExecutor"
import { DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"

const KUBERNETES = "kubernetes/highlighting"
const DOCKER = "docker/highlighting"
const METRIC = "firstCodeAnalysis"

// One series per name, each with the given values; the series-name parts are joined as DataQueryExecutor joins them.
function series(entries: [parts: string[], values: number[]][]): { data: DataQueryResult; configuration: DataQueryExecutorConfiguration } {
  const configuration = {
    seriesNames: entries.map(([parts]) => parts.join(SERIES_NAME_SEPARATOR)),
    measureNames: entries.map(() => METRIC),
    isCompanion: () => false,
  } as unknown as DataQueryExecutorConfiguration
  const data = entries.map(([, values]) => [values.map((_, index) => 1_000 + index), values]) as unknown as DataQueryResult
  return { data, configuration }
}

describe("compare table series index", () => {
  it("takes the mode names as the columns when a single branch is selected", () => {
    const { data, configuration } = series([
      [["master", "default", KUBERNETES], [100]],
      [["master", "legacy_and_new", KUBERNETES], [200]],
      [["master", "default", DOCKER], [50]],
    ])
    const indexed = indexSeries(data, configuration, ["default", "legacy_and_new"], [KUBERNETES, DOCKER], [METRIC])
    expect(indexed.byKey.get(seriesKey("default", KUBERNETES, METRIC))).toStrictEqual([100])
    expect(indexed.byKey.get(seriesKey("legacy_and_new", KUBERNETES, METRIC))).toStrictEqual([200])
  })

  it("keeps the runs of each mode apart when two branches and two modes are selected", () => {
    const { data, configuration } = series([
      [["master", "default", KUBERNETES], [100]],
      [["master", "legacy_and_new", KUBERNETES], [200]],
      [["261", "default", KUBERNETES], [90]],
      [["261", "legacy_and_new", KUBERNETES], [180]],
      [["master", "default", DOCKER], [50]],
    ])
    const indexed = indexSeries(data, configuration, ["261", "master"], [KUBERNETES, DOCKER], [METRIC], ["default", "legacy_and_new"])
    expect(indexed.byKey.get(seriesKey(columnName("master", "default"), KUBERNETES, METRIC))).toStrictEqual([100])
    expect(indexed.byKey.get(seriesKey(columnName("master", "legacy_and_new"), KUBERNETES, METRIC))).toStrictEqual([200])
    expect(indexed.byKey.get(seriesKey(columnName("261", "default"), KUBERNETES, METRIC))).toStrictEqual([90])
  })
})
