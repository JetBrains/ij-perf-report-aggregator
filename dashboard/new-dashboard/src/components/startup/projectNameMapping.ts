import { DimensionConfigurator, dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { PersistentStateManager } from "../common/PersistentStateManager"

const projectNameToTitle = new Map<string, string>()
// noinspection SpellCheckingInspection
projectNameToTitle.set("/q9N7EHxr8F1NHjbNQnpqb0Q0fs", "joda-time")
// noinspection SpellCheckingInspection
projectNameToTitle.set("1PbxeQ044EEghMOG9hNEFee05kM", "light edit (IJ)")

// noinspection SpellCheckingInspection
projectNameToTitle.set("j1a8nhKJexyL/zyuOXJ5CFOHYzU", "simple for PS")
// noinspection SpellCheckingInspection
projectNameToTitle.set("JeNLJFVa04IA+Wasc+Hjj3z64R0", "simple for WS")
Object.seal(projectNameToTitle)

export function getProjectName(value: string): string {
  return projectNameToTitle.get(value) ?? value
}

export function createProjectConfigurator(
  productConfigurator: DimensionConfigurator,
  serverConfigurator: ServerConfigurator,
  persistentStateManager: PersistentStateManager,
  filters: FilterConfigurator[] = []
): DimensionConfigurator {
  return dimensionConfigurator("project", serverConfigurator, persistentStateManager, false /* doesn't matter */, [productConfigurator, ...filters], (a, b) => {
    const t1 = getProjectName(a)
    const t2 = getProjectName(b)
    if (t1.startsWith("simple ") && !t2.startsWith("simple ")) {
      return -1
    }
    if (t2.startsWith("simple ") && !t1.startsWith("simple ")) {
      return 1
    }
    return t1.localeCompare(t2)
  })
}
