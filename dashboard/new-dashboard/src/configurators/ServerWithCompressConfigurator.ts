import { combineLatest, map, Observable, shareReplay } from "rxjs"
import { DataQuery, DataQueryExecutorConfiguration, serializeQuery, ServerConfigurator } from "../components/common/dataQuery"
import { getCompressor, getZstdObservable } from "../components/common/zstd"
import { dbTypeStore } from "../shared/dbTypes"
import { injectOrError, serverUrlObservableKey } from "../shared/injectionKeys"

export class ServerWithCompressConfigurator implements ServerConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"
  static readonly INTRANET_SERVER_URL = "https://ij-perf-api.labs.jb.gg"

  private readonly observable: Observable<null>
  private _serverUrl: string = this.determineServerUrl()

  constructor(
    readonly db: string,
    readonly table: string,
    serverUrlObservable: Observable<string> | null = null
  ) {
    dbTypeStore().setDbType(db, table)
    serverUrlObservable ??= injectOrError(serverUrlObservableKey)
    this.observable = combineLatest([serverUrlObservable, getZstdObservable()]).pipe(
      map(([url, _]) => {
        this._serverUrl = url
        return null
      }),
      shareReplay(1)
    )
  }

  get serverUrl(): string {
    return this._serverUrl
  }

  private determineServerUrl(): string {
    console.log(window.location.hostname)
    if (window.location.hostname === "ij-perf-api.labs.jb.gg") {
      return ServerWithCompressConfigurator.INTRANET_SERVER_URL
    } else {
      return ServerWithCompressConfigurator.DEFAULT_SERVER_URL
    }
  }

  compressString(params: string): string {
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-template-expression
    return `${getCompressor().compress(params)}`
  }

  computeQueryUrl(query: DataQuery): string {
    return `${this._serverUrl}/api/q/${this.compressString(serializeQuery(query))}`
  }

  computeSerializedQueryUrl(url: string): string {
    return `${this._serverUrl}/api/q/${this.compressString(url)}`
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.db = this.db
    query.table = this.table
    return true
  }

  configureFilter(_: DataQuery): boolean {
    return true
  }
}
