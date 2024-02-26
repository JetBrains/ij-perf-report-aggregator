import { computed, Ref, ShallowRef, shallowRef } from "vue"
import { Accident } from "../../../configurators/AccidentsConfigurator"
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
  dbType: DBType
}

class Description {
  constructor(
    readonly project: string,
    readonly branch: string,
    readonly url: string,
    readonly methodName: string,
    readonly description: string
  ) {}
}

export interface DataSeries {
  value: string
  metricName: string | undefined
  color: string
  nameToShow: string
}

export enum DBType {
  FLEET = "fleet",
  JBR = "jbr",
  DEV_FLEET = "dev_fleet",
  INTELLIJ = "intellij",
  INTELLIJ_DEV = "intellij_dev",
  QODANA = "qodana",
  BAZEL = "bazel",
  PERF_UNIT_TESTS = "perfUnitTests",
  IJENT_TESTS = "perfintDev_ijent",
  STARTUP_TESTS = "startupTests",
  STARTUP_TESTS_DEV = "startupTests_dev",
  UNKNOWN = "unknown",
}

export interface InfoSidebar {
  data: ShallowRef<InfoData | null>
  visible: ShallowRef<boolean>
  type: DBType

  show(data: InfoData): void

  close(): void
}

export class InfoSidebarImpl implements InfoSidebar {
  readonly data = shallowRef<InfoData | null>(null)
  readonly visible = computed(() => this.data.value != null)
  type = DBType.INTELLIJ

  constructor(type: DBType) {
    this.type = type
  }

  show(data: InfoData): void {
    this.data.value = data
  }

  close() {
    this.data.value = null
  }
}
