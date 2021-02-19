import { InjectionKey , Ref } from "vue"
import { DataQueryExecutor } from "./DataQueryExecutor"
import { ChartTooltipLinkProvider } from "./LineChartManager"
import { AggregationOperatorConfigurator } from "./configurators/AggregationOperatorConfigurator"
import { TimeRange } from "./configurators/TimeRangeConfigurator"

export const dataQueryExecutorKey: InjectionKey<DataQueryExecutor> = Symbol("dataQueryExecutor")
export const aggregationOperatorConfiguratorKey: InjectionKey<AggregationOperatorConfigurator> = Symbol("aggregationOperatorConfigurator")
export const timeRangeKey: InjectionKey<Ref<TimeRange>> = Symbol("timeRange")
export const tooltipUrlProviderKey: InjectionKey<ChartTooltipLinkProvider> = Symbol("tooltipUrlProvider")
export const serverUrlKey: InjectionKey<Ref<string>> = Symbol("serverUrl")
