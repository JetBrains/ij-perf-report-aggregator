import { describe, expect, it } from "vitest"
import { getMinYAxisTop } from "../../src/components/charts/yAxisRange"

describe("minimum y-axis top", () => {
  it.each([["ui.lagging#average"], ["ui.lagging#max"], ["ui.lagging#max_value"], ["ui.lagging#sum"]])("frames the lagging duration %s up to 5 s", (measure) => {
    expect(getMinYAxisTop([measure])).toBe(3000)
  })

  it("frames a chart combining several lagging durations", () => {
    expect(getMinYAxisTop(["ui.lagging#average", "ui.lagging#max"])).toBe(3000)
  })

  it("frames the lagging percentage share up to 50", () => {
    expect(getMinYAxisTop(["ui.lagging#percentage_share"])).toBe(10)
  })

  it("leaves the lagging count autoscaled", () => {
    expect(getMinYAxisTop(["ui.lagging#count"])).toBeUndefined()
  })

  it("leaves a chart mixing lagging durations with the percentage share autoscaled", () => {
    expect(getMinYAxisTop(["ui.lagging#average", "ui.lagging#max", "ui.lagging#percentage_share"])).toBeUndefined()
  })

  it.each([[["typing#max_awt_delay"]], [["test#max_awt_delay"]], [["ui.latency#max_value"]], [[]]])("leaves %s autoscaled", (measures) => {
    expect(getMinYAxisTop(measures)).toBeUndefined()
  })
})
