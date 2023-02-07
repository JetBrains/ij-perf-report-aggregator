import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { ValueUnit } from "shared/src/chart"
import { durationAxisPointerFormatter, nsToMs, timeFormatWithoutSeconds } from "shared/src/formatter"
import { computed, ShallowRef, shallowRef } from "vue"

export interface InfoSidebarVm {
  data: ShallowRef<InfoData | null>
  visible: ShallowRef<boolean>

  show(data: InfoData): void

  close(): void
}

export interface InfoData {
  color: string
  title: string
  projectName: string
  changesUrl: string
  artifactsUrl: string
  installerUrl: string|undefined
  machineName: string
  value: string
  build: string|undefined
  date: string
}

const buildUrl = (id: number) => `https://buildserver.labs.intellij.net/viewLog.html?buildId=${id}`

export function getInfoDataFrom(params: CallbackDataParams, valueUnit: ValueUnit): InfoData {
  const dataSeries = params.value as OptionDataValue[]
  const dateMs = dataSeries[0] as number
  const value: number = dataSeries[1]  as number
  let machineName: string|undefined
  let buildId: number|undefined
  let installerId: number|undefined
  let buildVersion: number|undefined
  let buildNum1: number|undefined
  let buildNum2: number|undefined
  let type: ValueUnit|undefined = valueUnit
  let buildNumber: string|undefined
  //dev builds
  if(dataSeries.length == 4){
    machineName = dataSeries[2] as string
    buildId = dataSeries[3] as number
  }
  //fleet
  if(dataSeries.length == 8){
    machineName = dataSeries[2] as string
    buildId = dataSeries[3] as number
    installerId = dataSeries[4] as number
    buildVersion = dataSeries[5] as number
    buildNum1 = dataSeries[6] as number
    buildNum2 = dataSeries[7] as number
  }
  //jbr
  if (dataSeries.length == 6) {
    if (dataSeries[2] == "c") {
      type = "counter"
    }
    machineName = dataSeries[3] as string
    buildId = dataSeries[4] as number
    buildNumber = dataSeries[5] as string
  }
  if (dataSeries.length == 9) {
    if (dataSeries[2] == "c") {
      type = "counter"
    }
    machineName = dataSeries[3] as string
    buildId = dataSeries[4] as number
    installerId = dataSeries[5] as number
    buildVersion = dataSeries[6] as number
    buildNum1 = dataSeries[7] as number
    buildNum2 = dataSeries[8] as number
  }

  const fullBuildId = buildVersion == undefined ? buildNumber :`${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
  const changesUrl = installerId == undefined  ? `${buildUrl(buildId as number)}&tab=changes` : `${buildUrl(installerId)}&tab=changes`
  const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
  const installerUrl = installerId == undefined ? undefined :`${buildUrl(installerId)}&tab=artifacts`

  let showValue = value.toString()
  if(type != "counter"){
    showValue = durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(value) : value )
  }

  return {
    build: fullBuildId,
    artifactsUrl,
    changesUrl,
    installerUrl,
    color: params.color as string,
    date: timeFormatWithoutSeconds.format(dateMs),
    value: showValue,
    machineName: machineName as string,
    projectName: params.seriesName as string,
    title: "Details",
  }
}

export class InfoSidebarVmImpl implements InfoSidebarVm {
  readonly data = shallowRef<InfoData | null>(null)
  readonly visible = computed(() => this.data.value != null)

  show(data: InfoData): void {
    this.data.value = data
  }

  close() {
    this.data.value = null
  }
}