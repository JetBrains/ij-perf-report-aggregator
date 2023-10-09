import { combineLatest, concat, debounceTime, filter, forkJoin, map, Observable, of, shareReplay, switchMap } from "rxjs"
import { measureNameToLabel } from "../../configurators/MeasureConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { defaultBodyConsumer, fromFetchWithRetryAndErrorHandling } from "../../configurators/rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, DataQueryFilter } from "./dataQuery"

export declare type DataQueryResult = (string | number)[][][]
export declare type DataQueryConsumer = (data: DataQueryResult | null, configuration: DataQueryExecutorConfiguration, isLoading: boolean) => void

interface Result {
  readonly isLoading: boolean
  readonly data: DataQueryResult | null
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
  constructor(configurators: DataQueryConfigurator[]) {
    const serverConfigurator = configurators.find((it) => it instanceof ServerConfigurator) as ServerConfigurator
    let abortController = new AbortController()
    this.observable = combineLatest(
      configurators.map((configurator) => {
        // combineLatest will not emit an initial value until each observable emits at least one value, so, null observer simply emits one null value
        return (configurator.createObservable() ?? of(null)).pipe(map(() => configurator))
      })
    ).pipe(
      debounceTime(100),
      switchMap((configurators) => {
        abortController.abort()
        const configuration = new DataQueryExecutorConfiguration()
        const query = new DataQuery()

        for (const configurator of configurators) {
          if (!configurator.configureQuery(query, configuration)) {
            return of(null)
          }
        }

        const queries = generateQueries(query, configuration)
        const mergedQueries = mergeQueries(queries)

        const loadingResults = of({ query, configuration, data: null, isLoading: true })
        abortController = new AbortController()
        return concat(
          loadingResults,
          forkJoin(
            mergedQueries.map((it) => {
              const stringQuery = JSON.stringify(it)
              return fromFetchWithRetryAndErrorHandling<DataQueryResult>(serverConfigurator.computeSerializedQueryUrl(`[${stringQuery}]`), defaultBodyConsumer, abortController)
            })
          ).pipe(
            // pass context along with data and flatten result
            map((data): Result => ({ query, configuration, data: data.flat(1), isLoading: false }))
          )
        )
      }),
      filter((it: Result | null): it is Result => it !== null),
      shareReplay(1)
    )
  }

  subscribe(listener: DataQueryConsumer): () => void {
    const subscription = this.observable.subscribe(({ configuration, query, data, isLoading }) => {
      this._lastQuery = query
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      // console.debug(`[queryExecutor] loaded (listenerAdded=${this.listener != null}, query=${JSON.stringify(risonDecode(serializedQuery))})`)
      listener(data, configuration, isLoading)
    })
    return () => {
      subscription.unsubscribe()
    }
  }
}

// https://stackoverflow.com/a/43053803
function computeCartesian<T>(input: T[][]): T[][] {
  // eslint-disable-next-line @typescript-eslint/ban-ts-comment
  // @ts-expect-error
  // eslint-disable-next-line unicorn/no-array-reduce
  return input.reduce((a, b) => {
    return a.flatMap((d) => b.map((e) => [d, e].flat()))
  })
}

export function generateQueries(query: DataQuery, configuration: DataQueryExecutorConfiguration): DataQuery[] {
  let producers = configuration.queryProducers
  if (producers.length === 0) {
    producers = [
      {
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
      },
    ]
  }

  const cartesian: number[][] =
    producers.length === 1
      ? Array.from({ length: producers[0].size() }).map((_, i) => [i])
      : computeCartesian(
          producers.map((it) => {
            return Array.from({ length: it.size() }).map((_, i) => i)
          })
        )

  let serializedQuery = []

  // https://en.wikipedia.org/wiki/Cartesian_product

  const result: DataQuery[] = []
  const last = Array.from({ length: producers.length - 1 }).map((_, i) => cartesian[0][i])
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

    serializedQuery.push(JSON.parse(JSON.stringify(query)) as DataQuery)

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
        result.push(...serializedQuery)
        serializedQuery = []
      }
    }

    configuration.seriesNames.push(seriesName)
    configuration.measureNames.push(measureName)
  }

  if (serializedQuery.length > 0) {
    result.push(...serializedQuery)
  }
  return result
}

