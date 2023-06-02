import { switchMap } from "rxjs"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter } from "../dataQuery"
import { DimensionConfigurator, filterSelected, loadDimension } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"
import { updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"

export class BranchConfigurator extends DimensionConfigurator {

  constructor() {
    super("branch", true)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return true
    }

    const filter: DataQueryFilter = {f: this.name, v: ""}
    const values = Array.isArray(value) ? [...value] : [value]
    configuration.queryProducers.push({
        size(): number {
          return values.length
        },
        mutate(index: number) {
          if(/\d+/.test(values[index])){
            filter.v = values[index] + "%"
            filter.o = "like"
          } else {
            filter.v = values[index]
          }
        },
        getSeriesName(index: number): string {
          return values.length > 1 ? values[index] : ""
        },
        getMeasureName(_index: number): string {
          return values[_index]
        },
      },
    )
    query.addFilter(filter)
    return true
  }
}

export function createBranchConfigurator(serverConfigurator: ServerConfigurator,
                                         persistentStateManager: PersistentStateManager | null,
                                         filters: FilterConfigurator[] = []): DimensionConfigurator {
  const configurator = new BranchConfigurator()
  const name = "branch"
  persistentStateManager?.add(name, configurator.selected)

  createFilterObservable(serverConfigurator, filters)
    .pipe(
      switchMap(() => loadDimension(name, serverConfigurator, filters, configurator.state)),
      updateComponentState(configurator.state),
    )
    .subscribe(data => {
      if (data == null) {
        return
      }

      data.sort(a => a.includes("/") ? 1 : -1)

      configurator.values.value = data.filter((value, _n, _a) => !/\d+\.\d+/.test(value))
      filterSelected(configurator, data, name)
    })
  return configurator
}
