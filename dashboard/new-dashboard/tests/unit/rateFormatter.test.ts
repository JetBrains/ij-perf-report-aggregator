import { describe, expect, it } from "vitest"
import { formatMeasureValue, reduceToAxisUnit, resolveMeasureUnit } from "../../src/components/common/formatter"

describe("rate units", () => {
  it.each([
    ["fps", "framesPerSecond"],
    ["pixelsPerSecond", "pixelsPerSecond"],
    ["jOPS", "javaOperationsPerSecond"],
  ] as const)("resolves the value unit %s to %s, whatever the metric name", (valueUnit, unit) => {
    expect(resolveMeasureUnit("dc_boulder_2013-13-30-06-13-17.ser", { valueUnit })).toBe(unit)
  })

  it("formats frames, pixels and Java operations per second", () => {
    expect(formatMeasureValue(59.94, "framesPerSecond")).toBe("59.9 FPS")
    expect(formatMeasureValue(2_700_000_000, "pixelsPerSecond")).toBe("2.7 Gpixels/s")
    expect(formatMeasureValue(48_000, "pixelsPerSecond")).toBe("48 kpixels/s")
    expect(formatMeasureValue(12_345.6, "javaOperationsPerSecond")).toBe("12,346 jOPS")
  })

  it("keeps a rate as the axis unit", () => {
    expect(reduceToAxisUnit(["framesPerSecond", "framesPerSecond"])).toBe("framesPerSecond")
  })

  it("renders a requested rate as a plain ratio while scaling", () => {
    expect(resolveMeasureUnit("CircleTests.ser", { valueUnit: "fps", scaling: true })).toBe("counter")
  })
})
