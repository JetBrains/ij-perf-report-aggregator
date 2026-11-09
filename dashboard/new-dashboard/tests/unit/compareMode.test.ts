import { describe, expect, it } from "vitest"
import { pickBaseAndCompared, pickBaseAndComparedModes, pickCompareColumns } from "../../src/components/charts/compareMode"

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

describe("compareMode base/compared mode selection", () => {
  it("returns null when fewer than two modes are selected", () => {
    expect(pickBaseAndComparedModes([])).toBeNull()
    expect(pickBaseAndComparedModes(["default"])).toBeNull()
  })

  it("takes the default mode as the base wherever it is in the selection", () => {
    expect(pickBaseAndComparedModes(["types2", "default"])).toStrictEqual({ base: "default", compared: ["types2"] })
  })

  it("takes the first selected mode as the base without the default mode", () => {
    expect(pickBaseAndComparedModes(["split", "wsl"])).toStrictEqual({ base: "split", compared: ["wsl"] })
  })
})

describe("compareMode columns", () => {
  it("compares the branches when a single mode is selected", () => {
    expect(pickCompareColumns(["master", "261"], ["default"])).toStrictEqual({ base: "261", compared: ["master"], branches: ["261", "master"], modes: [] })
  })

  it("compares the modes when a single branch is selected", () => {
    expect(pickCompareColumns(["master"], ["types2", "default"])).toStrictEqual({
      base: "default",
      compared: ["types2"],
      branches: ["default", "types2"],
      modes: [],
    })
  })

  it("takes every branch in every mode as a column, with the base branch in the base mode as the base", () => {
    expect(pickCompareColumns(["master", "261"], ["types2", "default"])).toStrictEqual({
      base: "261 – default",
      compared: ["261 – types2", "master – default", "master – types2"],
      branches: ["261", "master"],
      modes: ["default", "types2"],
    })
  })

  it("returns null when neither two branches nor two modes are selected", () => {
    expect(pickCompareColumns(["master"], ["default"])).toBeNull()
  })
})
