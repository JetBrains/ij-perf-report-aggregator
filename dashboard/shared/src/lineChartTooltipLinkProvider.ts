import { provide } from "vue"
import { DataQuery, DataQueryFilter, serializeQuery } from "./dataQuery"
import { reportInfoProviderKey } from "./injectionKeys"

export function provideReportUrlProvider(): void {
  provide(reportInfoProviderKey, {
    infoFields: ["tc_installer_build_id", "tc_build_id", "build_c1", "build_c2", "build_c3", "machine"],
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
  q["filters"] = (q["filters"] as Array<DataQueryFilter>).filter(it => it.f === "product" || it.f === "project" || it.f === "machine" || it.f === "generated_time")
  const filters = q["filters"] as Array<DataQueryFilter>
  for (let i = 0; i < filters.length; i++){
    if (filters[i].f === "generated_time") {
      filters[i] = {f: "generated_time", v: generatedTime / 1000}
      break
    }
  }
  return `/report?reportUrl=${encodeURIComponent(`/api/v1/report/${serializeQuery(q as never)}`)}`
}
