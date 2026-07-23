import { describe, expect, it } from "vitest"
import { LlmAnalysisListItem, LlmAnalysisState } from "../../src/components/common/llmAnalysis/LlmAnalysisClient"
import { AnalysesFilterState, currentUserLabel, distinctUsers, emptyAnalysesFilterState, filterAnalyses, userLabel } from "../../src/components/analyses/analysesFilter"

function item(overrides: Partial<LlmAnalysisListItem>): LlmAnalysisListItem {
  return {
    id: 1,
    createdAt: "2026-06-01T12:00:00Z",
    runBuildId: "100",
    state: LlmAnalysisState.Success,
    project: "intellij/indexing",
    metric: "indexing",
    currentBuildId: "200",
    prevBuildId: "199",
    feedbackCount: 0,
    ...overrides,
  }
}

function filter(overrides: Partial<AnalysesFilterState>): AnalysesFilterState {
  return { ...emptyAnalysesFilterState(), ...overrides }
}

// Helpers: 1 = item matches filter, 0 = item does not match
function matches(it: LlmAnalysisListItem, f: AnalysesFilterState): number {
  return filterAnalyses([it], f).length
}

describe("userLabel helper", () => {
  it("prefers userName", () => {
    expect(userLabel(item({ userName: "Jane Doe", userEmail: "jane@x.com" }))).toBe("Jane Doe")
  })

  it("falls back to email local part", () => {
    expect(userLabel(item({ userName: undefined, userEmail: "jane@jetbrains.com" }))).toBe("jane")
  })

  it("empty when nothing", () => {
    expect(userLabel(item({ userName: undefined, userEmail: undefined }))).toBe("")
  })
})

describe("currentUserLabel helper", () => {
  it("prefers name", () => {
    expect(currentUserLabel({ name: "Jane Doe", email: "jane@x.com" })).toBe("Jane Doe")
  })

  it("falls back to email local part", () => {
    expect(currentUserLabel({ name: "", email: "jane@jetbrains.com" })).toBe("jane")
  })

  it("empty when no user", () => {
    expect(currentUserLabel(null)).toBe("")
    expect(currentUserLabel({})).toBe("")
  })

  it("matches the label used by the user filter for the same identity", () => {
    const label = currentUserLabel({ name: "Jane Doe", email: "jane@x.com" })
    expect(matches(item({ userName: "Jane Doe" }), filter({ users: [label] }))).toBe(1)
  })
})

describe("filterAnalyses: search filter", () => {
  it("empty filter matches everything", () => {
    expect(matches(item({}), emptyAnalysesFilterState())).toBe(1)
  })

  it("matches project, metric, or user case-insensitively", () => {
    const it1 = item({ project: "goland/dataflow", metric: "completion", userName: "Bob" })
    expect(matches(it1, filter({ search: "DATAFLOW" }))).toBe(1)
    expect(matches(it1, filter({ search: "completion" }))).toBe(1)
    expect(matches(it1, filter({ search: "bob" }))).toBe(1)
    expect(matches(it1, filter({ search: "kotlin" }))).toBe(0)
  })
})

describe("filterAnalyses: state filter", () => {
  it("state filter is an OR set; empty means all", () => {
    expect(matches(item({ state: LlmAnalysisState.Failed }), filter({ states: [LlmAnalysisState.Success, LlmAnalysisState.InProgress] }))).toBe(0)
    expect(matches(item({ state: LlmAnalysisState.Failed }), filter({ states: [LlmAnalysisState.Failed] }))).toBe(1)
    expect(matches(item({ state: LlmAnalysisState.Failed }), filter({ states: [] }))).toBe(1)
  })
})

describe("filterAnalyses: date range filter", () => {
  it("date range is inclusive on whole days", () => {
    const it1 = item({ createdAt: "2026-06-10T08:30:00Z" })
    expect(matches(it1, filter({ dateFrom: new Date(2026, 5, 10), dateTo: new Date(2026, 5, 10) }))).toBe(1)
    expect(matches(it1, filter({ dateFrom: new Date(2026, 5, 11) }))).toBe(0)
    expect(matches(it1, filter({ dateTo: new Date(2026, 5, 9) }))).toBe(0)
  })
})

