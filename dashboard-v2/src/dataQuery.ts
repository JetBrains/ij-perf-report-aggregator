import { encode as risonEncode } from "rison-node"
import { Ref, ref } from "vue"
import { PersistentStateManager } from "./PersistentStateManager"
import { LineSeriesOption } from "echarts"
import { DimensionDefinition } from "echarts/types/src/util/types"

export function encodeQuery(query: DataQuery): string {
  // console.debug(`data-query: ${JSON.stringify(query, null, 2)}`)
  return risonEncode(query)
}

export class DataQueryExecutorConfiguration {
  serverUrl?: string

  readonly series: Array<LineSeriesOption> = []
  readonly dimensions: Array<DimensionDefinition> = []

  addSeries(series: LineSeriesOption, dimension: DimensionDefinition): void {
    this.series.push(series)
    this.dimensions.push(dimension)
  }
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
  readonly dimensions: Array<DataQueryDimension> = []

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

  addFilter(filter: DataQueryFilter): void {
    let filters = this.filters
    if (filters == null) {
      filters = []
      this.filters = filters
    }
    filters.push(filter)
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
  sql?: string
  resultKey?: string
}

export interface DataQueryConfigurator {
  readonly value?: Ref<unknown>

  scheduleLoad?(immediately: boolean): void

  configure(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean
}

export class ServerConfigurator implements DataQueryConfigurator {
  public readonly server = ref("https://ij-perf.labs.jb.gg")

  constructor(public readonly databaseName: string, persistentStateManager: PersistentStateManager) {
    persistentStateManager.add("serverUrl", this.server)
  }

  configure(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    configuration.serverUrl = this.server.value
    query.db = this.databaseName
    return true
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