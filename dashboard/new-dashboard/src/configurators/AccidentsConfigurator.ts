import { combineLatest, Observable } from "rxjs"
import { Ref, ref } from "vue"
import { Chart } from "../components/charts/DashboardCharts"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { DBType } from "../components/common/sideBar/InfoSidebar"
import { dbTypeStore } from "../shared/dbTypes"
import { ServerWithCompressConfigurator } from "./ServerWithCompressConfigurator"
import { TimeRange, TimeRangeConfigurator } from "./TimeRangeConfigurator"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

class AccidentFromServer {
  constructor(
    readonly id: number,
    readonly affectedTest: string,
    readonly date: string,
    readonly reason: string,
    readonly buildNumber: string,
    readonly kind: string,
    readonly stacktrace: string = ""
  ) {}
}

export enum AccidentKind {
  Regression = "Regression",
  Exception = "Exception",
  Improvement = "Improvement",
  Investigation = "Investigation",
  InferredImprovement = "InferredImprovement",
  InferredRegression = "InferredRegression",
}

export class Accident {
  constructor(
    readonly id: number,
    readonly affectedTest: string,
    readonly date: string,
    readonly reason: string,
    readonly buildNumber: string,
    readonly kind: AccidentKind,
    readonly stacktrace: string = ""
  ) {}
}

