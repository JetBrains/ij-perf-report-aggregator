<template>
  <Dashboard>
    <template #toolbar>
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
      <DimensionSelect
        label="Branch"
        :dimension="branchConfigurator"
      />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
    </template>
    <template
      v-for="metric in aggregatedMetrics"
      :key="metric"
    >
      <Divider :label="metric" />
      <div class="grid grid-cols-12 gap-4">
        <div class="col-span-12">
          <LineChartCard
            :skip-zero-values="false"
            :aggregated-measure="metric"
            trigger="item"
          />
        </div>
      </div>
    </template>
  </Dashboard>
</template>

<script lang="ts" setup>
import Divider from "tailwind-ui/src/Divider.vue"
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

const aggregatedMetrics = ["test#max_awt_delay", "test#average_awt_delay"]

const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  timeAverageConfigurator,
]
initDataComponent(configurators)
</script>