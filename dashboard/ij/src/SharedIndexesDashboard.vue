<template>
  <div class="w-full">
    <Toolbar>
      <template #start>
        <FilterInput :configurator="scenarioConfigurator"/>
        <DimensionSelect
          label="Scenarios"
          :dimension="scenarioConfigurator"
        />
        <DimensionSelect
          label="Branch"
          :dimension="branchConfigurator"
        />
        <MeasureSelect :configurator="measureConfigurator"/>
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        />
        <TimeRangeSelect :configurator="timeRangeConfigurator"/>
      </template>
      <template #end>
        <ReloadButton/>
      </template>
    </Toolbar>
  </div>

  <template
    v-for="metric in measureConfigurator.value.value"
    :key="metric"
  >
    <Divider align="center">
      {{ metric }}
    </Divider>
    <div class="grid grid-cols-12 gap-4">
      <div class="col-span-7">
        <LineChartCard
          :measures="[metric]"
        />
      </div>
      <div class="col-span-5">
        <BarChartCard
          :height="chartHeight"
          :measures="[metric]"
        />
      </div>
    </div>
  </template>
</template>

<script lang="ts" setup>
import { DEFAULT_LINE_CHART_HEIGHT } from "shared/src/chart"
import BarChartCard from "shared/src/components/BarChartCard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import MeasureSelect from "shared/src/components/MeasureSelect.vue"
import ReloadButton from "shared/src/components/ReloadButton.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { AggregationOperatorConfigurator } from "shared/src/configurators/AggregationOperatorConfigurator"
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { DataQueryExecutor, initDataComponent } from "shared/src/DataQueryExecutor"
import { aggregationOperatorConfiguratorKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { provide } from "vue"
import { useRouter } from "vue-router"
import FilterInput from "../../shared/src/components/FilterInput.vue"

provideReportUrlProvider()
// eslint-disable-next-line no-undef
const props = defineProps<{
  dbName: string
  defaultMeasures: Array<string>
}>()

const persistentStateManager = new PersistentStateManager(`${(props.dbName)}-dashboard`, {
  machine: "linux-blade",
  project: [],
  measure: props.defaultMeasures.slice(),
}, useRouter())

const serverConfigurator = new ServerConfigurator(props.dbName)
const scenarioConfigurator = new DimensionConfigurator("project", serverConfigurator, persistentStateManager, true)
const branchConfigurator = new DimensionConfigurator("branch", serverConfigurator, persistentStateManager, true)

const machineConfigurator = new MachineConfigurator(new DimensionConfigurator("machine", serverConfigurator, persistentStateManager),
  persistentStateManager)

const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, scenarioConfigurator)
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

// median by default, no UI control to change is added (insert <AggregationOperatorSelect /> if needed)
provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))

const dataQueryExecutor = new DataQueryExecutor([
  serverConfigurator,
  scenarioConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
], true)

const chartHeight = DEFAULT_LINE_CHART_HEIGHT
initDataComponent(persistentStateManager, dataQueryExecutor)
</script>