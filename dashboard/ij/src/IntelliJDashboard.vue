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
      :value-to-label="getProjectName"
      :dimension="projectConfigurator"
    />
    <DimensionHierarchicalSelect
      label="Machine"
      :dimension="machineConfigurator"
    />
    <TimeRangeSelect :configurator="timeRangeConfigurator" />

    <ReloadButton />
  </el-form>

  <el-divider style="margin-top: 5px">
    Progress Over Time
  </el-divider>
  <AggregationOperatorSelect />
  <el-row :gutter="5">
    <el-col :span="12">
      <BarChartCard
        :measures='[
          "bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d",
          "appComponentCreation_d", "projectComponentCreation_d",
        ]'
      />
    </el-col>
    <el-col :span="12">
      <BarChartCard
        :measures='["splash_i", "startUpCompleted_i"]'
      />
    </el-col>
  </el-row>
  <el-row
    :gutter="5"
    style="margin-top: 5px;"
  >
    <!-- todo "moduleLoading_d" when it will be fixed -->
    <el-col :span="12">
      <BarChartCard
        :measures='[
          "appStarter_d",
          "serviceSyncPreloading_d", "serviceAsyncPreloading_d",
          "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d",
        ]'
      />
    </el-col>
    <el-col :span="12">
      <BarChartCard
        :measures='[
          "projectDumbAware_d", "editorRestoring_d", "editorRestoringTillPaint_d"
        ]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Bootstrap
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='["pluginDescriptorInitV18_d", "appComponentCreation_d", "projectComponentCreation_d"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Class and Resource Loading
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["classLoadingTime_i", "classLoadingSearchTime_i", "classLoadingDefineTime_i"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='["classLoadingCount_i", "resourceLoadingCount_i"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Services
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["appComponentCreation_d", "serviceSyncPreloading_d", "serviceAsyncPreloading_d"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :skip-zero-values="false"
        :measures='["projectComponentCreation_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d", "moduleLoading_d"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Post-opening
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["projectDumbAware_d", "editorRestoring_d", "editorRestoringTillPaint_d"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='["splash_i", "startUpCompleted_i"]'
      />
    </el-col>
  </el-row>
</template>

<script lang="ts">
import { DataQueryExecutor, initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import AggregationOperatorSelect from "shared/src/components/AggregationOperatorSelect.vue"
import BarChartCard from "shared/src/components/BarChartCard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import ReloadButton from "shared/src/components/ReloadButton.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { AggregationOperatorConfigurator } from "shared/src/configurators/AggregationOperatorConfigurator"
import { DimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { SubDimensionConfigurator } from "shared/src/configurators/SubDimensionConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyle } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { defineComponent, provide } from "vue"
import { useRouter } from "vue-router"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

export default defineComponent({
  name: "IntelliJDashboard",
  components: {
    ReloadButton, LineChartCard, BarChartCard, DimensionHierarchicalSelect, DimensionSelect, TimeRangeSelect,
    AggregationOperatorSelect,
  },
  setup() {
    provideReportUrlProvider()
    provide(chartStyle, {
      ...chartDefaultStyle,
      // a lot of bars, as result, height of bar is not enough to make label readable
      barSeriesLabelPosition: "right",
    })

    // noinspection SpellCheckingInspection
    const persistentStateManager = new PersistentStateManager("ij-dashboard", {
      product: "IU",
      project: "73YWaW9bytiPDGuKvwNIYMK5CKI",
      machine: "macMini 2018",
    }, useRouter())
    const serverConfigurator = new ServerConfigurator("ij")
    const productConfigurator = new DimensionConfigurator("product", serverConfigurator, persistentStateManager)
    const projectConfigurator = createProjectConfigurator(productConfigurator, persistentStateManager)
    const machineConfigurator = new MachineConfigurator(
      new SubDimensionConfigurator("machine", productConfigurator),
      persistentStateManager,
    )
    const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

    provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))

    const dataQueryExecutor = new DataQueryExecutor([
      serverConfigurator,
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
    ], true)

    initDataComponent(persistentStateManager, dataQueryExecutor)

    return {
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
      getProjectName,
    }
  },
})
</script>
