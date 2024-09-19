<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <TimeRangeSelect :timerange-configurator="timeRangeConfigurator" />
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
        <MeasureSelect
          :configurator="measureConfigurator"
          title="Metrics"
          :selected-label="metricsSelectLabelFormat"
        />
        <MachineSelect :machine-configurator="machineConfigurator" />
      </template>
      <template #end>
        <PlotSettings @update:configurators="updateConfigurators" />
      </template>
    </StickyToolbar>

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-5 overflow-hidden pt-5"
      >
        <template
          v-for="measure in measureConfigurator.selected.value"
          :key="measure"
        >
          <LineChart
            :title="measure"
            :measures="[measure]"
            :configurators="configurators"
            :skip-zero-values="false"
            :with-measure-name="true"
          />
        </template>
      </div>
      <InfoSidebar :timerange-configurator="timeRangeConfigurator" />
    </main>
  </div>
</template>
<script setup lang="ts">
import { provide, ref, useTemplateRef } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { configuratorListKey } from "../../shared/injectionKeys"
import { accidentsConfiguratorKey, containerKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import { metricsSelectLabelFormat } from "../../shared/labels"
import DimensionSelect from "../charts/DimensionSelect.vue"
import LineChart from "../charts/LineChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import MachineSelect from "../common/MachineSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import StickyToolbar from "../common/StickyToolbar.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { DataQueryConfigurator } from "../common/dataQuery"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "../common/sideBar/InfoSidebar"
import InfoSidebar from "../common/sideBar/InfoSidebar.vue"
import PlotSettings from "../settings/PlotSettings.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"
import { AccidentsConfiguratorForStartup } from "../../configurators/accidents/AccidentsConfiguratorForStartup"

const { withInstaller } = defineProps<{
  withInstaller: boolean
}>()

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

provideReportUrlProvider(withInstaller)

const dbName = withInstaller ? "ij" : "ijDev"
const dbTable = "report"
const container = useTemplateRef<HTMLElement>("container")

const sidebarVm = new InfoSidebarImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, dbTable)
provide(serverConfiguratorKey, serverConfigurator)
const persistentStateManager = new PersistentStateManager(
  "startup-explore",
  {
    machine: "Linux Munich i7-13700, 64Gb",
    branch: "master",
  },
  useRouter()
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager, false, [timeRangeConfigurator, branchConfigurator, machineConfigurator])
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager, [
  timeRangeConfigurator,
  branchConfigurator,
  productConfigurator,
  machineConfigurator,
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

const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [
  timeRangeConfigurator,
  branchConfigurator,
  productConfigurator,
  projectConfigurator,
  machineConfigurator,
])
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
</script>
