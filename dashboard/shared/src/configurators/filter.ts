import { combineLatest, Observable } from "rxjs"
import { DataQuery } from "../dataQuery"
import { initZstdObservable } from "../zstd"
import { ServerConfigurator } from "./ServerConfigurator"

export function createFilterObservable(serverConfigurator: ServerConfigurator, filters: Array<FilterConfigurator>): Observable<unknown> {
  const filterObservables: Array<Observable<unknown>> = [initZstdObservable, serverConfigurator.createObservable()]
  if (filters.length !== 0) {
    filterObservables.push(...filters.map(it => it.createObservable()).filter((it: Observable<unknown> | null): it is Observable<unknown> => it !== null))
  }
  return combineLatest(filterObservables)
}

/**
 * Configurator may be used as a filter for dimension configurator - to limit dimension values. It is used when one dimension depends on another one.
 */
export interface FilterConfigurator {
  configureFilter(query: DataQuery): boolean

  /**
   * See {@link DataQueryConfigurator#createObservable}
   */
  createObservable(): Observable<unknown> | null
}

export function configureQueryFilters(query: DataQuery, filters: Array<FilterConfigurator>): boolean {
  for (const filter of filters) {
    if (!filter.configureFilter(query)) {
      return false
    }
  }
  return true
}