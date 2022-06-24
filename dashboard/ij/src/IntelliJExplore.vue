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
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import Dashboard from "shared/src/components/Dashboard.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import DimensionHierarchicalSelect from "../../shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../../shared/src/components/DimensionSelect.vue"
import MeasureSelect from "../../shared/src/components/MeasureSelect.vue"
import TimeRangeSelect from "../../shared/src/components/TimeRangeSelect.vue"
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
