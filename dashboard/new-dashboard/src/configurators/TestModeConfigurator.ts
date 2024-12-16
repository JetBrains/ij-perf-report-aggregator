import { switchMap } from "rxjs"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryExecutorConfiguration, DataQueryFilter, ServerConfigurator } from "../components/common/dataQuery"
import { DimensionConfigurator, filterSelected, loadDimension } from "./DimensionConfigurator"
import { updateComponentState } from "./componentState"
import { createFilterObservable, FilterConfigurator } from "./filter"

export const defaultModeName = "default"

export class TestModeConfigurator extends DimensionConfigurator {
  constructor(multiple: boolean) {
    super("mode", multiple)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    const filter: DataQueryFilter = { f: this.name, v: "" }
    if (value == null || value.length === 0 || value == defaultModeName) {
      query.addFilter(filter)
      return true
    }

    const values = Array.isArray(value) ? value.map((v) => (v == defaultModeName ? "" : v)) : [value]
    configuration.queryProducers.push({
      size(): number {
        return values.length
      },
      mutate(index: number) {
        filter.v = values[index]
      },
      getSeriesName(index: number): string {
        return values[index] == "" ? defaultModeName : values[index]
      },
      getMeasureName(_index: number): string {
        return configuration.measures[0]
      },
    })
    query.addFilter(filter)
    return true
  }

  configureFilter(query: DataQuery): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0 || value == defaultModeName) {
      query.addFilter({ f: this.name, v: "" })
    } else {
      if (Array.isArray(value)) {
        const withDefaultValue = value.map((v) => (v == defaultModeName ? "" : v))
        query.addFilter({ f: this.name, v: [...withDefaultValue] })
      } else {
        query.addFilter({ f: this.name, v: [value] })
      }
    }
    return true
  }
}

export function createTestModeConfigurator(
  serverConfigurator: ServerConfigurator,
  persistentStateManager: PersistentStateManager | null,
  filters: FilterConfigurator[] = [],
  persistentName: string = "mode",
  multiple: boolean = true
): DimensionConfigurator {
  const configurator = new TestModeConfigurator(multiple)
  const name = "mode"
  persistentStateManager?.add(persistentName, configurator.selected)

  createFilterObservable(serverConfigurator, filters)
    .pipe(
      switchMap(() => loadDimension(name, serverConfigurator, filters, configurator.state)),
      updateComponentState(configurator.state)
    )
    .subscribe((data) => {
      if (data == null) {
        return
      }

      const fetchedValues = data.filter((value, _n, _a) => value != "")
      configurator.values.value = [defaultModeName, ...fetchedValues]

      filterSelected(configurator, [...data, defaultModeName])
    })
  return configurator
}
