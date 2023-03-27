import { combineLatest, concat, debounceTime, filter, forkJoin, map, Observable, of, shareReplay, switchMap } from "rxjs"
import { provide } from "vue"
import { measureNameToLabel } from "./configurators/MeasureConfigurator"
import { ReloadConfigurator } from "./configurators/ReloadConfigurator"
import { ServerConfigurator } from "./configurators/ServerConfigurator"
import { fromFetchWithRetryAndErrorHandling } from "./configurators/rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, serializeQuery } from "./dataQuery"
import { configuratorListKey } from "./injectionKeys"

export declare type DataQueryResult = Array<Array<Array<string | number>>>
export declare type DataQueryConsumer = (data: DataQueryResult|null, configuration: DataQueryExecutorConfiguration, isLoading: boolean) => void

interface Result {
  readonly isLoading: boolean
  readonly data: DataQueryResult|null
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
  constructor(configurators: Array<DataQueryConfigurator>) {
    const serverConfigurator = configurators.find(it => it instanceof ServerConfigurator) as ServerConfigurator
    let abortController = new AbortController()
    this.observable = combineLatest(configurators.map(configurator => {
      // combineLatest will not emit an initial value until each observable emits at least one value, so, null observer simply emits one null value
      return (configurator.createObservable() ?? of(null))
        .pipe(
          map(() => configurator),
        )
    }))
      .pipe(
        debounceTime(100),
        switchMap(configurators => {
          abortController.abort()
          const configuration = new DataQueryExecutorConfiguration()
          const query = new DataQuery()

          for (const configurator of configurators) {
            if (!configurator.configureQuery(query, configuration)) {
              return of(null)
            }
          }

          const queries = generateQueries(query, configuration)
          const loadingResults = of({query, configuration, data: null, isLoading: true})
          abortController = new AbortController()
          return concat(loadingResults, forkJoin(queries.map(it => {
            return fromFetchWithRetryAndErrorHandling<DataQueryResult>(serverConfigurator.computeSerializedQueryUrl(`[${it}]`), null, it => it.json() , abortController)
          })).pipe(
            // pass context along with data and flatten result
            map((data): Result => ({query, configuration, data: data.flat(1), isLoading: false})),
          ))

        }),
        filter((it: Result | null): it is Result => it !== null),
        shareReplay(1),
      )
  }

  subscribe(listener: DataQueryConsumer): () => void {
    const subscription = this.observable.subscribe(({configuration, query, data, isLoading}) => {
      this._lastQuery = query
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      // console.debug(`[queryExecutor] loaded (listenerAdded=${this.listener != null}, query=${JSON.stringify(risonDecode(serializedQuery))})`)
      listener(data, configuration, isLoading)
    })
    return () => subscription.unsubscribe()
  }
}

// https://stackoverflow.com/a/43053803
function computeCartesian<T>(input: Array<Array<T>>): Array<Array<T>> {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-ignore
  // eslint-disable-next-line unicorn/no-array-reduce
  return input.reduce((a, b) => {
    return a.flatMap(d => b.map(e => [d, e].flat()))
  })
}

export function generateQueries(query: DataQuery, configuration: DataQueryExecutorConfiguration): Array<string> {
  let producers = configuration.queryProducers
  if (producers.length === 0) {
    producers = [{
      size(): number {
        return 1
      },
      getMeasureName(_index: number): string {
        return configuration.measures[0]
      },
      getSeriesName(_index: number): string {
        return measureNameToLabel(configuration.measures[0])
      },
      mutate(): void {
        // the only value
      },
    }]
  }

  const cartesian: Array<Array<number>> = producers.length === 1 ? Array.from({length: producers[0].size()}).map((_, i) => [i]) : computeCartesian(producers.map(it => {
    return Array.from({length: it.size()}).map((_, i) => i)
  }))

  let serializedQuery = ""

  // https://en.wikipedia.org/wiki/Cartesian_product

  const result: Array<string> = []
  const last = Array.from({length: producers.length - 1}).map((_, i) => cartesian[0][i])
  for (const item of cartesian) {
    // each column it is a producer
    // each row it is a combination
    // each column value it is an index of producer value
    let seriesName = ""
    let measureName = ""
    for (const [i, index] of item.entries()) {
      const producer = producers[i]
      producer.mutate(index)
      if (i !== 0) {
        measureName += " – "
      }
      measureName += producer.getMeasureName(index)

      const title = producer.getSeriesName(index)
      if (title.length > 0) {
        if (seriesName.length > 0) {
          seriesName += " – "
        }
        seriesName += title
      }
    }

    if (serializedQuery.length > 0) {
      serializedQuery += ","
    }
    serializedQuery += serializeQuery(query)

    if (item.length > 1) {
      let equal = true
      for (let i = item.length - 2; i >= 0; i--) {
        if (last[i] != item[i]) {
          for (let j = 0, n = item.length - 1; j < n; j++) {
            last[j] = item[j]
          }
          equal = false
          break
        }
      }

      if (!equal) {
        result.push(serializedQuery)
        serializedQuery = ""
      }
    }

    configuration.seriesNames.push(seriesName)
    configuration.measureNames.push(measureName)
  }

  if (serializedQuery.length > 0) {
    result.push(serializedQuery)
  }
  return result
}

export function initDataComponent(configurators: Array<DataQueryConfigurator>): void {
  provide(configuratorListKey, [...configurators, new ReloadConfigurator()])
}