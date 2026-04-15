import { describe, expect, it } from "vitest"
import { scaleToMedian } from "../../src/components/settings/configurators/ScalingConfigurator"

describe("median scaling", () => {
  it("scaling of empty array", () => {
    expect(scaleToMedian([])).toStrictEqual([])
  })

  it("scaling of to median", () => {
    expect(scaleToMedian([2, 1, 1, 1, 0, 0, 1])).toStrictEqual([100, 50, 50, 50, 0, 0, 50])
  })
})
