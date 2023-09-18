import { useStorage } from "@vueuse/core"
import { Observable } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export class ScalingConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly value = useStorage("scalingEnabled", false)

  createObservable(): Observable<unknown> {
    return refToObservable(this.value)
  }

  configureFilter(_: DataQuery): boolean {
    return true
  }

  configureQuery(_: DataQuery, _configuration: DataQueryExecutorConfiguration | null): boolean {
    return true
  }
}

function median(arr: number[]): number {
  const sortedArr = [...arr].sort((a, b) => a - b)
  const mid = Math.floor(sortedArr.length / 2)
  return sortedArr.length % 2 === 0 ? (sortedArr[mid - 1] + sortedArr[mid]) / 2 : sortedArr[mid]
}

export function scaleToMedian(arr: number[]): number[] {
  const medianValue = median(arr)
  return medianValue === 0 ? arr : arr.map((value) => (value / medianValue) * 50)
}
