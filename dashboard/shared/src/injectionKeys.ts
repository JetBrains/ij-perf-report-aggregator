import { InjectionKey, Ref } from "vue"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { ChartStyle } from "./chart"
import { ReportInfoProvider } from "./components/ChartToolTipManager"
import ChartTooltip from "./components/ChartTooltip.vue"
import { AggregationOperatorConfigurator } from "./configurators/AggregationOperatorConfigurator"
import { TimeRange } from "./configurators/TimeRangeConfigurator"

// inject is used instead of prop because on dashboard page there are a lot of chart cards and it is tedious to set property for each
export const dataQueryExecutorKey: InjectionKey<DataQueryExecutor> = Symbol("dataQueryExecutor")
export const aggregationOperatorConfiguratorKey: InjectionKey<AggregationOperatorConfigurator> = Symbol("aggregationOperatorConfigurator")
export const timeRangeKey: InjectionKey<Ref<TimeRange>> = Symbol("timeRange")
export const reportInfoProviderKey: InjectionKey<ReportInfoProvider> = Symbol("tooltipUrlProvider")
export const serverUrlKey: InjectionKey<Ref<string>> = Symbol("serverUrl")

export const chartStyleKey: InjectionKey<ChartStyle> = Symbol("chartStyle")

export const chartToolTipKey: InjectionKey<Ref<typeof ChartTooltip>> = Symbol("chartToolTip")
