import { computedAsync } from "@vueuse/core"
import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { computed, ref } from "vue"
import { AccidentsConfigurator } from "../../../configurators/AccidentsConfigurator"
import { measureNameToLabel } from "../../../configurators/MeasureConfigurator"
import { ServerWithCompressConfigurator } from "../../../configurators/ServerWithCompressConfigurator"
import { findDeltaInData, getDifferenceString } from "../../../util/Delta"
import { useSettingsStore } from "../../settings/settingsStore"
import { ValueUnit } from "../chart"
import { durationAxisPointerFormatter, getValueFormatterByMeasureName, isDurationFormatterApplicable, nsToMs, timeFormatWithoutSeconds } from "../formatter"
import { encodeRison } from "../rison"
import { buildUrl, DataSeries, DBType, InfoData } from "./InfoSidebar"

function filterUniqueByName(objects: CallbackDataParams[] | null): CallbackDataParams[] {
  const seen = new Set()
  return objects?.filter((item) => {
    const duplicate = seen.has(item.seriesName)
    seen.add(item.seriesName)
    return !duplicate
  }) as CallbackDataParams[]
}

export function getInfoDataFrom(
  dbType: DBType,
  params: CallbackDataParams | CallbackDataParams[],
  valueUnit: ValueUnit,
  accidentsConfigurator: AccidentsConfigurator | null
): InfoData {
  const accidents = accidentsConfigurator?.value
  if (Array.isArray(params) && params.length > 1) {
    const filteredParams = filterUniqueByName(params)

    const dataSeries = filteredParams[0].value as OptionDataValue[]
    const dateMs = dataSeries[0] as number
    const machineName: string = dataSeries[3] as string
    const buildId: number = dataSeries[4] as number
    const projectName: string = dataSeries[5] as string
    const installerId: number = dataSeries[6] as number
    const buildVersion: number = dataSeries[7] as number
    const buildNum1: number = dataSeries[8] as number
    const buildNum2: number = dataSeries[9] as number
    const branch: string = dataSeries[10] as string

    const series: DataSeries[] = []
    const prefixes = filteredParams.map((param) => getPrefix(param.seriesName as string))
    const commonPrefix = getCommonPrefix(prefixes)
    for (const param of filteredParams) {
      const currentSeriesData = param.value as OptionDataValue[]
      param.seriesName = removePrefix(param.seriesName as string, commonPrefix)
      const value = getValueFormatterByMeasureName(param.seriesName)(currentSeriesData[1] as number)
      series.push({ metricName: param.seriesName, value, color: param.color as string })
    }

    const fullBuildId = `${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
    const changesUrl = `${buildUrl(installerId)}&tab=changes`
    const artifactsUrl = `${buildUrl(buildId)}&tab=artifacts`
    const installerUrl = `${buildUrl(installerId)}&tab=artifacts`

    const filteredAccidents = computed(() => {
      return accidents?.value?.get(projectName + "_" + fullBuildId)
    })

    return {
      build: fullBuildId,
      accidents: filteredAccidents,
      description: ref(null),
      artifactsUrl,
      changesUrl,
      installerUrl,
      series,
      date: timeFormatWithoutSeconds.format(dateMs),
      machineName,
      projectName,
      deltaNext: undefined,
      deltaPrevious: undefined,
      title: "Details",
      installerId,
      buildId,
      branch,
    }
  } else {
    if (Array.isArray(params)) {
      params = params[0]
    }
    const dataSeries = params.value as OptionDataValue[]
    const dateMs = dataSeries[0] as number
    const value: number = useSettingsStore().scaling ? (dataSeries.at(-1) as number) : (dataSeries[1] as number)
    let projectName: string = params.seriesName as string
    let machineName: string | undefined
    let metricName: string | undefined
    let buildId: number | undefined
    let installerId: number | undefined
    let buildVersion: number | undefined
    let buildNum1: number | undefined
    let buildNum2: number | undefined
    let type: ValueUnit | undefined = valueUnit
    let buildNumber: string | undefined
    let accidentBuild: string | undefined
    let branch: string | undefined
    if (dbType == DBType.DEV_FLEET) {
      machineName = dataSeries[2] as string
      buildId = dataSeries[3] as number
      projectName = dataSeries[4] as string
      branch = dataSeries[5] as string
    }
    if (dbType == DBType.INTELLIJ_DEV || dbType == DBType.PERF_UNIT_TESTS) {
      metricName = dataSeries[2] as string
      if (dataSeries[3] == "c") {
        type = "counter"
      }
      machineName = dataSeries[4] as string
      buildId = dataSeries[5] as number
      projectName = dataSeries[6] as string
      branch = dataSeries[7] as string
      accidentBuild = buildId.toString()
    }
    if (dbType == DBType.FLEET || dbType == DBType.STARTUP_TESTS) {
      metricName = measureNameToLabel(dataSeries[2] as string)
      if (!isDurationFormatterApplicable(metricName)) {
        type = "counter"
      }
      machineName = dataSeries[3] as string
      buildId = dataSeries[4] as number
      projectName = dataSeries[5] as string
      installerId = dataSeries[6] as number
      buildVersion = dataSeries[7] as number
      buildNum1 = dataSeries[8] as number
      buildNum2 = dataSeries[9] as number
      branch = dataSeries[10] as string
      accidentBuild = `${buildVersion}.${buildNum1}`
    }
    if (dbType == DBType.JBR) {
      metricName = dataSeries[2] as string
      if (dataSeries[3] == "c") {
        type = "counter"
      }
      machineName = dataSeries[4] as string
      buildId = dataSeries[5] as number
      projectName = dataSeries[6] as string
      buildNumber = dataSeries[8] as string
      branch = dataSeries[9] as string
    }
    if (dbType == DBType.INTELLIJ) {
      metricName = dataSeries[2] as string
      if (dataSeries[3] == "c") {
        type = "counter"
      }
      machineName = dataSeries[4] as string
      buildId = dataSeries[5] as number
      projectName = dataSeries[6] as string
      installerId = dataSeries[7] as number
      buildVersion = dataSeries[8] as number
      buildNum1 = dataSeries[9] as number
      buildNum2 = dataSeries[10] as number
      branch = dataSeries[11] as string
      accidentBuild = `${buildVersion}.${buildNum1}`
    }
    if (dbType == DBType.BAZEL) {
      metricName = dataSeries[2] as string
      if (dataSeries[3] == "c") {
        type = "counter"
      }
      machineName = dataSeries[4] as string
      buildId = dataSeries[5] as number
      projectName = dataSeries[6] as string
      branch = dataSeries[7] as string
      console.log(dataSeries)
    }
    if (dbType == DBType.QODANA) {
      machineName = dataSeries[2] as string
      buildId = dataSeries[3] as number
    }
    if (dbType == DBType.UNKNOWN) {
      console.error("Unknown type of DB")
    }
    const fullBuildId = buildVersion == undefined ? buildNumber : `${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
    const changesUrl = installerId == undefined ? `${buildUrl(buildId as number)}&tab=changes` : `${buildUrl(installerId)}&tab=changes`
    const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
    const installerUrl = installerId == undefined ? undefined : `${buildUrl(installerId)}&tab=artifacts`
    const showValue: string = durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(value) : value, type)
    const filteredAccidents = computed(() => {
      return accidents?.value?.get(projectName + "_" + accidentBuild) ?? accidents?.value?.get(projectName + "/" + metricName + "_" + accidentBuild)
    })
    const description = computedAsync(async () => {
      return await getDescriptionFromMetaDb(projectName, "master")
    })
    const delta = findDeltaInData(dataSeries)
    let deltaPrevious: string | undefined
    let deltaNext: string | undefined
    if (delta != undefined) {
      if (delta.prev != null) {
        deltaPrevious = getDifferenceString(value, delta.prev, valueUnit == "ms", type as string)
      }
      if (delta.next != null) {
        deltaNext = getDifferenceString(value, delta.next, valueUnit == "ms", type as string)
      }
    }
    return {
      series: [{ metricName, value: showValue, color: params.color as string }],
      build: fullBuildId,
      artifactsUrl,
      changesUrl,
      installerUrl,
      date: timeFormatWithoutSeconds.format(dateMs),
      machineName: machineName as string,
      projectName,
      title: "Details",
      installerId,
      accidents: filteredAccidents,
      buildId: buildId as number,
      description,
      branch,
      deltaPrevious,
      deltaNext,
    }
  }
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

async function getDescriptionFromMetaDb(project: string | undefined, branch: string): Promise<Description | null> {
  const description_url = ServerWithCompressConfigurator.DEFAULT_SERVER_URL + "/api/meta/description/"
  const response = await fetch(description_url + encodeRison({ project, branch }))
  return response.ok ? response.json() : null
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
