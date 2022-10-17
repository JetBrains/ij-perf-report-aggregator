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
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        >
          <template v-slot:icon>
            <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionHierarchicalSelect>
        <DimensionSelect
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
        <section class="flex gap-6">
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'completion'"
              :title="'Completion'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'typing'"
              :title="'Typing'"
              :chart-color="'#9B51E0'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'indexing'"
              :title="'Indexing'"
              :chart-color="'#219653'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'test#max_awt_delay'"
              :title="'UI responsiveness'"
              :chart-color="'#F2994A'"
            />
          </div>
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing Long"
            measure="indexing"
            :projects="['community/indexing', 'lock-free-vfs-record-storage-intellij_sources/indexing', 'intellij_sources/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>

        <section class="flex gap-x-6">
          <div class="flex-1">
            <GroupProjectsChart
              label="Kotlin Builder Long"
              measure="kotlin_builder_time"
              :projects="['community/rebuild','intellij_sources/rebuild']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1">
            <GroupProjectsChart
              label="Rebuild Long"
              measure="build_compilation_duration"
              :projects="['community/rebuild','intellij_sources/rebuild']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>
      </div>
      <InfoSidebar />
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
import { provide, ref, shallowRef } from "vue"
import { useRouter } from "vue-router"
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import InfoSidebar from "../InfoSidebar.vue"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import { branchesSelectLabelFormat } from "../../shared/labels"
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
const persistenceForDashboard = new PersistentStateManager("idea_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const branchConfigurator = dimensionConfigurator(
  "branch",
  serverConfigurator,
  persistenceForDashboard,
  true,
  [timeRangeConfigurator],
  (a, _) => {
    return a.includes("/") ? 1 : -1
  })
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator, branchConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistenceForDashboard)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [branchConfigurator, timeRangeConfigurator],
)

const averagesConfigurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
]

const dashboardConfigurators = [
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