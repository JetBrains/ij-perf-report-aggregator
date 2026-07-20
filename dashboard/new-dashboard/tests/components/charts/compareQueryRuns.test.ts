import { describe, expect, it } from "vitest"
import { indexSeriesRuns, seriesKey } from "../../../src/components/charts/compareQuery"
import { DataQueryResult } from "../../../src/components/common/DataQueryExecutor"
import { DataQueryExecutorConfiguration } from "../../../src/components/common/dataQuery"

// indexSeriesRuns only reads seriesNames/measureNames off the configuration, so a plain object suffices.
function configOf(seriesNames: string[], measureNames: string[]): DataQueryExecutorConfiguration {
  return { seriesNames, measureNames } as unknown as DataQueryExecutorConfiguration
}

describe("index series runs", () => {
  it("keeps timestamp/value pairs, sorted by time, and resolves keys", () => {
    // One branch + one measure (both seeded), two projects distinguished by series name.
    const data: DataQueryResult = [
      [
        [2000, 1000],
        [20, 10],
      ],
      [
        [3000],
        [30],
      ],
    ]
    const indexed = indexSeriesRuns(data, configOf(["p1", "p2"], ["m1", "m1"]), ["b1"], ["p1", "p2"], ["m1"])
    expect(indexed.get(seriesKey("b1", "p1", "m1"))).toStrictEqual([
      { t: 1000, v: 10 },
      { t: 2000, v: 20 },
    ])
    expect(indexed.get(seriesKey("b1", "p2", "m1"))).toStrictEqual([{ t: 3000, v: 30 }])
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
