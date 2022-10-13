<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect
          :ranges="TimeRangeConfigurator.timeRanges"
          :value="timeRangeConfigurator.value.value"
          :on-change="onChangeRange"
        />
        <DimensionSelect
          label="Branch"
          :dimension="branchConfigurator"
        />
        <DimensionSelect
          label="Scenarios"
          tooltip="Scenarios"
          :dimension="scenarioConfigurator"
        />
        <MeasureSelect
          :configurator="measureConfigurator"
        />
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        />
        <DimensionSelect
          v-if="releaseConfigurator != null"
          label="Nightly/Release"
          :dimension="releaseConfigurator"
        />
        <DimensionSelect
          label="Triggered by"
          :dimension="triggeredByConfigurator"
        />
      </template>
    </Toolbar>

    <main class="flex">
      <div class="flex flex-1 flex-col gap-6 overflow-hidden" ref="container">

      </div>
      <InfoSidebar/>
    </main>
  </div>
</template>

<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import InfoSidebar from "../InfoSidebar.vue"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import MeasureSelect from "shared/src/components/MeasureSelect.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "idea"
const initialMachine = "macMini Intel 3.2, 16GB"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager(
  `${dbName}-${dbTable}-dashboard`,
  {
    machine: initialMachine,
    project: [],
    branch: "master",
  },
  router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator(
  "branch",
  serverConfigurator,
  persistentStateManager,
  true,
  [timeRangeConfigurator],
  (a, _) => {
    return a.includes("/") ? 1 : -1
  })
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator, branchConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistentStateManager)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistentStateManager,
  [branchConfigurator, timeRangeConfigurator],
)
const scenarioConfigurator = dimensionConfigurator(
  "project",
  serverConfigurator,
  persistentStateManager,
  true,
  [branchConfigurator, timeRangeConfigurator]
)
const measureConfigurator = new MeasureConfigurator(
  serverConfigurator,
  persistentStateManager,
  [scenarioConfigurator, branchConfigurator],
  true,
  "line",
)

const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator,
  triggeredByConfigurator,
]

function onChangeRange(value: string) {
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