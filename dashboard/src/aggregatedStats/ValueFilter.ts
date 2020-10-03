// Copyright 2000-2020 JetBrains s.r.o. Use of this source code is governed by the Apache 2.0 license that can be found in the LICENSE file.
import { ChartSettings } from "@/aggregatedStats/ChartSettings"

export function asArray(value: string | Array<string> | null): Array<string> {
  return (value == null || value === "") ? [] : (Array.isArray(value) ? value : [value])
}

export interface ValueFilter {
  getValue(chartSettings: ChartSettings, values: Array<string>): Array<string>
}

// export class SingleValueFilter implements ValueFilter {
//   constructor(private name: "selectedProjects") {
//   }
//
//   getValue(chartSettings: ChartSettings, values: Array<string>): Array<string> {
//     let result = chartSettings[this.name]
//     let changed = false
//     if (values.length === 0) {
//       result = ""
//       changed = true
//     }
//     else if (result == null || result.length === 0 || !values.includes(result)) {
//       result = values[0]
//       changed = true
//     }
//
//     if (changed) {
//       chartSettings[this.name] = result
//     }
//     return [result]
//   }
// }

export class MultiValueFilter implements ValueFilter {
  constructor(private name: "selectedProjects" | "selectedMachine") {
  }

  getValue(chartSettings: ChartSettings, values: Array<string>): Array<string> {
    let result = chartSettings[this.name]
    // if someone save it as string value instead of singleton array
    if (result != null && !Array.isArray(result)) {
      result = result === "" ? [] : [result]
    }

    let changed: boolean
    if (values.length === 0) {
      result = []
      changed = true
    }
    else if (result == null || result.length === 0) {
      result = [values[0]]
      changed = true
    }
    else {
      const intersection = result.filter(it => values.includes(it))
      changed = intersection.length != result.length
      if (intersection.length === 0) {
        result = [values[0]]
      }
      else {
        result = intersection
      }
    }

    if (changed) {
      chartSettings[this.name] = result
    }
    return result
  }
}