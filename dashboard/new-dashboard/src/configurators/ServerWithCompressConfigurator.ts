import { map } from "rxjs/map"
import { shareReplay } from "rxjs/share-replay"
import { combineLatest } from "./rxjs"
import { DataQuery, DataQueryExecutorConfiguration, serializeQuery, ServerConfigurator } from "../components/common/dataQuery"
import { getCompressor, getZstdObservable } from "../components/common/zstd"
import { dbTypeStore } from "../shared/dbTypes"
import { injectOrError, injectOrNull, serverUrlKey, serverUrlObservableKey } from "../shared/injectionKeys"
import type { Ref } from "vue"

export class ServerWithCompressConfigurator implements ServerConfigurator {
  static getDefaultServerUrl(): string {
    const hostname = window.location.hostname
    if (hostname === "ij-perf-api.labs.jb.gg" || hostname === "localhost") {
      return "https://ij-perf-api.labs.jb.gg"
    }
    return "https://ij-perf.labs.jb.gg"
  }

  static readonly DEFAULT_SERVER_URL = ServerWithCompressConfigurator.getDefaultServerUrl()

  private readonly observable: Observable<null>
  private _serverUrl: string = ServerWithCompressConfigurator.DEFAULT_SERVER_URL
  // the setting itself: _serverUrl is updated only while the observable is subscribed, which a page without data queries never does
  private readonly serverUrlSetting: Ref<string> | null = null

  constructor(
    readonly db: string,
    readonly table: string,
    serverUrlObservable: Observable<string> | null = null
  ) {
    dbTypeStore().setDbType(db, table)
    if (serverUrlObservable == null) {
      serverUrlObservable = injectOrError(serverUrlObservableKey)
      this.serverUrlSetting = injectOrNull(serverUrlKey)
    }
    this.observable = combineLatest([serverUrlObservable, getZstdObservable()])
      [map](([url, _]) => {
        this._serverUrl = url
        return null
      })
      [shareReplay](1)
  }

  get serverUrl(): string {
    const setting = this.serverUrlSetting?.value
    return setting == null || setting === "" ? this._serverUrl : setting
  }

  compressString(params: string): string {
    // eslint-disable-next-line @typescript-eslint/no-unnecessary-template-expression
    return `${getCompressor().compress(params)}`
  }

  computeQueryUrl(query: DataQuery): string {
    return `${this.serverUrl}/api/q/${this.compressString(serializeQuery(query))}`
  }

  computeSerializedQueryUrl(url: string): string {
    return `${this.serverUrl}/api/q/${this.compressString(url)}`
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
