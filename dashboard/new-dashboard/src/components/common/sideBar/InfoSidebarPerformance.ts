import { CallbackDataParams, OptionDataValue } from "echarts/types/src/util/types"
import { Observable } from "rxjs"
import { computed, ref, Ref } from "vue"
import { ServerConfigurator } from "../../../configurators/ServerConfigurator"
import { Accident, Description, getDescriptionFromMetaDb } from "../../../util/meta"
import { DataQueryExecutor } from "../DataQueryExecutor"
import { ValueUnit } from "../chart"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration, SimpleQueryProducer } from "../dataQuery"
import { durationAxisPointerFormatter, nsToMs, timeFormatWithoutSeconds } from "../formatter"
import { buildUrl, DataSeries, DBType, InfoData } from "./InfoSidebar"

export interface InfoDataPerformance extends DataSeries, InfoData {
  accidents: Ref<Accident[] | undefined> | undefined
  description: Description | null
  buildType: string | undefined
}

class AdditionalData {
  constructor(
    public readonly machine: string,
    public readonly tc_build_type?: string,
    public readonly tc_installer_build_id?: number
  ) {}
}

function getAdditionalData(serverConfigurator: ServerConfigurator, buildID: number, hasInstaller: boolean, hasBuildType: boolean): Promise<AdditionalData> {
  return new Promise((resolve) => {
    new DataQueryExecutor([
      serverConfigurator,
      new (class implements DataQueryConfigurator {
        configureQuery(query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
          configuration.queryProducers.push(new SimpleQueryProducer())
          query.addField("machine")
          if (hasInstaller) {
            query.addField("tc_installer_build_id")
          }
          if (hasBuildType) {
            query.addField("tc_build_type")
          }
          query.addFilter({ f: "tc_build_id", v: buildID })
          query.order = "tc_build_id"
          return true
        }

        createObservable(): Observable<unknown> | null {
          return null
        }
      })(),
    ]).subscribe((data, _configuration, isLoading) => {
      if (isLoading || data == null) {
        return
      }
      const dataFlattened = data.flat(1)
      const transposed = dataFlattened[0].map((_, i) => dataFlattened.map((row) => row[i]))[0]
      const [machine, tc_installer_build_id, tc_build_type] = transposed
      let result: AdditionalData
      if (hasInstaller) {
        result = hasBuildType
          ? new AdditionalData(machine as string, tc_build_type as string, tc_installer_build_id as number)
          : new AdditionalData(machine as string, undefined, tc_installer_build_id as number)
      } else {
        result = hasBuildType ? new AdditionalData(machine as string, tc_build_type as string) : new AdditionalData(machine as string)
      }
      resolve(result)
    })
  })
}

export async function getInfoDataFrom(
  serverConfigurator: ServerConfigurator,
  dbType: DBType,
  params: CallbackDataParams,
  valueUnit: ValueUnit,
  accidents: Ref<Accident[]> | null
): Promise<InfoDataPerformance> {
  const dataSeries = params.value as OptionDataValue[]
  const dateMs = dataSeries[0] as number
  const value: number = dataSeries[1] as number
  let projectName: string = params.seriesName as string
  let metricName: string | undefined
  let buildId: number | undefined
  let type: ValueUnit | undefined = valueUnit
  let buildNumber: string | undefined
  let buildVersion: number | undefined
  let buildNum1: number | undefined
  let buildNum2: number | undefined
  const accidentBuild: Ref<string | undefined> = ref()
  if (dbType == DBType.DEV_FLEET) {
    buildId = dataSeries[2] as number
    projectName = dataSeries[3] as string
  }
  if (dbType == DBType.INTELLIJ_DEV) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    buildId = dataSeries[4] as number
    projectName = dataSeries[5] as string
    accidentBuild.value = buildId.toString()
  }
  if (dbType == DBType.FLEET) {
    buildId = dataSeries[2] as number
    projectName = dataSeries[3] as string
    buildVersion = dataSeries[4] as number
    buildNum1 = dataSeries[5] as number
    buildNum2 = dataSeries[6] as number
  }
  if (dbType == DBType.JBR) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    buildId = dataSeries[4] as number
    projectName = dataSeries[5] as string
    buildNumber = dataSeries[6] as string
  }
  if (dbType == DBType.INTELLIJ) {
    metricName = dataSeries[2] as string
    if (dataSeries[3] == "c") {
      type = "counter"
    }
    buildId = dataSeries[4] as number
    projectName = dataSeries[5] as string
    buildVersion = dataSeries[6] as number
    buildNum1 = dataSeries[7] as number
    buildNum2 = dataSeries[8] as number
    accidentBuild.value = `${buildVersion}.${buildNum1}`
  }
  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }

  const hasInstaller = dbType == DBType.INTELLIJ || dbType == DBType.FLEET
  const hasBuildType = dbType == DBType.INTELLIJ || dbType == DBType.INTELLIJ_DEV

  const additionalData = await getAdditionalData(serverConfigurator, buildId as number, hasInstaller, hasBuildType)

  const installerId = additionalData.tc_installer_build_id
  const installerUrl = installerId == undefined ? undefined : `${buildUrl(installerId)}&tab=artifacts`
  const artifactsUrl = `${buildUrl(buildId as number)}&tab=artifacts`
  const changesUrl = installerId == undefined ? `${buildUrl(buildId as number)}&tab=changes` : `${buildUrl(installerId)}&tab=changes`
  const fullBuildId = buildVersion == undefined ? buildNumber : `${buildVersion}.${buildNum1}${buildNum2 == 0 ? "" : `.${buildNum2}`}`
  const machineName = additionalData.machine

  if (dbType == DBType.INTELLIJ) {
    accidentBuild.value = `${buildVersion}.${buildNum1}`
  }
  if (dbType == DBType.INTELLIJ_DEV && buildId) {
    accidentBuild.value = buildId.toString()
  }
  if (dbType == DBType.JBR && buildNumber) {
    accidentBuild.value = buildNumber
  }

  if (dbType == DBType.UNKNOWN) {
    console.error("Unknown type of DB")
  }

  let showValue = value.toString()
  if (type != "counter") {
    showValue = durationAxisPointerFormatter(valueUnit == "ns" ? nsToMs(value) : value)
  }

  const filteredAccidents = computed(() => {
    return accidents?.value.filter((accident: Accident) => accident.affectedTest == projectName && accident.buildNumber == accidentBuild.value)
  })

  const description = await getDescriptionFromMetaDb(projectName, "master")

  const buildType = additionalData.tc_build_type

  return {
    build: fullBuildId,
    artifactsUrl,
    changesUrl,
    installerUrl,
    color: params.color as string,
    date: timeFormatWithoutSeconds.format(dateMs),
    value: showValue,
    machineName,
    projectName,
    title: "Details",
    installerId,
    accidents: filteredAccidents,
    buildId: buildId as number,
    description,
    metricName,
    buildType,
  }
}
