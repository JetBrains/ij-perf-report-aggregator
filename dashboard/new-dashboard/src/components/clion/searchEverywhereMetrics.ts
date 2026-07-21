import { computed, inject, Ref } from "vue"
import { selectedToArray } from "../../configurators/DimensionConfigurator"
import { branchConfiguratorKey } from "../../shared/keys"

export const SE_MEASURES = new Set(["searchEverywhere", "searchEverywhere_first_elements_added"])

export function toNewSeProjects(projects: string[]): string[] {
  return projects.map((project) => project.replace("/go-to-", "/new-se-go-to-"))
}

export function containsSeMeasure(measure: string | string[]): boolean {
  return (Array.isArray(measure) ? measure : [measure]).some((m) => SE_MEASURES.has(m))
}

function isLegacySeBranch(branch: string): boolean {
  return branch === "261" || branch.startsWith("261.")
}

export function useNewSearchEverywhere(): Ref<boolean> {
  const branchConfigurator = inject(branchConfiguratorKey, null)
  return computed(() => {
    const branches = selectedToArray(branchConfigurator?.selected.value)
    return branches.length === 0 || !branches.every((branch) => isLegacySeBranch(branch))
  })
}
