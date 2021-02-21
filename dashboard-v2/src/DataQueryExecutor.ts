import { loadJson } from "@/util/httpUtil"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery } from "@/dataQuery"
import { Ref, watch } from "vue"
import { debounce } from "@/util/debounce"

declare type DataQueryConsumer = (data: Array<Array<never>>, configuration: DataQueryExecutorConfiguration) => void

export class DataQueryExecutor {
  private readonly debouncedExecution = debounce(() => this.execute())

  listener: DataQueryConsumer | null = null

  constructor(private readonly configurators: Array<DataQueryConfigurator> = []) {
    for (const configurator of configurators) {
      const value = configurator.value
      if (value != null) {
        watch(value, this.debouncedExecution)
      }
    }
  }

  setListener(listener: DataQueryConsumer | null): void {
    this.listener = listener
  }

  addConfigurator(configurator: DataQueryConfigurator): void {
    this.configurators.push(configurator)

    const value = configurator.value
    if (value != null) {
      watch(value, this.debouncedExecution)
    }
  }

  watch(value: Ref<unknown>): void {
    watch(value, this.debouncedExecution)
  }

  scheduleLoad(): void {
    this.debouncedExecution()
  }

  execute(): void {
    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()
    for (const configurator of this.configurators) {
      if (!configurator.configure(query, configuration)) {
        return
      }
    }

    loadJson<Array<Array<never>>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, null, data => {
      // console.debug(data)
      if (this.listener != null) {
        this.listener(data, configuration)
      }
    })
  }
}