<template>
  <Dashboard>
    <template #toolbar>
      <DimensionSelect
        label="Scenarios"
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
import { initDataComponent } from "src/DataQueryExecutor"
import { PersistentStateManager } from "src/PersistentStateManager"
import { DEFAULT_LINE_CHART_HEIGHT } from "src/chart"
import BarChartCard from "shared/src/components/BarChartCard.vue"
import Dashboard from "shared/src/components/Dashboard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import MeasureSelect from "shared/src/components/MeasureSelect.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { AggregationOperatorConfigurator } from "src/configurators/AggregationOperatorConfigurator"
import { DimensionConfigurator } from "src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "src/configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey } from "src/injectionKeys"
import { provideReportUrlProvider } from "src/lineChartTooltipLinkProvider"
import { provide, withDefaults } from "vue"
import { useRouter } from "vue-router"

provideReportUrlProvider()

// eslint-disable-next-line no-undef
const props = withDefaults(defineProps<{
  dbName: string
  defaultMeasures: Array<string>
  urlEnabled: boolean
}>(), {
  urlEnabled: true
})

if(props.urlEnabled){
  provideReportUrlProvider()
}

const persistentStateManager = new PersistentStateManager(`${(props.dbName)}-dashboard`, {
  machine: "linux-blade",
  project: [],
  branch: "master",
  measure: props.defaultMeasures.slice(),
}, useRouter())

const serverConfigurator = new ServerConfigurator(props.dbName)
const scenarioConfigurator = new DimensionConfigurator("project", serverConfigurator, persistentStateManager, true)
const branchConfigurator = new DimensionConfigurator("branch", serverConfigurator, persistentStateManager, true)

const machineConfigurator = new MachineConfigurator(new DimensionConfigurator("machine", serverConfigurator, persistentStateManager), persistentStateManager)

const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator, branchConfigurator, machineConfigurator])
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