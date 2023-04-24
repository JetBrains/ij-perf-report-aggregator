import { Observable } from "rxjs"
import { DataQueryExecutor } from "shared/src/DataQueryExecutor"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "shared/src/dataQuery"
import { shallowRef } from "vue"
import { InfoData, InfoSidebarVm } from "../InfoSidebarVm"

function stringToHex(string: string): string {
  return [...string].map(it => it.codePointAt(0)?.toString(16).slice(-4)).join("")
}

export function showSideBar(sidebarVm: InfoSidebarVm | undefined, infoData: InfoData) {
  if (infoData.installerId == undefined) {
    sidebarVm?.show(infoData)
    return
  }
  const serverUrlObservable = refToObservable(shallowRef(ServerConfigurator.DEFAULT_SERVER_URL))
  const separator = ".."
  //todo make db configurable
  new DataQueryExecutor([new ServerConfigurator("perfint", "installer", serverUrlObservable), new class implements DataQueryConfigurator {
    configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
      configuration.queryProducers.push(new SimpleQueryProducer())
      query.addField({n: "changes", sql: `concat(toString(arrayElement(changes, -1)),'${separator}',toString(arrayElement(changes, 1)))`})
      query.addFilter({f: "id", v: infoData.installerId})
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
      infoData.changes = changes.split(separator).map(it => stringToHex(atob(it)).slice(0, 7)).join(" .. ")
    }
    sidebarVm?.show(infoData)
  })
}