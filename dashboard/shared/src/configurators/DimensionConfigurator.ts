import { finalize, Observable, of, shareReplay, switchMap } from "rxjs"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter, encodeQuery } from "../dataQuery"
import { ServerConfigurator } from "./ServerConfigurator"
import { fromFetchWithRetryAndErrorHandling, refToObservable } from "./rxjs"

export abstract class BaseDimensionConfigurator implements DataQueryConfigurator {
  readonly selected = shallowRef<string | Array<string> | null>(null)
  readonly values = shallowRef<Array<string>>([])
  readonly loading = shallowRef(false)

  private readonly observable: Observable<string | Array<string> | null>

  protected constructor(readonly name: string, readonly multiple: boolean) {
    this.observable = refToObservable(this.selected, true).pipe(
      shareReplay(1),
    )
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.selected.value
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

    persistentStateManager?.add(name, this.selected)

    this.serverConfigurator.createObservable()
      .pipe(
        switchMap(() => {
          const query = this.createQuery()
          const configuration = new DataQueryExecutorConfiguration()
          if (!serverConfigurator.configureQuery(query, configuration)) {
            return of(null)
          }

          this.loading.value = true
          return fromFetchWithRetryAndErrorHandling<Array<string>>(`${configuration.getServerUrl()}/api/v1/load/${encodeQuery(query)}`).pipe(
            finalize(() => {
              this.loading.value = false
            }),
          )
        }),
      )
      .subscribe(data => {
        console.log("data for " + name)
        if (data != null) {
          this.values.value = data
        }
      })
  }
}

function configureQueryProducer(configuration: DataQueryExecutorConfiguration, filter: DataQueryFilter, values: Array<string>): void {
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
    },
  )
}