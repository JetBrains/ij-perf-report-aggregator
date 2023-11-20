import { Observable } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { FilterConfigurator } from "./filter"

export class NoOpConfigurator implements DataQueryConfigurator, FilterConfigurator {
  configureFilter(_: DataQuery): boolean {
    return true
  }

  configureQuery(_: DataQuery, _1: DataQueryExecutorConfiguration): boolean {
    return true
  }

  createObservable(): Observable<unknown> | null {
    return null
  }
}
