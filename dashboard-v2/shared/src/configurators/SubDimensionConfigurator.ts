import { watch } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQueryExecutorConfiguration, encodeQuery } from "../dataQuery"
import { TaskHandle } from "../util/debounce"
import { loadJson } from "../util/httpUtil"
import { BaseDimensionConfigurator, DimensionConfigurator } from "./DimensionConfigurator"

/**
 * Dimension, that depends on another dimension to get filtered data.
 */
export class SubDimensionConfigurator extends BaseDimensionConfigurator {
  constructor(name: string,
              private readonly parentDimensionConfigurator: DimensionConfigurator,
              persistentStateManager: PersistentStateManager | null = null,
              private readonly customValueSort: ((a: string, b: string) => number) | null = null) {
    super(name, false)

    if (persistentStateManager != null) {
      persistentStateManager.add(name, this.value)
    }

    watch(parentDimensionConfigurator.value, this.debouncedLoad.executeFunctionReference)
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

  load(taskHandle: TaskHandle): Promise<unknown> {
    const parentDimensionConfigurator = this.parentDimensionConfigurator
    const filterValue = parentDimensionConfigurator.value.value
    if (filterValue == null || filterValue.length === 0) {
      return Promise.resolve()
    }

    const query = this.createQuery()
    query.addFilter({field: parentDimensionConfigurator.name, value: filterValue})
    const configuration = new DataQueryExecutorConfiguration()
    if (!parentDimensionConfigurator.serverConfigurator.configureQuery(query, configuration)) {
      return Promise.resolve()
    }

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