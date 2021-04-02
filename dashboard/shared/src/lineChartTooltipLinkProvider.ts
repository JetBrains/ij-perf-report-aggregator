import { provide } from "vue"
import { DataQuery, DataQueryFilter, encodeQuery } from "./dataQuery"
import { reportInfoProviderKey } from "./injectionKeys"

export function provideReportUrlProvider(): void {
  provide(reportInfoProviderKey, {
    infoFields: ["tc_installer_build_id"],
    createReportUrl,
  })
}

function createReportUrl(generatedTime: number, query: DataQuery): string {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const q: Record<string, unknown> = {
    ...query
  }
  delete q["fields"]
  delete q["order"]
  delete q["flat"]
  q["filters"] = (q["filters"] as Array<DataQueryFilter>).filter(it => it.field === "product" || it.field === "project" || it.field === "machine" || it.field === "generated_time")
  const filters = q["filters"] as Array<DataQueryFilter>
  for (let i = 0; i < filters.length; i++){
    if (filters[i].field === "generated_time") {
      filters[i] = {field: "generated_time", value: generatedTime / 1000}
      break
    }
  }
  return `/report?reportUrl=${encodeURIComponent(`/api/v1/report/${encodeQuery(q as never)}`)}`
}
