import { switchMap } from "rxjs"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter, ServerConfigurator } from "../components/common/dataQuery"
import { DimensionConfigurator, filterSelected, loadDimension } from "./DimensionConfigurator"
import { updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"

export class BranchConfigurator extends DimensionConfigurator {
  constructor() {
    super("branch", true)
  }

  configureFilter(query: DataQuery): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return false
    }

    const values = Array.isArray(value) ? [...value] : [value]
    const OR_SEPARATOR = " or "

    const sqlClauses = values.map((val) => {
      return /\d+$/.test(val) ? `branch like '${val}%'` : `branch = '${val}'`
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
      size(): number {
        return values.length
      },
      mutate(index: number) {
        if (/\d+$/.test(values[index])) {
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
): DimensionConfigurator {
  const configurator = new BranchConfigurator()
  persistentStateManager?.add(persistentName, configurator.selected)

  createFilterObservable(serverConfigurator, filters)
    .pipe(
      switchMap(() => loadDimension("branch", serverConfigurator, filters, configurator.state)),
      updateComponentState(configurator.state)
    )
    .subscribe((data) => {
      if (data == null) {
        return
      }

      configurator.values.value = [
        ...new Set(
          data.map((value, _n, _a) => {
            const match = `${value}`.match(/^(\d+)\.\d+$/)
            return match ? match[1] : value
          })
        ),
      ]
      filterSelected(configurator, data, "branch")
    })
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
