import { describe, it, expect } from "vitest"
import { AccidentKind } from "../../src/configurators/accidents/AccidentsConfigurator"
import type { InfoData } from "../../src/components/common/sideBar/InfoSidebar"
import { detectPossibleMisclick } from "../../src/components/common/sideBar/MisclickHeuristic"

// Builds a minimal InfoData around the selected dot: `before` are the values preceding it
// (the last one is the previous point), `after` the ones following it.
function makeData(before: number[], value: number, after: number[] = []): InfoData {
  return {
    series: [{ rawValue: value }],
    previousValue: before.at(-1),
    nextValue: after[0],
    seriesValues: [...before, value, ...after],
    pointIndex: before.length,
  } as unknown as InfoData
}

describe("misclick heuristic: noise", () => {
  it("warns when the point lies within the range of the noisy points before it", () => {
    // diaspora firstCodeAnalysis on Linux: a 4% drop in a series swinging between 342 and 507.
    const warning = detectPossibleMisclick(makeData([383, 342, 355, 473, 424, 507, 474, 457, 455, 465], 447), AccidentKind.Improvement)
    expect(warning?.title).toBe("The selected point doesn't stand out from the noise")
    expect(warning?.detail).toContain("changed by 3.9%")
    expect(warning?.detail).toContain("the 10 points before it already vary by 32.5%")
  })

  it("does not warn when the point leaves the range of the points before it", () => {
    expect(detectPossibleMisclick(makeData([100, 102, 99, 101, 100, 98], 120), AccidentKind.Regression)).toBeNull()
  })

  it("does not judge noise from too few points", () => {
    expect(detectPossibleMisclick(makeData([100, 130, 110], 120), AccidentKind.Regression)).toBeNull()
  })

  it("points to a later, larger change when the selected point is within the noise", () => {
    const warning = detectPossibleMisclick(makeData([100, 120, 105, 115, 110], 112, [150]), AccidentKind.Regression)
    expect(warning?.title).toBe("The selected point doesn't look like the actual event")
    expect(warning?.detail).toContain("1 point later")
  })

  it("keeps the wrong-direction warning ahead of the noise warning", () => {
    const warning = detectPossibleMisclick(makeData([100, 120, 105, 115, 110], 104), AccidentKind.Regression)
    expect(warning?.title).toBe("This change looks like an improvement, not a regression")
  })
})
