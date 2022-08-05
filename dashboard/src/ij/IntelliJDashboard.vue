<template>
  <Dashboard>
    <template #toolbar>
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

    <TabMenu
      v-once
      class="mb-2"
      :model="tabs"
    />
    <router-view v-slot="{ Component }">
      <keep-alive>
        <component :is="Component" />
      </keep-alive>
    </router-view>
  </Dashboard>
</template>

<script setup lang="ts">
import { initDataComponent } from "../shared/DataQueryExecutor"
import { PersistentStateManager } from "../shared/PersistentStateManager"
import { chartDefaultStyle } from "../shared/chart"
import Dashboard from "../shared/components/Dashboard.vue"
import DimensionHierarchicalSelect from "../shared/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../shared/components/DimensionSelect.vue"
import TimeRangeSelect from "../shared/components/TimeRangeSelect.vue"
import { AggregationOperatorConfigurator } from "../shared/configurators/AggregationOperatorConfigurator"
import { dimensionConfigurator } from "../shared/configurators/DimensionConfigurator"
import { MachineConfigurator } from "../shared/configurators/MachineConfigurator"
import { ServerConfigurator } from "../shared/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "../shared/configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey } from "../shared/injectionKeys"
import { provideReportUrlProvider } from "../shared/lineChartTooltipLinkProvider"
import TabMenu from "../tailwind-ui/TabMenu.vue"
import { TabItem } from "../tailwind-ui/tabModel"
import { provide } from "vue"
import { useRouter } from "vue-router"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

provideReportUrlProvider()
provide(chartStyleKey, {
  ...chartDefaultStyle,
  // a lot of bars, as result, height of bar is not enough to make label readable
  barSeriesLabelPosition: "right",
})

// noinspection SpellCheckingInspection
const persistentStateManager = new PersistentStateManager("ij-dashboard", {
  product: "IU",
  project: "simple for IJ",
  machine: "macMini M1, 16GB",
}, useRouter())
const serverConfigurator = new ServerConfigurator("ij")
const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager)
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager)
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [productConfigurator])
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)

provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))

initDataComponent([
  serverConfigurator,
  productConfigurator,
  projectConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
])

const tabs: Array<TabItem> = [
  {
    label: "Pulse",
    to: "pulse",
  },
  {
    label: "Progress Over Time",
    to: "progressOverTime",
  },
  {
    label: "Module Loading",
    to: "moduleLoading",
  },
].map(it => {
  return {...it, to: `/ij/${it.to}`}
})
</script>
