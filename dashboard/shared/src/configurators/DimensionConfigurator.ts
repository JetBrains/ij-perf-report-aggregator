import { Observable } from "rxjs"
import { shallowRef } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter, encodeQuery } from "../dataQuery"
import { DebouncedTask, TaskHandle } from "../util/debounce"
import { loadJson } from "../util/httpUtil"
import { ServerConfigurator } from "./ServerConfigurator"
import { refToObservable } from "./rxjs"

export abstract class BaseDimensionConfigurator implements DataQueryConfigurator {
  readonly value = shallowRef<string | Array<string> | null>(null)
  readonly values = shallowRef<Array<string>>([])
  readonly loading = shallowRef(false)

  protected readonly debouncedLoad = new DebouncedTask(taskHandle => this.load(taskHandle))

  protected constructor(readonly name: string, readonly multiple: boolean) {
  }

  createObservable(): Observable<unknown> {
    return refToObservable(this.value, true)
  }

  scheduleLoadMetadata(immediately: boolean): void {
    this.debouncedLoad.execute(immediately)
  }

  abstract load(taskHandle: TaskHandle): Promise<unknown>

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.value.value
    if (value == null || value.length === 0) {
      // console.debug(`[dimensionConfigurator(name=${this.name})] value is not set`)
      return false
    }

    const filter: DataQueryFilter = {field: this.name, value}
    if (this.multiple && Array.isArray(value)) {
      filter.value = value[0]
      if (value.length > 1) {
        configureQueryProducer(configuration, filter, value)
      }
    }
    query.addFilter(filter)
    return true
  }

  protected createQuery(): DataQuery {
    const query = new DataQuery()
    query.addField({name: this.name, sql: `distinct ${this.name}`})
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
  }

  load(taskHandle: TaskHandle): Promise<unknown> {
    const query = this.createQuery()
    const configuration = new DataQueryExecutorConfiguration()
    if (!this.serverConfigurator.configureQuery(query, configuration)) {
      return Promise.resolve()
    }

    this.loading.value = true
    return loadJson<Array<string>>(`${configuration.getServerUrl()}/api/v1/load/${encodeQuery(query)}`, this.loading, taskHandle, data => {
      this.values.value = data
    })
  }
}

function configureQueryProducer(configuration: DataQueryExecutorConfiguration, filter: DataQueryFilter, values: Array<string>): void {
  configuration.extraQueryProducers.push({
      size(): number {
        return values.length
      },
      mutate(index: number) {
        filter.value = values[index]
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