import { inject, InjectionKey, Ref } from "vue"
import { DataQueryConfigurator } from "../components/common/dataQuery"
import { TimeRange } from "../configurators/TimeRangeConfigurator"
import { ReportInfoProvider } from "../components/charts/PopupProvider"

// inject is used instead of prop because on dashboard page there are a lot of chart cards and it is tedious to set property for each
export const configuratorListKey: InjectionKey<DataQueryConfigurator[]> = Symbol("configuratorList")
export const timeRangeKey: InjectionKey<Ref<TimeRange>> = Symbol("timeRange")
export const reportInfoProviderKey: InjectionKey<ReportInfoProvider> = Symbol("tooltipUrlProvider")

export const serverUrlObservableKey: InjectionKey<Observable<string>> = Symbol("serverUrlObservable")

export function injectOrError<T>(key: InjectionKey<T> | string): T {
  const value = inject(key)
  if (value === undefined) {
    throw new Error(`${key.toString()} is not provided`)
  }
  return value
}

export function injectOrNull<T>(key: InjectionKey<T> | string): T | null {
  const value = inject(key)
  if (value === undefined) {
    return null
  }
  return value
}
