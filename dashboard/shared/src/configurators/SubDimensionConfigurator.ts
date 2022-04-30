import { debounceTime, of, switchMap } from "rxjs"
import { watch } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQueryExecutorConfiguration, serializeAndEncodeQueryForUrl } from "../dataQuery"
import { BaseDimensionConfigurator, DimensionConfigurator } from "./DimensionConfigurator"
import { fromFetchWithRetryAndErrorHandling } from "./rxjs"

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
      persistentStateManager.add(name, this.selected)
    }

    parentDimensionConfigurator.createObservable()
      .pipe(
        debounceTime(100),
        switchMap(() => {
          const filterValue = parentDimensionConfigurator.selected.value
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
          return fromFetchWithRetryAndErrorHandling<Array<string>>(`${configuration.getServerUrl()}/api/q/${serializeAndEncodeQueryForUrl(query)}`)
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
      const selectedRef = this.selected
      if (values.length === 0) {
        // do not update value - don't unset if values temporary not set
        console.debug(`[subDimensionConfigurator(name=${name})] value list is empty`)
      }
      else {
        const selected = selectedRef.value
        if (Array.isArray(selected) && selected.length !== 0) {
          const filtered = selected.filter(it => values.includes(it))
          if (filtered.length !== selected.length) {
            selectedRef.value = filtered
          }
        }
        else if (selected == null || selected.length === 0 || !values.includes(selected as string)) {
          this.selected.value = values[0]
        }
      }
    })
  }
}