import { computed, Ref, ShallowRef, shallowRef } from "vue"
import { Accident } from "../../../configurators/accidents/AccidentsConfigurator"
import { getTeamcityBuildType } from "../../../util/artifacts"
import { ServerConfigurator } from "../dataQuery"
import { dbTypeStore } from "../../../shared/dbTypes"
import { getMachineGroupName } from "../../../configurators/MachineConfigurator"
import { Router } from "vue-router"
import { calculateChanges } from "../../../util/changes"

export const tcUrl = "https://buildserver.labs.intellij.net/"
export const buildUrl = (id: number) => `${tcUrl}viewLog.html?buildId=${id}`

export interface InfoData {
  seriesName: string
  build: string | undefined
  artifactsUrl: string
  changesUrl: string
  installerUrl: string | undefined
  buildId: number
  machineName: string
  projectName: string
  title: string
  installerId: number | undefined
  date: string
  branch: string | undefined
  series: DataSeries[]
  accidents: Ref<Accident[] | undefined> | undefined
  description: Ref<Description | null>
  deltaPrevious: string | undefined
  deltaNext: string | undefined
  chartDataUrl: string
  buildIdPrevious: number | undefined
  buildIdNext: number | undefined
  mode: string | undefined
}

class Description {
  constructor(
    readonly project: string,
    readonly branch: string,
    readonly url: string,
    readonly methodName: string | null,
    readonly description: string
  ) {}
}

export interface DataSeries {
  value: string
  metricName: string | undefined
  color: string
  rawValue: number
}

export enum DBType {
  FLEET = "fleet",
  FLEET_PERF = "fleet_perf",
  JBR = "jbr",
  INTELLIJ = "intellij",
  INTELLIJ_DEV = "intellij_dev",
  QODANA = "qodana",
  BAZEL = "bazel",
  PERF_UNIT_TESTS = "perfUnitTests",
  STARTUP_TESTS = "startupTests",
  STARTUP_TESTS_DEV = "startupTests_dev",
  UNKNOWN = "unknown",
}

export interface InfoSidebar {
  data: ShallowRef<InfoData | null>
  visible: ShallowRef<boolean>

  show(data: InfoData): void

  close(): void
}

export class InfoSidebarImpl implements InfoSidebar {
  readonly data = shallowRef<InfoData | null>(null)
  readonly visible = computed(() => this.data.value != null)

  show(data: InfoData): void {
    this.data.value = data
  }

  close() {
    this.data.value = null
  }
}

export async function getArtifactsUrl(data: InfoData | null, serverConfigurator: ServerConfigurator | null): Promise<string> {
  let url = ""
  if (serverConfigurator?.table == null) {
    url = data?.artifactsUrl ?? ""
  } else if (data?.installerId ?? data?.buildId) {
    const db = serverConfigurator.db
    if (db == "perfint" || db == "perfintDev") {
      const type = await getTeamcityBuildType(db, serverConfigurator.table, data.buildId)
      url = `${tcUrl}buildConfiguration/${type}/${data.buildId}?buildTab=artifacts#${replaceUnderscore("/" + data.projectName)}`
    } else {
      url = data.artifactsUrl
    }
  }
  return url
}

export function getNavigateToTestUrl(data: InfoData | null, router: Router) {
  const currentRoute = router.currentRoute.value
  let parts = currentRoute.path.split("/")
  if (parts.at(-1) == "startup" || parts.at(1) == "ij") {
    parts = ["", "ij", "explore"]
  } else if (parts.at(1) == "fleet" && parts.at(2) == "startupDashboard") {
    parts = ["", "fleet", "startupExplore"]
  } else {
    parts[parts.length - 1] = dbTypeStore().dbType == DBType.INTELLIJ_DEV ? "testsDev" : "tests"
  }
  const branch = data?.branch ?? ""
  const machineGroup = getMachineGroupName(data?.machineName ?? "")
  const majorBranch = /\d+\.\d+/.test(branch) ? branch.slice(0, branch.indexOf(".")) : branch
  const testURL = parts.join("/")

  const queryParams: string = new URLSearchParams({
    ...currentRoute.query,
    project: data?.projectName ?? "",
    branch: majorBranch,
    machine: machineGroup,
  }).toString()

  const measures =
    data?.series
      .map((s) => s.metricName)
      .filter((m): m is string => m != undefined)
      .map((m) => (dbTypeStore().isIJStartup() && m.includes("/") ? "metrics." + m : m))
      .map((m) => encodeURIComponent(m))
      .map((m) => "&measure=" + m)
      .join("") ?? ""

  return router.resolve(testURL + "?" + queryParams + measures).href
}

export async function getSpaceUrl(data: InfoData | null, serverConfigurator: ServerConfigurator | null): Promise<string | undefined> {
  const db = serverConfigurator?.db
  if (db != null && (data?.installerId ?? data?.buildId)) {
    const decodedChanges = await calculateChanges(db, data.installerId ?? data.buildId)
    if (decodedChanges == null || decodedChanges.length === 0) {
      console.log("No changes found")
      return data.changesUrl
    } else {
      return `https://code.jetbrains.team/p/ij/repositories/ultimate/commits?query=%22${decodedChanges}%22&tab=changes`
    }
  }
  return undefined
}

function replaceUnderscore(project: string) {
  return project.replaceAll("_", "-")
}
