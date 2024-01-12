<template>
  <div class="flex flex-col gap-5">
    <Toolbar :class="isSticky ? 'stickyToolbar' : 'customToolbar'">
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
        <CopyLink :timerange-configurator="timeRangeConfigurator" />
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
        <slot :configurators="configurators" />
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, provide, ref } from "vue"
import { useRouter } from "vue-router"
import { AccidentsConfiguratorForStartup } from "../../configurators/AccidentsConfigurator"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getDBType } from "../../shared/dbTypes"
import { configuratorListKey } from "../../shared/injectionKeys"
import { accidentsConfiguratorKey, containerKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import DimensionSelect from "../charts/DimensionSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import MachineSelect from "../common/MachineSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { DataQueryConfigurator } from "../common/dataQuery"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "../common/sideBar/InfoSidebar"
import { InfoDataPerformance } from "../common/sideBar/InfoSidebarPerformance"
import InfoSidebar from "../common/sideBar/InfoSidebarPerformance.vue"
import CopyLink from "../settings/CopyLink.vue"
import PlotSettings from "../settings/PlotSettings.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

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

const dbName = "ij"
const dbTable = "report"
const container = ref<HTMLElement>()

const sidebarVm = new InfoSidebarImpl<InfoDataPerformance>(getDBType(dbName, dbTable))

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, dbTable)
provide(serverConfiguratorKey, serverConfigurator)
const persistentStateManager = new PersistentStateManager(
  "startup-pulse",
  {
    machine: "Linux Munich i7-3770, 32Gb",
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

const accidentsConfigurator = new AccidentsConfiguratorForStartup(
  serverConfigurator.serverUrl,
  productConfigurator.selected,
  projectConfigurator.selected,
  ref(null),
  timeRangeConfigurator
)
provide(accidentsConfiguratorKey, accidentsConfigurator)

const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  productConfigurator,
  projectConfigurator,
  branchConfigurator,
  triggeredByConfigurator,
  accidentsConfigurator,
] as DataQueryConfigurator[]

provide(configuratorListKey, configurators)

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const isSticky = ref(false)
const checkIfSticky = () => (isSticky.value = window.scrollY > 100)
onMounted(() => {
  window.addEventListener("scroll", checkIfSticky)
})
onUnmounted(() => {
  window.removeEventListener("scroll", checkIfSticky)
})
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}

.stickyToolbar {
  top: 0rem;
  padding: 0.7rem 0.7rem 0.7rem 0.7rem;
  border-radius: 0;
  position: sticky;
  z-index: 100;
}
</style>