export abstract class AccidentsConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly value: Ref<Map<string, Accident[]> | undefined> = ref()
  protected dbType: DBType = DBType.UNKNOWN

  createObservable(): Observable<Map<string, Accident[]> | undefined> {
    return refToObservable(this.value)
  }

  configureFilter(_: DataQuery): boolean {
    return true
  }

  configureQuery(_: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    return true
  }

  protected abstract getAccidentUrl(): string

  async reloadAccidentData(id: number) {
    const response = await fetch(this.getAccidentUrl() + `accidents?id=${id}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    })

    if (!response.ok) {
      const errorMessage = `Cannot get accident by id ${id}. Response: ${response}`
      console.log(errorMessage)
      throw new Error(errorMessage)
    }

    const accident: Accident = await response.json()

    if (accident.id != undefined) {
      const updatedMap = new Map(this.value.value)
      for (const [_, value] of updatedMap) {
        const index = value.findIndex((obj) => obj.id === accident.id)
        if (index !== -1) {
          value.splice(index, 1, accident)
        }
      }
      this.value.value = updatedMap //we need to update value in reference to trigger the change
    }
  }

  async writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string | undefined, stacktrace: string = "") {
    try {
      let response = await fetch(this.getAccidentUrl() + "accidents/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ date, affected_test, reason, build_number: build_number.toString(), kind, stacktrace }),
      })

      if (!response.ok) {
        throw new Error("The accident wasn't created")
      }
      let idString: string = await response.text()
      const id = Number(idString)
      if (this.value.value == undefined) {
        this.value.value = new Map<string, Accident[]>()
      }
      const updatedMap = new Map(this.value.value)
      updatedMap.set(`${affected_test}_${build_number}`, [{ id, affectedTest: affected_test, date, reason, buildNumber: build_number, kind: kind as AccidentKind, stacktrace }])
      this.value.value = updatedMap //we need to update value in reference to trigger the change
      return id
    } catch (error) {
      console.error(error)
      return undefined
    }
  }

  async removeAccidentFromMetaDb(id: number) {
    const response = await fetch(this.getAccidentUrl() + "accidents/", {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id }),
    })

    if (!response.ok) {
      throw new Error("The accident wasn't deleted")
    }

    const updatedMap = new Map(this.value.value)
    for (const [key, value] of updatedMap) {
      if (value[0].id == id) {
        updatedMap.delete(key)
        break
      }
    }
    this.value.value = updatedMap //we need to update value in reference to trigger the change
  }

  public getAccidentsAroundDate(date: string): Promise<Accident[]> {
    return new Promise((resolve, reject) => {
      fetch(this.getAccidentUrl() + "accidentsAroundDate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ date }),
      })
        .then((response) => {
          if (!response.ok) {
            throw new Error("The accidents weren't fetched")
          }
          resolve(response.json())
        })
        .catch((error: unknown) => {
          console.error(error)
          reject(new Error("The accidents around date weren't fetched"))
        })
    })
  }

  public getAccidents(value: string[] | number[] | null): Accident[] | null {
    const accidents = this.value.value
    if (accidents != undefined && value != null) {
      let build = ""
      let key = ""
      let keyWithMetric = ""
      if (this.dbType == DBType.DEV_FLEET) {
        build = `${value[3]}`
        key = `${value[4]}_${build}`
        keyWithMetric = `${value[5]}/${value[2]}_${build}`
      }
      if (this.dbType == DBType.FLEET) {
        build = value[9] == 0 ? `${value[7]}.${value[8]}` : `${value[7]}.${value[8]}.${value[9]}`
        key = `${value[5]}_${build}`
        keyWithMetric = `${value[5]}/${value[2]}_${build}`
      }
      if (this.dbType == DBType.STARTUP_TESTS) {
        build = `${value[7]}.${value[8]}`
        key = `${value[5]}_${build}`
        keyWithMetric = `${value[5]}/${value[2]}_${build}`
      }
      if (this.dbType == DBType.INTELLIJ) {
        build = value[10] == 0 ? `${value[8]}.${value[9]}` : `${value[8]}.${value[9]}.${value[10]}`
        key = `${value[6]}_${build}`
        keyWithMetric = `${value[6]}/${value[2]}_${build}`
      }
      if (this.dbType == DBType.INTELLIJ_DEV || this.dbType == DBType.PERF_UNIT_TESTS || this.dbType == DBType.BAZEL) {
        build = `${value[5]}`
        key = `${value[6]}_${build}`
        keyWithMetric = `${value[6]}/${value[2]}_${build}`
      }
      if (this.dbType == DBType.STARTUP_TESTS_DEV) {
        build = `${value[4]}`
        key = `${value[5]}_${build}`
        keyWithMetric = `${value[5]}/${value[2]}_${build}`
      }
      const buildAccident = accidents.get(`_${build}`) ?? []
      const testAccident = accidents.get(key) ?? []
      const metricAccident = accidents.get(keyWithMetric) ?? []
      return [...testAccident, ...buildAccident, ...metricAccident]
    }
    return null
  }
}

function combineProjectsAndMetrics(projects: string | string[] | null, measures: string | string[] | null): string[] {
  const projectAndMetrics: string[] = []
  if (projects != null && measures != null) {
    if (Array.isArray(projects)) {
      projectAndMetrics.push(...projects)
    } else {
      projectAndMetrics.push(projects)
    }

    if (Array.isArray(projects)) {
      if (Array.isArray(measures)) {
        projectAndMetrics.push(...projects.map((project) => measures.map((metric) => `${project}/${metric}`)).flat(100))
      } else {
        projectAndMetrics.push(...projects.map((project) => `${project}/${measures}`))
      }
    } else {
      if (Array.isArray(measures)) {
        projectAndMetrics.push(...measures.map((metric) => `${projects}/${metric}`))
      } else {
        projectAndMetrics.push(`${projects}/${measures}`)
      }
    }
  }
  return projectAndMetrics
}

export class AccidentsConfiguratorForStartup extends AccidentsConfigurator {
  constructor(
    private serverUrl: string,
    private product: Ref<string | string[] | null>,
    projects: Ref<string | string[] | null>,
    metrics: Ref<string[] | string | null>,
    timeRangeConfigurator: TimeRangeConfigurator
  ) {
    super()
    this.dbType = dbTypeStore().dbType
    combineLatest([refToObservable(projects), refToObservable(metrics), timeRangeConfigurator.createObservable(), refToObservable(product)]).subscribe(
      ([projects, measures, [timeRange, customRange], product]) => {
        if (product == null) return
        if (Array.isArray(product)) return
        const projectAndMetrics = combineProjectsAndMetrics(projects, measures)
        const projectAndMetricsWithProduct = projectAndMetrics.map((it) => `${product}/${it}`)
        getAccidentsFromMetaDb(projectAndMetricsWithProduct, timeRange, customRange)
          .then((value) => {
            this.value.value = this.removeProductPrefix(product, value)
          })
          .catch((error: unknown) => {
            console.error(error)
          })
      }
    )
  }

  protected getAccidentUrl(): string {
    return this.serverUrl + "/api/meta/"
  }

  private removeProductPrefix(product: string, response: Map<string, Accident[]>): Map<string, Accident[]> {
    const map = new Map<string, Accident[]>()
    for (const [key, value] of response) {
      const keyWithoutProduct = key.replace(`${product}/`, "")
      map.set(keyWithoutProduct, value)
    }
    return map
  }

  async writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string | undefined, stacktrace: string = "") {
    if (this.product.value == null || Array.isArray(this.product.value)) return Promise.resolve(undefined)
    const test = `${this.product.value}/${affected_test}`
    try {
      let response = await fetch(this.getAccidentUrl() + "accidents/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ date, affected_test: test, reason, build_number: build_number.toString(), kind, stacktrace }),
      })

      if (!response.ok) {
        throw new Error("The accident wasn't created")
      }
      let idString: string = await response.text()
      const id = Number(idString)
      if (this.value.value == undefined) {
        this.value.value = new Map<string, Accident[]>()
      }
      const updatedMap = new Map(this.value.value)
      updatedMap.set(`${affected_test}_${build_number}`, [{ id, affectedTest: affected_test, date, reason, buildNumber: build_number, kind: kind as AccidentKind, stacktrace }])
      this.value.value = updatedMap //we need to update value in reference to trigger the change
      return id
    } catch (error) {
      console.error(error)
      return undefined
    }
  }
}

export class AccidentsConfiguratorForTests extends AccidentsConfigurator {
  constructor(
    private serverUrl: string,
    projects: Ref<string | string[] | null>,
    metrics: Ref<string[] | string | null>,
    timeRangeConfigurator: TimeRangeConfigurator
  ) {
    super()
    this.dbType = dbTypeStore().dbType
    combineLatest([refToObservable(projects), refToObservable(metrics), timeRangeConfigurator.createObservable()]).subscribe(([projects, measures, [timeRange, customRange]]) => {
      const projectAndMetrics = combineProjectsAndMetrics(projects, measures)
      getAccidentsFromMetaDb(projectAndMetrics, timeRange, customRange)
        .then((value) => {
          this.value.value = value
        })
        .catch((error: unknown) => {
          console.error(error)
        })
    })
  }

  protected getAccidentUrl(): string {
    return this.serverUrl + "/api/meta/"
  }
}

export class AccidentsConfiguratorForDashboard extends AccidentsConfigurator {
  constructor(
    private serverUrl: string,
    charts: Chart[] | null,
    timeRangeConfigurator: TimeRangeConfigurator
  ) {
    super()
    this.dbType = dbTypeStore().dbType
    const tests = this.getProjectAndProjectWithMetrics(charts)
    combineLatest([timeRangeConfigurator.createObservable()]).subscribe(([[timeRange, customRange]]) => {
      getAccidentsFromMetaDb(tests, timeRange, customRange)
        .then((value) => {
          this.value.value = value
        })
        .catch((error: unknown) => {
          console.error(error)
        })
    })
  }

  protected getAccidentUrl(): string {
    return this.serverUrl + "/api/meta/"
  }

  private getProjectAndProjectWithMetrics(charts: Chart[] | null): string[] {
    const projectsWithMetrics =
      charts?.flatMap((chart) => {
        const measures = Array.isArray(chart.definition.measure) ? chart.definition.measure : [chart.definition.measure]
        return chart.projects.flatMap((project) => {
          return measures.map((measure) => project + "/" + measure)
        })
      }) ?? []
    const projects = new Set(charts?.map((it) => it.projects).flat(Number.POSITIVE_INFINITY) as string[])
    return [...projectsWithMetrics, ...projects]
  }
}

function intervalToPostgresInterval(interval: TimeRange, customRange: string): string {
  const intervalMapping: Record<TimeRange, string> = {
    "1w": "7 DAYS",
    "2w": "14 DAYS",
    "1M": "1 MONTH",
    "3M": "3 MONTHS",
    "1y": "1 YEAR",
    all: "100 YEAR",
    custom: "",
  }
  if (customRange != "") {
    const delimiter = customRange.indexOf(":")
    let days = 365
    if (delimiter > 0) {
      days = getDaysDifference(customRange.slice(0, delimiter))
    }
    return days.toString() + " DAYS"
  }
  return intervalMapping[interval]
}

function getDaysDifference(dateString: string) {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const givenDate = new Date(dateString)
  givenDate.setHours(0, 0, 0, 0)

  const differenceInTime = today.getTime() - givenDate.getTime()
  const differenceInDays = differenceInTime / (1000 * 3600 * 24)

  return Math.abs(differenceInDays)
}

/**
 * This is needed for optimization since we search for accidents on each point on the plot.
 * @param accidents
 */
function convertAccidentsToMap(accidents: Accident[]): Map<string, Accident[]> {
  const accidentsMap = new Map<string, Accident[]>()
  for (const accident of accidents) {
    const key = `${accident.affectedTest}_${accident.buildNumber}` // assuming accident has a property 'value8'
    if (accidentsMap.get(key) == null) {
      accidentsMap.set(key, [accident])
    } else {
      accidentsMap.get(key)?.push(accident)
    }
  }
  return accidentsMap
}

async function getAccidentsFromMetaDb(tests: string[], timeRange: TimeRange, customRange: string): Promise<Map<string, Accident[]>> {
  const interval = intervalToPostgresInterval(timeRange, customRange)
  const params = tests.length === 0 ? { interval } : { interval, tests }
  try {
    const response = await fetch(ServerWithCompressConfigurator.DEFAULT_SERVER_URL + "/api/meta/getAccidents", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(params),
    })
    const data: AccidentFromServer[] = (await response.json()) as AccidentFromServer[]
    const accidents = data.map((value) => {
      return { ...value, kind: capitalizeFirstLetter(value.kind) }
    })
    return convertAccidentsToMap(accidents)
  } catch (error) {
    console.error(error)
    return new Map<string, Accident[]>()
  }
}

function capitalizeString(str: string): string {
  // Check if the string is all uppercase
  if (str === str.toUpperCase()) {
    // If all uppercase, capitalize the first letter, and make the rest lowercase
    return str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()
  }

  // If not all uppercase, handle as camel case
  return str
    .replaceAll(/([a-z])([A-Z])/g, "$1 $2") // Add space before each capital in a camel case word
    .split(" ") // Split the string into words
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase()) // Capitalize the first letter of each word
    .join("") // Rejoin the words without spaces
}

function capitalizeFirstLetter(str: string): AccidentKind {
  const result = capitalizeString(str)
  if (isAccidentKind(result)) {
    return result
  } else {
    throw new Error("Unsupported AccidentKind " + str)
  }
}

function isAccidentKind(str: string): str is AccidentKind {
  return Object.values(AccidentKind).includes(str as AccidentKind)
}
