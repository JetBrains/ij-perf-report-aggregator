import { describe, expect, it } from "vitest"
import { computeBaseStats, computeBranchStats, percentile } from "../../src/components/charts/compareStats"

describe("compareStats", () => {
  it("percentile interpolates between values", () => {
    expect(percentile([10, 20, 30, 40, 50], 50)).toBe(30)
    expect(percentile([10, 20, 30, 40, 50], 10)).toBe(14)
    expect(percentile([10, 20, 30, 40, 50], 90)).toBe(46)
  })

  it("flat baseline + small branch diff yields a large disparity", () => {
    const base = computeBaseStats([100, 100.1, 99.9, 100, 100, 99.95, 100.05])
    const branch = computeBranchStats(base, [105])
    expect(base.count).toBe(7)
    expect(branch.count).toBe(1)
    expect(base.center).toBeCloseTo(100, 5)
    expect(branch.diffPercent).toBeCloseTo(5, 5)
    expect(Math.abs(branch.disparity)).toBeGreaterThan(10)
  })

  it("noisy baseline + same branch diff yields a small disparity", () => {
    const base = computeBaseStats([80, 90, 100, 110, 120, 95, 105, 115, 85])
    const branch = computeBranchStats(base, [105])
    expect(Math.abs(branch.disparity)).toBeLessThan(2)
  })

  it("zero variance baseline with matching branch returns disparity=0", () => {
    const base = computeBaseStats([100, 100, 100])
    const branch = computeBranchStats(base, [100])
    expect(base.spread).toBe(0)
    expect(branch.disparity).toBe(0)
  })

  it("zero variance baseline with branch above returns +Infinity", () => {
    const base = computeBaseStats([100, 100, 100])
    const branch = computeBranchStats(base, [101])
    expect(base.spread).toBe(0)
    expect(branch.disparity).toBe(Number.POSITIVE_INFINITY)
  })

  it("Shamos spread survives quantized baselines that collapse MAD", () => {
    // 60 % of values equal the median; MAD = 0 here. Shamos spread uses pairwise differences,
    // so any two distinct values contribute — yielding a finite, sensible disparity instead of −∞.
    const masterValues: number[] = [...Array.from({ length: 60 }, () => 13), ...Array.from({ length: 20 }, () => 12), ...Array.from({ length: 20 }, () => 14)]
    const base = computeBaseStats(masterValues)
    const branch = computeBranchStats(base, [10])
    expect(base.center).toBe(13)
    expect(base.p10).toBeCloseTo(12, 5)
    expect(base.p90).toBeCloseTo(14, 5)
    expect(base.spread).toBeGreaterThan(0)
    expect(branch.disparity).toBeLessThan(-2)
  })

  it("single-element baseline has no spread and falls back to ±Infinity", () => {
    const base = computeBaseStats([100])
    const branch = computeBranchStats(base, [101])
    expect(base.spread).toBe(0)
    expect(branch.disparity).toBe(Number.POSITIVE_INFINITY)
  })

  it("tie-dominant baseline (Shamos collapses) falls back to p10–p90-based dispersion", () => {
    // 70 % of values identical: Shamos's median pairwise difference is 0, so pragmastat throws.
    // p10–p90 still carries dispersion info, so the disparity must stay finite — without this
    // fallback, near-identical branches show up as ±∞D (the bug from the indexedFiles screenshot).
    const baseValues = [...Array.from({ length: 70 }, () => 1000), ...Array.from({ length: 30 }, () => 990)]
    const base = computeBaseStats(baseValues)
    const branch = computeBranchStats(base, [1010])
    expect(base.spread).toBeGreaterThan(0)
    expect(branch.disparity).toBeGreaterThan(0)
    expect(branch.disparity).toBeLessThan(Number.POSITIVE_INFINITY)
  })

  it("missing branch data leaves diff & disparity as NaN", () => {
    const base = computeBaseStats([100, 100, 100])
    const branch = computeBranchStats(base, [])
    expect(branch.center).toBeNaN()
    expect(branch.diffPercent).toBeNaN()
    expect(branch.disparity).toBeNaN()
  })
})
