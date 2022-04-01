import { decode as risonDecode } from "rison-node"
import { provide, watch } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery } from "./dataQuery"
import { dataQueryExecutorKey } from "./injectionKeys"
import { PersistentStateManager } from "./PersistentStateManager"
import { DebouncedTask, TaskHandle } from "./util/debounce"
import { loadJson } from "./util/httpUtil"

export declare type DataQueryResult = Array<Array<Array<string | number>>>
export declare type DataQueryConsumer = (data: DataQueryResult, configuration: DataQueryExecutorConfiguration) => void

export class DataQueryExecutor {
  private readonly debouncedExecution = new DebouncedTask(taskHandle => this.execute(taskHandle))
  public readonly scheduleLoadIncludingConfiguratorsFunctionReference = (): void => this.scheduleLoadIncludingConfigurators()

  listener: DataQueryConsumer | null = null
  // data loading maybe started before html element is ready
  private notConsumedData: { data: Array<Array<Array<never>>>; configuration: DataQueryExecutorConfiguration } | null = null

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

  scheduleLoad(): void {
    this.debouncedExecution.execute()
  }

  init(): void {
    for (const configurator of this.configurators) {
      const value = configurator.value
      if (value == null) {
        continue
      }

      let debouncedExecution = this.debouncedExecution
      const valueChangeDelay = configurator.valueChangeDelay
      if (valueChangeDelay !== undefined) {
        debouncedExecution = new DebouncedTask(taskHandle => this.execute(taskHandle), valueChangeDelay)
      }
      watch(value, v => {
        // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
        console.debug(`[queryExecutor] schedule execution on configurator value change (new value=${v})`)
        debouncedExecution.execute()
      }, {deep: typeof value.value === "object" && value.value !== null})
    }

    this.scheduleLoadIncludingConfigurators(true)
  }

  scheduleLoadIncludingConfigurators(immediately: boolean = false): void {
    // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
    console.debug(`${this.getDebugTag()} scheduleLoadIncludingConfigurators (immediately=${immediately})`)

    for (const configurator of this.configurators) {
      if (configurator.scheduleLoadMetadata != null) {
        configurator.scheduleLoadMetadata(immediately)
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
        if (!configurator.configureQuery(query, configuration)) {
          return
        }
      }
    }

    if (!this.isGroup) {
      for (const configurator of this.configurators) {
        if (!configurator.configureQuery(query, configuration)) {
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

    let serializedQuery = "!("
    serializedQuery += encodeQuery(query)
    //TODO: implement a general algorithm based on recursion to cover any number of multivalued dimensions
    if(configuration.extraQueryProducers.length == 0){
      //empty
    } else if(configuration.extraQueryProducers.length === 1){
      const extraQueryProducer = configuration.extraQueryProducers[0]
      let done = false
      do {
        done = !extraQueryProducer.mutate()
        serializedQuery += "," + encodeQuery(query)
      }
      while (!done)
    } else if(configuration.extraQueryProducers.length === 2){
      const extraQueryProducer1 = configuration.extraQueryProducers[0]
      const extraQueryProducer2 = configuration.extraQueryProducers[1]
      let done1 = false
      let done2 = false
      let firstSerieIndex = 0
      extraQueryProducer2.setSecondSerieName(extraQueryProducer1.getSeriesName(firstSerieIndex))
      do {
        done2 = !extraQueryProducer2.mutate()
        serializedQuery += "," + encodeQuery(query)
      }
      while (!done2)
      extraQueryProducer2.reset()
      do {
        done1 = !extraQueryProducer1.mutate()
        firstSerieIndex++
        extraQueryProducer2.setSecondSerieName(extraQueryProducer1.getSeriesName(firstSerieIndex))
        serializedQuery += "," + encodeQuery(query)
        let done2 = false
        do {
          done2 = !extraQueryProducer2.mutate()
          serializedQuery += "," + encodeQuery(query)
        }
        while (!done2)
        extraQueryProducer2.reset()
        done2 = false
      }
      while (!done1)
    }
    serializedQuery += ")"

    return await loadJson<Array<Array<Array<never>>>>(`${configuration.getServerUrl()}/api/v1/load/${serializedQuery}`, null, taskHandle, data => {
      if (taskHandle.isCancelled) {
        return
      }

      this._lastQuery = query

      // console.debug(`[queryExecutor] loaded (listenerAdded=${this.listener != null}, query=${JSON.stringify(risonDecode(serializedQuery), null, 2)})`)
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      console.debug(`[queryExecutor] loaded (listenerAdded=${this.listener != null}, query=${JSON.stringify(risonDecode(serializedQuery))})`)
      if (this.listener == null) {
        this.notConsumedData = {configuration, data}
      }
      else {
        this.listener(data, configuration)
      }
    })
  }
}

export function initDataComponent(persistentStateManager: PersistentStateManager, dataQueryExecutor: DataQueryExecutor): void {
  provide(dataQueryExecutorKey, dataQueryExecutor)
  persistentStateManager.init()
  dataQueryExecutor.init()
}