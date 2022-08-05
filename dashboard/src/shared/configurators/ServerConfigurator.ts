import { combineLatest, map, Observable, shareReplay } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, serializeQuery } from "../dataQuery"
import { injectOrError, serverUrlObservableKey } from "../injectionKeys"
import { getCompressor, initZstdObservable } from "../zstd"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  private readonly observable: Observable<null>
  private _serverUrl: string = ServerConfigurator.DEFAULT_SERVER_URL

  constructor(readonly db: string, readonly table: string | null = null) {
    this.observable = combineLatest([injectOrError(serverUrlObservableKey), initZstdObservable]).pipe(
      map(([url, _]) => {
        this._serverUrl = url
        return null
      }),
      shareReplay(1),
    )
  }

  get serverUrl(): string {
    return this._serverUrl
  }

  computeQueryUrl(query: DataQuery): string {
    return `${this._serverUrl}/api/q/${getCompressor().compress(serializeQuery(query))}`
  }

  computeSerializedQueryUrl(url: string): string {
    return `${this._serverUrl}/api/q/${getCompressor().compress(url)}`
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    if (this._serverUrl == null || this._serverUrl.length === 0) {
      return false
    }

    query.db = this.db
    if (this.table != null) {
      query.table = this.table
    }
    return true
  }
}