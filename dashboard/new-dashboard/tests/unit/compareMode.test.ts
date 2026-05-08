import { describe, expect, it } from "vitest"
import { pickBaseAndCompared } from "../../src/components/charts/compareMode"

describe("compareMode base/compared selection", () => {
  it("returns null when nothing is selected", () => {
    expect(pickBaseAndCompared([])).toBeNull()
  })

  it("returns null when only one branch is selected", () => {
    expect(pickBaseAndCompared(["253"])).toBeNull()
    expect(pickBaseAndCompared(["master"])).toBeNull()
  })

  it("prefers the lowest-numbered release branch as the base, putting master in compared", () => {
    expect(pickBaseAndCompared(["253", "master", "261"])).toStrictEqual({ base: "253", compared: ["master", "261"] })
  })

  it("picks the lowest-numbered release branch when several releases are selected", () => {
    expect(pickBaseAndCompared(["261", "253", "262"])).toStrictEqual({ base: "253", compared: ["261", "262"] })
  })

  it("uses the first selected branch when neither master nor a release branch is present", () => {
    expect(pickBaseAndCompared(["my-feature", "their-feature"])).toStrictEqual({ base: "my-feature", compared: ["their-feature"] })
  })

  it("prefers a release branch over master", () => {
    expect(pickBaseAndCompared(["253", "master"])).toStrictEqual({ base: "253", compared: ["master"] })
  })

  it("falls back to master when no release branch is selected", () => {
    expect(pickBaseAndCompared(["master", "feature-b"])).toStrictEqual({ base: "master", compared: ["feature-b"] })
  })

  it("treats release branches as numeric, not lexicographic", () => {
    // "9" sorts after "253" lexicographically but should be picked as the lowest numeric base.
    expect(pickBaseAndCompared(["253", "9"])).toStrictEqual({ base: "9", compared: ["253"] })
  })

  it("preserves the user's selection order in compared[]", () => {
    expect(pickBaseAndCompared(["master", "feature-b", "feature-a"])).toStrictEqual({ base: "master", compared: ["feature-b", "feature-a"] })
  })
})
