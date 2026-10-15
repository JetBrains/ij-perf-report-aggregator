import { map } from "rxjs/map"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, ServerConfigurator } from "../components/common/dataQuery"
import { DimensionConfigurator, filterSelected, selectedToArray } from "./DimensionConfigurator"
import { updateComponentState } from "./componentState"
import { combineLatest, fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"

/**
 * Code owners live in the meta (Postgres) DB, not in ClickHouse, so the owner list and the owner -> projects
 * mapping are loaded from /api/meta/projectOwners. As a filter, the selected owners are translated into a
 * `project in (...)` ClickHouse filter, which lets the Tests dimension depend on the selected code owners.
 */
export class CodeOwnerConfigurator extends DimensionConfigurator implements DataQueryConfigurator {
  // null until the mapping is loaded
  private readonly ownerToProjects = shallowRef<Map<string, string[]> | null>(null)
  private readonly ownerObservable: Observable<string | string[] | null>

  constructor(serverConfigurator: ServerConfigurator, persistentStateManager: PersistentStateManager | null, dbName: string, table: string) {
    super("owner", true)
    persistentStateManager?.add("owner", this.selected)
    this.ownerObservable = combineLatest([refToObservable(this.selected, true), refToObservable(this.ownerToProjects)])[map](([selected]) => selected)

    const params = new URLSearchParams({ db: dbName, table })
    this.state.loading = true
    updateComponentState(fromFetchWithRetryAndErrorHandling<Record<string, string> | null>(`${serverConfigurator.serverUrl}/api/meta/projectOwners?${params.toString()}`), this.state).subscribe(
      (data) => {
        const mapping = new Map<string, string[]>()
        for (const [project, owner] of Object.entries(data ?? {})) {
          const projects = mapping.get(owner)
          if (projects == null) {
            mapping.set(owner, [project])
          } else {
            projects.push(project)
          }
        }
        const owners = [...mapping.keys()].toSorted((a, b) => a.localeCompare(b))
        this.values.value = owners
        filterSelected(this, owners)
        this.ownerToProjects.value = mapping
      }
    )
  }

  createObservable(): Observable<string | string[] | null> {
    return this.ownerObservable
  }

  // owner is not a ClickHouse column — it only narrows down dependent dimensions
  configureQuery(_query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    return true
  }

  configureFilter(query: DataQuery): boolean {
    const owners = selectedToArray(this.selected.value)
    if (owners.length === 0) {
      return true
    }
    const ownerToProjects = this.ownerToProjects.value
    if (ownerToProjects == null) {
      // wait for the mapping, the observable emits again once it is loaded
      return false
    }
    const projects = owners.flatMap((owner) => ownerToProjects.get(owner) ?? [])
    // an owner without projects must produce an empty result, not an unfiltered one
    query.addFilter({ f: "project", v: projects.length === 0 ? [""] : projects })
    return true
  }
}
