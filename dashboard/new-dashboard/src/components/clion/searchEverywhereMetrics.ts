import { computed, inject, Ref } from "vue"
import { selectedToArray } from "../../configurators/DimensionConfigurator"
import { TestModeConfigurator } from "../../configurators/TestModeConfigurator"
import { branchConfiguratorKey, dashboardConfiguratorsKey } from "../../shared/keys"

export const SE_MEASURES = new Set(["searchEverywhere", "searchEverywhere_first_elements_added"])

const SPLIT_MODE = "split"

export function toNewSeProjects(projects: string[]): string[] {
  return projects.map((project) => project.replace("/go-to-", "/new-se-go-to-"))
}

export function containsSeMeasure(measure: string | string[]): boolean {
  return (Array.isArray(measure) ? measure : [measure]).some((m) => SE_MEASURES.has(m))
}

function isLegacySeBranch(branch: string): boolean {
  return branch === "261" || branch.startsWith("261.")
}

function useTestModeConfigurator(): TestModeConfigurator | null {
  const configurators: unknown[] = inject(dashboardConfiguratorsKey, null) ?? []
  return configurators.find((configurator): configurator is TestModeConfigurator => configurator instanceof TestModeConfigurator) ?? null
}

export function useNewSearchEverywhere(): Ref<boolean> {
  const branchConfigurator = inject(branchConfiguratorKey, null)
  const testModeConfigurator = useTestModeConfigurator()
  return computed(() => {
    const branches = selectedToArray(branchConfigurator?.selected.value)
    const modes = selectedToArray(testModeConfigurator?.selected.value)
    return branches.length === 0 || !branches.every((branch) => isLegacySeBranch(branch)) || modes.includes(SPLIT_MODE)
  })
}
