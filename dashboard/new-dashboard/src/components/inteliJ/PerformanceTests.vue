<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect
          :ranges="TimeRangeConfigurator.timeRanges"
          :value="timeRangeConfigurator.value.value"
          :on-change="onChangeRange"
        >
          <template v-slot:icon>
            <CalendarIcon class="w-4 h-4 text-gray-500" />
          </template>
        </TimeRangeSelect>
        <DimensionSelect
          label="Branch"
          :selected-label="branchesSelectLabelFormat"
          :dimension="branchConfigurator"
        >
          <template v-slot:icon>
            <!--Temporary use custom icon, heroicons or primevue don't have such-->
            <div class="w-4 h-4 text-gray-500">
              <BranchIcon />
            </div>
          </template>
        </DimensionSelect>
        <DimensionSelect
          label="Metrics"
          :selected-label="metricsSelectLabelFormat"
          :dimension="scenarioConfigurator"
        >
          <template v-slot:icon>
            <ChartBarIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionSelect>
        <MeasureSelect
          title="Tests"
          :selected-label="testsSelectLabelFormat"
          :configurator="measureConfigurator"
        >
          <template v-slot:icon>
            <BeakerIcon class="w-4 h-4 text-gray-500" />
          </template>
        </MeasureSelect>
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        >
          <template v-slot:icon>
            <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionHierarchicalSelect>
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
import LineChart from '../charts/LineChart.vue'
import { branchesSelectLabelFormat, testsSelectLabelFormat, metricsSelectLabelFormat } from '../../shared/labels'
import BranchIcon from "../common/BranchIcon.vue"

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
    branch: "master",
    project: [],
    measure: [],
  }, router)

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