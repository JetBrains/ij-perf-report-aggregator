import { inject, ref, Ref } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { serverUrlKey } from "../injectionKeys"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  readonly value: Ref<string>
  readonly valueChangeDelay = 900

  constructor(readonly databaseName: string) {
    const serverUrl = inject(serverUrlKey, ref(ServerConfigurator.DEFAULT_SERVER_URL))
    if (serverUrl === undefined) {
      throw new Error("Server URL is not provided")
    }
    this.value = serverUrl
  }

  configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const serverUrl = this.value.value ?? ServerConfigurator.DEFAULT_SERVER_URL
    // noinspection HttpUrlsUsage
    if (serverUrl.length === 0 || !(serverUrl.startsWith("http://") || serverUrl.startsWith("https://"))) {
      console.debug(`[serverConfigurator] server url is not correct (url=${serverUrl})`)
      return false
    }

    configuration.serverUrl = serverUrl
    query.db = this.databaseName
    return true
  }
}