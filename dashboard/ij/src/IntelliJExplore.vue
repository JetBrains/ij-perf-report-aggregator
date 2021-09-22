<template>
  <el-form
    :inline="true"
    size="small"
  >
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

    <ReloadButton />
  </el-form>

  <el-form
    :inline="true"
    size="small"
  >
    <MeasureSelect :configurator="measureConfigurator" />
    <TimeRangeSelect :configurator="timeRangeConfigurator" />
  </el-form>

  <div class="grid grid-cols-2 gap-4">
    <LineChartCard />
  </div>
</template>

<script lang="ts">
import { DataQueryExecutor, initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import MeasureSelect from "shared/src/components/MeasureSelect.vue"
import ReloadButton from "shared/src/components/ReloadButton.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { SubDimensionConfigurator } from "shared/src/configurators/SubDimensionConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { defineComponent } from "vue"

import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

export default defineComponent({
  name: "IntelliJExplore",
  components: {ReloadButton, LineChartCard, DimensionHierarchicalSelect, DimensionSelect, MeasureSelect, TimeRangeSelect},
  setup() {
    const persistentStateManager = new PersistentStateManager("ij-explore")
    const serverConfigurator = new ServerConfigurator("ij")
    const productConfigurator = new DimensionConfigurator("product", serverConfigurator, persistentStateManager)
    const projectConfigurator = createProjectConfigurator(productConfigurator, persistentStateManager)
    const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager)
    const machineConfigurator = new MachineConfigurator(
      new SubDimensionConfigurator("machine", productConfigurator),
      persistentStateManager,
    )
    const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

    const dataQueryExecutor = new DataQueryExecutor([
      serverConfigurator,
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
      measureConfigurator,
    ])

    initDataComponent(persistentStateManager, dataQueryExecutor)

    return {
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      measureConfigurator,
      timeRangeConfigurator,
      dataQueryExecutor,
      getProjectName,
    }
  },
})
</script>
