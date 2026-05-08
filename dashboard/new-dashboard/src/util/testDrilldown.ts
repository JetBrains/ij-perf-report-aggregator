import { Router } from "vue-router"
import { DBType } from "../components/common/sideBar/InfoSidebar"
import { dbTypeStore } from "../shared/dbTypes"

interface OpenTestDrilldownOptions {
  project: string
  measure: string
  // If non-empty, replaces the existing `branch` query parameter. Each entry is appended
  // separately so multi-branch URLs round-trip as `branch=A&branch=B`.
  branches?: readonly string[] | null
}

// Drills down from a "compare" view to the per-test page by swapping the dashboard's last path
// segment with `tests` (or `testsDev` for IntelliJ-dev DB), preserving the current query string,
// then overriding project/measure and the branch list. Iterates `currentRoute.query` manually
// because Vue Router's LocationQuery may carry array values that URLSearchParams' string coercion
// would collapse into `key=a,b` instead of `key=a&key=b`.
export function openTestDrilldown(router: Router, options: OpenTestDrilldownOptions): void {
  const currentRoute = router.currentRoute.value
  const parts = currentRoute.path.split("/")
  parts[parts.length - 1] = dbTypeStore().dbType == DBType.INTELLIJ_DEV ? "testsDev" : "tests"
  const testURL = parts.join("/")

  const queryParams = new URLSearchParams()
  for (const [key, val] of Object.entries(currentRoute.query)) {
    if (val == null) continue
    if (Array.isArray(val)) {
      for (const v of val) {
        if (v != null) queryParams.append(key, v)
      }
    } else {
      queryParams.append(key, val)
    }
  }
  queryParams.set("project", options.project)
  queryParams.set("measure", options.measure)
  queryParams.delete("metrics")
  queryParams.delete("tests")
  if (options.branches != null && options.branches.length > 0) {
    queryParams.delete("branch")
    for (const b of options.branches) queryParams.append("branch", b)
  }

  window.open(router.resolve(`${testURL}?${queryParams.toString()}`).href, "_blank")
}
