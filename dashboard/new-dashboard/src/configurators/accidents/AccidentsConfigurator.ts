import { Observable } from "rxjs"
import { Ref, ref } from "vue"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../components/common/dataQuery"
import { DBType } from "../../components/common/sideBar/InfoSidebar"
import { ServerWithCompressConfigurator } from "../ServerWithCompressConfigurator"
import { TimeRange } from "../TimeRangeConfigurator"
import { FilterConfigurator } from "../filter"
import { refToObservable } from "../rxjs"
import { useUserStore } from "../../shared/useUserStore"

class AccidentFromServer {
  constructor(
    readonly id: number,
    readonly affectedTest: string,
    readonly date: string,
    readonly reason: string,
    readonly buildNumber: string,
    readonly kind: string,
    readonly stacktrace: string = "",
    readonly userName: string = ""
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
    readonly stacktrace: string = "",
    readonly userName: string = ""
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
      const errorMessage = `Cannot get accident by id ${id}`
      console.error(errorMessage)
      throw new Error(errorMessage)
    }

    const accident = (await response.json()) as Accident

    const updatedMap = new Map(this.value.value)
    for (const [_, value] of updatedMap) {
      const index = value.findIndex((obj) => obj.id === accident.id)
      if (index !== -1) {
        value.splice(index, 1, accident)
      }
    }
    this.value.value = updatedMap //we need to update value in reference to trigger the change
  }

  async writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string | undefined, stacktrace: string = "") {
    const userName = useUserStore().user?.name ?? ""
    const response = await fetch(this.getAccidentUrl() + "accidents/", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ date, affected_test, reason, build_number: build_number.toString(), kind, stacktrace, user_name: userName }),
    })

    if (!response.ok) {
      if (response.status === 409) {
        throw new Error("An accident for this test and build already exists")
      }
      throw new Error(`The accident wasn't created (HTTP ${response.status}: ${response.statusText})`)
    }
    const idString: string = await response.text()
    const id = Number(idString)
    this.value.value ??= new Map<string, Accident[]>()
    const updatedMap = new Map(this.value.value)
    updatedMap.set(`${affected_test}_${build_number}`, [
      { id, affectedTest: affected_test, date, reason, buildNumber: build_number, kind: kind as AccidentKind, stacktrace, userName },
    ])
    this.value.value = updatedMap //we need to update value in reference to trigger the change
    return id
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
      if (
        this.dbType == DBType.INTELLIJ_DEV ||
        this.dbType == DBType.PERF_UNIT_TESTS ||
        this.dbType == DBType.BAZEL ||
        this.dbType == DBType.FLEET_PERF ||
        this.dbType == DBType.DIOGEN
      ) {
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

  combineProjectsAndMetrics(projects: string | string[] | null, measures: string | string[] | null): string[] {
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

  async getAccidentsFromMetaDb(tests: string[], timeRange: TimeRange, customRange: string): Promise<Map<string, Accident[]>> {
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
        return new Accident(value.id, value.affectedTest, value.date, value.reason, value.buildNumber, capitalizeFirstLetter(value.kind), value.stacktrace, value.userName)
      })
      return convertAccidentsToMap(accidents)
    } catch (error) {
      console.error(error)
      return new Map<string, Accident[]>()
    }
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
