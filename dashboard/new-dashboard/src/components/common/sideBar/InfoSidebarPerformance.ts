import { computedAsync } from "@vueuse/core"
import type { DefaultLabelFormatterCallbackParams as CallbackDataParams } from "echarts"
import type { OptionDataValue } from "../../../shared/echarts-types"
import { computed, Ref } from "vue"
import { Accident, AccidentsConfigurator } from "../../../configurators/accidents/AccidentsConfigurator"
import { ServerWithCompressConfigurator } from "../../../configurators/ServerWithCompressConfigurator"
import { dbTypeStore, resolveMeasureUnitForDb } from "../../../shared/dbTypes"
import { findDeltaInData, getDifferenceString } from "../../../util/Delta"
import { useSettingsStore } from "../../settings/settingsStore"
import { ValueUnit } from "../chart"
import { formatMeasureValue, timeFormatWithoutSeconds } from "../formatter"
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

export function getBuildId(dataSeries: (number | string)[][]): number[] | undefined
export function getBuildId(dataSeries: (number | string)[]): number | undefined
export function getBuildId(dataSeries: (number | string)[] | (number | string)[][]): number | number[] | undefined {
  const dbType = dbTypeStore().dbType

  let buildId: number | undefined

  if (
    dbType == DBType.INTELLIJ_DEV ||
    dbType == DBType.PERF_UNIT_TESTS ||
    dbType == DBType.FLEET_PERF ||
    dbType == DBType.DIOGEN ||
    dbType == DBType.QODANA ||
    dbType == DBType.TOOLBOX
  ) {
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
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }
  return buildId
}

function getAccidentBuild(params: CallbackDataParams): string | undefined {
  const dbType = dbTypeStore().dbType
  if (dbType == DBType.INTELLIJ_DEV || dbType == DBType.PERF_UNIT_TESTS || dbType == DBType.FLEET_PERF || dbType == DBType.DIOGEN || dbType == DBType.TOOLBOX) {
    return getBuildId(params.value as number[])?.toString()
  }
  if (dbType == DBType.FLEET || dbType == DBType.STARTUP_TESTS) {
    return getFullBuildId(params)
  }
  if (dbType == DBType.STARTUP_TESTS_DEV) {
    return getBuildId(params.value as number[])?.toString()
  }
  if (dbType == DBType.INTELLIJ) {
    return getFullBuildId(params)
  }
  if (dbType == DBType.BAZEL) {
    return getBuildId(params.value as number[])?.toString()
  }
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }
  return undefined
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

export function getBasicInfo(params: CallbackDataParams, valueUnit: ValueUnit) {
  const seriesName = params.seriesName as string
  const dataSeries = params.value as OptionDataValue[]
  const dateMs = dataSeries[0] as number
  let projectName: string = params.seriesName as string
  let machineName: string | undefined
  let metricName: string | undefined
  let installerId: number | undefined
  let type: ValueUnit | undefined = valueUnit
  let branch: string | undefined
  let mode: string | undefined
  const dbType = dbTypeStore().dbType
  if (dbType == DBType.FLEET_PERF) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    machineName = dataSeries[4] as string
    projectName = dataSeries[6] as string
    branch = dataSeries[7] as string
  }
  if (dbType == DBType.INTELLIJ_DEV || dbType == DBType.PERF_UNIT_TESTS || dbType == DBType.DIOGEN || dbType == DBType.QODANA || dbType == DBType.TOOLBOX) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    machineName = dataSeries[4] as string
    projectName = dataSeries[6] as string
    branch = dataSeries[7] as string
    if (dbType == DBType.INTELLIJ_DEV) {
      mode = dataSeries[8] as string
    }
  }
  if (dbType == DBType.FLEET || dbType == DBType.STARTUP_TESTS) {
    metricName = dataSeries[2] as string
    machineName = dataSeries[3] as string
    projectName = dataSeries[5] as string
    installerId = dataSeries[6] as number
    branch = dataSeries[10] as string
  }
  if (dbType == DBType.STARTUP_TESTS_DEV) {
    metricName = dataSeries[2] as string
    machineName = dataSeries[3] as string
    projectName = dataSeries[5] as string
    branch = dataSeries[6] as string
  }
  if (dbType == DBType.JBR) {
    metricName = dataSeries[2] as string
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
    mode = dataSeries[12] as string
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
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }

  // An explicit value-unit on the chart wins over the per-series type detected from the data.
  if (valueUnit !== "auto") {
    type = valueUnit
  }

  const buildId: number | undefined = getBuildId(params.value as (string | number)[])
  const changesUrl = installerId == undefined ? `${buildUrl(buildId as number)}&buildTab=changes` : `${buildUrl(installerId)}&buildTab=changes`
  const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
  const installerUrl = installerId == undefined ? undefined : `${buildUrl(installerId)}&tab=artifacts`

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
    buildId: buildId as number,
    branch,
    metricName,
    type,
    mode,
  }
}

