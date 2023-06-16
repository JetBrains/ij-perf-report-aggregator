import { computed, ShallowRef, shallowRef } from "vue"

export const buildUrl = (id: number) => `https://buildserver.labs.intellij.net/viewLog.html?buildId=${id}`

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
}

export interface DataSeries {
  value: string
  metricName: string | undefined
  color: string
}

export interface InfoSidebar<T extends InfoData> {
  data: ShallowRef<T | null>
  visible: ShallowRef<boolean>

  show(data: T): void

  close(): void
}

export class InfoSidebarImpl<D extends InfoData> implements InfoSidebar<D> {
  readonly data = shallowRef<D | null>(null)
  readonly visible = computed(() => this.data.value != null)

  show(data: D): void {
    this.data.value = data
  }

  close() {
    this.data.value = null
  }
}