describe("filterAnalyses: boolean toggles", () => {
  it("hasTicket toggle", () => {
    expect(matches(item({ ytIssueId: undefined }), filter({ hasTicket: true }))).toBe(0)
    expect(matches(item({ ytIssueId: "IJPL-1" }), filter({ hasTicket: true }))).toBe(1)
  })

  it("hasFeedback toggle", () => {
    expect(matches(item({ feedbackCount: 0 }), filter({ hasFeedback: true }))).toBe(0)
    expect(matches(item({ feedbackCount: 3 }), filter({ hasFeedback: true }))).toBe(1)
  })

  it("hasCommits toggle", () => {
    expect(matches(item({ llmGuiltyCommits: [] }), filter({ hasCommits: true }))).toBe(0)
    expect(matches(item({ llmGuiltyCommits: ["abc"] }), filter({ hasCommits: true }))).toBe(1)
  })
})

describe("filterAnalyses: commit filter", () => {
  it("empty hash matches all items regardless of commits", () => {
    expect(matches(item({}), filter({ commit: "" }))).toBe(1)
    expect(matches(item({ llmGuiltyCommits: ["abc123"] }), filter({ commit: "" }))).toBe(1)
  })

  it("matches by substring (contains), case-insensitively, across found commits", () => {
    const it1 = item({ llmGuiltyCommits: ["4bd9c5d518ecb58238b44c0645488cc424f5264d", "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"] })
    expect(matches(it1, filter({ commit: "4bd9c5d5" }))).toBe(1)
    expect(matches(it1, filter({ commit: "4BD9C5D5" }))).toBe(1)
    expect(matches(it1, filter({ commit: "88cc424f" }))).toBe(1)
    expect(matches(it1, filter({ commit: "deadbeef" }))).toBe(1)
    expect(matches(it1, filter({ commit: "abc123" }))).toBe(0)
  })

  it("non-empty hash does not match items with no commits", () => {
    expect(matches(item({ llmGuiltyCommits: undefined }), filter({ commit: "4bd9" }))).toBe(0)
  })
})

describe("filterAnalyses: user filter", () => {
  it("user filter is an OR set over user labels; empty means all", () => {
    const jane = item({ userName: "Jane Doe" })
    const bob = item({ userName: undefined, userEmail: "bob@jetbrains.com" })
    expect(matches(jane, filter({ users: ["Jane Doe"] }))).toBe(1)
    expect(matches(jane, filter({ users: ["bob"] }))).toBe(0)
    expect(matches(bob, filter({ users: ["bob", "Jane Doe"] }))).toBe(1)
    expect(matches(jane, filter({ users: [] }))).toBe(1)
  })
})

describe("filterAnalyses: combined filters", () => {
  it("combines all predicates with AND", () => {
    const it1 = item({ state: LlmAnalysisState.Success, ytIssueId: "IJPL-1", project: "kotlin/highlighting" })
    expect(matches(it1, filter({ states: [LlmAnalysisState.Success], hasTicket: true, search: "kotlin" }))).toBe(1)
    expect(matches(it1, filter({ states: [LlmAnalysisState.Failed], hasTicket: true, search: "kotlin" }))).toBe(0)
  })

  it("returns only matching items from a list", () => {
    const items = [item({ id: 1, state: LlmAnalysisState.Success }), item({ id: 2, state: LlmAnalysisState.Failed }), item({ id: 3, state: LlmAnalysisState.Success })]
    expect(filterAnalyses(items, filter({ states: [LlmAnalysisState.Success] })).map((i) => i.id)).toStrictEqual([1, 3])
  })
})

describe("distinctUsers helper", () => {
  it("returns sorted unique non-empty user labels", () => {
    const items = [
      item({ userName: "Bob" }),
      item({ userName: "Anna" }),
      item({ userName: "Bob" }),
      item({ userName: undefined, userEmail: undefined }),
      item({ userName: undefined, userEmail: "carol@x.com" }),
    ]
    expect(distinctUsers(items)).toStrictEqual(["Anna", "Bob", "carol"])
  })
})
