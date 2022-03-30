<template>
  <Toolbar>
    <template #start>
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
    </template>
    <template #end>
      <ReloadButton />
    </template>
  </Toolbar>

  <TabView v-model:active-index="activeTab">
    <TabPanel
      v-for="tab in tabs"
      :key="tab.title"
      :header="tab.title"
    />
  </TabView>
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

export interface Tab {
  route: string
  title: string
}

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
    const activeTab = ref(0)
    const tabs = ref<Array<Tab>>([
      {
        title: "Pulse",
        route: "pulse",
      },
      {
        title: "Progress Over Time",
        route: "progressOverTime",
      },
      {
        title: "Module Loading",
        route: "moduleLoading",
      },
    ])

    function updateActiveTab(newPath: string): void {
      activeTab.value = tabs.value.findIndex(tab => newPath.endsWith("/" + tab.route))
    }

    updateActiveTab(route.path)

    const router = useRouter()
    watch(() => route.path, updateActiveTab)
    watch(activeTab, async value => {
      console.log(tabs.value[value].route)
      await router.push({...route, path: `/ij/${tabs.value[value].route}`})
    })

    return {
      productConfigurator,
      projectConfigurator,
      machineConfigurator,
      timeRangeConfigurator,
      getProjectName,
      activeTab,
      tabs,
    }
  },
})
</script>
