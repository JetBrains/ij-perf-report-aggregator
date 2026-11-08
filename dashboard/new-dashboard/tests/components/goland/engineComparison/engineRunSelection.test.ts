import { describe, expect, it } from "vitest"
import { compareCell, countMissing, pairRuns, pickRunValues } from "../../../../src/components/goland/engineComparison/engineRunSelection"

const DAY0 = 0
const DAY1 = 86_400_000

describe("engine run selection", () => {
  it("pairRuns keeps only the builds measured with both engines, sorted by time", () => {
    const legacy = [
      { t: DAY1, v: 110 },
      { t: DAY0, v: 100 },
      { t: DAY0 + 100, v: 105 },
    ]
    const brandNew = [
      { t: DAY0, v: 50 },
      { t: DAY1, v: 40 },
    ]
    expect(pairRuns(legacy, brandNew)).toStrictEqual([
      { t: DAY0, legacy: 100, new: 50 },
      { t: DAY1, legacy: 110, new: 40 },
    ])
  })

  it("pairRuns pairs a build that straddles UTC midnight only by its exact timestamp", () => {
    expect(pairRuns([{ t: DAY1 - 1, v: 1 }], [{ t: DAY1, v: 2 }])).toStrictEqual([])
  })

  it("aggregated mode returns every paired build", () => {
    const pairs = [
      { t: DAY0, legacy: 100, new: 50 },
      { t: DAY1, legacy: 110, new: 40 },
    ]
    expect(pickRunValues(false, null, pairs)).toStrictEqual({ legacy: [100, 110], new: [50, 40] })
  })

  it("aggregated mode returns null when no build has both engines", () => {
    expect(pickRunValues(false, null, [])).toBeNull()
  })

  it("single mode picks the selected build, also among several builds of one day", () => {
    const pairs = [
      { t: DAY0, legacy: 100, new: 50 },
      { t: DAY0 + 100, legacy: 101, new: 51 },
      { t: DAY1, legacy: 110, new: 40 },
    ]
    expect(pickRunValues(true, DAY0, pairs)).toStrictEqual({ legacy: [100], new: [50] })
    expect(pickRunValues(true, DAY0 + 100, pairs)).toStrictEqual({ legacy: [101], new: [51] })
  })

  it("single mode returns null for a run without a paired build", () => {
    expect(pickRunValues(true, DAY0 + 1, [{ t: DAY0, legacy: 100, new: 50 }])).toBeNull()
    expect(pickRunValues(true, null, [{ t: DAY0, legacy: 100, new: 50 }])).toBeNull()
  })

  it("single mode does not flag a change inside the build-to-build spread of the ratio", () => {
    const pairs = [1, 1.1, 0.9, 1.05, 0.95].map((ratio, i) => ({ t: DAY0 + i, legacy: 100, new: 100 * ratio }))
    pairs.push({ t: DAY0 + 10, legacy: 100, new: 106 })
    const cell = compareCell(true, DAY0 + 10, pairs)
    expect(cell?.diffPercent).toBeCloseTo(6, 10)
    expect(cell).toMatchObject({ significant: false })
  })

  it("single mode flags a change far outside the spread of the ratio", () => {
    const pairs = [1, 1.01, 0.99, 1, 1.01].map((ratio, i) => ({ t: DAY0 + i, legacy: 100, new: 100 * ratio }))
    pairs.push({ t: DAY0 + 10, legacy: 100, new: 150 })
    expect(compareCell(true, DAY0 + 10, pairs)).toMatchObject({ significant: true })
  })

  it("single mode does not count the drift both engines share as noise", () => {
    // Both engines slow down build by build, NEW stays at 0.6 of LEGACY, and the selected build is 10% off.
    const pairs = [100, 200, 300, 400, 500].map((legacy, i) => ({ t: DAY0 + i, legacy, new: legacy * 0.6 }))
    pairs.push({ t: DAY0 + 10, legacy: 100, new: 66 })
    expect(compareCell(true, DAY0 + 10, pairs)).toMatchObject({ significant: true })
  })

  it("aggregated mode takes the ratio of each build, not the ratio of the two centers", () => {
    // NEW is 10% faster in two builds and 10% slower in one slow build. The ratio of the centers
    // (342.5 / 325) would let that build decide: NEW 5% slower.
    const pairs = [
      { t: DAY0, legacy: 100, new: 90 },
      { t: DAY0 + 1, legacy: 100, new: 90 },
      { t: DAY0 + 2, legacy: 1000, new: 1100 },
    ]
    // The Hodges-Lehmann center of the three log ratios: the median of their six pairwise averages.
    const expected = Math.exp((Math.log(0.9) + (Math.log(0.9) + Math.log(1.1)) / 2) / 2)
    expect(compareCell(false, null, pairs)?.ratio).toBeCloseTo(expected, 10)
    expect(compareCell(false, null, pairs)?.ratio).toBeLessThan(1)
  })

  it("countMissing counts the builds of the scope that published one engine only", () => {
    const legacy = [
      { t: DAY0, v: 100 },
      { t: DAY0 + 1, v: 100 },
      { t: DAY1, v: 100 },
    ]
    const brandNew = [
      { t: DAY0, v: 50 },
      { t: DAY0 + 2, v: 50 },
    ]
    // DAY0: paired; DAY0 + 1: NEW missing; DAY0 + 2: LEGACY missing; DAY1: outside the scope; DAY1 + 1: neither engine.
    expect(countMissing([DAY0, DAY0 + 1, DAY0 + 2, DAY1 + 1], legacy, brandNew)).toStrictEqual({ legacy: 1, new: 1 })
  })

  it("compareCell returns null without a paired build of the run", () => {
    expect(compareCell(true, DAY1, [{ t: DAY0, legacy: 100, new: 50 }])).toBeNull()
  })
})
