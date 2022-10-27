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
          :branchConfigurator="branchConfigurator"
          :releaseConfigurator="releaseConfigurator"
          :triggeredByConfigurator="triggeredByConfigurator"
        />
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
        <section class="flex gap-6">
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'processingSpeed#JAVA'"
              :title="'Indexing Java (kB/s)'"
              :chart-color="'#bb2083'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'processingSpeed#Kotlin'"
              :title="'Indexing Kotlin (kB/s)'"
              :chart-color="'#9B51E0'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'completion\_%'"
              :is-like="true"
              :title="'Completion'"
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
              :configurators="[...averagesConfigurators, typingOnlyConfigurator]"
              :aggregated-measure="'test#average_awt_delay'"
              :title="'UI responsiveness during typing'"
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
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { DataQuery, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

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

const typingOnlyConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({f: "project", v: "%typing", o: "like"})
    return true
  },
  createObservable() {
    return null
  },
}

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