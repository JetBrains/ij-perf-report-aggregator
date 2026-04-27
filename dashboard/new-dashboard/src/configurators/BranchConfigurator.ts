import { getActivePinia } from "pinia"
import { map, merge, switchMap } from "rxjs"
import { ref, Ref, toRef, watch } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter, ServerConfigurator } from "../components/common/dataQuery"
import { useSettingsStore } from "../components/settings/settingsStore"
import { DimensionConfigurator, filterSelected, loadDimension } from "./DimensionConfigurator"
import { updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export class BranchConfigurator extends DimensionConfigurator {
  groupBranches: Ref<boolean> = ref(true)
  hasGroupableBranches: Ref<boolean> = ref(false)

  constructor() {
    super("branch", true)
  }

  createObservable() {
    return merge(super.createObservable(), refToObservable(this.groupBranches).pipe(map(() => null)))
  }

  configureFilter(query: DataQuery): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return false
    }

    const values = Array.isArray(value) ? [...value] : [value]
    const OR_SEPARATOR = " or "

    const sqlClauses = values.map((val) => {
      return this.groupBranches.value && /\d+$/.test(val) ? `branch like '${val}%'` : `branch = '${val}'`
    })

    const sql = sqlClauses.join(OR_SEPARATOR).replace(/^branch/, "")
    query.addFilter({ f: "branch", q: sql })
    return true
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return true
    }

    const filter: DataQueryFilter = { f: this.name, v: "" }
    const values = Array.isArray(value) ? [...value] : [value]
    configuration.queryProducers.push({
      isBranchDimension: true,
      size(): number {
        return values.length
      },
      mutate: (index: number) => {
        if (this.groupBranches.value && /^\d{3}$/.test(values[index])) {
          filter.v = values[index] + "%"
          filter.o = "like"
        } else {
          filter.v = values[index]
          filter.o = undefined
        }
      },
      getSeriesName(index: number): string {
        return values.length > 1 ? values[index] : ""
      },
      getMeasureName(_index: number): string {
        return values[_index]
      },
    })
    query.addFilter(filter)
    return true
  }
}

export function createBranchConfigurator(
  serverConfigurator: ServerConfigurator,
  persistentStateManager: PersistentStateManager | null,
  filters: FilterConfigurator[] = [],
  persistentName: string = "branch"
): BranchConfigurator {
  const configurator = new BranchConfigurator()
  persistentStateManager?.add(persistentName, configurator.selected)

  const settingsStore = getActivePinia() != null ? useSettingsStore() : null
  if (settingsStore != null) {
    configurator.groupBranches = toRef(settingsStore, "groupBranches")
  }

  let rawData: string[] = []

  function processData(data: string[]) {
    configurator.hasGroupableBranches.value = data.some((v) => /^\d+\.\d+$/.test(v))
    configurator.values.value = configurator.groupBranches.value
      ? [
          ...new Set(
            data.map((value) => {
              const match = /^(\d+)\.\d+$/.exec(value)
              return match ? match[1] : value
            })
          ),
        ]
      : [...new Set(data)]
    filterSelected(configurator, configurator.values.value as string[])
  }

  createFilterObservable(serverConfigurator, filters)
    .pipe(
      switchMap(() => loadDimension("branch", serverConfigurator, filters, configurator.state)),
      updateComponentState(configurator.state)
    )
    .subscribe((data) => {
      if (data == null) {
        return
      }
      rawData = data
      processData(data)
    })

  if (settingsStore != null) {
    watch(
      () => settingsStore.groupBranches,
      () => {
        if (rawData.length > 0) {
          processData(rawData)
        }
      }
    )
  }

  return configurator
}

export function sortBranches(a: string, b: string): number {
  if (a === "master") return -1
  if (b === "master") return 1

  // Then numbers
  const isANumber = !Number.isNaN(Number(a))
  const isBNumber = !Number.isNaN(Number(b))

  if (isANumber && !isBNumber) return -1
  if (!isANumber && isBNumber) return 1
  if (isANumber && isBNumber) return Number(b) - Number(a)

  // Then strings that contain "/"
  const hasASlash = a.includes("/")
  const hasBSlash = b.includes("/")

  if (hasASlash && !hasBSlash) return -1
  if (!hasASlash && hasBSlash) return 1

  // Then everything else
  return a.localeCompare(b)
}
