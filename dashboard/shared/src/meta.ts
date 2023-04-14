import { Ref } from "vue"
import { ServerConfigurator } from "./configurators/ServerConfigurator"
import { TimeRange } from "./configurators/TimeRangeConfigurator"
import { encodeRison } from "./rison"

const url = ServerConfigurator.DEFAULT_SERVER_URL + "/api/meta/"

export class Accident {
  constructor(readonly affectedTest: string, readonly date: string, readonly reason: string, readonly buildNumber: string) {}
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

export function writeRegressionToMetaDb(date: string, affected_test: string, reason: string, build_number: string) {
  fetch(url, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({date, affected_test, reason, build_number: build_number.toString()}),
  }).catch(error => console.error(error))
}

export function getWarningFromMetaDb(warnings: Ref<Array<Accident> | undefined>, tests: Array<string> | string | null, timeRange: TimeRange) {
  if (tests != null && !Array.isArray(tests)) {
    tests = [tests]
  }
  warnings.value = []
  const interval = intervalToPostgresInterval(timeRange)
  const params = tests == null ? {interval} : {interval, tests}
  fetch(url + encodeRison(params))
    .then(response => response.json())
    .then((data: Array<Accident>) => {
      if(data != null) {
        warnings.value?.push(...data)
      }
    })
    .catch(error => console.error(error))
}

export function isValueShouldBeMarked(accidents: Array<Accident> | null, value: Array<string>): boolean {
  return getAccident(accidents, value) != null
}

export function getAccident(accidents: Array<Accident> | null, value: Array<string>): Accident | null {
  if (accidents != null) {
    //perf db
    if (value.length == 10) {
      return accidents.find(accident => accident.affectedTest == value[5] && accident.buildNumber == value[7] + "." + value[8]) ?? null
    }
    //perf dev db
    if (value.length == 6) {
      return accidents.find(accident => accident.affectedTest == value[5] && accident.buildNumber == value[4]) ?? null
    }
  }
  return null
}