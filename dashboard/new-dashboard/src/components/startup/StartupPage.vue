<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="TimeRangeConfigurator.timeRanges"
        :value="timeRangeConfigurator.value.value"
        :on-change="onChangeRange"
      >
        <template #icon>
          <CalendarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </TimeRangeSelect>
      <BranchSelect
        :branch-configurator="branchConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
      />
      <DimensionSelect
        label="Product"
        :value-to-label="(it: string ) => productCodeToName.get(it) ?? it"
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
      >
        <template #icon>
          <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
        </template>
      </DimensionHierarchicalSelect>
      <slot name="toolbar" />
    </template>
  </Toolbar>
  <slot />
  <ChartTooltip
    ref="tooltip"
  />
</template>
<script setup lang="ts">

import { provide, Ref, ref } from "vue"
import { useRouter } from "vue-router"
import { AggregationOperatorConfigurator } from "../../configurators/AggregationOperatorConfigurator"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey, chartToolTipKey, configuratorListKey } from "../../shared/injectionKeys"
import ChartTooltip from "../charts/ChartTooltip.vue"
import DimensionHierarchicalSelect from "../charts/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../charts/DimensionSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { chartDefaultStyle } from "../common/chart"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

const productCodeToName = new Map([
  ["DB", "DataGrip"],
  ["IU", "IntelliJ IDEA"],
  ["PS", "PhpStorm"],
  ["WS", "WebStorm"],
  ["GO", "GoLand"],
  ["PY", "PyCharm Professional"],
  ["RM", "RubyMine"],
])


provideReportUrlProvider()
provide(chartStyleKey, {
  ...chartDefaultStyle,
  // a lot of bars, as result, height of bar is not enough to make label readable
  barSeriesLabelPosition: "right",
})
const tooltip = ref<typeof ChartTooltip>()
provide(chartToolTipKey, tooltip as Ref<typeof ChartTooltip>)

const dbName = "ij"
const dbTable = "report"

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager("ij-dashboard", {
  product: "IU",
  project: "simple for IJ",
  machine: "macMini M1, 16GB",
  branch: "master",
}, useRouter())

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator, branchConfigurator],
)
const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager, false, [timeRangeConfigurator, branchConfigurator])
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistentStateManager,
  [branchConfigurator, timeRangeConfigurator],
)
const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  productConfigurator,
  projectConfigurator,
  branchConfigurator,
  triggeredByConfigurator
]

provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))
provide(configuratorListKey, configurators)

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}
</script>