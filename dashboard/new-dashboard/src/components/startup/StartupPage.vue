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
      Sidebar:
      <InputSwitch v-model="sidebarEnabled" />
    </template>
  </Toolbar>
  <main class="flex">
    <div
      ref="container"
      class="flex flex-1 flex-col gap-6 overflow-hidden"
    >
      <slot />
    </div>
    <InfoSidebarStartup />
  </main>
  <ChartTooltip ref="tooltip" />
</template>
<script setup lang="ts">
import { useStorage } from "@vueuse/core"
import { provide, Ref, ref, watch } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { chartStyleKey, chartToolTipKey, configuratorListKey, sidebarEnabledKey } from "../../shared/injectionKeys"
import { containerKey, sidebarStartupKey } from "../../shared/keys"
import ChartTooltip from "../charts/ChartTooltip.vue"
import DimensionSelect from "../charts/DimensionSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import MachineSelect from "../common/MachineSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { chartDefaultStyle } from "../common/chart"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "../common/sideBar/InfoSidebar"
import { InfoDataFromStartup } from "../common/sideBar/InfoSidebarStartup"
import InfoSidebarStartup from "../common/sideBar/InfoSidebarStartup.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

const container = ref<HTMLElement>()
const sidebarVm = new InfoSidebarImpl<InfoDataFromStartup>()

const tooltip = ref<typeof ChartTooltip>()
provide(chartToolTipKey, tooltip as Ref<typeof ChartTooltip>)
provide(containerKey, container)
provide(sidebarStartupKey, sidebarVm)

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
})

const dbName = "ij"
const dbTable = "report"

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
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
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
const configurators = [serverConfigurator, machineConfigurator, timeRangeConfigurator, productConfigurator, projectConfigurator, branchConfigurator, triggeredByConfigurator]

provide(configuratorListKey, configurators)

const sidebarEnabled = useStorage("sidebarEnabled", true)
watch(sidebarEnabled, (value) => {
  if (!value) {
    sidebarVm.close()
  }
})
provide(sidebarEnabledKey, sidebarEnabled)

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
