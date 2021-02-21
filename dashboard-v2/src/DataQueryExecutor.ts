import { loadJson } from "./util/httpUtil"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery } from "./dataQuery"
import { Ref, watch } from "vue"
import { DebouncedTask, TaskHandle } from "./util/debounce"

declare type DataQueryConsumer = (data: Array<Array<never>>, configuration: DataQueryExecutorConfiguration) => void

export class DataQueryExecutor {
  private readonly debouncedExecution = new DebouncedTask(taskHandle => this.execute(taskHandle))

  listener: DataQueryConsumer | null = null
  // data loading maybe started before html element is ready
  private notConsumedData: {data: Array<Array<never>>; configuration: DataQueryExecutorConfiguration} | null = null

  constructor(public readonly configurators: Array<DataQueryConfigurator> = []) {
    for (const configurator of configurators) {
      const value = configurator.value
      if (value != null) {
        watch(value, this.debouncedExecution.executeFunctionReference)
      }
    }
  }

  setListener(listener: DataQueryConsumer | null): void {
    this.listener = listener
    if (listener != null) {
      const notConsumedData = this.notConsumedData
      if (notConsumedData != null) {
        this.notConsumedData = null
        listener(notConsumedData.data, notConsumedData.configuration)
      }
    }
  }

  addConfigurator(configurator: DataQueryConfigurator): void {
    this.configurators.push(configurator)

    const value = configurator.value
    if (value != null) {
      watch(value, this.debouncedExecution.executeFunctionReference)
    }
  }

  watch(value: Ref<unknown>): void {
    watch(value, this.debouncedExecution.executeFunctionReference)
  }

  scheduleLoad(): void {
    this.debouncedExecution.execute()
  }

  scheduleLoadIncludingConfigurators(immediately: boolean = false): void {
    for (const configurator of this.configurators) {
      if (configurator.scheduleLoad != null) {
        configurator.scheduleLoad(immediately)
      }
    }

    // Will be asked to load data if some configurator is changed,
    // but on route navigation all configurators are not changed, so, schedule load to make sure that chart will be not empty
    this.scheduleLoad()
  }

  execute(taskHandle: TaskHandle): Promise<unknown> {
    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()
    for (const configurator of this.configurators) {
      if (!configurator.configure(query, configuration)) {
        return Promise.resolve()
      }
    }

    if (taskHandle.isCancelled) {
      return Promise.resolve()
    }

    // console.debug(query)
    return loadJson<Array<Array<never>>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, null, taskHandle, data => {
      if (taskHandle.isCancelled) {
        return
      }

      console.debug(`[queryExecutor] loaded (listenerAdded=${this.listener != null})`)
      if (this.listener == null) {
        this.notConsumedData = {configuration, data}
      }
      else {
        this.listener(data, configuration)
      }
    })
  }
}