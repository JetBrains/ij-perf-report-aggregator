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
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="indexing"
            :projects="['intellij_sources/indexing', 'intellij_commit/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning"
            measure="scanning"
            :projects="['intellij_sources/indexing', 'intellij_commit/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Rebuild"
            measure="build_compilation_duration"
            :projects="['intellij_sources/rebuild']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Find Usages Java"
            measure="findUsages"
            :projects="['intellij_sources/findUsages/Application_runReadAction', 'intellij_sources/findUsages/LocalInspectionTool_getID',
                        'intellij_sources/findUsages/PsiManager_getInstance', 'intellij_sources/findUsages/PropertyMapping_value',
                        'intellij_commit/findUsages/Application_runReadAction', 'intellij_commit/findUsages/LocalInspectionTool_getID',
                        'intellij_commit/findUsages/PsiManager_getInstance', 'intellij_commit/findUsages/PropertyMapping_value']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Find Usages Kotlin"
            measure="findUsages"
            :projects="['intellij_sources/findUsages/ActionsKt_runReadAction', 'intellij_sources/findUsages/DynamicPluginListener_TOPIC', 'intellij_sources/findUsages/Path_div',
                        'intellij_sources/findUsages/Persistent_absolutePath', 'intellij_sources/findUsages/RelativeTextEdit_rangeTo',
                        'intellij_sources/findUsages/TemporaryFolder_invoke', 'intellij_sources/findUsages/Project_guessProjectDir',
                        'intellij_commit/findUsages/ActionsKt_runReadAction', 'intellij_commit/findUsages/DynamicPluginListener_TOPIC', 'intellij_commit/findUsages/Path_div',
                        'intellij_commit/findUsages/Persistent_absolutePath', 'intellij_commit/findUsages/RelativeTextEdit_rangeTo',
                        'intellij_commit/findUsages/TemporaryFolder_invoke', 'intellij_commit/findUsages/Project_guessProjectDir']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Local Inspection"
            measure="localInspections"
            :projects="['intellij_sources/localInspection/java_file','intellij_sources/localInspection/kotlin_file',
                        'intellij_commit/localInspection/java_file','intellij_commit/localInspection/kotlin_file']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Completion: execution time"
            measure="completion"
            :projects="['intellij_sources/completion/java_file','intellij_sources/completion/kotlin_file',
                        'intellij_commit/completion/java_file','intellij_commit/completion/kotlin_file']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Completion: awt delay"
            measure="test#average_awt_delay"
            :projects="['intellij_sources/completion/java_file','intellij_sources/completion/kotlin_file',
                        'intellij_commit/completion/java_file','intellij_commit/completion/kotlin_file']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
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
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

const dbName = "perfintDev"
const dbTable = "idea"
const initialMachine = "linux-blade-hetzner"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("ideaDev_dashboard", {
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


const dashboardConfigurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
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