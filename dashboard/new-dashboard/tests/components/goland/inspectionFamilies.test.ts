import { describe, expect, it } from "vitest"
import { GO_INSPECTION_FAMILIES, GOLANGCI_LINT_FAMILY } from "../../../src/components/goland/inspectionFamilies"

const ALL_FAMILIES = [GOLANGCI_LINT_FAMILY, ...GO_INSPECTION_FAMILIES]

describe("Go inspection families", () => {
  // Within the light palette, an inspection keeps its color across charts.
  it.each(ALL_FAMILIES)("keeps $key within the light palette", (family) => {
    expect(family.inspections.length).toBeLessThanOrEqual(9)
  })

  it("puts each inspection in one family", () => {
    const inspections = ALL_FAMILIES.flatMap((family) => family.inspections)
    expect(inspections.filter((it, index) => inspections.indexOf(it) !== index)).toStrictEqual([])
  })
})
