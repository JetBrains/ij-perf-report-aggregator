import { computedAsync } from "@vueuse/core"
import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { computed, Ref } from "vue"
import { Accident, AccidentsConfigurator } from "../../../configurators/AccidentsConfigurator"
import { ServerWithCompressConfigurator } from "../../../configurators/ServerWithCompressConfigurator"
import { dbTypeStore } from "../../../shared/dbTypes"
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

export function getBuildId(params: CallbackDataParams): number | undefined {
  const dbType = dbTypeStore().dbType
  const dataSeries = params.value as OptionDataValue[]

  let buildId: number | undefined

  if (dbType == DBType.DEV_FLEET) {
    buildId = dataSeries[3] as number
  }
  if (dbType == DBType.INTELLIJ_DEV || dbType == DBType.PERF_UNIT_TESTS) {
    buildId = dataSeries[5] as number
  }
  if (dbType == DBType.FLEET || dbType == DBType.STARTUP_TESTS) {
    buildId = dataSeries[4] as number
  }
  if (dbType == DBType.STARTUP_TESTS_DEV) {
    buildId = dataSeries[4] as number
  }
  if (dbType == DBType.JBR) {
    buildId = dataSeries[5] as number
  }
  if (dbType == DBType.INTELLIJ) {
    buildId = dataSeries[5] as number
  }
  if (dbType == DBType.BAZEL) {
    buildId = dataSeries[5] as number
  }
  if (dbType == DBType.QODANA) {
    buildId = dataSeries[3] as number
  }
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }
  return buildId
}

export function getAccidentBuild(params: CallbackDataParams): string | undefined {
  const dbType = dbTypeStore().dbType
  if (dbType == DBType.INTELLIJ_DEV || dbType == DBType.PERF_UNIT_TESTS) {
    return getBuildId(params)?.toString()
  }
  if (dbType == DBType.FLEET || dbType == DBType.STARTUP_TESTS) {
    return getFullBuildId(params)
  }
  if (dbType == DBType.DEV_FLEET) {
    return getBuildId(params)?.toString()
  }
  if (dbType == DBType.STARTUP_TESTS_DEV) {
    return getBuildId(params)?.toString()
  }
  if (dbType == DBType.INTELLIJ) {
    return getFullBuildId(params)
  }
  if (dbType == DBType.BAZEL) {
    return getBuildId(params)?.toString()
  }
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }
  return
}