/**
 * Returns the name of the filter that can be combined between the two queries, or null if no such filter exists.
 * If there are multiple filters that are different, returns null.
 */
function getFilterNameForMerge(query1: DataQuery, query2: DataQuery): string | null {
  if (query1.filters == undefined || query2.filters == undefined) return null
  const differingFilters = query1.filters.filter((filter1) => !query2.filters?.some((filter2) => deepEqual(filter1, filter2)))
  if (differingFilters.length !== 1) return null
  const targetFilter1 = differingFilters[0]
  const targetFilter2 = query2.filters.find((filter2) => filter2.f === targetFilter1.f)
  if (!targetFilter2 || !isFilterCanBeMerged(targetFilter1, targetFilter2)) return null
  return targetFilter1.f
}

function deepEqual(obj1: DataQueryFilter, obj2: DataQueryFilter): boolean {
  return JSON.stringify(obj1) === JSON.stringify(obj2)
}

function isFilterCanBeMerged(filter1: DataQueryFilter, filter2: DataQueryFilter): boolean {
  // Check if both filters have a field 'f' and they are equal
  if (filter1.f !== filter2.f) return false

  // Check if filter1.v and filter2.v are either strings or arrays of strings
  if (!(typeof filter1.v === "string" ? true : Array.isArray(filter1.v)) || !(typeof filter2.v === "string" ? true : Array.isArray(filter2.v))) return false

  //Check that filters are different
  if (filter1.v == filter2.v) return false

  //We only support combining filters with no operator
  // noinspection RedundantIfStatementJS
  if (filter1.o !== undefined || filter2.o !== undefined) return false
  return true
}

export function mergeQueries(queries: DataQuery[]): DataQuery[] {
  if (queries.length === 1) return queries
  const resultQueries: DataQuery[] = [...queries]

  let currentFilterField: string | null = null
  for (let i = 0; i < resultQueries.length; i++) {
    for (let j = i + 1; j < resultQueries.length; j++) {
      const matchingFilterField = getFilterNameForMerge(resultQueries[i], resultQueries[j])
      if ((currentFilterField == null && matchingFilterField != null) || (matchingFilterField === currentFilterField && matchingFilterField != null)) {
        currentFilterField = matchingFilterField
        const mergedQuery = { ...resultQueries[i] } as DataQuery
        mergedQuery.filters = mergeFilters(resultQueries[i].filters, resultQueries[j].filters)
        resultQueries.push(mergedQuery)

        // Remove the original queries
        resultQueries.splice(j, 1) // remove j first since it's higher index
        resultQueries.splice(i, 1)

        // Reset indices to re-evaluate with new list
        i = -1
        break // exit the inner loop
      }
    }
  }
  return resultQueries
}

function mergeFilters(filters1?: DataQueryFilter[], filters2?: DataQueryFilter[]): DataQueryFilter[] {
  if (!filters1 || !filters2) return filters1 ?? filters2 ?? []
  return filters1.map((filter1) => {
    const filter2 = filters2.find((f2) => isFilterCanBeMerged(filter1, f2))
    //@ts-expect-error - filter1 and filter2 are strings or arrays of strings this is checked in isFilterCanBeMerged
    return filter2 ? { f: filter1.f, v: mergeValues(filter1.v, filter2.v), s: true } : filter1
  })
}

function mergeValues(value1: string | string[], value2: string | string[]): string[] {
  const array1 = Array.isArray(value1) ? value1 : [value1]
  const array2 = Array.isArray(value2) ? value2 : [value2]
  return [...new Set([...array1, ...array2])]
}
