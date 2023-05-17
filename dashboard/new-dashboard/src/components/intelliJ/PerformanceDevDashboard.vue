<template>
  <div class="flex flex-col gap-5">
    <DashboardToolbar
      :branch-configurator="branchConfigurator"
      :machine-configurator="machineConfigurator"
      :on-change-range="onChangeRange"
      :time-range-configurator="timeRangeConfigurator"
      :triggered-by-configurator="triggeredByConfigurator"
    />
    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <section>
          <GroupProjectsChart
            v-for="chart in charts"
            :key="chart.definition.label"
            :label="chart.definition.label"
            :measure="chart.definition.measure"
            :projects="chart.projects"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
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
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardToolbar from "../common/DashboardToolbar.vue"

provideReportUrlProvider(false)

const chartsDeclaration: Array<ChartDefinition> = [{
  labels: ["Indexing", "Scanning", "Updating time"],
  measures: ["indexing", "scanning", "updatingTime"],
  projects: ["intellij_sources/indexing", "intellij_commit/indexing"],
}, {
  labels: ["Find Usages Java"],
  measures: ["findUsages"],
  projects: ["intellij_sources/findUsages/Application_runReadAction", "intellij_sources/findUsages/LocalInspectionTool_getID",
    "intellij_sources/findUsages/PsiManager_getInstance", "intellij_sources/findUsages/PropertyMapping_value",
    "intellij_commit/findUsages/Application_runReadAction", "intellij_commit/findUsages/LocalInspectionTool_getID",
    "intellij_commit/findUsages/PsiManager_getInstance", "intellij_commit/findUsages/PropertyMapping_value"],
}, {
  labels: ["Find Usages Kotlin"],
  measures: ["findUsages"],
  projects: ["intellij_sources/findUsages/ActionsKt_runReadAction", "intellij_sources/findUsages/DynamicPluginListener_TOPIC", "intellij_sources/findUsages/Path_div",
    "intellij_sources/findUsages/Persistent_absolutePath", "intellij_sources/findUsages/RelativeTextEdit_rangeTo",
    "intellij_sources/findUsages/TemporaryFolder_invoke", "intellij_sources/findUsages/Project_guessProjectDir",
    "intellij_commit/findUsages/ActionsKt_runReadAction", "intellij_commit/findUsages/DynamicPluginListener_TOPIC", "intellij_commit/findUsages/Path_div",
    "intellij_commit/findUsages/Persistent_absolutePath", "intellij_commit/findUsages/RelativeTextEdit_rangeTo",
    "intellij_commit/findUsages/TemporaryFolder_invoke", "intellij_commit/findUsages/Project_guessProjectDir"],
}, {
  labels: ["Local Inspection"],
  measures: ["localInspections"],
  projects: ["intellij_sources/localInspection/java_file", "intellij_sources/localInspection/kotlin_file",
    "intellij_commit/localInspection/java_file", "intellij_commit/localInspection/kotlin_file"],
}, {
  labels: ["Completion: execution time"],
  measures: ["completion"],
  projects: ["intellij_sources/completion/java_file", "intellij_sources/completion/kotlin_file",
    "intellij_commit/completion/java_file", "intellij_commit/completion/kotlin_file"],
}, {
  labels: ["Completion: awt delay"],
  measures: ["test#average_awt_delay"],
  projects: ["intellij_sources/completion/java_file", "intellij_sources/completion/kotlin_file",
    "intellij_commit/completion/java_file", "intellij_commit/completion/kotlin_file"],
}]

const dbName = "perfintDev"
const dbTable = "idea"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
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
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator,
]

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const charts = combineCharts(chartsDeclaration)
const projects = chartsDeclaration.map(it => it.projects).flat(Number.POSITIVE_INFINITY) as Array<string>
const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getAccidentsFromMetaDb(warnings, projects, data)
})
</script>