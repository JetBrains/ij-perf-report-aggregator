import { describe, expect, it } from "vitest"
import { DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"
import { mergeSeries } from "../../src/configurators/MeasureConfigurator"

describe("mergeSeries()", () => {
  it("sorts a merged alias series by time instead of concatenation order", () => {
    // Simulates aliasing an old project name into a new one (e.g. "wideTree" -> "wide tree"):
    // the new name's query runs first and currently only has a couple of recent rows, the old
    // name's query returns the bulk of the (earlier) history. Both share the aliased series name.
    const newNameChunk: (string | number)[][] = [[2000], [120]]
    const oldNameChunk: (string | number)[][] = [
      [1000, 1001, 1002],
      [200, 190, 180],
    ]

    const configuration = new DataQueryExecutorConfiguration()
    configuration.seriesNames = ["wide tree", "wide tree"]
    configuration.measureNames.push("fleet.test", "fleet.test")

    const result = mergeSeries([newNameChunk, oldNameChunk], configuration)

    expect(result.data).toHaveLength(1)
    expect(result.data[0][0]).toStrictEqual([1000, 1001, 1002, 2000])
    expect(result.data[0][1]).toStrictEqual([200, 190, 180, 120])
  })

  it("leaves a single, already time-sorted series untouched", () => {
    const chunk: (string | number)[][] = [
      [1000, 1001, 1002],
      [200, 190, 180],
    ]

    const configuration = new DataQueryExecutorConfiguration()
    configuration.seriesNames = ["wide tree"]
    configuration.measureNames.push("fleet.test")

    const result = mergeSeries([chunk], configuration)

    expect(result.data[0][0]).toStrictEqual([1000, 1001, 1002])
    expect(result.data[0][1]).toStrictEqual([200, 190, 180])
  })
})
