import { describe, expect, it } from "vitest"
import { computeEngineAggregates, EngineComparePoint } from "../../../../src/components/goland/engineComparison/engineCompareAggregates"

function point(overrides: Partial<EngineComparePoint>): EngineComparePoint {
  return {
    base: "caddy",
    title: "caddy",
    metricType: "coldStartHighlighting_fast",
    phase: "coldStartHighlighting",
    bucket: "fast",
    ratio: 1,
    diffPercent: 0,
    significant: false,
    ...overrides,
  }
}

describe("engine comparison aggregates", () => {
  it("returns the geometric mean of positive ratios", () => {
    const aggregates = computeEngineAggregates([point({ ratio: 0.5 }), point({ ratio: 2 })])
    // sqrt(0.5 * 2) === 1
    expect(aggregates.grandGeomean).toBeCloseTo(1, 10)
    expect(aggregates.count).toBe(2)
  })

  it("skips zero, negative, and NaN ratios", () => {
    const aggregates = computeEngineAggregates([point({ ratio: 2 }), point({ ratio: 0 }), point({ ratio: -1 }), point({ ratio: Number.NaN })])
    expect(aggregates.count).toBe(1)
    expect(aggregates.grandGeomean).toBeCloseTo(2, 10)
  })

  it("groups per bucket in first-appearance order", () => {
    const rows = [point({ bucket: "fast", ratio: 0.5 }), point({ bucket: "medium", ratio: 2 }), point({ bucket: "slow", ratio: 4 })]
    const aggregates = computeEngineAggregates(rows)
    expect(aggregates.perBucket.map((group) => group.key)).toStrictEqual(["fast", "medium", "slow"])
    expect(aggregates.perBucket[0].geomean).toBeCloseTo(0.5, 10)
  })

  it("groups per phase and averages within a phase", () => {
    const rows = [point({ phase: "coldStartHighlighting", ratio: 0.5 }), point({ phase: "coldStartHighlighting", ratio: 2 }), point({ phase: "warmStartHighlighting", ratio: 4 })]
    const aggregates = computeEngineAggregates(rows)
    expect(aggregates.perPhase.map((group) => group.key)).toStrictEqual(["coldStartHighlighting", "warmStartHighlighting"])
    // cold has ratios 0.5 and 2 -> geomean 1
    expect(aggregates.perPhase[0].geomean).toBeCloseTo(1, 10)
    expect(aggregates.perPhase[0].count).toBe(2)
  })

  it("counts improved / regressed / neutral by significance and sign", () => {
    const rows = [
      point({ significant: true, diffPercent: 30 }), // regressed
      point({ significant: true, diffPercent: -20 }), // improved
      point({ significant: false, diffPercent: 40 }), // neutral (not significant)
      point({ significant: true, diffPercent: 0 }), // neutral (no move)
      point({ significant: true, diffPercent: Number.NaN }), // neutral (unusable)
    ]
    const aggregates = computeEngineAggregates(rows)
    expect(aggregates.regressed).toBe(1)
    expect(aggregates.improved).toBe(1)
    expect(aggregates.neutral).toBe(3)
  })

  it("selects the worst regression among significant rows", () => {
    const rows = [
      point({ base: "a", significant: true, diffPercent: 10 }),
      point({ base: "b", significant: true, diffPercent: 45 }),
      point({ base: "e", significant: false, diffPercent: 99 }), // ignored: not significant
    ]
    const aggregates = computeEngineAggregates(rows)
    expect(aggregates.worstRegression?.base).toBe("b")
    expect(aggregates.worstRegression?.diffPercent).toBe(45)
  })

  it("selects the best improvement among significant rows", () => {
    const rows = [point({ base: "c", significant: true, diffPercent: -5 }), point({ base: "d", significant: true, diffPercent: -60 })]
    const aggregates = computeEngineAggregates(rows)
    expect(aggregates.bestImprovement?.base).toBe("d")
    expect(aggregates.bestImprovement?.diffPercent).toBe(-60)
  })

  it("handles empty input", () => {
    const aggregates = computeEngineAggregates([])
    expect(aggregates).toStrictEqual({
      count: 0,
      grandGeomean: Number.NaN,
      perBucket: [],
      perPhase: [],
      improved: 0,
      regressed: 0,
      neutral: 0,
      worstRegression: null,
      bestImprovement: null,
    })
  })
})
