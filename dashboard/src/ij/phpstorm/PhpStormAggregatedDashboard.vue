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
import { initDataComponent } from "../../shared/DataQueryExecutor"
import { PersistentStateManager } from "../../shared/PersistentStateManager"
import Dashboard from "../../shared/components/Dashboard.vue"
import DimensionSelect from "../../shared/components/DimensionSelect.vue"
import LineChartCard from "../../shared/components/LineChartCard.vue"
import TimeRangeSelect from "../../shared/components/TimeRangeSelect.vue"
import { dimensionConfigurator } from "../../shared/configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../shared/configurators/MachineConfigurator"
import { ServerConfigurator } from "../../shared/configurators/ServerConfigurator"
import { TimeAverageConfigurator } from "../../shared/configurators/TimeAverageConfigurator"
import { TimeRangeConfigurator } from "../../shared/configurators/TimeRangeConfigurator"
import Divider from "../../tailwind-ui/Divider.vue"
import { useRouter } from "vue-router"
import DimensionHierarchicalSelect from "/../../shared/components/DimensionHierarchicalSelect.vue"

const persistentStateManager = new PersistentStateManager("phpstorm_dashboard", {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "phpstorm")
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true, [timeRangeConfigurator], (a, _) => {
  return a.includes("/") ? 1 : -1
})
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])

const timeAverageConfigurator = new TimeAverageConfigurator()

const aggregatedMetrics = ["responsiveness_time", "completion_execution_time"]

const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  timeAverageConfigurator
]
initDataComponent(configurators)
</script>