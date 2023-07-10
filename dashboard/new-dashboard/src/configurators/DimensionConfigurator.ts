import { Observable, of, shareReplay, switchMap } from "rxjs"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter } from "../components/common/dataQuery"
import { ServerConfigurator } from "./ServerConfigurator"
import { ComponentState, createComponentState, updateComponentState } from "./componentState"
import { configureQueryFilters, createFilterObservable, FilterConfigurator } from "./filter"
import { fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"

export class DimensionConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly state = createComponentState()

  readonly selected = shallowRef<string | string[] | null>(null)
  readonly values = shallowRef<(string | boolean)[]>([])

  private readonly observable: Observable<string | string[] | null>

  constructor(
    readonly name: string,
    readonly multiple: boolean
  ) {
    this.observable = refToObservable(this.selected, true).pipe(shareReplay(1))
  }

  createObservable(): Observable<string | string[] | null> {
    return this.observable
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return false
    }

    const filter: DataQueryFilter = { f: this.name, v: value }
    if (this.multiple && Array.isArray(value)) {
      filter.v = value[0]
      if (value.length > 1) {
        configureQueryProducer(configuration, filter, value)
      }
    }
    query.addFilter(filter)
    return true
  }

  configureFilter(query: DataQuery): boolean {
    const value = this.selected.value
    if (value == null || value.length === 0) {
      return false
    }

    query.addFilter({ f: this.name, v: value })
    return true
  }
}

export function loadDimension(name: string, serverConfigurator: ServerConfigurator, filters: FilterConfigurator[], state: ComponentState) {
  const query = new DataQuery()
  query.addField({ n: name, sql: `distinct ${name}` })
  query.order = name
  query.flat = true

  const configuration = new DataQueryExecutorConfiguration()
  if (!serverConfigurator.configureQuery(query, configuration) || !configureQueryFilters(query, filters)) {
    return of(null)
  }

  state.loading = true
  return fromFetchWithRetryAndErrorHandling<string[]>(serverConfigurator.computeQueryUrl(query))
}

export function dimensionConfigurator(
  name: string,
  serverConfigurator: ServerConfigurator,
  persistentStateManager: PersistentStateManager | null,
  multiple: boolean = false,
  filters: FilterConfigurator[] = [],
  customValueSort: ((a: string, b: string) => number) | null = null
): DimensionConfigurator {
  const configurator = new DimensionConfigurator(name, multiple)
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

      if (customValueSort != null) {
        data.sort(customValueSort)
      }
      configurator.values.value = data

      filterSelected(configurator, data, name)
    })
  return configurator
}

export function filterSelected(configurator: DimensionConfigurator, data: string[], name: string) {
  const selectedRef = configurator.selected
  if (data.length === 0) {
    // do not update value - don't unset if values temporary not set
    console.debug(`[dimensionConfigurator(name=${name})] value list is empty`)
  } else {
    const selected = selectedRef.value
    if (Array.isArray(selected) && selected.length > 0) {
      const filtered = selected.filter((it) => data.includes(it))
      if (filtered.length !== selected.length) {
        selectedRef.value = filtered
      }
    } else if (selected == null || selected.length === 0 || !data.includes(selected as string)) {
      selectedRef.value = data[0]
    }
  }
}

export function configureQueryProducer(configuration: DataQueryExecutorConfiguration, filter: DataQueryFilter, values: string[]): void {
  configuration.queryProducers.push({
    size(): number {
      return values.length
    },
    mutate(index: number) {
      filter.v = values[index]
    },
    getSeriesName(index: number): string {
      return values[index]
    },
    getMeasureName(_index: number): string {
      return configuration.measures[0]
    },
  })
}
