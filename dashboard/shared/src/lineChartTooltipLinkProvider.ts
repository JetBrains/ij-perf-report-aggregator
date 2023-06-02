import { Observable } from "rxjs"
import { provide } from "vue"
import { DataQuery, DataQueryFilter, serializeQuery } from "./dataQuery"
import { injectOrError, reportInfoProviderKey, serverUrlObservableKey } from "./injectionKeys"

export function provideReportUrlProvider(isInstallerExists: boolean = true, isBuildNumberExists: boolean = false): void {
  const serverUrl = injectOrError(serverUrlObservableKey)
  const infoFields = ["machine", "tc_build_id", "project"]
  if (isInstallerExists) {
    infoFields.push("tc_installer_build_id", "build_c1", "build_c2", "build_c3")
  }
  if(isBuildNumberExists){
    infoFields.push("build_number")
  }
  provide(reportInfoProviderKey, {
    infoFields,
    createReportUrl: (generatedTime, query) => createReportUrl(generatedTime, query, serverUrl),
  })
}

function createReportUrl(generatedTime: number, query: DataQuery, serverUrl: Observable<string>): string {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const q: Record<string, unknown> = {
    ...query
  }
  delete q["fields"]
  delete q["order"]
  delete q["flat"]
  q["filters"] = (q["filters"] as DataQueryFilter[]).filter(it => it.f === "product" || it.f === "project" || it.f === "machine" || it.f === "generated_time")
  const filters = q["filters"] as DataQueryFilter[]
  for (let i = 0; i < filters.length; i++){
    if (filters[i].f === "generated_time") {
      filters[i] = {f: "generated_time", v: generatedTime / 1000}
      break
    }
  }

  let v
  serverUrl.subscribe(value => {
    v = value
  }).unsubscribe()
  return `/report?reportUrl=${encodeURIComponent(`${v}/api/v1/report/${serializeQuery(q as never)}`)}`
}
