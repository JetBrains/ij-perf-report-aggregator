import { BehaviorSubject, Observable } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"

export class ReloadConfigurator implements DataQueryConfigurator {
  readonly subject = new BehaviorSubject<number>(0)

  configureQuery(_query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    return true
  }

  createObservable(): Observable<unknown> {
    return this.subject
  }
}