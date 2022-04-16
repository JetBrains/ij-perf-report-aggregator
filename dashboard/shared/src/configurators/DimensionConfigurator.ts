import { deepEqual } from "fast-equals"
import { debounceTime, distinctUntilChanged, Observable, of, switchMap } from "rxjs"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter, encodeQuery } from "../dataQuery"
import { ServerConfigurator } from "./ServerConfigurator"
import { fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"

export abstract class BaseDimensionConfigurator implements DataQueryConfigurator {
  readonly value = shallowRef<string | Array<string> | null>(null)
  readonly values = shallowRef<Array<string>>([])
  readonly loading = shallowRef(false)

  protected constructor(readonly name: string, readonly multiple: boolean) {
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.value, true)
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.value.value
    if (value == null || value.length === 0) {
      return false
    }

    const filter: DataQueryFilter = {f: this.name, v: value}
    if (this.multiple && Array.isArray(value)) {
      filter.v = value[0]
      if (value.length > 1) {
        configureQueryProducer(configuration, filter, value)
      }
    }
    query.addFilter(filter)
    return true
  }

  protected createQuery(): DataQuery {
    const query = new DataQuery()
    query.addField({n: this.name, sql: `distinct ${this.name}`})
    query.order = [this.name]
    query.table = "report"
    query.flat = true
    return query
  }
}

export class DimensionConfigurator extends BaseDimensionConfigurator {
  constructor(name: string,
              readonly serverConfigurator: ServerConfigurator,
              persistentStateManager: PersistentStateManager | null,
              multiple: boolean = false) {
    super(name, multiple)

    persistentStateManager?.add(name, this.value)

    this.serverConfigurator.createObservable()
      .pipe(
        debounceTime(100),
        distinctUntilChanged(deepEqual),
        switchMap(() => {
          const query = this.createQuery()
          const configuration = new DataQueryExecutorConfiguration()
          if (!serverConfigurator.configureQuery(query, configuration)) {
            return of(null)
          }

          this.loading.value = true
          return fromFetchWithRetryAndErrorHandling<Array<string>>(`${configuration.getServerUrl()}/api/v1/load/${encodeQuery(query)}`)
        }),
      )
      .subscribe(data => {
        this.loading.value = false
        if (data != null) {
          this.values.value = data
        }
      })
  }
}

function configureQueryProducer(configuration: DataQueryExecutorConfiguration, filter: DataQueryFilter, values: Array<string>): void {
  configuration.extraQueryProducers.push({
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
    },
  )
}