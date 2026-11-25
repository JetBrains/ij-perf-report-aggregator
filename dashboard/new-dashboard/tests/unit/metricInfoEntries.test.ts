import { describe, expect, it } from "vitest"
import { metricInfoEntries } from "../../src/components/charts/metricInfoEntries"

const options = { valueUnit: "auto" as const, label: (measure: string) => measure }

describe("info popup entries", () => {
  it("names every measure that shares a description", () => {
    const entries = metricInfoEntries(["globalInspections#Annotator.ms", "globalInspections#GoMaybeNil.ms"], options)
    expect(entries.map((entry) => entry.labels)).toStrictEqual([["globalInspections#Annotator.ms", "globalInspections#GoMaybeNil.ms"]])
    expect(entries[0].unit).toBe("milliseconds")
  })

  it("gives measures with different descriptions their own entries", () => {
    const entries = metricInfoEntries(["globalInspections#GoMaybeNil.ms", "globalInspections#GoMaybeNil.problems"], options)
    expect(entries.map((entry) => [entry.labels, entry.unit, entry.betterDirection])).toStrictEqual([
      [["globalInspections#GoMaybeNil.ms"], "milliseconds", "lower"],
      [["globalInspections#GoMaybeNil.problems"], "counter", "stable"],
    ])
  })

  it("lets the description and direction of the chart override those of its measures", () => {
    const entries = metricInfoEntries(["globalInspections#GoMaybeNil.ms", "globalInspections#GoMaybeNil.problems"], {
      ...options,
      description: "chart",
      betterDirection: "none",
    })
    expect(entries.map((entry) => [entry.description, entry.betterDirection])).toStrictEqual([
      ["chart", "none"],
      ["chart", "none"],
    ])
  })
})
