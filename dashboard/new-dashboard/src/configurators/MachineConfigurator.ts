import { deepEqual } from "fast-equals"
import { combineLatest, distinctUntilChanged, Observable, shareReplay, switchMap } from "rxjs"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter, ServerConfigurator, toArray } from "../components/common/dataQuery"
import { createComponentState, updateComponentState } from "./componentState"
import { loadDimension } from "./DimensionConfigurator"
import { createFilterObservable, FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export class MachineConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly selected = shallowRef<string[]>([])
  readonly values = shallowRef<GroupedDimensionValue[]>([])

  private readonly observable: Observable<unknown>
  readonly state = createComponentState()
  private readonly groupNameToItem = new Map<string, GroupedDimensionValue>()

  readonly filters = shallowRef<FilterConfigurator[]>([])

  constructor(
    serverConfigurator: ServerConfigurator,
    persistentStateManager?: PersistentStateManager,
    initialFilters: FilterConfigurator[] = [],
    readonly multiple: boolean = true,
    predefinedMachines?: string[]
  ) {
    const name = "machine"
    persistentStateManager?.add(name, this.selected, (it) => toArray(it as never))
    if (predefinedMachines) {
      this.selected.value = predefinedMachines
    }

    this.filters.value = initialFilters
    const filterObservable = refToObservable(this.filters).pipe(
      switchMap((currentFilters) => {
        return createFilterObservable(serverConfigurator, currentFilters)
      }),
      shareReplay(1)
    )

    const listObservable = filterObservable.pipe(
      switchMap(() => {
        // The same distinct-machine query as any dimension, but answered by /api/machineGroups/
        // with the agents already bucketed into hardware-class groups by the backend (the sole
        // owner of the grouping — see pkg/machine).
        return loadDimension<MachineGroupResponseItem[]>(name, serverConfigurator, this.filters.value, this.state, "/api/machineGroups/")
      }),
      updateComponentState(this.state),
      shareReplay(1)
    )

    listObservable.subscribe((data) => {
      if (data == null) {
        return
      }

      this.groupNameToItem.clear()
      this.values.value = this.buildGroups(data)
      void this.normalizeSelectionToGroups(serverConfigurator)
    })

    // selected value may be a group name, so, we must re-execute query on machine list update
    this.observable = combineLatest([refToObservable(this.selected, true), listObservable]).pipe(distinctUntilChanged(deepEqual))
  }

  updateFilters(newFilters: FilterConfigurator[]) {
    this.filters.value = newFilters
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  // Drilldown links may carry a raw agent name in the selection; map it to its hardware-class
  // group so the dropdown highlights the group and the chart queries the whole class. The
  // backend owns the mapping — the group is not derivable from the loaded (filtered) member
  // lists, which may not contain that exact ephemeral instance.
  private async normalizeSelectionToGroups(serverConfigurator: ServerConfigurator): Promise<void> {
    const selected = this.selected.value
    const groupNames = new Set(this.values.value.map((group) => group.value))
    // A name present as a live agent in the loaded list is a deliberate single-agent choice
    // (leaves are selectable in the dropdown) — only names absent from the list (ephemeral
    // drilldown instances) are mapped to their group.
    const leafNames = new Set(this.values.value.flatMap((group) => group.children?.map((child) => child.value) ?? []))
    const unresolved = selected.filter((value) => !groupNames.has(value) && !leafNames.has(value))
    if (unresolved.length === 0) {
      return
    }

    const resolved = await Promise.all(unresolved.map((value) => fetchMachineGroup(serverConfigurator, value)))
    // The selection may have changed while the lookup was in flight (user picked another machine,
    // or a filter change re-ran this). Don't clobber the newer choice with a stale result.
    if (!deepEqual(this.selected.value, selected)) {
      return
    }
    const rawToGroup = new Map(unresolved.map((value, index) => [value, resolved[index]]))
    const next = [
      ...new Set(
        selected.map((value) => {
          // The backend answers "Unknown" for anything it can't map — including a selected group
          // name whose agents merely didn't run in the current filtered window. Never rewrite to
          // that bucket: the rewrite is persisted and would destroy the real selection.
          const group = rawToGroup.get(value)
          return group != null && group !== "Unknown" && groupNames.has(group) ? group : value
        })
      ),
    ]
    if (!deepEqual(next, selected)) {
      this.selected.value = next
    }
  }

  private buildGroups(data: MachineGroupResponseItem[]): GroupedDimensionValue[] {
    const grouped: GroupedDimensionValue[] = []
    for (const { group, machines } of data) {
      const item: GroupedDimensionValue = {
        value: group,
        children: machines.map((value) => ({ value })),
        icon: this.getIcons(group),
      }
      this.groupNameToItem.set(group, item)
      grouped.push(item)
    }
    grouped.sort((a, b) => a.value.localeCompare(b.value))
    return grouped
  }

  private getIcons(groupName: string): string {
    if (groupName.toLowerCase().startsWith("linux")) {
      return "pi icon-linux"
    } else if (groupName.toLowerCase().startsWith("mac")) {
      return "pi pi-apple"
    }
    return "pi pi-microsoft"
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const selected = this.selected.value
    if (selected.length === 0) {
      console.debug("machine is not configured")
      return false
    }

    if (!this.multiple) {
      this.configureQueryAsFilter(selected, query)
      return true
    }

    const groupNameToItem = this.groupNameToItem

    const values: string[] = []
    const filter: DataQueryFilter = { f: "machine", v: values }
    query.addFilter(filter)
    configuration.queryProducers.push({
      size(): number {
        return selected.length
      },

      mutate(index: number): void {
        const value = selected[index]
        const groupItem = groupNameToItem.get(value)
        values.length = 0
        if (groupItem == null) {
          values.push(value)
        } else {
          // it's group
          if (groupItem.children != null) {
            if (groupItem.children.length > 1) {
              filter.v = prefix(groupItem.children.map((it) => it.value)) + "%"
              filter.o = "like"
              return
            }
            for (const child of groupItem.children) {
              values.push(child.value)
            }
          }
          values.sort()
        }
      },

      getSeriesName(index: number): string {
        return selected.length > 1 ? selected[index] : ""
      },

      getMeasureName(index: number): string {
        return selected[index]
      },
    })
    return true
  }

  configureFilter(query: DataQuery): boolean {
    const value = this.selected.value
    if (value.length === 0) {
      return false
    }

    this.configureQueryAsFilter(value, query)
    return true
  }

  getMergedValue(): string {
    const values: string[] = this.getSelectedValues(this.selected.value)
    return prefix(values) + "%"
  }

  private configureQueryAsFilter(selected: string[], query: DataQuery) {
    const values = this.getSelectedValues(selected)

    if (values.length > 0) {
      // stable order of fields in query (caching)
      values.sort()
      if (values.length > 50) {
        query.addFilter({ f: "machine", v: prefix(values) + "%", o: "like" })
      } else {
        query.addFilter({ f: "machine", v: values })
      }
    }
  }

  private getSelectedValues(selected: string[]) {
    const values: string[] = []
    for (const value of selected) {
      const groupItem = this.groupNameToItem.get(value)
      if (groupItem == null) {
        values.push(value)
      } else {
        // it's group
        // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
        for (const child of groupItem.children!) {
          values.push(child.value)
        }
      }
    }
    return values
  }
}

export interface GroupedDimensionValue {
  value: string
  children?: GroupedDimensionValue[]
  icon?: string
}

function prefix(words: string[]): string {
  if (!words[0] || words.length == 1) return words[0] || ""
  let i = 0
  while (words[0][i] && words.every((w) => w[i] === words[0][i])) i++
  return words[0].slice(0, Math.max(0, i))
}

interface MachineGroupResponseItem {
  group: string
  machines: string[]
}

// Resolves a single raw agent name to its hardware-class group via the backend (the sole owner
// of the mapping). Falls back to the given name on error.
async function fetchMachineGroup(serverConfigurator: ServerConfigurator, machineName: string): Promise<string> {
  try {
    const response = await fetch(`${serverConfigurator.serverUrl}/api/machineGroup?machine=${encodeURIComponent(machineName)}`)
    const parsed = (await response.json()) as { group: string }
    return parsed.group
  } catch {
    return machineName
  }
}
