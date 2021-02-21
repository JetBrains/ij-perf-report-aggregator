import { ref, watch } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { loadJson } from "../util/httpUtil"
import {
  DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery, ServerConfigurator,
} from "../dataQuery"
import { debounce } from "../util/debounce"

export interface Item {
  label: string
  value: string
}

abstract class BaseDimensionConfigurator implements DataQueryConfigurator {
  readonly value = ref<string>("")
  readonly values = ref<Array<string>>([])
  readonly loading = ref(false)

  protected readonly debouncedLoad = debounce(() => this.load(), 100)

  protected constructor(readonly name: string) {
  }

  scheduleLoad(): void {
    this.debouncedLoad()
  }

  abstract load(): void

  configure(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    const value = this.value.value
    if (value == null || value.length === 0) {
      console.warn(`Value of dimension ${this.name} is empty`)
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

  load(): void {
    const query = this.createQuery()
    const configuration = new DataQueryExecutorConfiguration()
    this.serverConfigurator.configure(query, configuration)
    this.loading.value = true
    loadJson<Array<string>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, this.loading, data => {
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

    watch(parentDimensionConfigurator.value, this.debouncedLoad)
    watch(this.values, (values) => {
      const {value} = this.value
      if (values.length === 0) {
        if (value != null && value.length !== 0) {
          this.value.value = ""
        }
      }
      else if (value == null || value.length === 0 || !values.includes(value)) {
        this.value.value = values[0]
      }
    })
  }

  load(): void {
    const {parentDimensionConfigurator} = this
    const filterValue = parentDimensionConfigurator.value.value
    if (filterValue == null || filterValue.length === 0) {
      return
    }

    const query = this.createQuery()
    query.addFilter({field: parentDimensionConfigurator.name, value: filterValue})
    const configuration = new DataQueryExecutorConfiguration()
    parentDimensionConfigurator.serverConfigurator.configure(query, configuration)
    this.loading.value = true
    loadJson<Array<string>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, this.loading, data => {
      this.loading.value = false
      if (this.customValueSort != null) {
        data.sort(this.customValueSort)
      }
      this.values.value = data
    })
  }
}
