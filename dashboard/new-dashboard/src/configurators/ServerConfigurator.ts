import { combineLatest, map, Observable, shareReplay } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, serializeQuery } from "../components/common/dataQuery"
import { getCompressor, initZstdObservable } from "../components/common/zstd"
import { injectOrError, serverUrlObservableKey } from "../shared/injectionKeys"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  private readonly observable: Observable<null>
  private _serverUrl: string = ServerConfigurator.DEFAULT_SERVER_URL

  constructor(
    readonly db: string,
    readonly table: string,
    serverUrlObservable: Observable<string> | null = null
  ) {
    if (serverUrlObservable == null) {
      serverUrlObservable = injectOrError(serverUrlObservableKey)
    }
    this.observable = combineLatest([serverUrlObservable, initZstdObservable]).pipe(
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
    query.db = this.db
    query.table = this.table
    return true
  }
}
