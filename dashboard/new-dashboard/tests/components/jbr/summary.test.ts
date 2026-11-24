import { describe, expect, it } from "vitest"
import { compareConfigurationNames, RunValue, summarizeRow, SummaryRow } from "../../../src/components/jbr/summary"

const render: SummaryRow = { label: "Render onscreen", benchmark: "Render", route: "", unit: "milliseconds", betterDirection: "lower", measureSuffix: "_onscreen" }
const mapBench: SummaryRow = { label: "MapBench", benchmark: "Mapbench", route: "", unit: "framesPerSecond", betterDirection: "higher" }

// Twelve runs at each level, spread by a fixed ±1 % pattern so that the levels overlap where unchanged.
const NOISE = [1, 1.01, 0.99, 1.005, 0.995, 1.01, 0.99, 1, 1.004, 0.996, 1.008, 0.992]
const STEP = NOISE.length

function runs(project: string, measure: string, before: number, after: number): RunValue[] {
  return NOISE.flatMap((factor, run) => [
    { project, measure, time: run, value: before * factor },
    { project, measure, time: STEP + run, value: after * factor },
  ])
}

describe("JBR summary", () => {
  it("finds the step of shorter frame times as a significant improvement at its time", () => {
    const values = [...runs("Render_Windows11x86_64", "TextLCD_onscreen", 10, 8), ...runs("Render_Windows11x86_64", "Image_onscreen", 4, 3.2)]
    const cell = summarizeRow(render, values).get("Windows11x86_64")
    expect(cell?.latest?.time).toBe(STEP)
    expect(cell?.latest?.change).toBeCloseTo(0.25, 2)
    expect(cell?.latest).toMatchObject({ significant: true })
  })

  it("bounds the step of shorter frame times by its confidence interval", () => {
    const values = [...runs("Render_Windows11x86_64", "TextLCD_onscreen", 10, 8), ...runs("Render_Windows11x86_64", "Image_onscreen", 4, 3.2)]
    const cell = summarizeRow(render, values).get("Windows11x86_64")
    expect(cell?.latest?.low).toBeLessThan(0.25)
    expect(cell?.latest?.high).toBeGreaterThan(0.25)
    expect(cell?.worst).toBeNull()
  })

  it("leaves a change with too few runs after it pending", () => {
    // The step of the first test, but only five runs after it.
    const values = [...runs("Render_Windows11x86_64", "TextLCD_onscreen", 10, 8), ...runs("Render_Windows11x86_64", "Image_onscreen", 4, 3.2)].filter((it) => it.time < STEP + 5)
    const cell = summarizeRow(render, values).get("Windows11x86_64")
    expect(cell?.pending).toMatchObject({ time: STEP, runsAfter: 5 })
    expect(cell?.latest).toBeNull()
    expect(cell?.scenarios).toStrictEqual([])
  })

  it("finds no significant change in noise", () => {
    const cell = summarizeRow(mapBench, runs("Mapbench_Windows11x86_64", "CircleTests.ser", 100, 100)).get("Windows11x86_64")
    expect(cell?.latest).toBeNull()
    expect(cell?.runs).toBe(2 * STEP)
  })

  it("names the one regressed scenario that the geometric mean dilutes", () => {
    const values = [
      ...Array.from({ length: 7 }, (_, i) => runs("Mapbench_macOS15aarch64Metal", `scenario${i}.ser`, 100, 100)).flat(),
      ...runs("Mapbench_macOS15aarch64Metal", "spiralTest-fill.ser", 100, 75),
    ]
    const cell = summarizeRow(mapBench, values).get("macOS15aarch64Metal")
    expect(cell?.latest?.time).toBe(STEP)
    expect(cell?.latest?.change).toBeGreaterThan(-0.06)
    expect(cell?.worst?.scenario).toBe("spiralTest-fill.ser")
    expect(cell?.worst?.change).toBeCloseTo(-0.25, 2)
  })

  it("keeps only the row's mode, its benchmark and positive values", () => {
    const values = [
      ...runs("Render_Windows11x86_64", "TextLCD_volatile", 10, 5),
      ...runs("Render_Windows11x86_64", "TextLCD_onscreen#p90", 10, 5),
      ...runs("Dacapo_Windows11x86_64", "TextLCD_onscreen", 10, 5),
      ...runs("Render_Windows11x86_64", "Image_onscreen", 0, 0),
    ]
    expect(summarizeRow(render, values).size).toBe(0)
  })

  it("orders configurations by platform", () => {
    expect(["Windows11x86_64", "Ubuntu2404x86_64", "macOS15aarch64"].toSorted(compareConfigurationNames)).toStrictEqual(["macOS15aarch64", "Ubuntu2404x86_64", "Windows11x86_64"])
  })
})
