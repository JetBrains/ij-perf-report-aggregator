import { encode as risonEncode } from "rison-node"
import { Observable } from "rxjs"
import { ChartConfigurator } from "./chart"

export function encodeQuery(query: DataQuery): string {
  // console.debug(`data-query: ${JSON.stringify(query, null, 2)}`)
  return risonEncode(query)
}

export class DataQuery {
  db?: string
  table?: string

  // noinspection JSMismatchedCollectionQueryUpdate
  private readonly fields: Array<string | DataQueryDimension> = []
  private filters?: Array<DataQueryFilter>

  order?: Array<string>

  // used only for grouped query
  aggregator?: string
  // cube.js term (group by)
  private dimensions?: Array<DataQueryDimension>

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

  removeFilters(toRemove: Array<DataQueryFilter>): void {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    this.filters = this.filters!.filter(it => !toRemove.includes(it))
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
  f: string
  v?: number | string | Array<string>
  q?: string
  // `=` by default
  // operator
  o?: "=" | "!=" | ">"
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

export interface QueryProducer {
  size(): number

  /**
   * Mutate query and return false to stop producing.
   */
  mutate(index: number): void

  getSeriesName(index: number): string

  getMeasureName(index: number): string
}

export class DataQueryExecutorConfiguration {
  serverUrl: string | null = null

  readonly seriesNames: Array<string> = []
  readonly measureNames: Array<string> = []

  getServerUrl(): string {
    const result = this.serverUrl
    if (result == null) {
      throw new Error("serverUrl is not configured")
    }
    return result
  }

  readonly queryProducers: Array<QueryProducer> = []

  private _chartConfigurator: ChartConfigurator | null = null
  get chartConfigurator(): ChartConfigurator {
    const result = this._chartConfigurator
    if (result == null) {
      throw new Error("measure list is not yet set")
    }
    return result
  }

  set chartConfigurator(value: ChartConfigurator) {
    if (this._chartConfigurator != null) {
      throw new Error("measure list is already set")
    }
    this._chartConfigurator = value
  }

  private _measures: Array<string> | null = null
  get measures(): Array<string> {
    const result = this._measures
    if (result == null) {
      throw new Error("measure list is not yet set")
    }
    return result
  }
  set measures(value: Array<string>) {
    if (this._measures != null) {
      throw new Error("measure list is already set")
    }
    this._measures = value
  }
}

export interface Machine {
  readonly name: string
}

export function toMutableArray(value: string | Array<string> | null): Array<string> {
  return (value == null || value === "") ? [] : (Array.isArray(value) ? value.slice() : [value])
}

export function toArray(value: string | Array<string> | null): Array<string> {
  return (value == null || value === "") ? [] : (Array.isArray(value) ? value : [value])
}