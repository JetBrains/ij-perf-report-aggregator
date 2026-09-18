import { describe, expect, it } from "vitest"
import { buildStdDevBand, getStdDevBandSeries, getStdDevMeasureName } from "../../src/components/charts/stdDevInterval"
import { parseSeriesId } from "../../src/components/charts/seriesId"

const SERIES_ID = "completion#mean_value@caddy/completion/variable"

/** The deviations are already paired with the values by the caller, so the helper takes them index by index. */
function band(values: number[], stdDevs: (number | undefined)[], valueScale = 1) {
  const timestamps = values.map((_, index) => 1_000 + index)
  return buildStdDevBand({
    timestamps,
    values,
    stdDevs: stdDevs.map((stdDev) => stdDev ?? 0),
    valueScale,
  })
}

describe("standard deviation companion measure", () => {
  it("is derived from the mean it belongs to", () => {
    expect(getStdDevMeasureName("completion#mean_value")).toBe("completion#standard_deviation")
    expect(getStdDevMeasureName("completion#firstElementShown#mean_value")).toBe("completion#firstElementShown#standard_deviation")
  })

  it.each([["completion#median_value"], ["completion#standard_deviation"], ["completion"], ["mean_value"], ["typing"]])("is not derived from %s", (measure) => {
    expect(getStdDevMeasureName(measure)).toBeNull()
  })
})

describe("standard deviation band", () => {
  it("spans one deviation either side of the mean", () => {
    expect(band([100, 200], [10, 20])).toStrictEqual({
      timestamps: [1000, 1001],
      lowerEdge: [90, 180],
      width: [20, 40],
    })
  })

  it("stops at zero rather than stretching the axis into negative values", () => {
    // A mean of 30 ms with a 50 ms deviation reaches below zero, where a duration cannot be.
    expect(band([30], [50])).toMatchObject({ lowerEdge: [0], width: [80] })
  })

  it("pinches to the line where no deviation was reported", () => {
    expect(band([100, 200, 300], [10, undefined, 30])).toMatchObject({
      lowerEdge: [90, 200, 270],
      width: [20, 0, 60],
    })
  })

  it("follows the values onto the scaled axis", () => {
    expect(band([100, 200], [10, 20], 0.5)).toMatchObject({ lowerEdge: [95, 190], width: [10, 20] })
  })

  it("skips points the series has no value for", () => {
    expect(band([100, Number.NaN, 300], [10, 20, 30])).toMatchObject({
      timestamps: [1000, 1002],
      lowerEdge: [90, 270],
    })
  })

  it("is not drawn when every deviation is zero or missing", () => {
    expect(band([100, 200], [0, undefined])).toBeNull()
    expect(band([], [])).toBeNull()
  })
})

describe("standard deviation band series", () => {
  const series = getStdDevBandSeries(band([100, 200], [10, 20])!, SERIES_ID, "caddy/completion/variable")

  it("stacks the width onto the lower edge", () => {
    expect(series.map((it) => it.data)).toStrictEqual([
      [
        [1000, 90],
        [1001, 180],
      ],
      [
        [1000, 20],
        [1001, 40],
      ],
    ])
    expect(new Set(series.map((it) => it.stack)).size).toBe(1)
    // "samesign" would refuse to stack onto a lower edge sitting at zero.
    expect(series.map((it) => it.stackStrategy)).toStrictEqual(["all", "all"])
  })

  it("only fills between the two edges", () => {
    expect(series[0].areaStyle).toBeUndefined()
    expect(series[1].areaStyle?.opacity).toBeGreaterThan(0)
    expect(series.map((it) => it.lineStyle?.opacity)).toStrictEqual([0, 0])
  })

  it("shares the name, and so the color and the legend entry, of the series it wraps", () => {
    expect(series.map((it) => it.name)).toStrictEqual(["caddy/completion/variable", "caddy/completion/variable"])
  })

  it("stays out of the way of the lines", () => {
    expect(series.map((it) => [it.silent, it.z])).toStrictEqual([
      [true, 1],
      [true, 1],
    ])
  })

  it("declares which half of the band each series is", () => {
    expect(series.map((it) => parseSeriesId(it.id as string).role)).toStrictEqual(["stdDevBandLowerEdge", "stdDevBandFill"])
  })

  it("is resolved back to the series it wraps", () => {
    expect(series.map((it) => parseSeriesId(it.id as string).ownerSeriesId)).toStrictEqual([SERIES_ID, SERIES_ID])
  })

  it("does not mistake a line, or a deviation the user selected, for a band", () => {
    expect([SERIES_ID, "completion#standard_deviation@project"].map((it) => parseSeriesId(it).role)).toStrictEqual([null, null])
  })

  it("is left out of the tooltip by the option itself, rather than by the formatter", () => {
    expect(series.map((it) => it.tooltip?.show)).toStrictEqual([false, false])
  })
})
