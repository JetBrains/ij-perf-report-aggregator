import { Observable, switchMap } from "rxjs"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, ServerConfigurator } from "../components/common/dataQuery"
import { FilterConfigurator } from "./filter"
import { DimensionConfigurator, loadDimension, filterSelected } from "./DimensionConfigurator"
import { createFilterObservable } from "./filter"
import { updateComponentState } from "./componentState"

class ProjectLikeFilter implements FilterConfigurator {
  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "", q: "project like '%/measureStartup%' or project like '%/warmup%'" })
    return true
  }

  createObservable(): Observable<unknown> {
    return new Observable((subscriber) => {
      subscriber.next(null)
      subscriber.complete()
    })
  }
}

export function startupProjectConfigurator(
  serverConfigurator: ServerConfigurator,
  persistentStateManager: PersistentStateManager | null,
  multiple: boolean = false,
  filters: FilterConfigurator[] = [],
  customValueSort: ((a: string, b: string) => number) | null = null,
  aliases: Map<string, string> | null = null
): DimensionConfigurator {
  const projectLikeFilter = new ProjectLikeFilter()
  const allFilters = [projectLikeFilter, ...filters]

  const configurator = new DimensionConfigurator("project", multiple, aliases)
  persistentStateManager?.add("project", configurator.selected)
  createFilterObservable(serverConfigurator, allFilters)
    .pipe(
      switchMap(() => loadDimension("project", serverConfigurator, allFilters, configurator.state)),
      updateComponentState(configurator.state)
    )
    .subscribe((data) => {
      if (data == null) {
        return
      }

      // Merge related projects by removing /warmup and /measureStartup suffixes
      const mergedProjects = new Set<string>()
      for (const project of data) {
        const warmupIndex = project.indexOf("/warmup")
        const measureStartupIndex = project.indexOf("/measureStartup")

        if (warmupIndex !== -1) {
          mergedProjects.add(project.substring(0, warmupIndex))
        } else if (measureStartupIndex !== -1) {
          mergedProjects.add(project.substring(0, measureStartupIndex))
        } else {
          mergedProjects.add(project)
        }
      }

      const mergedData = Array.from(mergedProjects)

      if (customValueSort != null) {
        mergedData.sort(customValueSort)
      }
      configurator.values.value = mergedData

      filterSelected(configurator, mergedData)
    })

  return configurator
}
