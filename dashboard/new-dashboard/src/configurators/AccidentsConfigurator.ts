import { combineLatest, Observable } from "rxjs"
import { Ref, ref } from "vue"
import { Chart } from "../components/charts/DashboardCharts"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { ServerConfigurator } from "./ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "./TimeRangeConfigurator"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export const accidents_url = ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/accidents/"

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

export class AccidentsConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly value: Ref<Map<string, Accident[]> | undefined> = ref()

  createObservable(): Observable<Map<string, Accident[]> | undefined> {
    return refToObservable(this.value)
  }

  configureFilter(_: DataQuery): boolean {
    return true
  }

  configureQuery(_: DataQuery, _configuration: DataQueryExecutorConfiguration | null): boolean {
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
}

export class AccidentsConfiguratorForTests extends AccidentsConfigurator {
  constructor(projects: Ref<string | string[] | null>, metrics: Ref<string[] | string | null>, timeRangeConfigurator: TimeRangeConfigurator) {
    super()
    combineLatest([refToObservable(projects), refToObservable(metrics), timeRangeConfigurator.createObservable(), timeRangeConfigurator.createCustomRangeObservable()]).subscribe(
      ([projects, measures, timeRange, customRange]) => {
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
        getAccidentsFromMetaDb(projectAndMetrics, timeRange, customRange)
          .then((value) => {
            this.value.value = value
          })
          .catch((error) => {
            console.error(error)
          })
      }
    )
  }
}

export class AccidentsConfiguratorForDashboard extends AccidentsConfigurator {
  constructor(charts: Chart[] | null, timeRangeConfigurator: TimeRangeConfigurator) {
    super()
    const tests = this.getProjectAndProjectWithMetrics(charts)
    combineLatest([timeRangeConfigurator.createObservable(), timeRangeConfigurator.createCustomRangeObservable()]).subscribe(([timeRange, customRange]) => {
      console.log(timeRange)
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
    const response = await fetch(ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/getAccidents", {
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

function capitalizeFirstLetter(str: string): AccidentKind {
  const result = str.charAt(0).toUpperCase() + str.slice(1).toLowerCase()
  if (isAccidentKind(result)) {
    return result
  } else {
    throw new Error("Unsupported AccidentKind")
  }
}

function isAccidentKind(str: string): str is AccidentKind {
  return Object.values(AccidentKind).includes(str as AccidentKind)
}

export function getAccidents(accidents: Map<string, Accident[]> | undefined, value: string[] | null): Accident[] | null {
  if (accidents != undefined) {
    //perf db
    if (value?.length == 12 || value?.length == 11) {
      const key = `${value[6]}_${value[8]}.${value[9]}`
      const keyWithMetric = `${value[6]}/${value[2]}_${value[8]}.${value[9]}`
      return accidents.get(key) ?? accidents.get(keyWithMetric) ?? null
    }
    //perf dev db
    if (value?.length == 8 || value?.length == 7) {
      const key = `${value[6]}_${value[5]}`
      const keyWithMetric = `${value[6]}/${value[2]}_${value[5]}`
      return accidents.get(key) ?? accidents.get(keyWithMetric) ?? null
    }
  }
  return null
}
