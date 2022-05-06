import { Observable } from "rxjs"
import { ref } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../dataQuery"
import { refToObservable } from "./rxjs"

export class AggregationOperatorConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_OPERATOR = "median"
  readonly value = ref({operator: AggregationOperatorConfigurator.DEFAULT_OPERATOR, quantile: 50})

  createObservable(): Observable<unknown> {
    return refToObservable(this.value, true)
  }

  constructor(persistentStateManager: PersistentStateManager) {
    persistentStateManager.add("aggregationOperator", this.value)
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    const aggregator = this.value.value
    const operator = aggregator.operator
    if (operator === "median" || operator.length === 0) {
      query.aggregator = "quantileTDigest(0.5)"
    }
    else if (operator === "quantile") {
      query.aggregator = `quantileTDigest(${(aggregator.quantile ?? 50) / 100})`
    }
    else {
      query.aggregator = operator
    }
    return true
  }
}