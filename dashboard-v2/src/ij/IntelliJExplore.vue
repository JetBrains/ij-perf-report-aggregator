<template>
  <el-form
    :inline="true"
    size="small"
  >
    <ServerSelect v-model="serverUrl" />

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

    <ReloadButton :load="loadData" />
  </el-form>

  <el-form
    :inline="true"
    size="small"
  >
    <MeasureSelect :configurator="measureConfigurator" />
    <TimeRangeSelect :configurator="timeRangeConfigurator" />
  </el-form>

  <el-row>
    <el-col :span="12">
      <ChartCard :provider="dataQueryExecutor" />
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { ServerConfigurator } from "../dataQuery"
import { defineComponent } from "vue"
import { PersistentStateManager } from "../PersistentStateManager"
import { DimensionConfigurator, SubDimensionConfigurator } from "../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../configurators/MachineConfigurator"
import DimensionSelect from "../components/DimensionSelect.vue"
import TimeRangeSelect from "../components/TimeRangeSelect.vue"
import MeasureSelect from "../components/MeasureSelect.vue"
import ServerSelect from "../components/ServerSelect.vue"
import DimensionHierarchicalSelect from "../components/DimensionHierarchicalSelect.vue"
import { DataQueryExecutor, initDataComponent } from "../DataQueryExecutor"
import { MeasureConfigurator } from "../configurators/MeasureConfigurator"
import { TimeRangeConfigurator } from "../configurators/TimeRangeConfigurator"
import ChartCard from "../components/ChartCard.vue"
import ReloadButton from "../components/ReloadButton.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

export default defineComponent({
  name: "IntelliJExplore",
  components: {ReloadButton, ChartCard, ServerSelect, DimensionHierarchicalSelect, DimensionSelect, MeasureSelect, TimeRangeSelect},
  setup() {
    const persistentStateManager = new PersistentStateManager("ij-explore")
    const serverConfigurator = new ServerConfigurator("ij", persistentStateManager)
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

    initDataComponent(serverConfigurator, persistentStateManager, dataQueryExecutor)

    return {
      serverUrl: serverConfigurator.server,
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      measureConfigurator,
      timeRangeConfigurator,
      dataQueryExecutor,
      getProjectName,
      loadData: dataQueryExecutor.scheduleLoadIncludingConfiguratorsFunctionReference,
    }
  },
})
</script>
