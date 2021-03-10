import { InjectionKey, Ref } from "vue"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { ChartTooltipLinkProvider } from "./LineChartManager"
import { AggregationOperatorConfigurator } from "./configurators/AggregationOperatorConfigurator"
import { TimeRange } from "./configurators/TimeRangeConfigurator"

// inject is used instead of prop because on dashboard page there are a lot of chart cards and it is tedious to set property for each
export const dataQueryExecutorKey: InjectionKey<DataQueryExecutor> = Symbol("dataQueryExecutor")
export const aggregationOperatorConfiguratorKey: InjectionKey<AggregationOperatorConfigurator> = Symbol("aggregationOperatorConfigurator")
export const timeRangeKey: InjectionKey<Ref<TimeRange>> = Symbol("timeRange")
export const tooltipUrlProviderKey: InjectionKey<ChartTooltipLinkProvider> = Symbol("tooltipUrlProvider")
export const serverUrlKey: InjectionKey<Ref<string>> = Symbol("serverUrl")
