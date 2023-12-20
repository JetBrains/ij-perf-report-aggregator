import { Observable } from "rxjs"
import { ChartConfigurator } from "./chart"

export function serializeQuery(query: DataQuery): string {
  // const encoded = encodeRison(query)
  return JSON.stringify(query)
}

export class DataQuery {
  db?: string
  table?: string

  // noinspection JSMismatchedCollectionQueryUpdate
  public fields: (string | DataQueryDimension)[] = []
  public filters?: DataQueryFilter[]

  order?: string[] | string

  // used only for grouped query
  aggregator?: string
  // cube.js term (group by)
  private dimensions?: DataQueryDimension[]

  // used only for grouped query
  timeDimensionFormat?: string

  /**
   * Whether to return flat array . Applicable only if the only field is specified.
   * `false`: array of objects ([{field: value},...])
   * `true`: array of values ([value,...])
   */
  flat?: boolean

  addField(field: string | DataQueryDimension): void {
    this.fields.push(field)
  }

  insertField(field: string | DataQueryDimension, index: number): void {
    this.fields.splice(index, 0, field)
  }

  // removeField(field: DataQueryDimension) {
  //   const index = this.fields.indexOf(field)
  //   if (index > -1) {
  //     this.fields.splice(index, 1)
  //   }
  // }

  addFilter(filter: DataQueryFilter): void {
    let filters = this.filters
    if (filters == null) {
      filters = []
      this.filters = filters
    }
    filters.push(filter)
  }

  removeFilters(toRemove: DataQueryFilter[]): void {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    this.filters = this.filters!.filter((it) => !toRemove.includes(it))
  }

  addDimension(dimension: DataQueryDimension): void {
    let dimensions = this.dimensions
    if (dimensions == null) {
      dimensions = []
      this.dimensions = dimensions
    }
    dimensions.push(dimension)
  }
}

export interface DataQueryFilter {
  f?: string
  v?: number | string | string[] | boolean
  q?: string
  // `=` by default
  // operator
  o?: "=" | "!=" | ">" | "like"
  // query is combined and should be split
  s?: boolean
}

export interface DataQueryDimension {
  n: string
  // for nested
  subName?: string
  sql?: string
  resultKey?: string
}

export interface DataQueryConfigurator {
  /**
   * Return null if no need to observe — a fully static configurator that doesn't trigger query reconfiguration.
   * Observable must emit at least one value.
   *
   * Do not emit duplicated values - use distinctUntilChanged operator if needed.
   */
  createObservable(): Observable<unknown> | null

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean

  createDimension?(): QueryProducer
}

export interface ServerConfigurator extends DataQueryConfigurator {
  get serverUrl(): string

  compressString(params: string): string

  computeQueryUrl(query: DataQuery): string

  computeSerializedQueryUrl(url: string): string

  createObservable(): Observable<unknown>

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean

  createDimension?(): QueryProducer
}

export interface QueryProducer {
  size(): number

  /**
   * Mutate query and return false to stop producing.
   */
  mutate(index: number): void

  getSeriesName(index: number): string

  getMeasureName(index: number): string
}

export class SimpleQueryProducer implements QueryProducer {
  getMeasureName(_: number): string {
    return ""
  }

  getSeriesName(_: number): string {
    return ""
  }

  // eslint-disable-next-line @typescript-eslint/no-empty-function
  mutate(_: number): void {}

  size(): number {
    return 1
  }
}

export class DataQueryExecutorConfiguration {
  public seriesNames: string[] = []
  readonly measureNames: string[] = []

  readonly queryProducers: QueryProducer[] = []

  private chartConfigurator: ChartConfigurator[] = []

  public getChartConfigurators(): ChartConfigurator[] {
    return this.chartConfigurator
  }

  public addChartConfigurator(configurator: ChartConfigurator): void {
    this.chartConfigurator.push(configurator)
  }

  private _measures: string[] | null = null
  get measures(): string[] {
    const result = this._measures
    if (result == null) {
      throw new Error("measure list is not yet set")
    }
    return result
  }

  set measures(value: string[]) {
    if (this._measures != null) {
      throw new Error("measure list is already set")
    }
    this._measures = value
  }
}

export interface Machine {
  readonly name: string
}

export function toMutableArray(value: string | string[] | null): string[] {
  return value == null || value === "" ? [] : Array.isArray(value) ? [...value] : [value]
}

export function toArray(value: string | string[] | null): string[] {
  return value == null || value === "" ? [] : Array.isArray(value) ? value : [value]
}
