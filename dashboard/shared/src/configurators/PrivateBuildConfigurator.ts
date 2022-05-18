import { PersistentStateManager } from "../PersistentStateManager"
import { configureQueryProducer, DimensionConfigurator, loadDimension } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"
import { createFilterObservable, FilterConfigurator } from "./filter"
import { switchMap } from "rxjs"
import { updateComponentState } from "./componentState"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter } from "../dataQuery"

export class PrivateBuildConfigurator extends DimensionConfigurator {

  constructor() {
    super("triggeredBy", true)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return true
    }

    const filter: DataQueryFilter = {f: this.name, v: ""}
    const values: Array<string> = ["", ...value]
    configureQueryProducer(configuration, filter, values)
    query.addFilter(filter)
    return true
  }
}

export function privateBuildConfigurator(serverConfigurator: ServerConfigurator,
                                         persistentStateManager: PersistentStateManager | null,
                                         filters: Array<FilterConfigurator> = []): DimensionConfigurator {
  const configurator = new PrivateBuildConfigurator()
  persistentStateManager?.add("triggeredBy", configurator.selected)

  createFilterObservable(serverConfigurator, filters)
    .pipe(
      switchMap(() => loadDimension("triggeredBy", serverConfigurator, filters, configurator.state)),
      updateComponentState(configurator.state),
    )
    .subscribe(data => {
      if (data == null) {
        return
      }

      configurator.values.value = data.filter((value, _n, _a) => !(value == ""))

      const selectedRef = configurator.selected
      if (data.length === 0) {
        // do not update value - don't unset if values temporary not set
        console.debug("[dimensionConfigurator(name=triggeredBy)] value list is empty")
      }
      else {
        const selected = selectedRef.value
        if (selected instanceof Array && selected.length !== 0) {
          const filtered = selected.filter(it => data.includes(it))
          if (filtered.length !== selected.length) {
            selectedRef.value = filtered
          }
        }
        else if (selected == null || selected.length === 0 || !data.includes(selected as string)) {
          selectedRef.value = data[0]
        }
      }
    })
  return configurator
}