export function getFullBuildId(params: CallbackDataParams): string | undefined {
  const dbType = dbTypeStore().dbType
  const dataSeries = params.value as OptionDataValue[]

  let buildVersion: number | undefined
  let buildNum1: number | undefined
  let buildNum2: number | undefined
  let buildNumber: string | undefined

  if (dbType == DBType.FLEET || dbType == DBType.STARTUP_TESTS) {
    buildVersion = dataSeries[7] as number
    buildNum1 = dataSeries[8] as number
    buildNum2 = dataSeries[9] as number
  }
  if (dbType == DBType.JBR) {
    buildNumber = dataSeries[7] as string
  }
  if (dbType == DBType.INTELLIJ) {
    buildVersion = dataSeries[8] as number
    buildNum1 = dataSeries[9] as number
    buildNum2 = dataSeries[10] as number
  }
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }
  return buildVersion == undefined ? buildNumber : `${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
}

function getInfo(params: CallbackDataParams, valueUnit: ValueUnit, accidents: Ref<Map<string, Accident[]> | undefined> | undefined) {
  const seriesName = params.seriesName as string
  const dataSeries = params.value as OptionDataValue[]
  const dateMs = dataSeries[0] as number
  let projectName: string = params.seriesName as string
  let machineName: string | undefined
  let metricName: string | undefined
  let installerId: number | undefined
  let type: ValueUnit | undefined = valueUnit
  let branch: string | undefined
  const dbType = dbTypeStore().dbType
  if (dbType == DBType.DEV_FLEET) {
    machineName = dataSeries[2] as string
    projectName = dataSeries[4] as string
    branch = dataSeries[5] as string
  }
  if (dbType == DBType.INTELLIJ_DEV || dbType == DBType.PERF_UNIT_TESTS) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    machineName = dataSeries[4] as string
    projectName = dataSeries[6] as string
    branch = dataSeries[7] as string
  }
  if (dbType == DBType.FLEET || dbType == DBType.STARTUP_TESTS) {
    metricName = dataSeries[2] as string
    if (!isDurationFormatterApplicable(metricName)) {
      type = "counter"
    }
    machineName = dataSeries[3] as string
    projectName = dataSeries[5] as string
    installerId = dataSeries[6] as number
    branch = dataSeries[10] as string
  }
  if (dbType == DBType.STARTUP_TESTS_DEV) {
    metricName = dataSeries[2] as string
    if (!isDurationFormatterApplicable(metricName)) {
      type = "counter"
    }
    machineName = dataSeries[3] as string
    projectName = dataSeries[5] as string
    branch = dataSeries[6] as string
  }
  if (dbType == DBType.JBR) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    machineName = dataSeries[4] as string
    projectName = dataSeries[6] as string
    branch = dataSeries[8] as string
  }
  if (dbType == DBType.INTELLIJ) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    machineName = dataSeries[4] as string
    projectName = dataSeries[6] as string
    installerId = dataSeries[7] as number
    branch = dataSeries[11] as string
  }
  if (dbType == DBType.BAZEL) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    machineName = dataSeries[4] as string
    projectName = dataSeries[6] as string
    branch = dataSeries[7] as string
  }
  if (dbType == DBType.QODANA) {
    machineName = dataSeries[2] as string
    branch = dataSeries[5] as string
    projectName = dataSeries[4] as string
  }
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }

  const buildId: number | undefined = getBuildId(params)
  const changesUrl = installerId == undefined ? `${buildUrl(buildId as number)}&buildTab=changes` : `${buildUrl(installerId)}&buildTab=changes`
  const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
  const installerUrl = installerId == undefined ? undefined : `${buildUrl(installerId)}&tab=artifacts`
  const accidentBuild = getAccidentBuild(params)

  const filteredAccidents = computed(() => {
    if (accidentBuild == undefined) return []
    const testAccident = accidents?.value?.get(projectName + "_" + accidentBuild) ?? []
    const metricAccident = metricName == undefined ? [] : (accidents?.value?.get(projectName + "/" + metricName + "_" + accidentBuild) ?? [])
    const buildAccident = accidents?.value?.get(`_${accidentBuild}`) ?? []
    return [...testAccident, ...buildAccident, ...metricAccident]
  })
  const description = computedAsync(async () => {
    return await getDescriptionFromMetaDb(projectName, "master")
  })
  return {
    seriesName,
    build: getFullBuildId(params),
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
    metricName,
    type,
  }
}

export function getInfoDataFrom(
  params: CallbackDataParams | CallbackDataParams[],
  valueUnit: ValueUnit,
  accidentsConfigurator: AccidentsConfigurator | null,
  chartDataUrl: string
): InfoData {
  const accidents = accidentsConfigurator?.value
  if (Array.isArray(params) && params.length > 1) {
    const filteredParams = filterUniqueByName(params)
    const info = getInfo(params[0], valueUnit, accidents)
    const series: DataSeries[] = []
    for (const param of filteredParams) {
      const currentSeriesData = param.value as OptionDataValue[]
      const value = getValueFormatterByMeasureName(param.seriesName as string)(currentSeriesData[1] as number)
      series.push({ metricName: param.seriesName as string, value, color: param.color as string })
    }

    return { ...info, series, deltaPrevious: undefined, deltaNext: undefined, chartDataUrl, buildIdPrevious: undefined, buildIdNext: undefined }
  } else {
    if (Array.isArray(params)) {
      params = params[0]
    }
    const info = getInfo(params, valueUnit, accidents)
    const dataSeries = params.value as OptionDataValue[]
    const value: number = useSettingsStore().scaling ? (dataSeries.at(-1) as number) : (dataSeries[1] as number)
    const showValue: string = durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(value) : value, info.type)
    const delta = findDeltaInData(dataSeries)
    let deltaPrevious: string | undefined
    let deltaNext: string | undefined
    let buildIdPrevious: number | undefined
    let buildIdNext: number | undefined
    if (delta != undefined) {
      if (delta.prev != null) {
        deltaPrevious = getDifferenceString(value, delta.prev, valueUnit == "ms", info.type)
        buildIdPrevious = delta.prevBuildId
      }
      if (delta.next != null) {
        deltaNext = getDifferenceString(value, delta.next, valueUnit == "ms", info.type)
        buildIdNext = delta.nextBuildId
      }
    }
    return {
      ...info,
      deltaNext,
      deltaPrevious,
      series: [{ metricName: info.metricName, value: showValue, color: params.color as string }],
      chartDataUrl,
      buildIdPrevious,
      buildIdNext,
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
  return response.ok ? ((await response.json()) as Description) : null
}
