import { Observable } from "rxjs"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "shared/src/dataQuery"
import { shallowRef } from "vue"
import { InfoData, InfoSidebarVm } from "../InfoSidebarVm"

function base64ToHex(base64: string): string {
  const decodedArray = new Uint8Array([...atob(base64)].map(c => c.codePointAt(0) ?? 0))
  let hex = ""
  for (const byte of decodedArray) {
    hex += byte.toString(16).padStart(2, "0")
  }
  return hex
}

export function showSideBar(sidebarVm: InfoSidebarVm | undefined, infoData: InfoData) {
  const serverUrlObservable = refToObservable(shallowRef(ServerConfigurator.DEFAULT_SERVER_URL))
  const separator = ".."
  const db = infoData.installerId ? "perfint" : "perfintDev"
  const id = infoData.installerId ?? infoData.buildId
  new DataQueryExecutor([new ServerConfigurator(db, "installer", serverUrlObservable), new class implements DataQueryConfigurator {
    configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
      configuration.queryProducers.push(new SimpleQueryProducer())
      query.addField({n: "changes", sql: `arrayStringConcat(changes,'${separator}')`})
      query.addFilter({f: "id", v: id})
      query.order = "changes"
      return true
    }

    createObservable(): Observable<unknown> | null {
      return null
    }
  }]).subscribe((data, _configuration,isLoading) => {
    if(isLoading || data == null){
      return
    }
    const changes = data.flat(3)[0]
    if(typeof changes === "string"){
      //commit has to be decoded as base64 and converted to hex
      infoData.changes = changes.split(separator).map(it => base64ToHex(it)).join("%2C")
    }
    sidebarVm?.show(infoData)
  })
}