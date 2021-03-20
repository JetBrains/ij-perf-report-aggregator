import { encode as risonEncode } from "rison-node"
import { Ref } from "vue"
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

  addFilter(filter: DataQueryFilter): void {
    let filters = this.filters
    if (filters == null) {
      filters = []
      this.filters = filters
    }
    filters.push(filter)
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
  field: string
  value?: number | string | Array<string>
  sql?: string
  // `=` by default
  operator?: "=" | "!=" | ">"
}

export interface DataQueryDimension {
  name: string
  // for nested
  subName?: string
  sql?: string
  resultKey?: string
}

export interface DataQueryConfigurator {
  readonly value?: Ref<unknown>
  readonly valueChangeDelay?: number

  scheduleLoadMetadata?(immediately: boolean): void

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean
}

interface ExtraQueryProducer {
  /**
   * Mutate query and return false to stop producing.
   */
  mutate(): boolean

  getSeriesName(index: number): string

  getMeasureName(index: number): string
}

export class DataQueryExecutorConfiguration {
  serverUrl: string | null = null

  getServerUrl(): string {
    const result = this.serverUrl
    if (result == null) {
      throw new Error("serverUrl is not configured")
    }
    return result
  }

  // returns false if done
  extraQueryProducer: ExtraQueryProducer | null = null

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