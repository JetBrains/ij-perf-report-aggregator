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
        <section>
          <GroupProjectsChart
            v-for="chart in charts"
            :key="chart.definition.label"
            :label="chart.definition.label"
            :measure="chart.definition.measure"
            :projects="chart.projects"
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
import { ChartDefinition, combineCharts, extractUniqueProjects } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "idea"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const chartsDeclaration: Array<ChartDefinition> = [{
  labels: ["Indexing", "Scanning", "Number of indexing runs"],
  measures: ["indexing", "scanning", "numberOfIndexingRuns"],
  projects: ["keycloak_release_20/indexing"],
}, {
  labels: ["Local Inspection", "First Code Analysis"],
  measures: ["localInspections", "firstCodeAnalysis"],
  projects: ["keycloak_release_20/localInspection/AuthenticationManagementResource", "keycloak_release_20/localInspection/IdentityBrokerService",
    "keycloak_release_20/localInspection/JpaUserProvider", "keycloak_release_20/localInspection/RealmAdminResource",
    "keycloak_release_20/localInspection/QuarkusRuntimePomXml", "keycloak_release_20/localInspection/RootPomXml",
    "keycloak_release_20/localInspection/CorePomXml"],
}, {
  labels: ["Completion"],
  measures: ["completion"],
  projects: ["keycloak_release_20/completion/AuthenticationManagementResource", "keycloak_release_20/completion/IdentityBrokerService",
    "keycloak_release_20/completion/JpaUserProvider", "keycloak_release_20/completion/RealmAdminResource",
    "keycloak_release_20/completion/QuarkusRuntimePomXml", "keycloak_release_20/completion/RootPomXml",
    "keycloak_release_20/completion/CorePomXml"],
}, {
  labels: ["Show Intentions (average awt delay)", "Show Intentions (showQuickFixes)"],
  measures: ["test#average_awt_delay", "showQuickFixes"],
  projects: ["keycloak_release_20/showIntentions/AuthenticationManagementResource", "keycloak_release_20/showIntentions/IdentityBrokerService",
    "keycloak_release_20/showIntentions/JpaUserProvider", "keycloak_release_20/showIntentions/RealmAdminResource"],
}, {
  labels: ["Typing", "Typing (firstCodeAnalysis)", "Typing (typingCodeAnalyzing)", "Typing (average_awt_delay)", "Typing (max_awt_delay)"],
  measures: ["typing", "firstCodeAnalysis", "typingCodeAnalyzing", "test#average_awt_delay", "test#max_awt_delay"],
  projects: ["keycloak_release_20/typing/ClientEntity", "keycloak_release_20/typing/PolicyEntity"],
}]

const charts = combineCharts(chartsDeclaration)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("idea_ultimate_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)
const scenarioConfigurator = dimensionConfigurator(
  "project",
  serverConfigurator,
  null,
  true,
)
scenarioConfigurator.selected.value = extractUniqueProjects(chartsDeclaration)

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator, scenarioConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator, branchConfigurator, scenarioConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistenceForDashboard)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [branchConfigurator, timeRangeConfigurator, scenarioConfigurator],
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