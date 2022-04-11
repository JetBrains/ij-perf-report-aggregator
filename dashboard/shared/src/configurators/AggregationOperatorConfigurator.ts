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
    const aggregator = this.value.value ?? AggregationOperatorConfigurator.DEFAULT_OPERATOR
    if (aggregator.operator === "median" || aggregator.operator.length === 0) {
      query.aggregator = "quantileTDigest(0.5)"
    }
    else if (aggregator.operator === "quantile") {
      query.aggregator = `quantileTDigest(${aggregator.quantile / 100})`
    }
    else {
      query.aggregator = aggregator.operator
    }
    return true
  }
}