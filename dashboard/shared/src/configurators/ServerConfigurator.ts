import { forkJoin, map, Observable, shareReplay, switchMap } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, serializeQuery } from "../dataQuery"
import { injectOrError, serverUrlObservableKey } from "../injectionKeys"
import { CompressorUsingDictionary, initZstdObservable } from "../zstd"
import { fromFetchWithRetryAndErrorHandling } from "./rxjs"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  private readonly observable: Observable<null>
  private _serverUrl: string = ServerConfigurator.DEFAULT_SERVER_URL

  private compressor: CompressorUsingDictionary | null = null

  constructor(readonly db: string, readonly table: string | null = null) {
    this.observable = injectOrError(serverUrlObservableKey).pipe(
      switchMap(url => {
        this._serverUrl = url
        return forkJoin([
          fromFetchWithRetryAndErrorHandling<ArrayBuffer>(`${url}/api/zstd-dictionary`, null, it => it.arrayBuffer()),
          initZstdObservable,
        ])
      }),
      map(([dictionaryData, _]) => {
        if (this.compressor !== null) {
          this.compressor.dispose()
        }
        this.compressor = new CompressorUsingDictionary(dictionaryData)
        return null
      }),
      shareReplay(1)
    )
  }

  get serverUrl(): string {
    return this._serverUrl
  }

  computeQueryUrl(query: DataQuery): string {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return `${this._serverUrl}/api/q/${this.compressor!.compress(serializeQuery(query))}`
  }

  computeSerializedQueryUrl(url: string): string {
    // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
    return `${this._serverUrl}/api/q/${this.compressor!.compress(url)}`
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