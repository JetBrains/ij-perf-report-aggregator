import { Observable } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, serializeAndEncodeQueryForUrl } from "../dataQuery"
import { injectOrError, serverUrlObservableKey } from "../injectionKeys"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  private readonly observable: Observable<string>
  private _serverUrl: string = ServerConfigurator.DEFAULT_SERVER_URL

  constructor(readonly db: string, readonly table: string | null = null) {
    this.observable = injectOrError(serverUrlObservableKey)
    this.observable.subscribe(value => {
      this._serverUrl = value
    })
  }

  get serverUrl(): string {
    return this._serverUrl
  }

  computeQueryUrl(query: DataQuery): string {
    return `${this._serverUrl}/api/q/${serializeAndEncodeQueryForUrl(query)}`
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