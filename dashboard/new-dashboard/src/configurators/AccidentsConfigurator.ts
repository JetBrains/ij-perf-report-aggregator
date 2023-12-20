import { combineLatest, Observable } from "rxjs"
import { Ref, ref } from "vue"
import { Chart } from "../components/charts/DashboardCharts"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { DBType } from "../components/common/sideBar/InfoSidebar"
import { ServerWithCompressConfigurator } from "./ServerWithCompressConfigurator"
import { TimeRange, TimeRangeConfigurator } from "./TimeRangeConfigurator"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export const accidents_url = ServerWithCompressConfigurator.DEFAULT_SERVER_URL + "/api/meta/accidents/"

class AccidentFromServer {
  constructor(
    readonly id: number,
    readonly affectedTest: string,
    readonly date: string,
    readonly reason: string,
    readonly buildNumber: string,
    readonly kind: string
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
    readonly kind: AccidentKind
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

  writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string | undefined) {
    fetch(accidents_url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ date, affected_test, reason, build_number: build_number.toString(), kind }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("The accident wasn't created")
        }
        return response.text()
      })
      .then((idString) => {
        const id = Number(idString)
        if (this.value.value == undefined) {
          this.value.value = new Map<string, Accident[]>()
        }
        const updatedMap = new Map(this.value.value)
        updatedMap.set(`${affected_test}_${build_number}`, [{ id, affectedTest: affected_test, date, reason, buildNumber: build_number, kind: kind as AccidentKind }])
        this.value.value = updatedMap //we need to update value in reference to trigger the change
      })
      .catch((error) => {
        console.error(error)
      })
  }

  removeAccidentFromMetaDb(id: number) {
    fetch(accidents_url, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("The accident wasn't deleted")
        }
      })
      .then((_) => {
        const updatedMap = new Map(this.value.value)
        for (const [key, value] of updatedMap) {
          if (value[0].id == id) {
            updatedMap.delete(key)
            break
          }
        }
        this.value.value = updatedMap //we need to update value in reference to trigger the change
      })
      .catch((error) => {
        console.error(error)
      })
  }

  public getAccidents(value: string[] | null): Accident[] | null {
    const accidents = this.value.value
    if (accidents != undefined && value != null) {
      if (this.dbType == DBType.STARTUP_TESTS) {
        const key = `${value[5]}_${value[7]}.${value[8]}`
        const keyWithMetric = `${value[5]}/${value[2]}_${value[7]}.${value[8]}`
        return accidents.get(key) ?? accidents.get(keyWithMetric) ?? null
      }
      if (this.dbType == DBType.INTELLIJ) {
        const key = `${value[6]}_${value[8]}.${value[9]}`
        const keyWithMetric = `${value[6]}/${value[2]}_${value[8]}.${value[9]}`
        return accidents.get(key) ?? accidents.get(keyWithMetric) ?? null
      }
      if (this.dbType == DBType.INTELLIJ_DEV || this.dbType == DBType.PERF_UNIT_TESTS) {
        const key = `${value[6]}_${value[5]}`
        const keyWithMetric = `${value[6]}/${value[2]}_${value[5]}`
        return accidents.get(key) ?? accidents.get(keyWithMetric) ?? null
      }
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
    private product: string,
    projects: Ref<string | string[] | null>,
    metrics: Ref<string[] | string | null>,
    timeRangeConfigurator: TimeRangeConfigurator
  ) {
    super()
    this.dbType = DBType.STARTUP_TESTS
    combineLatest([refToObservable(projects), refToObservable(metrics), timeRangeConfigurator.createObservable()]).subscribe(([projects, measures, [timeRange, customRange]]) => {
      const projectAndMetrics = combineProjectsAndMetrics(projects, measures)
      const projectAndMetricsWithProduct = projectAndMetrics.map((it) => `${product}/${it}`)
      getAccidentsFromMetaDb(projectAndMetricsWithProduct, timeRange, customRange)
        .then((value) => {
          this.value.value = this.removeProductPrefix(product, value)
        })
        .catch((error) => {
          console.error(error)
        })
    })
  }

  private removeProductPrefix(product: string, response: Map<string, Accident[]>): Map<string, Accident[]> {
    const map = new Map<string, Accident[]>()
    for (const [key, value] of response) {
      const keyWithoutProduct = key.replace(`${product}/`, "")
      map.set(keyWithoutProduct, value)
    }
    return map
  }

  writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string | undefined) {
    const test = `${this.product}/${affected_test}`
    fetch(accidents_url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ date, test, reason, build_number: build_number.toString(), kind }),
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("The accident wasn't created")
        }
        return response.text()
      })
      .then((idString) => {
        const id = Number(idString)
        if (this.value.value == undefined) {
          this.value.value = new Map<string, Accident[]>()
        }
        const updatedMap = new Map(this.value.value)
        updatedMap.set(`${affected_test}_${build_number}`, [{ id, affectedTest: affected_test, date, reason, buildNumber: build_number, kind: kind as AccidentKind }])
        this.value.value = updatedMap //we need to update value in reference to trigger the change
      })
      .catch((error) => {
        console.error(error)
      })
  }
}

export class AccidentsConfiguratorForTests extends AccidentsConfigurator {
  constructor(projects: Ref<string | string[] | null>, metrics: Ref<string[] | string | null>, timeRangeConfigurator: TimeRangeConfigurator, dbType: DBType) {
    super()
    this.dbType = dbType
    combineLatest([refToObservable(projects), refToObservable(metrics), timeRangeConfigurator.createObservable()]).subscribe(([projects, measures, [timeRange, customRange]]) => {
      const projectAndMetrics = combineProjectsAndMetrics(projects, measures)
      getAccidentsFromMetaDb(projectAndMetrics, timeRange, customRange)
        .then((value) => {
          this.value.value = value
        })
        .catch((error) => {
          console.error(error)
        })
    })
  }
}

export class AccidentsConfiguratorForDashboard extends AccidentsConfigurator {
  constructor(charts: Chart[] | null, timeRangeConfigurator: TimeRangeConfigurator, dbType: DBType) {
    super()
    this.dbType = dbType
    const tests = this.getProjectAndProjectWithMetrics(charts)
    combineLatest([timeRangeConfigurator.createObservable()]).subscribe(([[timeRange, customRange]]) => {
      getAccidentsFromMetaDb(tests, timeRange, customRange)
        .then((value) => {
          this.value.value = value
        })
        .catch((error) => {
          console.error(error)
        })
    })
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
    return days + " DAYS"
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
