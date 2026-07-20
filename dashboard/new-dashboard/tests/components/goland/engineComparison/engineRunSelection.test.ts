import { describe, expect, it } from "vitest"
import { dayKey, latestRunByDay, MS_PER_DAY, pickRunValues } from "../../../../src/components/goland/engineComparison/engineRunSelection"

const DAY0 = 0
const DAY1 = MS_PER_DAY
const DAY2 = 2 * MS_PER_DAY

describe("engine run selection", () => {
  it("dayKey floors a timestamp to its UTC day", () => {
    expect(dayKey(DAY1 + 5_000)).toBe(1)
    expect(dayKey(DAY0)).toBe(0)
  })

  it("latestRunByDay keeps the last run of each day", () => {
    const byDay = latestRunByDay([
      { t: DAY0, v: 1 },
      { t: DAY0 + 100, v: 2 },
    ])
    expect(byDay.get(0)?.v).toBe(2)
  })

  it("aggregated mode returns every run's value", () => {
    const values = pickRunValues(
      false,
      null,
      [
        { t: DAY0, v: 100 },
        { t: DAY1, v: 110 },
      ],
      [{ t: DAY0, v: 50 }]
    )
    expect(values).toStrictEqual({ legacy: [100, 110], new: [50] })
  })

  it("single mode with no day picks the latest common day", () => {
    const legacy = [
      { t: DAY0, v: 100 },
      { t: DAY1, v: 110 },
    ]
    const brandNew = [
      { t: DAY0, v: 50 },
      { t: DAY1, v: 40 },
    ]
    expect(pickRunValues(true, null, legacy, brandNew)).toStrictEqual({ legacy: [110], new: [40] })
  })

  it("single mode honors an explicit day", () => {
    const legacy = [
      { t: DAY0, v: 100 },
      { t: DAY1, v: 110 },
    ]
    const brandNew = [
      { t: DAY0, v: 50 },
      { t: DAY1, v: 40 },
    ]
    expect(pickRunValues(true, 0, legacy, brandNew)).toStrictEqual({ legacy: [100], new: [50] })
  })

  it("single mode skips a day missing one engine", () => {
    expect(pickRunValues(true, DAY2 / MS_PER_DAY, [{ t: DAY2, v: 1 }], [{ t: DAY0, v: 2 }])).toBeNull()
  })

  it("single mode returns null when the engines share no day", () => {
    expect(pickRunValues(true, null, [{ t: DAY0, v: 1 }], [{ t: DAY1, v: 2 }])).toBeNull()
  })
})
