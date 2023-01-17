<template>
  <div class="flex flex-col gap-5">
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
          :release-configurator="releaseConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <DimensionSelect
          label="Tests"
          :selected-label="testsSelectLabelFormat"
          :dimension="scenarioConfigurator"
        >
          <template #icon>
            <ChartBarIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionSelect>
        <MeasureSelect
          title="Metrics"
          :selected-label="metricsSelectLabelFormat"
          :configurator="measureConfigurator"
        >
          <template #icon>
            <BeakerIcon class="w-4 h-4 text-gray-500" />
          </template>
        </MeasureSelect>
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        >
          <template #icon>
            <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionHierarchicalSelect>
      </template>
    </Toolbar>

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
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
          />
        </template>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import MeasureSelect from "shared/src/components/MeasureSelect.vue"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import { testsSelectLabelFormat, metricsSelectLabelFormat } from "../../shared/labels"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import LineChart from "../charts/LineChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "kotlin"
const initialMachine = "linux-blade-hetzner"
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
    branch: "master",
    project: [],
    measure: [],
  }, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator, branchConfigurator],
)
const scenarioConfigurator = dimensionConfigurator(
  "project",
  serverConfigurator,
  persistentStateManager,
  true,
  [branchConfigurator, timeRangeConfigurator],
)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistentStateManager,
  [branchConfigurator, timeRangeConfigurator],
)
const measureConfigurator = new MeasureConfigurator(
  serverConfigurator,
  persistentStateManager,
  [scenarioConfigurator, branchConfigurator],
  true,
  "line",
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistentStateManager)

const configurators = [
  serverConfigurator,
  scenarioConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator,
  releaseConfigurator,
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