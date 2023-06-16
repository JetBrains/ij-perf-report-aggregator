import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { computed, ShallowRef, shallowRef } from "vue"
import { calculateChanges } from "../util/changes"
import { timeFormatWithoutSeconds } from "./common/formatter"

export interface InfoSidebarStartup {
  data: ShallowRef<InfoDataFromStartup | null>
  visible: ShallowRef<boolean>

  show(data: InfoDataFromStartup | null): void

  close(): void
}

export interface DataSerie {
  value: number
  name: string
  color: string
}

export interface InfoDataFromStartup {
  build: string
  artifactsUrl: string
  changesUrl: string
  installerUrl: string
  date: string
  series: DataSerie[]
  machineName: string
  projectName: string
  title: string
  installerId: number
  buildId: number
  changes: string | undefined
}
const buildUrl = (id: number) => `https://buildserver.labs.intellij.net/viewLog.html?buildId=${id}`

function filterUniqueByName(objects: CallbackDataParams[]| null): CallbackDataParams[]|undefined {
  const seen = new Set()
  return objects?.filter(item => {
    const duplicate = seen.has(item.seriesName)
    seen.add(item.seriesName)
    return !duplicate
  })
}


export function getInfoDataForStartup(originalParams: CallbackDataParams[] | null): InfoDataFromStartup | null {
  const params = filterUniqueByName(originalParams)
  if (params && params.length > 0) {
    const dataSeries = params[0].value as OptionDataValue[]
    const dateMs = dataSeries[0] as number
    const machineName: string = dataSeries[2] as string
    const buildId: number = dataSeries[3] as number
    const projectName: string = dataSeries[4] as string
    const installerId: number = dataSeries[5] as number
    const buildVersion: number = dataSeries[6] as number
    const buildNum1: number = dataSeries[7] as number
    const buildNum2: number = dataSeries[8] as number

    const series: DataSerie[] = []
    const prefixes = params.map(param => getPrefix(param.seriesName as string))
    const commonPrefix = getCommonPrefix(prefixes)
    for (const param of params) {
      const currentSeriesData = param.value as OptionDataValue[]
      param.seriesName = removePrefix(param.seriesName as string, commonPrefix)
      series.push({ name: param.seriesName, value: currentSeriesData[1] as number, color: param.color as string })
    }

    const fullBuildId = `${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
    const changesUrl = `${buildUrl(installerId)}&tab=changes`
    const artifactsUrl = `${buildUrl(buildId)}&tab=artifacts`
    const installerUrl = `${buildUrl(installerId)}&tab=artifacts`



    return {
      build: fullBuildId,
      artifactsUrl,
      changesUrl,
      installerUrl,
      series,
      date: timeFormatWithoutSeconds.format(dateMs),
      machineName,
      projectName,
      title: "Details",
      changes: undefined,
      installerId,
      buildId,
    }
  }
  return null
}

function getPrefix(name: string): string {
  const lastDot = name.lastIndexOf(".")
  const lastSlash = name.lastIndexOf("/")
  const lastIndex = Math.max(lastDot, lastSlash)
  return name.slice(0, Math.max(0, lastIndex))
}

function getCommonPrefix(names: string[]): string {
  if (names.length === 0) {
    return ""
  }
  let commonPrefix = names[0]
  for (let i = 1; i < names.length; i++) {
    while (!names[i].startsWith(commonPrefix) && commonPrefix) {
      const lastIndex = Math.max(commonPrefix.lastIndexOf("/"), commonPrefix.lastIndexOf("."))
      commonPrefix = commonPrefix.slice(0, Math.max(0, lastIndex))
    }
  }
  if (commonPrefix.endsWith("/") || commonPrefix.endsWith(".")) {
    commonPrefix = commonPrefix.slice(0, Math.max(0, commonPrefix.length - 1))
  }
  return commonPrefix
}

function removePrefix(name: string, prefix: string): string {
  if (name.startsWith(prefix + ".") || name.startsWith(prefix + "/")) {
    return name.slice(Math.max(0, prefix.length + 1))
  }
  return name
}

export class InfoSidebarStartupImpl implements InfoSidebarStartup {
  readonly data = shallowRef<InfoDataFromStartup | null>(null)
  readonly visible = computed(() => this.data.value != null)

  show(data: InfoDataFromStartup): void {
    this.data.value = data
  }

  close() {
    this.data.value = null
  }
}
