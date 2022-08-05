<template>
  <Dashboard>
    <template #toolbar>
      <DimensionSelect
        label="Product"
        :dimension="productConfigurator"
      />
      <DimensionSelect
        label="Project"
        :value-label="getProjectName"
        :dimension="projectConfigurator"
      />

      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />

      <MeasureSelect :configurator="measureConfigurator" />
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>

    <div class="grid grid-cols-2 gap-4 mt-2">
      <LineChartCard />
    </div>
  </Dashboard>
</template>

<script setup lang="ts">
import { initDataComponent } from "../shared/DataQueryExecutor"
import { PersistentStateManager } from "../shared/PersistentStateManager"
import Dashboard from "../shared/components/Dashboard.vue"
import DimensionHierarchicalSelect from "../shared/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../shared/components/DimensionSelect.vue"
import LineChartCard from "../shared/components/LineChartCard.vue"
import MeasureSelect from "../shared/components/MeasureSelect.vue"
import TimeRangeSelect from "../shared/components/TimeRangeSelect.vue"
import { dimensionConfigurator } from "../shared/configurators/DimensionConfigurator"
import { MachineConfigurator } from "../shared/configurators/MachineConfigurator"
import { MeasureConfigurator } from "../shared/configurators/MeasureConfigurator"
import { ServerConfigurator } from "../shared/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "../shared/configurators/TimeRangeConfigurator"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

const persistentStateManager = new PersistentStateManager("ij-explore")
const serverConfigurator = new ServerConfigurator("ij")

const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager)
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager)
const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager)
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [productConfigurator])
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

initDataComponent([
  serverConfigurator,
  productConfigurator,
  projectConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  measureConfigurator,
])
</script>
