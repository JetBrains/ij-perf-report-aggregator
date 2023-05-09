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
          :triggered-by-configurator="triggeredByConfigurator"
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
          <div class="flex-1 min-w-0">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'completion\_%'"
              :aggregated-project="'%\_k1'"
              :is-like="true"
              :title="'completion K1'"
            />
          </div>
          <div class="flex-1 min-w-0">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'completion\_%'"
              :aggregated-project="'%\_k2'"
              :is-like="true"
              :title="'completion K2'"
            />
          </div>
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Completion kts"
              measure="completion#mean_value"
              :projects="['arrow/completion/completion_kts_with_library_cache_k1', 'kotlin_lang/completion/completion_kts_with_library_cache_k1']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Indexing K1"
              measure="indexing"
              :projects="['kotlin_empty/indexing_k1', 'intellij_commit/indexing_k1', 'kotlin_lang/indexing_k1']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Indexing K2"
              measure="indexing"
              :projects="['kotlin_empty/indexing_k2', 'intellij_commit/indexing_k2', 'kotlin_lang/indexing_k2']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>
        </section>

        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Completion mean value on hello-world"
              measure="completion#mean_value"
              :projects="['kotlin_empty/completion/empty_place_with_library_cache_k1', 'kotlin_empty/completion/empty_place_with_library_cache_k2']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>
        </section>

        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Completion mean value K1"
              measure="completion#mean_value"
              :projects="[
                'intellij_commit/completion/empty_place_with_library_cache_k1',
                'intellij_commit/completion/after_parameter_with_library_cache_k1',
                'kotlin_lang/completion/after_parameter_with_library_cache_k1',
                'kotlin_lang/completion/empty_place_with_library_cache_k1'
              ]"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Completion mean value K2"
              measure="completion#mean_value"
              :projects="[
                'intellij_commit/completion/empty_place_with_library_cache_k2',
                'intellij_commit/completion/after_parameter_with_library_cache_k2',
                'kotlin_lang/completion/after_parameter_with_library_cache_k2',
                'kotlin_lang/completion/empty_place_with_library_cache_k2'
              ]"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>
        </section>

        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Completion mean value with typing"
              measure="completion#mean_value"
              :projects="[
                'intellij_commit/completion/empty_place_typing_with_library_cache_k1',
                'intellij_commit/completion/empty_place_typing_with_library_cache_k2',
                'kotlin_empty/completion/empty_place_typing_with_library_cache_k1',
                'kotlin_empty/completion/empty_place_typing_with_library_cache_k2',
                'kotlin_lang/completion/empty_place_typing_with_library_cache_k1',
                'kotlin_lang/completion/empty_place_typing_with_library_cache_k2'
              ]"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Completion mean value Kotlin language server K1"
              measure="completion#mean_value"
              :projects="[
                'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k1',
                'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k1',
                'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k1',
                'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k1',
              ]"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
            />
          </div>
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Completion mean value Kotlin language server K2"
              measure="completion#mean_value"
              :projects="[
                'kotlin_language_server/completion/Completions_emptyPlace_completions_typing_with_library_cache_k2',
                'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_typing_with_library_cache_k2',
                'kotlin_language_server/completion/Completions_emptyPlace_completions_with_library_cache_k2',
                'kotlin_language_server/completion/QuickFixesTest_emptyPlace_completions_with_library_cache_k2',
              ]"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
              :accidents="warnings"
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
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { Accident, getAccidentsFromMetaDb } from "shared/src/meta"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../../shared/keys"
import InfoSidebar from "../../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../../InfoSidebarVm"
import AggregationChart from "../../charts/AggregationChart.vue"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import BranchSelect from "../../common/BranchSelect.vue"
import TimeRangeSelect from "../../common/TimeRangeSelect.vue"

provideReportUrlProvider(false)

const dbName = "perfintDev"
const dbTable = "kotlin"
const initialMachine = "linux-blade-hetzner"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("kotlin_completion_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator, branchConfigurator],
)
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
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator,
]

const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getAccidentsFromMetaDb(warnings, null, data as TimeRange)
})

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