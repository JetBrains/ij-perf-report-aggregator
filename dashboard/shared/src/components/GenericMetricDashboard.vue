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
      <DimensionSelect
        v-if="supportReleases"
        label="Nightly/Release"
        :dimension="releaseConfigurator"
      />
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>
    <template
      v-for="metric in measureConfigurator.selected.value"
      :key="metric"
    >
      <Divider :label="metric" />
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
import Divider from "tailwind-ui/src/Divider.vue"
import { provide, withDefaults } from "vue"
import { useRouter } from "vue-router"
import { initDataComponent } from "../DataQueryExecutor"
import { PersistentStateManager } from "../PersistentStateManager"
import { chartDefaultStyle, DEFAULT_LINE_CHART_HEIGHT, ValueUnit } from "../chart"
import { AggregationOperatorConfigurator } from "../configurators/AggregationOperatorConfigurator"
import { dimensionConfigurator } from "../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../configurators/MachineConfigurator"
import { ChartType, MeasureConfigurator } from "../configurators/MeasureConfigurator"
import { ReleaseNightlyConfigurator } from "../configurators/ReleaseNightlyConfigurator"
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
  installerExists?: boolean
  supportReleases?: boolean
  valueUnit?: ValueUnit
}>(), {
  compoundTooltip: true,
  urlEnabled: true,
  table: undefined,
  chartType: "line",
  valueUnit: "ms",
  installerExists: true,
  supportReleases: false,
})

provide(chartStyleKey, {
  ...chartDefaultStyle,
  valueUnit: props.valueUnit,
})

if (props.urlEnabled) {
  provideReportUrlProvider(props.installerExists)
}

const persistentStateManager = new PersistentStateManager(`${(props.dbName)}-${(props.table == null ? "" : `${props.table}-`)}dashboard`, {
  machine: "linux-blade",
  project: [],
  branch: "master",
  measure: props.defaultMeasures.slice(),
}, useRouter())

const serverConfigurator = new ServerConfigurator(props.dbName, props.table)

const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true)

const measureConfigurator = new MeasureConfigurator(
  serverConfigurator,
  persistentStateManager,
  [scenarioConfigurator, branchConfigurator],
  true,
  props.chartType,
)

const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator])

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
let configurators = [
  serverConfigurator,
  scenarioConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
]
if (props.supportReleases) {
  configurators.push(new ReleaseNightlyConfigurator(persistentStateManager))
}
// median by default, no UI control to change is added (insert <AggregationOperatorSelect /> if needed)
provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
initDataComponent(configurators)
</script>