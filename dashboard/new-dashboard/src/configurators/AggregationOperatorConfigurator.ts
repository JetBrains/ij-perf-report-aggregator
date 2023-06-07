import { combineLatest, Observable, shareReplay } from "rxjs"
import { computed, ref } from "vue"
import { PersistentStateManager } from "../components/common/PersistentStateManager"
import { DataQuery, DataQueryConfigurator, DataQueryExecutorConfiguration } from "../components/common/dataQuery"
import { refToObservable } from "./rxjs"

export class AggregationOperatorConfigurator implements DataQueryConfigurator {
  static readonly DEFAULT_OPERATOR = "median"

  readonly operator = ref(AggregationOperatorConfigurator.DEFAULT_OPERATOR)
  readonly quantile = ref(50)

  private readonly observable: Observable<unknown>

  constructor(persistentStateManager: PersistentStateManager) {
    persistentStateManager.add(
      "aggregationOperator",
      computed(() => ({ operator: this.operator.value, quantile: this.quantile.value }))
    )
    this.observable = combineLatest([refToObservable(this.operator), refToObservable(this.quantile)]).pipe(shareReplay(1))
  }

  createObservable(): Observable<unknown> {
    return this.observable
  }

  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    const operator = this.operator.value
    if (operator === "median" || operator.length === 0) {
      query.aggregator = "quantileTDigest(0.5)"
    } else if (operator === "quantile") {
      query.aggregator = `quantileTDigest(${this.quantile.value / 100})`
    } else {
      query.aggregator = operator
    }
    return true
  }
}
