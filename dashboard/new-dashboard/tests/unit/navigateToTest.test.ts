import { createPinia, setActivePinia } from "pinia"
import { beforeEach, describe, expect, it } from "vitest"
import { Router } from "vue-router"
import { valueUnitParamName } from "../../src/components/common/chart"
import { getNavigateToTestUrl, InfoData } from "../../src/components/common/sideBar/InfoSidebar"
import { dbTypeStore } from "../../src/shared/dbTypes"

function router(path: string, query: Record<string, string> = {}): Router {
  return {
    currentRoute: { value: { path, query } },
    resolve: (to: string) => ({ href: to }),
  } as unknown as Router
}

function sidebarData(data: Partial<InfoData>): InfoData {
  return {
    projectName: "intellij_commit/indexing",
    machineName: "intellij-linux-hw-hetzner-agent-17",
    branch: "master",
    buildId: 458553579,
    series: [{ metricName: "JVM.heapUsageMb/afterExecution", value: "216", color: "#000", rawValue: 216 }],
    ...data,
  } as InfoData
}

describe("Navigate to test URL", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    dbTypeStore().setDbType("perfintDev", "phpstorm")
  })

  it("carries the unit of the clicked chart", () => {
    const url = getNavigateToTestUrl(sidebarData({ metricType: "counter" }), router("/phpstorm/memoryDashboard"))
    expect(url).toContain(`${valueUnitParamName}=counter`)
  })

  it("omits the unit when the chart left it to the stored metric type", () => {
    const url = getNavigateToTestUrl(sidebarData({ metricType: "auto" }), router("/phpstorm/memoryDashboard"))
    expect(url).not.toContain(valueUnitParamName)
  })

  it("drops a unit inherited from the source page query", () => {
    const url = getNavigateToTestUrl(sidebarData({ metricType: undefined }), router("/phpstorm/memoryDashboard", { [valueUnitParamName]: "counter" }))
    expect(url).not.toContain(valueUnitParamName)
  })
})
