import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../common/dataQuery"

export class SeriesNameConfigurator implements DataQueryConfigurator {
  constructor(private readonly measureName: string) {}
  createObservable() {
    return null
  }
  configureQuery(_query: DataQuery, configuration: DataQueryExecutorConfiguration): boolean {
    const measureName = this.measureName
    configuration.queryProducers.push({
      getSeriesName(_index: number): string {
        return measureName
      },
      // eslint-disable-next-line @typescript-eslint/no-empty-function
      mutate(_index: number): void {},
      size(): number {
        return 1
      },
      getMeasureName(_index: number): string {
        return ""
      },
    })
    return true
  }
}
