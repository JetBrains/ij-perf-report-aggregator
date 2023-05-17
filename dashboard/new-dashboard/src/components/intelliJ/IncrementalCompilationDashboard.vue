<template>
  <div class="flex flex-col gap-5">
    <DashboardToolbar
      :branch-configurator="branchConfigurator"
      :machine-configurator="machineConfigurator"
      :release-configurator="releaseConfigurator"
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
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
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
import { ChartDefinition, combineCharts, extractUniqueProjects } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardToolbar from "../common/DashboardToolbar.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "idea"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const chartsDeclaration: Array<ChartDefinition> = [
  {
  labels: ["Elastic Rebuild/Build time"],
  measures: ["build_compilation_duration"],
  projects: ["incremental-build-java/build_incremental", "incremental-build-java/rebuild_initial"],
  }, {
    labels: ["IntelliJ Rebuild/Build time"],
    measures: ["build_compilation_duration"],
    projects: ["incremental-build-intellij/build_incremental", "incremental-build-intellij/rebuild_initial"],
  }, {
    labels: ["Coroutines Rebuild/Build time"],
    measures: ["build_compilation_duration"],
    projects: ["incremental-build-kotlin/build_incremental", "incremental-build-kotlin/rebuild_initial"],
  }, {
    labels: ["Youtrack JPS Rebuild/Build time"],
    measures: ["build_compilation_duration"],
    projects: ["incremental-build-youtrack-jps/build_incremental", "incremental-build-youtrack-jps/rebuild_initial"],
  }, {
    labels: ["Youtrack Gradle Rebuild/Build time"],
    measures: ["build_compilation_duration"],
    projects: ["incremental-build-youtrack-gradle/build_incremental", "incremental-build-youtrack-gradle/rebuild_initial"],
  }, {
    labels: ["Hub JPS Rebuild/Build time"],
    measures: ["build_compilation_duration"],
    projects: ["incremental-build-hub-jps/build_incremental", "incremental-build-hub-jps/rebuild_initial"],
  }, {
    labels: ["Hub Gradle Rebuild/Build time"],
    measures: ["build_compilation_duration"],
    projects: ["incremental-build-hub-gradle/build_incremental", "incremental-build-hub-gradle/rebuild_initial"],
  }
]

const charts = combineCharts(chartsDeclaration)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("idea_incremental_build_dashboard", {
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
  [timeRangeConfigurator]
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

const dashboardConfigurators = [
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator,
  triggeredByConfigurator,
]

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getAccidentsFromMetaDb(warnings, null, data)
})
</script>