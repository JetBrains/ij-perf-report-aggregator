import { switchMap } from "rxjs"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter, ServerConfigurator } from "../components/common/dataQuery"
import { configureQueryProducer, DimensionConfigurator, filterSelected, loadDimension } from "./DimensionConfigurator"
import { updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"

export class PrivateBuildConfigurator extends DimensionConfigurator {
  constructor() {
    super("triggeredBy", true)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    const filter: DataQueryFilter = { f: this.name, v: "" }
    if (value == null || value.length === 0) {
      query.addFilter(filter)
      return true
    }
    const values = Array.isArray(value) ? ["", ...value] : ["", value]
    configureQueryProducer(configuration, filter, values)
    query.addFilter(filter)
    return true
  }

  configureFilter(query: DataQuery): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      query.addFilter({ f: this.name, v: "" })
    } else {
      if (Array.isArray(value)) {
        query.addFilter({ f: this.name, v: [...value, ""] })
      } else {
        query.addFilter({ f: this.name, v: [value, ""] })
      }
    }
    return true
  }
}

export function privateBuildConfigurator(
  serverConfigurator: ServerConfigurator,
  persistentStateManager: PersistentStateManager | null,
  filters: FilterConfigurator[] = []
): DimensionConfigurator {
  const configurator = new PrivateBuildConfigurator()
  const name = "triggeredBy"
  persistentStateManager?.add(name, configurator.selected)

  createFilterObservable(serverConfigurator, filters)
    .pipe(
      switchMap(() => loadDimension(name, serverConfigurator, filters, configurator.state)),
      updateComponentState(configurator.state)
    )
    .subscribe((data) => {
      if (data == null) {
        return
      }

      configurator.values.value = data.filter((value, _n, _a) => value != "")

      filterSelected(configurator, data, name)
    })
  return configurator
}
