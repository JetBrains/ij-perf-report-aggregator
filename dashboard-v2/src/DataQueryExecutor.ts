import { loadJson } from "./util/httpUtil"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery, ServerConfigurator } from "./dataQuery"
import { watch } from "vue"
import { DebouncedTask, debounceSync, TaskHandle } from "./util/debounce"
import { PersistentStateManager } from "./PersistentStateManager"

declare type DataQueryConsumer = (data: Array<Array<never>>, configuration: DataQueryExecutorConfiguration) => void

export class DataQueryExecutor {
  private readonly debouncedExecution = new DebouncedTask(taskHandle => this.execute(taskHandle))
  public readonly scheduleLoadIncludingConfiguratorsFunctionReference = (): void => this.scheduleLoadIncludingConfigurators()

  listener: DataQueryConsumer | null = null
  // data loading maybe started before html element is ready
  private notConsumedData: {data: Array<Array<never>>; configuration: DataQueryExecutorConfiguration} | null = null

  private readonly dependents: Array<DataQueryExecutor> = []

  private _lastQuery: DataQuery | null = null
  get lastQuery(): DataQuery | null {
    return this._lastQuery
  }

  /**
   * `isGroup = true` means that this DataQueryExecutor only manages dependent executors but doesn't load data itself.
   */
  constructor(private readonly configurators: Array<DataQueryConfigurator> = [],
              private readonly isGroup: boolean = false,
              private readonly parent: DataQueryExecutor | null = null) {
    for (const configurator of configurators) {
      const value = configurator.value
      if (value != null) {
        // watch(value, this.debouncedExecution.executeFunctionReference)
        watch(value, v => {
          console.debug(`[queryExecutor] schedule execution on configurator value change (new value=${v})`)
          this.debouncedExecution.execute()
        })
      }
    }
  }

  createSub(configurators: Array<DataQueryConfigurator>): DataQueryExecutor {
    const dataQueryExecutor = new DataQueryExecutor(configurators, false, this)
    this.dependents.push(dataQueryExecutor)
    return dataQueryExecutor
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

  // addConfigurator(configurator: DataQueryConfigurator): void {
  //   this.configurators.push(configurator)
  //
  //   const value = configurator.value
  //   if (value != null) {
  //     watch(value, this.debouncedExecution.executeFunctionReference)
  //   }
  // }

  // watch(value: Ref<unknown>): void {
  //   watch(value, this.debouncedExecution.executeFunctionReference)
  // }

  scheduleLoad(): void {
    this.debouncedExecution.execute()
  }

  scheduleLoadIncludingConfigurators(immediately: boolean = false): void {
    console.debug(`${this.getDebugTag()} scheduleLoadIncludingConfigurators (immediately=${immediately})`)

    for (const configurator of this.configurators) {
      if (configurator.scheduleLoad != null) {
        configurator.scheduleLoad(immediately)
      }
    }

    // Will be asked to load data if some configurator is changed,
    // but on route navigation all configurators are not changed, so, schedule load to make sure that chart will be not empty
    this.debouncedExecution.execute(immediately)
  }

  private getDebugTag() {
    if (this.isGroup) {
      return `[queryExecutor(isGroup=true, dependentCount=${this.dependents.length})]`
    }
    else {
      return `[queryExecutor${(this.parent == null ? "" : "(isDependent=true)")}]`
    }
  }

  async execute(taskHandle: TaskHandle): Promise<unknown> {
    const query = new DataQuery()
    const configuration = new DataQueryExecutorConfiguration()

    if (this.parent != null) {
      for (const configurator of this.parent.configurators) {
        if (!configurator.configure(query, configuration)) {
          return
        }
      }
    }

    if (!this.isGroup) {
      for (const configurator of this.configurators) {
        if (!configurator.configure(query, configuration)) {
          return
        }
      }
    }

    if (taskHandle.isCancelled) {
      return
    }

    if (this.dependents.length > 0) {
      await Promise.allSettled(this.dependents.map(it => it.execute(taskHandle)))
    }

    if (this.isGroup) {
      return
    }

    return await loadJson<Array<Array<never>>>(`${configuration.serverUrl}/api/v1/load/${encodeQuery(query)}`, null, taskHandle, data => {
      if (taskHandle.isCancelled) {
        return
      }

      this._lastQuery = query

      console.debug(`[queryExecutor] loaded (listenerAdded=${this.listener != null}, query=${JSON.stringify(query)})`)
      if (this.listener == null) {
        this.notConsumedData = {configuration, data}
      }
      else {
        this.listener(data, configuration)
      }
    })
  }
}

export function initDataComponent(serverConfigurator: ServerConfigurator,
                                  persistentStateManager: PersistentStateManager,
                                  dataQueryExecutor: DataQueryExecutor): void {
  persistentStateManager.init()
  dataQueryExecutor.scheduleLoadIncludingConfigurators(true)

  // after persistentStateManager.init to avoid double calling of scheduleLoadIncludingConfigurators
  watch(serverConfigurator.server, debounceSync(() => dataQueryExecutor.scheduleLoadIncludingConfiguratorsFunctionReference, 900))
}