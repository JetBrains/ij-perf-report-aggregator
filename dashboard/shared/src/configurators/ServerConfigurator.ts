import { Observable, shareReplay } from "rxjs"
import { inject, ref, Ref } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { serverUrlKey } from "../injectionKeys"
import { refToObservable } from "./rxjs"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  readonly value: Ref<string>

  private readonly observable: Observable<string>

  constructor(readonly db: string, readonly table: string | null = null, serverUrl: string | null = null) {
    if (serverUrl == null) {
      const value = inject(serverUrlKey, ref(ServerConfigurator.DEFAULT_SERVER_URL))
      if (value === undefined) {
        throw new Error("Server URL is not provided")
      }
      this.value = value
    }
    else {
      this.value = ref(serverUrl)
    }

    this.observable = refToObservable(this.value).pipe(
      shareReplay(1),
    )
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const serverUrl = this.value.value ?? ServerConfigurator.DEFAULT_SERVER_URL
    // noinspection HttpUrlsUsage
    if (serverUrl.length === 0 || !(serverUrl.startsWith("http://") || serverUrl.startsWith("https://"))) {
      console.debug(`[serverConfigurator] server url is not correct (url=${serverUrl})`)
      return false
    }

    configuration.serverUrl = serverUrl
    query.db = this.db
    if (this.table != null) {
      query.table = this.table
    }
    return true
  }
}