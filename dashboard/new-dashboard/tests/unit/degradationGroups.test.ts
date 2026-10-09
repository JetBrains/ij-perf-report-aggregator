import { describe, expect, it } from "vitest"
import { LlmAnalysisDegradationMatch } from "../../src/components/common/llmAnalysis/LlmAnalysisClient"
import { groupDegradations } from "../../src/components/common/llmAnalysis/degradationGroups"

function degradation(overrides: Partial<LlmAnalysisDegradationMatch>): LlmAnalysisDegradationMatch {
  return { project: "p", metric: "m", buildId: "1", date: "2026-09-28", matchedCommits: ["a"], rangeSize: 100, ...overrides }
}

describe("degradation matches grouping", () => {
  it("merges builds of the same chart and describes it by the strongest build only", () => {
    const result = groupDegradations([
      degradation({ buildId: "1", rangeSize: 200, matchedCommits: ["a", "b", "c"], date: "2026-09-30" }),
      degradation({ buildId: "2", rangeSize: 2, matchedCommits: ["a"], date: "2026-09-29" }),
    ])
    expect(result.chartCount).toBe(1)
    const chart = result.groups[0].charts[0]
    expect(chart.builds.map((b) => b.buildId)).toStrictEqual(["2", "1"])
    // not "3 of 2 commits" mixed from both builds
    expect([chart.matchedCommits.length, chart.rangeSize, chart.date]).toStrictEqual([1, 2, "2026-09-29"])
  })

  it("groups by metric, larger groups first, charts by range size", () => {
    const result = groupDegradations([
      degradation({ project: "a", metric: "related" }),
      degradation({ project: "b", metric: "same", rangeSize: 300 }),
      degradation({ project: "c", metric: "same", rangeSize: 20 }),
    ])
    expect(result.groups.map((g) => g.metric)).toStrictEqual(["same", "related"])
    expect(result.groups[0].charts.map((c) => c.project)).toStrictEqual(["c", "b"])
    expect(result.strongest?.project).toBe("c")
  })

  it("reports common commits only when every chart matched the same ones", () => {
    expect(groupDegradations([degradation({ project: "a", matchedCommits: ["x", "y"] }), degradation({ project: "b", matchedCommits: ["y", "x"] })]).commonCommits).toStrictEqual([
      "x",
      "y",
    ])
    expect(groupDegradations([degradation({ project: "a", matchedCommits: ["x"] }), degradation({ project: "b", matchedCommits: ["y"] })]).commonCommits).toBeNull()
    expect(groupDegradations([]).commonCommits).toBeNull()
  })
})
