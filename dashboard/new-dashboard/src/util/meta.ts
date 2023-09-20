import { Ref } from "vue"
import { encodeRison } from "../components/common/rison"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { TimeRange } from "../configurators/TimeRangeConfigurator"

const accidents_url = ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/accidents/"

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

export class AccidentFromServer {
  constructor(
    readonly id: number,
    readonly affectedTest: string,
    readonly date: string,
    readonly reason: string,
    readonly buildNumber: string,
    readonly kind: string
  ) {}
}

function intervalToPostgresInterval(interval: TimeRange): string {
  const intervalMapping = {
    "1w": "7 DAYS",
    "1M": "1 MONTH",
    "3M": "3 MONTHS",
    "1y": "1 YEAR",
    all: "100 YEAR",
  }
  return intervalMapping[interval]
}

export function getAccidentTypes(): string[] {
  return Object.values(AccidentKind)
}

export function removeAccidentFromMetaDb(id: number) {
  fetch(accidents_url, {
    method: "DELETE",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ id }),
  })
    .then((_) => {
      window.location.reload()
    })
    .catch((error) => {
      console.error(error)
    })
}

export function writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string | undefined) {
  fetch(accidents_url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ date, affected_test, reason, build_number: build_number.toString(), kind }),
  })
    .then((_) => {
      window.location.reload()
    })
    .catch((error) => {
      console.error(error)
    })
}

export function getAccidentsFromMetaDb(tests: string[] | string | null, timeRange: Ref<TimeRange>): Promise<Accident[]> {
  if (tests != null && !Array.isArray(tests)) {
    tests = [tests]
  }
  const interval = intervalToPostgresInterval(timeRange.value)
  const params = tests == null ? { interval } : { interval, tests }
  return fetch(ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/getAccidents", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params),
  })
    .then((response) => response.json())
    .then((data: AccidentFromServer[]) => {
      return data.map((value) => {
        return { ...value, kind: capitalizeFirstLetter(value.kind) }
      })
    })
    .catch((error) => {
      console.error(error)
      return []
    })
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

const description_url = ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/description/"

export class Description {
  constructor(
    readonly project: string,
    readonly branch: string,
    readonly url: string,
    readonly methodName: string,
    readonly description: string
  ) {}
}

export async function getDescriptionFromMetaDb(project: string | undefined, branch: string): Promise<Description | null> {
  const response = await fetch(description_url + encodeRison({ project, branch }))
  return response.ok ? response.json() : null
}

/**
 * This is needed for optimization since we search for accidents on each point on the plot.
 * @param accidents
 */
export function convertAccidentsToMap(accidents: Accident[] | undefined | null): Map<string, Accident> {
  const accidentsMap = new Map<string, Accident>()
  if (accidents) {
    for (const accident of accidents) {
      const key = `${accident.affectedTest}_${accident.buildNumber}` // assuming accident has a property 'value8'
      accidentsMap.set(key, accident)
    }
  }
  return accidentsMap
}

export function isValueShouldBeMarkedWithPin(accidents: Map<string, Accident> | undefined, value: string[]): boolean {
  const accident = getAccident(accidents, value)
  return accident != null && accident.kind != AccidentKind.Exception
}

export function getAccident(accidents: Map<string, Accident> | undefined, value: string[] | null): Accident | null {
  if (accidents != null) {
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
