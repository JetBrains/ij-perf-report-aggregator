import { describe, expect, it } from "vitest"
import { compareCell, pairRuns, pickRunValues } from "../../../../src/components/goland/engineComparison/engineRunSelection"

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

  it("single mode does not flag a change inside the LEGACY build-to-build spread", () => {
    const pairs = [100, 120, 90, 110, 80].map((legacy, i) => ({ t: DAY0 + i, legacy, new: legacy }))
    pairs.push({ t: DAY0 + 10, legacy: 100, new: 106 })
    const cell = compareCell(true, DAY0 + 10, pairs)
    expect(cell?.diffPercent).toBeCloseTo(6, 10)
    expect(cell).toMatchObject({ significant: false })
  })

  it("single mode flags a change far outside the LEGACY spread", () => {
    const pairs = [100, 101, 99, 100, 101].map((legacy, i) => ({ t: DAY0 + i, legacy, new: legacy }))
    pairs.push({ t: DAY0 + 10, legacy: 100, new: 150 })
    expect(compareCell(true, DAY0 + 10, pairs)).toMatchObject({ significant: true })
  })

  it("compareCell returns null without a paired build of the run", () => {
    expect(compareCell(true, DAY1, [{ t: DAY0, legacy: 100, new: 50 }])).toBeNull()
  })
})
