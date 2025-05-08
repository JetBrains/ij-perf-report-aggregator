<template>
  <div class="flex flex-col gap-5">
    <DashboardToolbar
      :branch-configurator="branchConfigurator"
      :machine-configurator="machineConfigurator"
      :release-configurator="releaseNightlyConfigurator"
      :on-change-range="onChangeRange"
      :time-range-configurator="timeRangeConfigurator"
      :triggered-by-configurator="triggeredByConfigurator"
      :test-mode-configurator="testModeConfigurator"
    >
      <template #configurator>
        <slot name="configurator" />
      </template>
      <template #toolbar>
        <PlotSettings @update:configurators="updateConfigurators" />
      </template>
    </DashboardToolbar>

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <slot :averages-configurators="averagesConfigurators" />
      </div>
      <InfoSidebar :timerange-configurator="timeRangeConfigurator" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { provide, useTemplateRef } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { nightly, ReleaseNightlyConfigurator, ReleaseType } from "../../configurators/ReleaseNightlyConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { accidentsConfiguratorKey, containerKey, dashboardConfiguratorsKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import { Chart, extractUniqueProjects } from "../charts/DashboardCharts"
import PlotSettings from "../settings/PlotSettings.vue"
import DashboardToolbar from "./DashboardToolbar.vue"
import { PersistentStateManager } from "./PersistentStateManager"
import { DataQueryConfigurator } from "./dataQuery"
import { provideReportUrlProvider } from "./lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "./sideBar/InfoSidebar"
import InfoSidebar from "./sideBar/InfoSidebar.vue"
import { AccidentsConfiguratorForDashboard } from "../../configurators/accidents/AccidentsConfiguratorForDashboard"
import { dbTypeStore } from "../../shared/dbTypes"
import { createTestModeConfigurator, defaultModeName } from "../../configurators/TestModeConfigurator"

interface PerformanceDashboardProps {
  dbName: string
  table: string
  initialMachine?: string | null
  persistentId: string
  withInstaller?: boolean
  charts?: Chart[] | null
  isBuildNumberExists?: boolean
  releaseConfigurator?: ReleaseType
  branch?: string | null
  initialMode?: string[]
}

const {
  dbName,
  table,
  initialMachine = null,
  persistentId,
  withInstaller = true,
  charts = null,
  isBuildNumberExists = false,
  releaseConfigurator = nightly,
  branch = "master",
  initialMode = defaultModeName,
} = defineProps<PerformanceDashboardProps>()

const container = useTemplateRef<HTMLElement>("container")
const router = useRouter()
const sidebarVm = new InfoSidebarImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, table)

provideReportUrlProvider(withInstaller, isBuildNumberExists)
provide(serverConfiguratorKey, serverConfigurator)

const persistenceForDashboard = new PersistentStateManager(
  persistentId,
  {
    machine: initialMachine ?? "",
    project: [],
    branch,
    releaseConfigurator,
    mode: initialMode,
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const scenarioConfigurator = charts == null ? null : dimensionConfigurator("project", serverConfigurator, null, true)
if (scenarioConfigurator != null && charts != null) {
  scenarioConfigurator.selected.value = extractUniqueProjects(charts)
}

const branchConfigurator = branch == null ? null : createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const filters = []
filters.push(timeRangeConfigurator)
if (scenarioConfigurator != null) {
  filters.push(scenarioConfigurator)
}
if (branchConfigurator != null) {
  filters.push(branchConfigurator)
}
const machineConfigurator = initialMachine == null ? undefined : new MachineConfigurator(serverConfigurator, persistenceForDashboard, filters)
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistenceForDashboard, filters)

const averagesConfigurators = [serverConfigurator, timeRangeConfigurator] as DataQueryConfigurator[]
if (machineConfigurator != null) {
  averagesConfigurators.push(machineConfigurator)
}
if (branchConfigurator != null) {
  averagesConfigurators.push(branchConfigurator)
}

const accidentsConfigurator = new AccidentsConfiguratorForDashboard(serverConfigurator.serverUrl, charts, timeRangeConfigurator)
provide(accidentsConfiguratorKey, accidentsConfigurator)

const dashboardConfigurators = [timeRangeConfigurator, triggeredByConfigurator, accidentsConfigurator] as FilterConfigurator[]
if (machineConfigurator != null) {
  dashboardConfigurators.push(machineConfigurator)
}
if (branchConfigurator != null) {
  dashboardConfigurators.push(branchConfigurator)
}

const releaseNightlyConfigurator = withInstaller ? new ReleaseNightlyConfigurator(persistenceForDashboard) : undefined
if (releaseNightlyConfigurator != null) {
  dashboardConfigurators.push(releaseNightlyConfigurator)
}

const testModeConfigurator = dbTypeStore().isModeSupported() ? createTestModeConfigurator(serverConfigurator, persistenceForDashboard, filters, "mode", true, initialMode) : null
if (testModeConfigurator != null) {
  dashboardConfigurators.push(testModeConfigurator)
}

provide(dashboardConfiguratorsKey, dashboardConfigurators)

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const updateConfigurators = (configurator: FilterConfigurator) => {
  dashboardConfigurators.push(configurator)
}
</script>
