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

  <el-tabs v-model="activeTab">
    <el-tab-pane
      label="Pulse"
      name="pulse"
    />
    <el-tab-pane
      label="Progress Over Time"
      name="progressOverTime"
    />
    <el-tab-pane
      label="Module Loading"
      name="moduleLoading"
    />
  </el-tabs>
  <router-view v-slot="{ Component }">
    <keep-alive>
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>

<script lang="ts">
import { DataQueryExecutor, initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
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
import { defineComponent, provide, ref, watch } from "vue"
import { useRoute, useRouter } from "vue-router"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

export default defineComponent({
  name: "IntelliJDashboard",
  components: {
    ReloadButton, DimensionHierarchicalSelect, DimensionSelect, TimeRangeSelect,
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
      project: "simple for IJ",
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

    const route = useRoute()
    const activeTab = ref("pulse")

    function updateActiveTab(newPath: string): void {
      if (newPath.endsWith("/pulse")) {
        activeTab.value = "pulse"
      }
      else if (newPath.endsWith("/progressOverTime")) {
        activeTab.value = "progressOverTime"
      }
    }

    updateActiveTab(route.path)

    const router = useRouter()
    watch(() => route.path, updateActiveTab)
    watch(activeTab, async value => {
      await router.push({...route, path: `/ij/${value}`})
    })

    return {
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
      getProjectName,
      activeTab,
    }
  },
})
</script>
