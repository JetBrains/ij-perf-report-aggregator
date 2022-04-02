<template>
  <Toolbar class="!py-0 !px-2">
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

  <TabMenu
    :model="tabs"
    class="mb-1"
  />
  <router-view v-slot="{ Component }">
    <keep-alive>
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>

<script setup lang="ts">
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
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

export interface Tab {
  label: string
  to: string
}

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

const tabs = ref<Array<Tab>>([
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
}))
</script>
