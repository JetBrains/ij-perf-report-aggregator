import { deepEqual } from "fast-equals"
import { ref, Ref } from "vue"
import { ValueUnit } from "../common/chart"
import { defaultModeName } from "../../configurators/TestModeConfigurator"
import { columnName } from "./compareQuery"

export type RenderMode = "charts" | "compare"

export interface BaseAndCompared {
  base: string
  compared: string[]
}

// Selection priority for the compare-table base:
//   1. the lowest-numbered release-style branch (matches RELEASE_BRANCH_RE), e.g. "253" beats "261"
//   2. master, if no release branch is selected
//   3. the first selected branch otherwise
// Release wins over master because users typically ask "did master regress vs the stable
// release?" — putting the release on the base axis makes +Δ% on master read as a regression.
// All other selected branches become the compared columns, preserving the user's selection order.
// Returns null when there's nothing to compare (no selection, or only the base is selected).
const RELEASE_BRANCH_RE = /^\d+$/

export function pickBaseAndCompared(selected: readonly string[]): BaseAndCompared | null {
  if (selected.length === 0) return null
  const releases = selected.filter((b) => RELEASE_BRANCH_RE.test(b))
  let base: string
  if (releases.length > 0) {
    base = releases.reduce((lo, b) => (Number(b) < Number(lo) ? b : lo))
  } else if (selected.includes("master")) {
    base = "master"
  } else {
    base = selected[0]
  }
  const compared = selected.filter((b) => b !== base)
  if (compared.length === 0) return null
  return { base, compared }
}

// The base of a mode comparison is the default mode when it is selected, else the first selected mode.
export function pickBaseAndComparedModes(selected: readonly string[]): BaseAndCompared | null {
  if (selected.length < 2) return null
  const base = selected.includes(defaultModeName) ? defaultModeName : selected[0]
  return { base, compared: selected.filter((m) => m !== base) }
}

export interface CompareColumns extends BaseAndCompared {
  // The column dimension that indexSeries matches against series-name parts, base first: the branches, or the modes
  // when a single branch is selected.
  branches: string[]
  // The modes, base first, when every column is one branch in one mode; otherwise empty.
  modes: string[]
}

// Two or more branches are compared with each other; with two or more modes as well, each column is one branch in one
// mode, and the base is the base branch in the base mode. A single branch compares the selected modes.
export function pickCompareColumns(branches: readonly string[], modes: readonly string[]): CompareColumns | null {
  const byBranch = pickBaseAndCompared(branches)
  const byMode = pickBaseAndComparedModes(modes)
  if (byBranch == null) {
    return byMode == null ? null : { ...byMode, branches: [byMode.base, ...byMode.compared], modes: [] }
  }
  const branchColumns = [byBranch.base, ...byBranch.compared]
  if (byMode == null) {
    return { ...byBranch, branches: branchColumns, modes: [] }
  }
  const modeColumns = [byMode.base, ...byMode.compared]
  const base = columnName(byBranch.base, byMode.base)
  const compared = branchColumns.flatMap((branch) => modeColumns.map((mode) => columnName(branch, mode))).filter((column) => column !== base)
  return { base, compared, branches: branchColumns, modes: modeColumns }
}

export interface CompareSectionConfig {
  id: string
  label: string
  measure: string | string[]
  projects: string[]
  aliases: string[] | null
  machines: string[] | null
  valueUnit: ValueUnit
}

export class CompareSectionsRegistry {
  private nextId = 0
  readonly sections: Ref<CompareSectionConfig[]> = ref<CompareSectionConfig[]>([])

  nextSectionId(): string {
    this.nextId++
    return `section-${this.nextId}`
  }

  register(config: CompareSectionConfig): void {
    const existing = this.sections.value.findIndex((s) => s.id === config.id)
    if (existing === -1) {
      this.sections.value = [...this.sections.value, config]
      return
    }
    // GroupProjectsChart re-runs its deep watcher on every prop tick during dashboard mount;
    // skip the array reassignment (and the downstream watcher fan-out) when nothing actually changed.
    if (deepEqual(this.sections.value[existing], config)) return
    const next = [...this.sections.value]
    next[existing] = config
    this.sections.value = next
  }

  unregister(id: string): void {
    this.sections.value = this.sections.value.filter((s) => s.id !== id)
  }
}
