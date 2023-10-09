import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { getValueFormatterByMeasureName, timeFormatWithoutSeconds } from "../formatter"
import { buildUrl, DataSeries, InfoData } from "./InfoSidebar"

export interface InfoDataFromStartup extends InfoData {
  series: DataSeries[]
}

function filterUniqueByName(objects: CallbackDataParams[] | null): CallbackDataParams[] | undefined {
  const seen = new Set()
  return objects?.filter((item) => {
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
    const machineName: string = dataSeries[3] as string
    const buildId: number = dataSeries[4] as number
    const projectName: string = dataSeries[5] as string
    const installerId: number = dataSeries[6] as number
    const buildVersion: number = dataSeries[7] as number
    const buildNum1: number = dataSeries[8] as number
    const buildNum2: number = dataSeries[9] as number

    const series: DataSeries[] = []
    const prefixes = params.map((param) => getPrefix(param.seriesName as string))
    const commonPrefix = getCommonPrefix(prefixes)
    for (const param of params) {
      const currentSeriesData = param.value as OptionDataValue[]
      param.seriesName = removePrefix(param.seriesName as string, commonPrefix)
      const value = getValueFormatterByMeasureName(param.seriesName)(currentSeriesData[1] as number)
      series.push({ metricName: param.seriesName, value, color: param.color as string })
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
