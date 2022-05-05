import { combineLatest, map, Observable, shareReplay } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, serializeQuery } from "../dataQuery"
import { compressorObservableKey, injectOrError, serverUrlObservableKey } from "../injectionKeys"
import { CompressorUsingDictionary } from "../zstd"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  private readonly observable: Observable<null>
  private _serverUrl: string = ServerConfigurator.DEFAULT_SERVER_URL

  private compressor: CompressorUsingDictionary | null = null

  constructor(readonly db: string, readonly table: string | null = null) {
    injectOrError(compressorObservableKey).subscribe(it => {
      this.compressor = it
    })

    this.observable = combineLatest([injectOrError(serverUrlObservableKey), injectOrError(compressorObservableKey)]).pipe(
      map(([url, compressor]) => {
        this._serverUrl = url
        this.compressor = compressor
        return null
      }),
      shareReplay(1),
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