import { describe, expect, it } from "vitest"
import { SERIES_NAME_SEPARATOR } from "../../../src/components/common/DataQueryExecutor"
import { resolveSeriesStyles, splitSeriesName } from "../../../src/components/charts/seriesStyle"

const MEASURES = ["completion#mean_value", "fus_completion_duration_90p"]

function name(project: string, measure: string): string {
  return project + SERIES_NAME_SEPARATOR + measure
}

describe("series name split", () => {
  it("separates the project from the measure", () => {
    expect(splitSeriesName(name("caddy/variable", "completion#mean_value"), MEASURES)).toStrictEqual({ group: "caddy/variable", measure: "completion#mean_value" })
  })

  it("prefers the longest measure that ends the name", () => {
    expect(splitSeriesName(name("act", "typing#latency"), ["latency", "typing#latency"])).toStrictEqual({ group: "act", measure: "typing#latency" })
  })

  it("a name without a project is the measure alone", () => {
    expect(splitSeriesName("completion#mean_value", MEASURES)).toStrictEqual({ group: "", measure: "completion#mean_value" })
  })

  it("a name that ends with no measure is not split", () => {
    expect(splitSeriesName("caddy/variable", MEASURES)).toBeNull()
  })
})

describe("series styles", () => {
  it("colors by project and dashes by measure", () => {
    const series = [
      { group: "caddy", measure: MEASURES[0] },
      { group: "caddy", measure: MEASURES[1] },
      { group: "permify", measure: MEASURES[0] },
      { group: "permify", measure: MEASURES[1] },
    ]
    const styles = resolveSeriesStyles(series, MEASURES, false)!
    expect(styles[0]?.color).toBe(styles[1]?.color)
    expect(styles[2]?.color).toBe(styles[3]?.color)
    expect(styles[0]?.color).not.toBe(styles[2]?.color)
    expect(styles.map((it) => it?.lineType)).toStrictEqual(["solid", "dashed", "solid", "dashed"])
  })

  it("leaves a companion series unstyled", () => {
    const styles = resolveSeriesStyles([{ group: "caddy", measure: MEASURES[0] }, null, { group: "permify", measure: MEASURES[1] }], MEASURES, false)!
    expect(styles[1]).toBeNull()
  })

  it("a measure without data does not take a line type", () => {
    const series = [
      { group: "caddy", measure: MEASURES[1] },
      { group: "permify", measure: "x" },
    ]
    expect(resolveSeriesStyles(series, [MEASURES[0], MEASURES[1], "x"], false)?.map((it) => it?.lineType)).toStrictEqual(["solid", "dashed"])
  })

  it("keeps the default palette for one project", () => {
    expect(
      resolveSeriesStyles(
        [
          { group: "caddy", measure: MEASURES[0] },
          { group: "caddy", measure: MEASURES[1] },
        ],
        MEASURES,
        false
      )
    ).toBeNull()
  })

  it("keeps the default palette for one measure", () => {
    expect(
      resolveSeriesStyles(
        [
          { group: "caddy", measure: MEASURES[0] },
          { group: "permify", measure: MEASURES[0] },
        ],
        MEASURES,
        false
      )
    ).toBeNull()
  })

  it("keeps the default palette for four measures", () => {
    const measures = ["a", "b", "c", "d"]
    const series = measures.flatMap((measure) => [
      { group: "caddy", measure },
      { group: "permify", measure },
    ])
    expect(resolveSeriesStyles(series, measures, false)).toBeNull()
  })
})

describe("series styles past the series cap", () => {
  it("keeps the default palette for more than six series", () => {
    const series = ["caddy", "permify", "act", "river"].flatMap((group) => MEASURES.map((measure) => ({ group, measure })))
    expect(resolveSeriesStyles(series, MEASURES, false)).toBeNull()
  })
})
