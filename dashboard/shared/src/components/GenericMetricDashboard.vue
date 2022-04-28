<template>
  <Dashboard>
    <template #toolbar>
      <DimensionSelect
        label="Scenarios"
        tooltip="Scenarios"
        :dimension="scenarioConfigurator"
      />
      <DimensionSelect
        label="Branch"
        :dimension="branchConfigurator"
      />
      <MeasureSelect :configurator="measureConfigurator" />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>
    <template
      v-for="metric in measureConfigurator.selected.value"
      :key="metric"
    >
      <Divider align="center">
        {{ metric }}
      </Divider>
      <div class="grid grid-cols-12 gap-4">
        <div class="col-span-8">
          <LineChartCard
            :compound-tooltip="compoundTooltip"
            :chart-type="chartType"
            :value-unit="valueUnit"
            :measures="[metric]"
          />
        </div>
        <div class="col-span-4">
          <BarChartCard
            :height="chartHeight"
            :measures="[metric]"
          />
        </div>
      </div>
    </template>
  </Dashboard>
</template>

<script lang="ts" setup>
import { provide, withDefaults } from "vue"
import { useRouter } from "vue-router"
import { initDataComponent } from "../DataQueryExecutor"
import { PersistentStateManager } from "../PersistentStateManager"
import { chartDefaultStyle, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../chart"
import { AggregationOperatorConfigurator } from "../configurators/AggregationOperatorConfigurator"
import { DimensionConfigurator } from "../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../configurators/MachineConfigurator"
import { ChartType, MeasureConfigurator } from "../configurators/MeasureConfigurator"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "../configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey } from "../injectionKeys"
import { provideReportUrlProvider } from "../lineChartTooltipLinkProvider"
import BarChartCard from "./BarChartCard.vue"
import Dashboard from "./Dashboard.vue"
import DimensionHierarchicalSelect from "./DimensionHierarchicalSelect.vue"
import DimensionSelect from "./DimensionSelect.vue"
import LineChartCard from "./LineChartCard.vue"
import MeasureSelect from "./MeasureSelect.vue"
import TimeRangeSelect from "./TimeRangeSelect.vue"

const props = withDefaults(defineProps<{
  dbName: string
  table?: string
  compoundTooltip?: boolean
  chartType?: ChartType
  defaultMeasures: Array<string>
  urlEnabled?: boolean
  valueUnit?: ValueUnit
}>(), {
  compoundTooltip: true,
  urlEnabled: true,
  table: undefined,
  chartType: "line",
  valueUnit: "ms",
})

provide(chartStyleKey, {
  ...chartDefaultStyle,
  valueUnit: props.valueUnit,
})

if (props.urlEnabled) {
  provideReportUrlProvider()
}

const persistentStateManager = new PersistentStateManager(`${(props.dbName)}-dashboard`, {
  machine: "linux-blade",
  project: [],
  branch: "master",
  measure: props.defaultMeasures.slice(),
}, useRouter())

const serverConfigurator = new ServerConfigurator(props.dbName, props.table)
const scenarioConfigurator = new DimensionConfigurator("project", serverConfigurator, persistentStateManager, true)
const branchConfigurator = new DimensionConfigurator("branch", serverConfigurator, persistentStateManager, true)

const machineConfigurator = new MachineConfigurator(
  new DimensionConfigurator("machine", serverConfigurator, persistentStateManager),
  persistentStateManager,
)

const measureConfigurator = new MeasureConfigurator(
  serverConfigurator,
  persistentStateManager,
  [scenarioConfigurator, branchConfigurator, machineConfigurator],
  true,
  props.chartType,
)
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

// median by default, no UI control to change is added (insert <AggregationOperatorSelect /> if needed)
provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
initDataComponent(persistentStateManager, [
  serverConfigurator,
  scenarioConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
])
</script>