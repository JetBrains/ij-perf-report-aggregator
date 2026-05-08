import { deepEqual } from "fast-equals"
import { ref, Ref } from "vue"
import { ValueUnit } from "../common/chart"

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
