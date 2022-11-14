<template>
  <Dashboard>
    <template #toolbar>
      <DimensionSelect
        label="Branch"
        :dimension="branchConfigurator"
      />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>
    <template
      v-for="metric in aggregatedMetrics"
      :key="metric"
    >
      <LineChartCard
        class="w-full"
        :label="metric"
        :skip-zero-values="false"
        :aggregated-measure="metric"
        trigger="item"
      />
    </template>
  </Dashboard>
</template>

<script lang="ts" setup>
import { useRouter } from "vue-router"
import { initDataComponent } from "../DataQueryExecutor"
import { PersistentStateManager } from "../PersistentStateManager"
import { dimensionConfigurator } from "../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../configurators/MachineConfigurator"
import { ServerConfigurator } from "../configurators/ServerConfigurator"
import { TimeAverageConfigurator } from "../configurators/TimeAverageConfigurator"
import { TimeRangeConfigurator } from "../configurators/TimeRangeConfigurator"
import Dashboard from "./Dashboard.vue"
import DimensionHierarchicalSelect from "./DimensionHierarchicalSelect.vue"
import DimensionSelect from "./DimensionSelect.vue"
import LineChartCard from "./LineChartCard.vue"
import TimeRangeSelect from "./TimeRangeSelect.vue"

const props = defineProps<{
  dbName: string
  table?: string
}>()

const persistentStateManager = new PersistentStateManager(`${props.dbName}_${props.table == null ? "" : props.table}_aggregated_dashboard`, {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator(props.dbName, props.table)
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true, [timeRangeConfigurator], (a, _) => {
  return a.includes("/") ? 1 : -1
})
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])

const timeAverageConfigurator = new TimeAverageConfigurator()

const aggregatedMetrics = ["indexing", "completion"]

const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  timeAverageConfigurator,
]
initDataComponent(configurators)
</script>