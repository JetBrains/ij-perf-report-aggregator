import { DimensionConfigurator, SubDimensionConfigurator } from "../configurators/DimensionConfigurator"
import { PersistentStateManager } from "../PersistentStateManager"

const projectNameToTitle = new Map<string, string>()
// noinspection SpellCheckingInspection
projectNameToTitle.set("/q9N7EHxr8F1NHjbNQnpqb0Q0fs", "joda-time")
// noinspection SpellCheckingInspection
projectNameToTitle.set("73YWaW9bytiPDGuKvwNIYMK5CKI", "simple for IJ")
// noinspection SpellCheckingInspection
projectNameToTitle.set("1PbxeQ044EEghMOG9hNEFee05kM", "light edit (IJ)")

// noinspection SpellCheckingInspection
projectNameToTitle.set("j1a8nhKJexyL/zyuOXJ5CFOHYzU", "simple for PS")
// noinspection SpellCheckingInspection
projectNameToTitle.set("JeNLJFVa04IA+Wasc+Hjj3z64R0", "simple for WS")
// noinspection SpellCheckingInspection
projectNameToTitle.set("nC4MRRFMVYUSQLNIvPgDt+B3JqA", "Idea")
Object.seal(projectNameToTitle)

export function getProjectName(value: string): string {
  return projectNameToTitle.get(value) || value
}

export function createProjectConfigurator(productConfigurator: DimensionConfigurator,
                                          persistentStateManager: PersistentStateManager): SubDimensionConfigurator {
  return new SubDimensionConfigurator("project", productConfigurator, persistentStateManager, (a, b) => {
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