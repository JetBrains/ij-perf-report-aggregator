import { computed, ShallowRef, shallowRef } from "vue"
export const tcUrl = "https://buildserver.labs.intellij.net/"
export const buildUrl = (id: number) => `${tcUrl}viewLog.html?buildId=${id}`

export interface InfoData {
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
}

export interface DataSeries {
  value: string
  metricName: string | undefined
  color: string
}

export enum DBType {
  FLEET = "fleet",
  JBR = "jbr",
  DEV_FLEET = "dev_fleet",
  INTELLIJ = "intellij",
  INTELLIJ_DEV = "intellij_dev",
  QODANA = "qodana",
  UNKNOWN = "unknown",
}

export interface InfoSidebar<T extends InfoData> {
  data: ShallowRef<T | null>
  visible: ShallowRef<boolean>
  type: DBType

  show(data: T): void

  close(): void
}

export class InfoSidebarImpl<D extends InfoData> implements InfoSidebar<D> {
  readonly data = shallowRef<D | null>(null)
  readonly visible = computed(() => this.data.value != null)
  type = DBType.INTELLIJ

  constructor(type: DBType) {
    this.type = type
  }

  show(data: D): void {
    this.data.value = data
  }

  close() {
    this.data.value = null
  }
}
