import { inject, provide } from "vue"
import { DataQuery, DataQueryFilter, encodeQuery } from "./dataQuery"
import { serverUrlKey, reportInfoProviderKey } from "./injectionKeys"

export function provideReportUrlProvider(): void {
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  const serverUrl = inject(serverUrlKey)!
  provide(reportInfoProviderKey, {
    infoFields: ["tc_installer_build_id"],
    createReportUrl: (generatedTime, query) => createReportUrl(generatedTime, query, serverUrl.value)
  })
}

function createReportUrl(generatedTime: number, query: DataQuery, serverUrl: string): string {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const q: Record<string, unknown> = {
    ...query
  }
  delete q["fields"]
  delete q["order"]
  delete q["flat"]
  q["filters"] = (q["filters"] as Array<DataQueryFilter>).filter(it => !it.field.startsWith("measures."))
  const filters = q["filters"] as Array<DataQueryFilter>
  for (let i = 0; i < filters.length; i++){
    if (filters[i].field === "generated_time") {
      filters[i] = {field: "generated_time", value: generatedTime / 1000}
      break
    }
  }
  return `/report?reportUrl=${encodeURIComponent(serverUrl)}/api/v1/report/${encodeQuery(q as never)}`
}
