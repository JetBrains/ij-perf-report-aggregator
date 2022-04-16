import { deepEqual } from "fast-equals"
import { debounceTime, distinctUntilChanged, of, switchMap } from "rxjs"
import { watch } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQueryExecutorConfiguration, encodeQuery } from "../dataQuery"
import { BaseDimensionConfigurator, DimensionConfigurator } from "./DimensionConfigurator"
import { fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"

/**
 * Dimension, that depends on another dimension to get filtered data.
 */
export class SubDimensionConfigurator extends BaseDimensionConfigurator {
  constructor(name: string,
              parentDimensionConfigurator: DimensionConfigurator,
              persistentStateManager: PersistentStateManager | null = null,
              private readonly customValueSort: ((a: string, b: string) => number) | null = null) {
    super(name, false)

    if (persistentStateManager != null) {
      persistentStateManager.add(name, this.value)
    }

    refToObservable(parentDimensionConfigurator.value, true)
      .pipe(
        debounceTime(100),
        distinctUntilChanged(deepEqual),
        switchMap(() => {
          const filterValue = parentDimensionConfigurator.value.value
          if (filterValue == null || filterValue.length === 0) {
            return of(null)
          }

          const query = this.createQuery()
          query.addFilter({f: parentDimensionConfigurator.name, v: filterValue})
          const configuration = new DataQueryExecutorConfiguration()
          if (!parentDimensionConfigurator.serverConfigurator.configureQuery(query, configuration)) {
            return of(null)
          }

          this.loading.value = true
          return fromFetchWithRetryAndErrorHandling<Array<string>>(`${configuration.getServerUrl()}/api/v1/load/${encodeQuery(query)}`)
        }),
      )
      .subscribe(data => {
        this.loading.value = false

        if (data == null) {
          return
        }

        if (this.customValueSort != null) {
          data.sort(this.customValueSort)
        }
        this.values.value = data
      })
    watch(this.values, values => {
      const value = this.value
      if (values.length === 0) {
        // do not update value - don't unset if values temporary not set
        console.debug(`[subDimensionConfigurator(name=${name})] value list is empty`)
      }
      else if (value.value == null || value.value.length === 0 ||
        !values.includes(value.value as string /* multiple selection not supported yet */)) {
        const newValue = values[0]
        console.debug(`[subDimensionConfigurator(name=${name})] values loaded and selected value is set to ${newValue}`)
        this.value.value = newValue
      }
    })
  }
}