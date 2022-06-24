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
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import Dashboard from "shared/src/components/Dashboard.vue"
import { AggregationOperatorConfigurator } from "shared/src/configurators/AggregationOperatorConfigurator"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import TabMenu from "tailwind-ui/src/TabMenu.vue"
import { provide } from "vue"
import { useRouter } from "vue-router"
import DimensionHierarchicalSelect from "../../shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../../shared/src/components/DimensionSelect.vue"
import TimeRangeSelect from "../../shared/src/components/TimeRangeSelect.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

export interface Tab {
  label: string
  to: string
}

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
  machine: "macMini 2018",
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

const tabs: Array<Tab> = [
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
