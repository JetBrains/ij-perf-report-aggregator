import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { durationAxisPointerFormatter, nsToMs, timeFormatWithoutSeconds } from "shared/src/formatter"
import { computed, ShallowRef, shallowRef } from "vue"
import { ValueUnit } from "shared/src/chart"

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
  duration: string
  build: string|undefined
  date: string
}

const buildUrl = (id: number) =>`https://buildserver.labs.intellij.net/viewLog.html?buildId=${id}`

export function getInfoDataFrom(params: CallbackDataParams, valueUnit: ValueUnit): InfoData {
  const dataSeries = params.value as OptionDataValue[]
  //dev builds
  if(dataSeries.length == 4){
    const [
      dateMs,
      durationMs,
      machineName,
      buildId,
    ]  = dataSeries
    const changesUrl = `${buildUrl(buildId as number)}&tab=changes`
    const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
    return {
      build: undefined,
      artifactsUrl,
      changesUrl,
      installerUrl: undefined,
      color: params.color as string,
      date: timeFormatWithoutSeconds.format(dateMs as number),
      duration: durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(durationMs) : durationMs as number),
      machineName: machineName as string,
      projectName: params.seriesName!,
      title: "Title",
    }
  }
  //fleet
  if(dataSeries.length == 8){
    const [
      dateMs,
      durationMs,
      machineName,
      buildId,
      installerId,
      buildVersion,
      buildNum1,
      buildNum2,
    ] = params.value as OptionDataValue[]
    const fullBuildId = `${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
    const changesUrl = `${buildUrl(buildId as number)}&tab=changes`
    const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
    const installerUrl = `${buildUrl(installerId as number)}&tab=artifacts`
    return {
      build: fullBuildId,
      artifactsUrl,
      changesUrl,
      installerUrl,
      color: params.color as string,
      date: timeFormatWithoutSeconds.format(dateMs as number),
      duration: durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(durationMs) : durationMs as number),
      machineName: machineName as string,
      projectName: params.seriesName!,
      title: "Title",
    }
  }
  //jbr
  if(dataSeries.length == 5){
    const [
      dateMs,
      durationMs,
      _,
      machineName,
      buildId,
    ]  = dataSeries
    const changesUrl = `${buildUrl(buildId as number)}&tab=changes`
    const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
    return {
      build: undefined,
      artifactsUrl,
      changesUrl,
      installerUrl: undefined,
      color: params.color as string,
      date: timeFormatWithoutSeconds.format(dateMs as number),
      duration: durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(durationMs) : durationMs as number),
      machineName: machineName as string,
      projectName: params.seriesName!,
      title: "Title",
    }
  }
  const [
    dateMs,
    durationMs,
    _,
    machineName,
    buildId,
    installerId,
    buildVersion,
    buildNum1,
    buildNum2,
  ] = params.value as OptionDataValue[]
  const fullBuildId = `${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
  const changesUrl = `${buildUrl(buildId as number)}&tab=changes`
  const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
  const installerUrl = `${buildUrl(installerId as number)}&tab=artifacts`

  return {
    build: fullBuildId,
    artifactsUrl,
    changesUrl,
    installerUrl,
    color: params.color as string,
    date: timeFormatWithoutSeconds.format(dateMs as number),
    duration: durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(durationMs) : durationMs as number),
    machineName: machineName as string,
    projectName: params.seriesName!,
    title: "Title",
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