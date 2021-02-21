import { ref, watch } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { loadJson } from "../util/httpUtil"
import {
  DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery, ServerConfigurator,
} from "../dataQuery"
import { DebouncedTask, TaskHandle } from "../util/debounce"

export interface Item {
  label: string
  value: string
}

abstract class BaseDimensionConfigurator implements DataQueryConfigurator {
  readonly value = ref<string>("")
  readonly values = ref<Array<string>>([])
  readonly loading = ref(false)

  protected readonly debouncedLoad = new DebouncedTask(taskHandle => this.load(taskHandle))

  protected constructor(readonly name: string) {
  }

  scheduleLoad(immediately: boolean): void {
    this.debouncedLoad.execute(immediately)
  }

  abstract load(taskHandle: TaskHandle): Promise<unknown>

  configure(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.value.value
    if (value == null || value.length === 0) {
      console.debug(`[dimensionConfigurator(name=${this.name})] value is not set`)
      return false
    }

    query.addFilter({field: this.name, value: value})
    return true
  }

  protected createQuery() {
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
              persistentStateManager: PersistentStateManager) {
    super(name)
    persistentStateManager.add(name, this.value)
  }

  load(taskHandle: TaskHandle): Promise<unknown> {
    const query = this.createQuery()
    const configuration = new DataQueryExecutorConfiguration()
    this.serverConfigurator.configure(query, configuration)
    this.loading.value = true
    return loadJson<Array<string>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, this.loading, taskHandle, data => {
      this.values.value = data
    })
  }
}

/**
 * Dimension, that depends on another dimension to get filtered data.
 */
export class SubDimensionConfigurator extends BaseDimensionConfigurator {
  constructor(name: string,
              private readonly parentDimensionConfigurator: DimensionConfigurator,
              persistentStateManager: PersistentStateManager,
              private readonly customValueSort?: (a: string, b: string) => number) {
    super(name)

    persistentStateManager.add(name, this.value)

    watch(parentDimensionConfigurator.value, this.debouncedLoad.executeFunctionReference)
    watch(this.values, values => {
      const value = this.value
      if (values.length === 0) {
        if (value.value != null && value.value.length !== 0) {
          this.value.value = ""
        }
      }
      else if (value.value == null || value.value.length === 0 || !values.includes(value.value)) {
        this.value.value = values[0]
      }
    })
  }

  load(taskHandle: TaskHandle): Promise<unknown> {
    const parentDimensionConfigurator = this.parentDimensionConfigurator
    const filterValue = parentDimensionConfigurator.value.value
    if (filterValue == null || filterValue.length === 0) {
      return Promise.resolve()
    }

    const query = this.createQuery()
    query.addFilter({field: parentDimensionConfigurator.name, value: filterValue})
    const configuration = new DataQueryExecutorConfiguration()
    parentDimensionConfigurator.serverConfigurator.configure(query, configuration)
    this.loading.value = true
    return loadJson<Array<string>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, this.loading, taskHandle, data => {
      this.loading.value = false
      if (this.customValueSort != null) {
        data.sort(this.customValueSort)
      }
      this.values.value = data
    })
  }
}
