import { Observable } from "rxjs"
import { inject, InjectionKey, Ref } from "vue"
import { ChartStyle } from "./chart"
import { ReportInfoProvider } from "./components/ChartToolTipManager"
import ChartTooltip from "./components/ChartTooltip.vue"
import { AggregationOperatorConfigurator } from "./configurators/AggregationOperatorConfigurator"
import { TimeRange } from "./configurators/TimeRangeConfigurator"
import { DataQueryConfigurator } from "./dataQuery"
import { CompressorUsingDictionary } from "./zstd"

// inject is used instead of prop because on dashboard page there are a lot of chart cards and it is tedious to set property for each
export const configuratorListKey: InjectionKey<Array<DataQueryConfigurator>> = Symbol("dataQueryExecutor")
export const aggregationOperatorConfiguratorKey: InjectionKey<AggregationOperatorConfigurator> = Symbol("aggregationOperatorConfigurator")
export const timeRangeKey: InjectionKey<Ref<TimeRange>> = Symbol("timeRange")
export const reportInfoProviderKey: InjectionKey<ReportInfoProvider> = Symbol("tooltipUrlProvider")

export const serverUrlObservableKey: InjectionKey<Observable<string>> = Symbol("serverUrlObservable")
export const compressorObservableKey: InjectionKey<Observable<CompressorUsingDictionary>> = Symbol("compressorObservable")

export const chartStyleKey: InjectionKey<ChartStyle> = Symbol("chartStyle")

export const chartToolTipKey: InjectionKey<Ref<typeof ChartTooltip>> = Symbol("chartToolTip")

export function injectOrError<T>(key: InjectionKey<T> | string): T {
  const value = inject(key)
  if (value === undefined) {
    throw new Error(`${key.toString()} is not provided`)
  }
  return value
}
