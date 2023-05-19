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
        <slot
          :dashboard-configurators="dashboardConfigurators"
          :averages-configurators="averagesConfigurators"
          :warnings="warnings"
        />
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
import { DataQueryConfigurator } from "shared/src/dataQuery"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { Accident, getAccidentsFromMetaDb } from "shared/src/meta"
import { provide, ref, withDefaults } from "vue"
import { useRouter } from "vue-router"
import { containerKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import { Chart, extractUniqueProjects } from "../charts/DashboardCharts"
import DashboardToolbar from "./DashboardToolbar.vue"


interface PerformanceDashboardProps {
  dbName: string
  table: string
  initialMachine?: string
  persistentId: string
  withInstaller?: boolean
  charts?: Array<Chart>
  isBuildNumberExists?: boolean
}

const props = withDefaults(defineProps<PerformanceDashboardProps>(), {
  withInstaller: true,
  isBuildNumberExists: false,
  charts: null,
  initialMachine: null,
})

provideReportUrlProvider(props.withInstaller, props.isBuildNumberExists)

const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(props.dbName, props.table)
provide(serverConfiguratorKey, serverConfigurator)

const persistenceForDashboard = new PersistentStateManager(props.persistentId, {
  machine: props.initialMachine ?? "",
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const scenarioConfigurator = props.charts == null ? null : dimensionConfigurator(
  "project",
  serverConfigurator,
  null,
  true,
  [timeRangeConfigurator],
)
if (scenarioConfigurator != null) {
  scenarioConfigurator.selected.value = extractUniqueProjects(props.charts)
}

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const machineConfigurator = props.initialMachine == null ? null : new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  scenarioConfigurator == null ? [timeRangeConfigurator, branchConfigurator] : [timeRangeConfigurator, branchConfigurator, scenarioConfigurator],
)

const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  scenarioConfigurator == null ? [branchConfigurator, timeRangeConfigurator] : [branchConfigurator, timeRangeConfigurator, scenarioConfigurator],
)

const averagesConfigurators = [
  serverConfigurator,
  branchConfigurator,
  timeRangeConfigurator,
] as DataQueryConfigurator[]
if (machineConfigurator != null) {
  averagesConfigurators.push(machineConfigurator)
}

const dashboardConfigurators = [
  branchConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator,
] as DataQueryConfigurator[]

if (machineConfigurator != null) {
  dashboardConfigurators.push(machineConfigurator)
}

const releaseConfigurator = (props.withInstaller) ? new ReleaseNightlyConfigurator(persistenceForDashboard) : null
if (releaseConfigurator != null) {
  dashboardConfigurators.push(releaseConfigurator)
}


function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const projects = props.charts?.map(it => it.projects).flat(Number.POSITIVE_INFINITY) as Array<string>
const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getAccidentsFromMetaDb(warnings, projects, data)
})
</script>