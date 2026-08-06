import { debounceTime, Observable } from "rxjs"
import { Ref } from "vue"
import { DataQuery } from "../components/common/dataQuery"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

// Charts register their projects while mounting, so the set arrives in a burst — and again, much
// later, when a lazy accordion expands. Coalesce each burst into a single dependent reload.
const settleMs = 100

// Narrows a dimension query (the machine list) to the projects a dashboard actually charts, so a
// hardware group that never ran any of them is not offered as a choice that plots nothing.
class DashboardProjectsFilter implements FilterConfigurator {
  constructor(private readonly projects: Ref<string[]>) {}

  configureFilter(query: DataQuery): boolean {
    const projects = this.projects.value
    // Empty until the first chart registers; filtering on an empty set would match nothing.
    if (projects.length > 0) {
      query.addFilter({ f: "project", v: projects })
    }
    return true
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.projects, true).pipe(debounceTime(settleMs))
  }
}

export function dashboardProjectsFilter(projects: Ref<string[]>): FilterConfigurator {
  return new DashboardProjectsFilter(projects)
}
