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
    <TimeRangeSelect :configurator="timeRangeConfigurator" />
  </el-form>

  <el-divider class="dividerAfterForm">
    Bootstrap
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <ChartCard
        :provider="dataQueryExecutor"
        :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
      />
    </el-col>
    <el-col :span="12">
      <ChartCard
        :provider="dataQueryExecutor"
        :measures='["pluginDescriptorInitV18_d", "appComponentCreation_d", "projectComponentCreation_d"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Class and Resource Loading
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <ChartCard
        :provider="dataQueryExecutor"
        :measures='["classLoadingTime_i", "classLoadingSearchTime_i", "classLoadingDefineTime_i"]'
      />
    </el-col>
    <el-col :span="12">
      <ChartCard
        :provider="dataQueryExecutor"
        :measures='["classLoadingCount_i", "resourceLoadingCount_i"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Services
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <ChartCard
        :provider="dataQueryExecutor"
        :measures='["appComponentCreation_d", "serviceSyncPreloading_d", "serviceAsyncPreloading_d"]'
      />
    </el-col>
    <el-col :span="12">
      <ChartCard
        :skip-zero-values="false"
        :provider="dataQueryExecutor"
        :measures='["projectComponentCreation_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d", "moduleLoading_d"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Post-opening
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <ChartCard
        :provider="dataQueryExecutor"
        :measures='["projectDumbAware_d", "editorRestoring_d", "editorRestoringTillPaint_d"]'
      />
    </el-col>
    <el-col :span="12">
      <ChartCard
        :provider="dataQueryExecutor"
        :measures='["splash_i", "startUpCompleted_i"]'
      />
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
import ServerSelect from "../components/ServerSelect.vue"
import DimensionHierarchicalSelect from "../components/DimensionHierarchicalSelect.vue"
import { DataQueryExecutor, initDataComponent } from "../DataQueryExecutor"
import { TimeRangeConfigurator } from "../configurators/TimeRangeConfigurator"
import ChartCard from "../components/ChartCard.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"
import ReloadButton from "../components/ReloadButton.vue"

export default defineComponent({
  name: "IntelliJDashboard",
  components: {ReloadButton, ChartCard, ServerSelect, DimensionHierarchicalSelect, DimensionSelect, TimeRangeSelect},
  setup() {
    const persistentStateManager = new PersistentStateManager("ij-dashboard")
    const serverConfigurator = new ServerConfigurator("ij", persistentStateManager)
    const productConfigurator = new DimensionConfigurator("product", serverConfigurator, persistentStateManager)
    const projectConfigurator = createProjectConfigurator(productConfigurator, persistentStateManager)
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
    ], true)

    initDataComponent(serverConfigurator, persistentStateManager, dataQueryExecutor)

    return {
      serverUrl: serverConfigurator.server,
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
      dataQueryExecutor,
      getProjectName,
      loadData: dataQueryExecutor.scheduleLoadIncludingConfiguratorsFunctionReference,
    }
  },
})
</script>
