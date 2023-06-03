import { of, switchMap } from "rxjs"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { DimensionConfigurator } from "./DimensionConfigurator"
import { ServerConfigurator } from "./ServerConfigurator"
import { ComponentState, updateComponentState } from "./componentState"
import { configureQueryFilters, createFilterObservable, FilterConfigurator } from "./filter"
import { fromFetchWithRetryAndErrorHandling } from "./rxjs"


export class BuildConfigurator extends DimensionConfigurator {

  constructor(name: string) {
    super(name, false)
  }

  configureQuery(_: DataQuery, _2: DataQueryExecutorConfiguration): boolean {
    return true
  }
}

function loadBuilds(serverConfigurator: ServerConfigurator, filters: FilterConfigurator[], state: ComponentState){
  const query = new DataQuery()
  query.addField({n: "build", sql: "distinct concat(toString(build_c1),'.',toString(build_c2),'.',toString(build_c3))"})
  query.order = "build"
  query.flat = true

  const configuration = new DataQueryExecutorConfiguration()
  if (!serverConfigurator.configureQuery(query, configuration) || !configureQueryFilters(query, filters)) {
    return of(null)
  }

  state.loading = true
  return fromFetchWithRetryAndErrorHandling<string[]>(serverConfigurator.computeQueryUrl(query))
}

export function buildConfigurator(name: string, serverConfigurator: ServerConfigurator,
                                  persistentStateManager: PersistentStateManager | null,
                                  filters: FilterConfigurator[] = []): DimensionConfigurator {
  const configurator = new BuildConfigurator(name)
  persistentStateManager?.add(name, configurator.selected)

  createFilterObservable(serverConfigurator, filters)
    .pipe(
      switchMap(() => loadBuilds(serverConfigurator, filters, configurator.state)),
      updateComponentState(configurator.state),
    )
    .subscribe(data => {
      if (data == null) {
        return
      }

      configurator.values.value = data.filter(value => value != "").filter(value => value.split(".").length == 3).map(value => {
        const buildParts = value.split(".")
        return buildParts[2] == "0" ? buildParts[0] + "." + buildParts[1] : value
      }).sort(compareBuilds)
    })
  return configurator
}

function compareBuilds(a: string, b: string) {
  const [branch1, build1] = a.split(".").map(value => Number.parseInt(value))
  const [branch2, build2] = b.split(".").map(value => Number.parseInt(value))
  if (branch1 < branch2) {
    return 1
  }
  else if (branch1 > branch2) {
    return -1
  }
  else {
    return build1 < build2 ? 1 : -1
  }
}
