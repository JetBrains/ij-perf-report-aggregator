import { inject, Ref, ref } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { serverUrlKey } from "../componentKeys"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"

export class ServerConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_SERVER_URL = "https://ij-perf.labs.jb.gg"

  readonly value: Ref<string>
  readonly valueChangeDelay = 900

  constructor(readonly databaseName: string, persistentStateManager: PersistentStateManager) {
    this.value = inject(serverUrlKey, ref(""))
    persistentStateManager.add("serverUrl", this.value)
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