<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="timeRangeConfigurator.timeRanges"
        :value="timeRangeConfigurator.value.value"
        :on-change="onChangeRange"
      />
      <BranchSelect
        :branch-configurator="branchConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
      />
      <DimensionSelect
        label="Product"
        :value-to-label="(it: string) => productCodeToName.get(it) ?? it"
        :dimension="productConfigurator"
      />
      <DimensionSelect
        label="Project"
        :value-to-label="getProjectName"
        :dimension="projectConfigurator"
      />
      <MachineSelect :machine-configurator="machineConfigurator" />
      <slot name="toolbar" />
    </template>
    <template #end>
      <PlotSettings @update:configurators="updateConfigurators" />
    </template>
  </Toolbar>
  <main class="flex">
    <div
      ref="container"
      class="flex flex-1 flex-col gap-5 overflow-hidden pt-5"
    >
      <slot />
    </div>
    <InfoSidebarStartup />
  </main>
  <ChartTooltip ref="tooltip" />
</template>
<script setup lang="ts">
import { provide, Ref, ref } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getDBType } from "../../shared/dbTypes"
import { chartStyleKey, chartToolTipKey, configuratorListKey } from "../../shared/injectionKeys"
import { containerKey, sidebarStartupKey } from "../../shared/keys"
import ChartTooltip from "../charts/ChartTooltip.vue"
import DimensionSelect from "../charts/DimensionSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import MachineSelect from "../common/MachineSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { chartDefaultStyle } from "../common/chart"
import { DataQueryConfigurator } from "../common/dataQuery"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "../common/sideBar/InfoSidebar"
import { InfoDataFromStartup } from "../common/sideBar/InfoSidebarStartup"
import InfoSidebarStartup from "../common/sideBar/InfoSidebarStartup.vue"
import PlotSettings from "../settings/PlotSettings.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

const container = ref<HTMLElement>()

const tooltip = ref<typeof ChartTooltip>()
provide(chartToolTipKey, tooltip as Ref<typeof ChartTooltip>)
provide(containerKey, container)

const productCodeToName = new Map([
  ["DB", "DataGrip"],
  ["IU", "IntelliJ IDEA"],
  ["PS", "PhpStorm"],
  ["WS", "WebStorm"],
  ["GO", "GoLand"],
  ["PY", "PyCharm Professional"],
  ["RM", "RubyMine"],
  ["CL", "CLion"],
])

provideReportUrlProvider()
provide(chartStyleKey, {
  ...chartDefaultStyle,
})

const dbName = "ij"
const dbTable = "report"

const sidebarVm = new InfoSidebarImpl<InfoDataFromStartup>(getDBType(dbName, dbTable))
provide(sidebarStartupKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager(
  "ij-dashboard",
  {
    product: "IU",
    project: "simple for IJ",
    machine: "macMini M1, 16GB",
    branch: "master",
  },
  useRouter()
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager, false, [timeRangeConfigurator, branchConfigurator])
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager, [
  productConfigurator,
  timeRangeConfigurator,
  branchConfigurator,
])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  productConfigurator,
  projectConfigurator,
  branchConfigurator,
  triggeredByConfigurator,
] as DataQueryConfigurator[]

provide(configuratorListKey, configurators)

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>
