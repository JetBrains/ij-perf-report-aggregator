import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { FilterConfigurator } from "./filter"

/**
 * Leaves out the points generated before [start], a date in the `YYYY-MM-DD` form. A dashboard whose tests changed how
 * they measure shows only the points that the current method produced, so a chart and its alerts compare like with like.
 */
export class HistoryStartConfigurator implements DataQueryConfigurator, FilterConfigurator {
  constructor(private readonly start: string) {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(start)) {
      throw new Error(`The history start must be a YYYY-MM-DD date, not ${start}`)
    }
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    return this.configureFilter(query)
  }

  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "generated_time", q: `>= toDateTime('${this.start} 00:00:00')` })
    return true
  }

  createObservable(): null {
    return null
  }
}
