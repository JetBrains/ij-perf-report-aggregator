import { describe, expect, it } from "vitest"
import { columnName, indexSeriesRuns, seriesKey } from "../../../src/components/charts/compareQuery"
import { DataQueryResult, SERIES_NAME_SEPARATOR } from "../../../src/components/common/DataQueryExecutor"
import { DataQueryExecutorConfiguration } from "../../../src/components/common/dataQuery"

// indexSeriesRuns reads the series/measure names and asks the configuration which of them are companions, so
// it gets a real one. No measure here is a companion: every one is its own owner.
function configOf(seriesNames: string[], measureNames: string[]): DataQueryExecutorConfiguration {
  const configuration = new DataQueryExecutorConfiguration()
  configuration.seriesNames = seriesNames
  configuration.measureNames.push(...measureNames)
  configuration.ownerMeasureNames.push(...measureNames)
  return configuration
}

describe("index series runs", () => {
  it("keeps timestamp/value pairs, sorted by time, and resolves keys", () => {
    // One branch + one measure (both seeded), two projects distinguished by series name.
    const data: DataQueryResult = [
      [
        [2000, 1000],
        [20, 10],
      ],
      [[3000], [30]],
    ]
    const indexed = indexSeriesRuns(data, configOf(["p1", "p2"], ["m1", "m1"]), ["b1"], ["p1", "p2"], ["m1"])
    expect(indexed.get(seriesKey("b1", "p1", "m1"))).toStrictEqual([
      { t: 1000, v: 10 },
      { t: 2000, v: 20 },
    ])
    expect(indexed.get(seriesKey("b1", "p2", "m1"))).toStrictEqual([{ t: 3000, v: 30 }])
  })

  it("keys the runs of one project by branch and mode when modes are given", () => {
    // The engine comparison: one branch and one project, so only the mode and the measure name each series.
    const data: DataQueryResult = [
      [[1000], [10]],
      [[1000], [7]],
      [[2000], [20]],
    ]
    const names = [`default${SERIES_NAME_SEPARATOR}m1`, `types2${SERIES_NAME_SEPARATOR}m1`, `default${SERIES_NAME_SEPARATOR}m2`]
    const indexed = indexSeriesRuns(data, configOf(names, ["m1", "m1", "m2"]), ["b1"], ["p1"], ["m1", "m2"], ["default", "types2"])
    expect(indexed.get(seriesKey(columnName("b1", "default"), "p1", "m1"))).toStrictEqual([{ t: 1000, v: 10 }])
    expect(indexed.get(seriesKey(columnName("b1", "types2"), "p1", "m1"))).toStrictEqual([{ t: 1000, v: 7 }])
    expect(indexed.get(seriesKey(columnName("b1", "default"), "p1", "m2"))).toStrictEqual([{ t: 2000, v: 20 }])
    expect(indexed.get(seriesKey(columnName("b1", "types2"), "p1", "m2"))).toBeUndefined()
  })

  it("drops points with a non-finite timestamp or value", () => {
    const data: DataQueryResult = [
      [
        [1000, 2000, 3000],
        [10, Number.NaN, 30],
      ],
    ]
    const indexed = indexSeriesRuns(data, configOf(["p1"], ["m1"]), ["b1"], ["p1"], ["m1"])
    expect(indexed.get(seriesKey("b1", "p1", "m1"))).toStrictEqual([
      { t: 1000, v: 10 },
      { t: 3000, v: 30 },
    ])
  })
})
