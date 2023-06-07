import { switchMap } from "rxjs"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter } from "../components/common/dataQuery"
import { configureQueryProducer, DimensionConfigurator, filterSelected, loadDimension } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"
import { updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"

export class PrivateBuildConfigurator extends DimensionConfigurator {
  constructor() {
    super("triggeredBy", true)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return true
    }

    const filter: DataQueryFilter = { f: this.name, v: "" }
    const values = Array.isArray(value) ? ["", ...value] : ["", value]
    configureQueryProducer(configuration, filter, values)
    query.addFilter(filter)
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
