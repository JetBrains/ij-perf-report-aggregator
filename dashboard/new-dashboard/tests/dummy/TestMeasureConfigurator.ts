import { Observable, of } from "rxjs"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../../src/components/common/dataQuery"

export class TestMeasureConfigurator implements DataQueryConfigurator {
  constructor(private readonly measures: string[] = ["testMeasure"]) {}

  createObservable(): Observable<unknown> {
    return of(null)
  }

  configureQuery(_query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    if (this.measures.length === 0) {
      return false
    }
    configuration.measures = this.measures
    return true
  }
}