function getInfoWithAccidentsAndDescription(params: CallbackDataParams, valueUnit: ValueUnit, accidents: Ref<Map<string, Accident[]> | undefined> | undefined) {
  const basicInfo = getBasicInfo(params, valueUnit)
  const accidentBuild = getAccidentBuild(params)

  const filteredAccidents = computed(() => {
    if (accidentBuild == undefined) return []
    const testAccident = accidents?.value?.get(basicInfo.projectName + "_" + accidentBuild) ?? []
    const metricAccident = basicInfo.metricName == undefined ? [] : (accidents?.value?.get(basicInfo.projectName + "/" + basicInfo.metricName + "_" + accidentBuild) ?? [])
    const buildAccident = accidents?.value?.get(`_${accidentBuild}`) ?? []
    return [...testAccident, ...buildAccident, ...metricAccident]
  })

  const description = computedAsync(async () => {
    return await getDescriptionFromMetaDb(basicInfo.projectName, basicInfo.branch ?? "master")
  }) as Ref<Description | null>

  const owner = computedAsync(async () => {
    return await getOwnerFromMetaDb(basicInfo.projectName)
  }) as Ref<string | null>

  return {
    ...basicInfo,
    accidents: filteredAccidents,
    description,
    owner,
  }
}

function getInfo(params: CallbackDataParams, valueUnit: ValueUnit, accidents: Ref<Map<string, Accident[]> | undefined> | undefined) {
  return getInfoWithAccidentsAndDescription(params, valueUnit, accidents)
}

function toDataSeries(params: CallbackDataParams, valueUnit: ValueUnit, scaling: boolean): DataSeries {
  const dataSeries = params.value as OptionDataValue[]
  // Take the metric from the data, not from the series label: a chart drawing one measure for
  // several projects labels its series by project, and consumers of `metricName` (bisect, accident
  // reports, drilldown links) would then get a project name as the metric.
  const metricName = getBasicInfo(params, valueUnit).metricName
  const unit = resolveMeasureUnitForDb(metricName ?? "", { storedType: dataSeries[3] as string, valueUnit, scaling })
  const rawValue = (scaling ? dataSeries.at(-1) : dataSeries[1]) as number
  return { metricName, label: params.seriesName ?? metricName ?? "", value: formatMeasureValue(rawValue, unit), color: params.color as string, rawValue }
}

export function getInfoDataFrom(
  params: CallbackDataParams | CallbackDataParams[],
  valueUnit: ValueUnit,
  accidentsConfigurator: AccidentsConfigurator | null,
  chartDataUrl: string,
  seriesContext?: { seriesValues: number[] | undefined; pointIndex: number | undefined }
): InfoData {
  const accidents = accidentsConfigurator?.value
  const paramsList = Array.isArray(params) ? params : [params]
  // Several points can be selected at once - all series of a chart at the hovered x, or every
  // series matching the `point` URL parameter. They are listed in `series`, while everything else
  // below (project, build, deltas) describes the first one, which is what the sidebar actions
  // (bisect, accident report, analysis) are launched for.
  const mainParams = paramsList[0]
  const scaling = useSettingsStore().scaling
  const info = getInfo(mainParams, valueUnit, accidents)
  const series = filterUniqueByName(paramsList).map((param) => toDataSeries(param, valueUnit, scaling))

  const mainSeriesData = mainParams.value as OptionDataValue[]
  const unit = resolveMeasureUnitForDb(info.metricName ?? "", { storedType: mainSeriesData[3] as string, valueUnit, scaling })
  const value = series[0].rawValue
  const delta = findDeltaInData(mainSeriesData)
  let deltaPrevious: string | undefined
  let deltaNext: string | undefined
  let buildIdPrevious: number | undefined
  let buildIdNext: number | undefined
  let formattedPreviousValue: string | undefined
  let previousValue: number | undefined
  let nextValue: number | undefined
  if (delta != undefined) {
    if (delta.prev != null) {
      deltaPrevious = getDifferenceString(value, delta.prev, unit)
      buildIdPrevious = delta.prevBuildId
      formattedPreviousValue = formatMeasureValue(delta.prev, unit)
      previousValue = delta.prev
    }
    if (delta.next != null) {
      deltaNext = getDifferenceString(value, delta.next, unit)
      buildIdNext = delta.nextBuildId
      nextValue = delta.next
    }
  }
  return {
    ...info,
    deltaNext,
    deltaPrevious,
    series,
    chartDataUrl,
    buildIdPrevious,
    buildIdNext,
    formattedCurrentValue: series[0].value,
    formattedPreviousValue,
    previousValue,
    nextValue,
    metricType: info.type,
    seriesValues: seriesContext?.seriesValues,
    pointIndex: seriesContext?.pointIndex,
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

async function getOwnerFromMetaDb(project: string | undefined): Promise<string | null> {
  if (project == undefined) return null
  const url = ServerWithCompressConfigurator.DEFAULT_SERVER_URL + "/api/meta/ownerByProject?project=" + encodeURIComponent(project)
  const response = await fetch(url)
  if (!response.ok) return null
  const data = (await response.json()) as { owner?: string }
  return data.owner ?? null
}
