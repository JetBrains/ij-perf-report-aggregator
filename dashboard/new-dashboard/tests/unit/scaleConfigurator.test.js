import { expect, test } from "vitest"
import { scaleToMedian } from "../../src/configurators/ScalingConfigurator"

test("scaling of empty array", () => {
  expect(scaleToMedian([])).toEqual([])
})

test("scaling of to median", () => {
  expect(scaleToMedian([2, 1, 1, 1, 0, 0, 1])).toEqual([100, 50, 50, 50, 0, 0, 50])
})
