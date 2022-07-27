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
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import Dashboard from "shared/src/components/Dashboard.vue"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeAverageConfigurator } from "shared/src/configurators/TimeAverageConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import Divider from "tailwind-ui/src/Divider.vue"
import { useRouter } from "vue-router"
import DimensionHierarchicalSelect from "/shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"

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