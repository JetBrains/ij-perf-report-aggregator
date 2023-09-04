import { useStorage } from "@vueuse/core"
import { Observable } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { FilterConfigurator } from "./filter"
import { refToObservable } from "./rxjs"

export class SmoothingConfigurator implements DataQueryConfigurator, FilterConfigurator {
  readonly value = useStorage("smoothingEnabled", false)

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
