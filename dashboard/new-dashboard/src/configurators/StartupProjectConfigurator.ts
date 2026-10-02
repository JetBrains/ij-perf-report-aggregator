import { ColdObservable } from "rxjs"
import { switchMap } from "rxjs/switch-map"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, ServerConfigurator } from "../components/common/dataQuery"
import { FilterConfigurator, createFilterObservable } from "./filter"
import { DimensionConfigurator, loadDimension, filterSelected, selectedToArray } from "./DimensionConfigurator"
import { updateComponentState } from "./componentState"

// Runs are reported under "<project>/measureStartup" and "<project>/warmup"; the selector shows
// the "<project>" stem, so queries here translate between the two forms.
const startupSuffixes = ["/measureStartup", "/warmup"]

// `projectFilter` narrows the projects to those whose name contains it (e.g. "JetBrains Light"
// runs share the "idea" table with the regular IDEA startup runs).
class ProjectLikeFilter implements FilterConfigurator {
  constructor(private readonly projectFilter: string | null) {}

  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "", q: startupSuffixes.map((suffix) => `project like '%${suffix}%'`).join(" or ") })
    if (this.projectFilter != null) {
      // Value form, not a predicate: the backend escapes it for us.
      query.addFilter({ f: "project", v: `%${this.projectFilter}%`, o: "like" })
    }
    return true
  }

  createObservable(): Observable<unknown> {
    return new ColdObservable((subscriber) => {
      subscriber.next(null)
      subscriber.complete()
    })
  }
}

// Narrows the machine list to the selected startup projects, so a hardware group that never ran
// them is not offered in the selector.
class SelectedStartupProjectsFilter implements FilterConfigurator {
  private readonly anyStartupProject: ProjectLikeFilter

  constructor(
    private readonly configurator: DimensionConfigurator,
    projectFilter: string | null
  ) {
    this.anyStartupProject = new ProjectLikeFilter(projectFilter)
  }

  configureFilter(query: DataQuery): boolean {
    const selected = selectedToArray(this.configurator.selected.value)
    if (selected.length === 0) {
      // A filter matching nothing would empty the machine list.
      return this.anyStartupProject.configureFilter(query)
    }

    // Value form, not a predicate: the backend escapes the names for us.
    const projects = selected.flatMap((project) => startupSuffixes.map((suffix) => project + suffix))
    // stable order of fields in query (caching)
    projects.sort()
    query.addFilter({ f: "project", v: projects })
    return true
  }

  createObservable(): Observable<unknown> {
    return this.configurator.createObservable()
  }
}

export function selectedStartupProjectsFilter(configurator: DimensionConfigurator, projectFilter: string | null = null): FilterConfigurator {
  return new SelectedStartupProjectsFilter(configurator, projectFilter)
}

export function startupProjectConfigurator(
  serverConfigurator: ServerConfigurator,
  persistentStateManager: PersistentStateManager | null,
  multiple: boolean = false,
  filters: FilterConfigurator[] = [],
  projectFilter: string | null = null,
  customValueSort: ((a: string, b: string) => number) | null = null,
  aliases: Map<string, string> | null = null
): DimensionConfigurator {
  const projectLikeFilter = new ProjectLikeFilter(projectFilter)
  const allFilters = [projectLikeFilter, ...filters]

  const configurator = new DimensionConfigurator("project", multiple, aliases)
  persistentStateManager?.add("project", configurator.selected)
  updateComponentState(
    createFilterObservable(serverConfigurator, allFilters)[switchMap](() => loadDimension("project", serverConfigurator, allFilters, configurator.state)),
    configurator.state
  ).subscribe((data) => {
    if (data == null) {
      return
    }

    const mergedProjects = new Set<string>()
    for (const project of data) {
      const suffix = startupSuffixes.find((it) => project.includes(it))
      mergedProjects.add(suffix == null ? project : project.slice(0, project.indexOf(suffix)))
    }

    const mergedData = [...mergedProjects]

    if (customValueSort != null) {
      mergedData.sort(customValueSort)
    }
    configurator.values.value = mergedData

    filterSelected(configurator, mergedData)
  })

  return configurator
}
