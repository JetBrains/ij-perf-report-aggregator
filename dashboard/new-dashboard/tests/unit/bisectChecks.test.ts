import { describe, it, expect } from "vitest"
import { checkGraphStability, checkTargetValue } from "../../src/components/common/sideBar/BisectChecks"
import type { InfoData } from "../../src/components/common/sideBar/InfoSidebar"

// Builds a minimal InfoData with the surrounding series values. `before` are the
// values preceding the selected dot, `after` the ones following it; the selected
// dot's own value sits between them and is excluded from both windows.
function makeData(before: number[], after: number[], metricType = "ms"): InfoData {
  const seriesValues = [...before, 0, ...after]
  return { seriesValues, pointIndex: before.length, metricType } as unknown as InfoData
}

describe("graph stability check", () => {
  it("degradation: no warning when before is entirely below after", () => {
    expect(checkGraphStability(makeData([10, 11, 12], [20, 21, 22]), "DEGRADATION")).toBeNull()
  })

  it("degradation: warns when the levels overlap (max before >= min after)", () => {
    expect(checkGraphStability(makeData([10, 25], [20, 30]), "DEGRADATION")).not.toBeNull()
  })

  it("optimization: no warning when before is entirely above after", () => {
    expect(checkGraphStability(makeData([30, 31], [10, 11]), "OPTIMIZATION")).toBeNull()
  })

  it("optimization: warns when the levels overlap (min before <= max after)", () => {
    expect(checkGraphStability(makeData([30, 15], [20, 10]), "OPTIMIZATION")).not.toBeNull()
  })

  it("returns null when there are no points after the dot", () => {
    const data = { seriesValues: [10, 11, 12], pointIndex: 2, metricType: "ms" } as unknown as InfoData
    expect(checkGraphStability(data, "DEGRADATION")).toBeNull()
  })

  it("returns null when series context is missing", () => {
    const data = { seriesValues: undefined, pointIndex: undefined } as unknown as InfoData
    expect(checkGraphStability(data, "DEGRADATION")).toBeNull()
  })
})

describe("target value check", () => {
  it("degradation: no warning when the target sits between before and after", () => {
    expect(checkTargetValue(makeData([10, 12], [20, 22]), "DEGRADATION", 15)).toBeNull()
  })

  it("degradation: warns when the highest before value exceeds the target", () => {
    expect(checkTargetValue(makeData([10, 18], [20]), "DEGRADATION", 15)).not.toBeNull()
  })

  it("degradation: warns when the lowest after value is below the target", () => {
    expect(checkTargetValue(makeData([10], [14, 22]), "DEGRADATION", 15)).not.toBeNull()
  })

  it("optimization: no warning when the target sits between before and after", () => {
    expect(checkTargetValue(makeData([30, 32], [10, 12]), "OPTIMIZATION", 20)).toBeNull()
  })

  it("optimization: warns when the lowest before value is below the target", () => {
    expect(checkTargetValue(makeData([18, 30], [10]), "OPTIMIZATION", 20)).not.toBeNull()
  })

  it("optimization: warns when the highest after value exceeds the target", () => {
    expect(checkTargetValue(makeData([30], [10, 25]), "OPTIMIZATION", 20)).not.toBeNull()
  })

  it("converts ns values to milliseconds before comparing with the target", () => {
    // 10ms and 12ms before, 20ms after, target 15ms -> no warning despite raw ns magnitudes.
    expect(checkTargetValue(makeData([10e6, 12e6], [20e6], "ns"), "DEGRADATION", 15)).toBeNull()
    // 14ms after dips below the 15ms target -> warning.
    expect(checkTargetValue(makeData([10e6], [14e6, 22e6], "ns"), "DEGRADATION", 15)).not.toBeNull()
  })

  it("returns null for a non-finite target value", () => {
    expect(checkTargetValue(makeData([10], [20]), "DEGRADATION", Number.NaN)).toBeNull()
  })
})
