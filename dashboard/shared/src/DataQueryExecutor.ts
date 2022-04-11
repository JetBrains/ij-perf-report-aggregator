import { deepEqual } from "fast-equals"
import { combineLatest, debounceTime, distinctUntilChanged, filter, forkJoin, map, mergeMap, Observable, of, shareReplay, switchMap, takeUntil, timer } from "rxjs"
import { fromFetch } from "rxjs/fetch"
import { provide } from "vue"
import { PersistentStateManager } from "./PersistentStateManager"
import { measureNameToLabel } from "./configurators/MeasureConfigurator"
import { ReloadConfigurator } from "./configurators/ReloadConfigurator"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, encodeQuery } from "./dataQuery"
import { configuratorListKey } from "./injectionKeys"

export declare type DataQueryResult = Array<Array<Array<string | number>>>
export declare type DataQueryConsumer = (data: DataQueryResult, configuration: DataQueryExecutorConfiguration) => void

interface Result {
  readonly data: DataQueryResult
  readonly query: DataQuery
  readonly configuration: DataQueryExecutorConfiguration
}

export class DataQueryExecutor {
  private _lastQuery: DataQuery | null = null

  get lastQuery(): DataQuery | null {
    return this._lastQuery
  }

  private readonly observable: Observable<Result>

  /**
   * `isGroup = true` means that this DataQueryExecutor only manages dependent executors but doesn't load data itself.
   */
  constructor(private readonly configurators: Array<DataQueryConfigurator>) {
    this.observable = combineLatest(this.configurators.map(configurator => {
      // combineLatest will not emit an initial value until each observable emits at least one value, so, null observer simply emits one null value
      return (configurator.createObservable() ?? of(null))
        .pipe(
          distinctUntilChanged(deepEqual),
          map(_ => configurator),
        )
    }))
      .pipe(
        debounceTime(100),
        switchMap(configurators => {
          const configuration = new DataQueryExecutorConfiguration()
          const query = new DataQuery()

          for (const configurator of configurators) {
            if (!configurator.configureQuery(query, configuration)) {
              return of(null)
            }
          }

          const queries = generateQueries(query, configuration)
          if (queries.length == 1) {
            return fetch(queries[0], configuration)
              .pipe(
                // pass context along with data
                map((data): Result => {
                  return {query, configuration, data}
                }),
              )
          }
          else {
            return forkJoin(queries.map(it => fetch(it, configuration)))
              .pipe(
                // pass context along with data and flatten result
                map((data): Result => {
                  return {query, configuration, data: data.flat(1)}
                }),
              )
          }
        }),
        filter((it: Result | null): it is Result => it !== null),
        shareReplay(1),
      )
  }

  subscribe(listener: DataQueryConsumer): () => void {
    const subscription = this.observable.subscribe(({configuration, query, data}) => {
      this._lastQuery = query
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      // console.debug(`[queryExecutor] loaded (listenerAdded=${this.listener != null}, query=${JSON.stringify(risonDecode(serializedQuery))})`)
      listener(data, configuration)
    })
    return () => subscription.unsubscribe()
  }

  scheduleLoadIncludingConfigurators(immediately: boolean = false): void {
    for (const configurator of this.configurators) {
      if (configurator.scheduleLoadMetadata != null) {
        configurator.scheduleLoadMetadata(immediately)
      }
    }
  }
}

function fetch(serializedQuery: string, configuration: DataQueryExecutorConfiguration) {
  return fromFetch(`${configuration.getServerUrl()}/api/v1/load/${serializedQuery}`)
    .pipe(
      // promise to result
      mergeMap(response => response.json() as Promise<DataQueryResult>),
      // timeout
      takeUntil(timer(8_000)),
    )
}

// https://stackoverflow.com/a/43053803
function computeCartesian<T>(input: Array<Array<T>>): Array<Array<T>> {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  return input.reduce((a, b) => {
    return a.flatMap(d => b.map(e => [d, e].flat()))
  })
}

export function generateQueries(query: DataQuery, configuration: DataQueryExecutorConfiguration): Array<string> {
  let producers = configuration.extraQueryProducers
  if (producers.length === 0) {
    producers = [{
      size(): number {
        return 1
      },
      getMeasureName(_index: number): string {
        return configuration.measures[0]
      },
      getSeriesName(_index: number): string {
        return measureNameToLabel(this.getMeasureName(-1))
      },
      mutate(): boolean {
        debugger
        throw new Error("Not applicable")
      }
    }]
  }

  let cartesian: Array<Array<number>>
  if (producers.length === 1) {
    cartesian = Array.from(new Array<number>(producers[0].size()), (_, i) => [i])
  }
  else {
    cartesian = computeCartesian(producers.map(it => {
      return Array.from(new Array<number>(it.size()), (_, i) => i)
    }))
  }

  let serializedQuery = ""

  const result: Array<string> = []
  let last = cartesian[0][0]
  for (let combinationIndex = 0; combinationIndex < cartesian.length; combinationIndex++) {
    const item = cartesian[combinationIndex]
    // each column it is a producer
    // each row it is a combination
    // each column value it is an index of producer value
    let seriesName = ""
    let measureName = ""
    for (let i = 0; i < item.length; i++) {
      const producer = producers[i]
      const index = item[i]
      producer.mutate(index)
      if (i !== 0) {
        seriesName += " - "
      }
      if (i !== 0) {
        measureName += " - "
      }
      seriesName += producer.getSeriesName(index)
      measureName += producer.getMeasureName(index)
    }

    if (serializedQuery.length !== 0) {
      serializedQuery += ","
    }
    serializedQuery += encodeQuery(query)

    if (item.length > 1 && last !== item[0]) {
      last = item[0]
      result.push("!(" + serializedQuery + ")")
      serializedQuery = ""
    }

    configuration.seriesNames.push(seriesName)
    configuration.measureNames.push(measureName)
  }

  if (serializedQuery.length !== 0) {
    result.push("!(" + serializedQuery + ")")
  }
  return result
}

export function initDataComponent(persistentStateManager: PersistentStateManager, configurators: Array<DataQueryConfigurator>): void {
  persistentStateManager.init()

  for (const configurator of configurators) {
    if (configurator.scheduleLoadMetadata != null) {
      configurator.scheduleLoadMetadata(true)
    }
  }

  provide(configuratorListKey, configurators.concat(new ReloadConfigurator()))
}