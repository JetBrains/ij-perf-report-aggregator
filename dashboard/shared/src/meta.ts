import { Ref } from "vue"
import { ServerConfigurator } from "./configurators/ServerConfigurator"
import { TimeRange } from "./configurators/TimeRangeConfigurator"
import { encodeRison } from "./rison"

const accidents_url = ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/accidents/"


export enum AccidentKind {
  Regression = "Regression",
  Exception = "Exception",
  Improvement = "Improvement",
}

export class Accident {
  constructor(readonly id: number, readonly affectedTest: string, readonly date: string, readonly reason: string, readonly buildNumber: string, readonly kind: AccidentKind) {}
}

export class AccidentFromServer {
  constructor(readonly id: number, readonly affectedTest: string, readonly date: string, readonly reason: string, readonly buildNumber: string, readonly kind: string) {}
}

function intervalToPostgresInterval(interval: TimeRange): string {
  const intervalMapping = {
    "1w": "7 DAYS",
    "1M": "1 MONTH",
    "3M": "3 MONTHS",
    "1y": "1 YEAR",
    "all": "100 YEAR",
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
    body: JSON.stringify({id}),
  }).then(
    _ => window.location.reload(),
  ).catch(error => console.error(error))
}

export function writeAccidentToMetaDb(date: string, affected_test: string, reason: string, build_number: string, kind: string|undefined) {
  fetch(accidents_url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({date, affected_test, reason, build_number: build_number.toString(), kind}),
  }).then(
    _ => window.location.reload(),
  ).catch(error => console.error(error))
}

export function getAccidentsFromMetaDb(accidents: Ref<Accident[] | undefined>, tests: string[] | string | null, timeRange: TimeRange) {
  if (tests != null && !Array.isArray(tests)) {
    tests = [tests]
  }
  accidents.value = []
  const interval = intervalToPostgresInterval(timeRange)
  const params = tests == null ? {interval} : {interval, tests}
  fetch(accidents_url + encodeRison(params))
    .then(response => response.json())
    .then((data: AccidentFromServer[]) => {
      if (data != null) {
        const mappedData = data.map(value => {
          return {...value, kind: capitalizeFirstLetter(value.kind)}
        })
        accidents.value?.push(...mappedData)
      }
    })
    .catch(error => console.error(error))
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
  constructor(readonly project: string, readonly branch: string, readonly url: string, readonly methodName: string, readonly description: string) {}
}

export function getDescriptionFromMetaDb(descriptionRef: Ref<Description|undefined>, project: string | undefined, branch: string) {
  if (project != undefined && branch != undefined) {
    fetch(description_url + encodeRison({project, branch}))
      .then(response => {
        return response.ok ? response.json() : null
      })
      .then((data: Description) => {
        if (data != null) {
          descriptionRef.value = data
        }
      })
      .catch(error => console.error(error))
  }
}

/**
 * This is needed for optimization since we search for accidents on each point on the plot.
 * @param accidents
 */
export function convertAccidentsToMap(accidents: Accident[] | null): Map<string, Accident> {
  const accidentsMap = new Map<string, Accident>()
  if (accidents) {
    for (const accident of accidents) {
      const key = `${accident.affectedTest}_${accident.buildNumber}` // assuming accident has a property 'value8'
      accidentsMap.set(key, accident)
    }
  }
  return accidentsMap
}

export function isValueShouldBeMarkedWithPin(accidents: Map<string, Accident> | null, value: string[]): boolean {
  const accident = getAccident(accidents, value)
  return accident != null && accident.kind != AccidentKind.Exception
}

export function getAccident(accidents: Map<string, Accident> | null, value: string[]): Accident | null {
  if (accidents != null) {
    //perf db
    if (value.length == 10) {
      const key = `${value[5]}_${value[7]}.${value[8]}`
      return accidents.get(key) ?? null
    }
    //perf dev db
    if (value.length == 6) {
      const key = `${value[5]}_${value[4]}`
      return accidents.get(key) ?? null
    }
  }
  return null
}