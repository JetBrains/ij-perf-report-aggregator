<template>
  <div class="flex flex-col gap-5">
    <DashboardToolbar
      :branch-configurator="branchConfigurator"
      :machine-configurator="machineConfigurator"
      :release-configurator="releaseConfigurator"
      :on-change-range="onChangeRange"
      :time-range-configurator="timeRangeConfigurator"
      :triggered-by-configurator="triggeredByConfigurator"
    >
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
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computedAsync } from "@vueuse/core"
import { provide, Ref, ref } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "../../configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { getDBType } from "../../shared/dbTypes"
import { accidentsKeys, containerKey, dashboardConfiguratorsKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import { Accident, getAccidentsFromMetaDb } from "../../util/meta"
import { Chart, extractUniqueProjects } from "../charts/DashboardCharts"
import PlotSettings from "../settings/PlotSettings.vue"
import DashboardToolbar from "./DashboardToolbar.vue"
import { PersistentStateManager } from "./PersistentStateManager"
import { DataQueryConfigurator } from "./dataQuery"
import { provideReportUrlProvider } from "./lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "./sideBar/InfoSidebar"
import { InfoDataPerformance } from "./sideBar/InfoSidebarPerformance"
import InfoSidebar from "./sideBar/InfoSidebarPerformance.vue"

interface PerformanceDashboardProps {
  dbName: string
  table: string
  initialMachine?: string | null
  persistentId: string
  withInstaller?: boolean
  charts?: Chart[] | null
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
const sidebarVm = new InfoSidebarImpl<InfoDataPerformance>(getDBType(props.dbName, props.table))

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(props.dbName, props.table)
provide(serverConfiguratorKey, serverConfigurator)

const persistenceForDashboard = new PersistentStateManager(
  props.persistentId,
  {
    machine: props.initialMachine ?? "",
    project: [],
    branch: "master",
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const scenarioConfigurator = props.charts == null ? null : dimensionConfigurator("project", serverConfigurator, null, true, [timeRangeConfigurator])
if (scenarioConfigurator != null && props.charts != null) {
  scenarioConfigurator.selected.value = extractUniqueProjects(props.charts)
}

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const machineConfigurator =
  props.initialMachine == null
    ? undefined
    : new MachineConfigurator(
        serverConfigurator,
        persistenceForDashboard,
        scenarioConfigurator == null ? [timeRangeConfigurator, branchConfigurator] : [timeRangeConfigurator, branchConfigurator, scenarioConfigurator]
      )

const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  scenarioConfigurator == null ? [branchConfigurator, timeRangeConfigurator] : [branchConfigurator, timeRangeConfigurator, scenarioConfigurator]
)

const averagesConfigurators = [serverConfigurator, branchConfigurator, timeRangeConfigurator] as DataQueryConfigurator[]
if (machineConfigurator != null) {
  averagesConfigurators.push(machineConfigurator)
}

const dashboardConfigurators = [branchConfigurator, timeRangeConfigurator, triggeredByConfigurator] as FilterConfigurator[]

if (machineConfigurator != null) {
  dashboardConfigurators.push(machineConfigurator)
}

const releaseConfigurator = props.withInstaller ? new ReleaseNightlyConfigurator(persistenceForDashboard) : undefined
if (releaseConfigurator != null) {
  dashboardConfigurators.push(releaseConfigurator)
}
provide(dashboardConfiguratorsKey, dashboardConfigurators)

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const updateConfigurators = (configurator: FilterConfigurator) => {
  dashboardConfigurators.push(configurator)
}

function getProjectAndProjectWithMetrics(charts: Chart[] | null): string[] {
  const projectsWithMetrics =
    charts?.flatMap((chart) => {
      const measures = Array.isArray(chart.definition.measure) ? chart.definition.measure : [chart.definition.measure]
      return chart.projects.flatMap((project) => {
        return measures.map((measure) => project + "/" + measure)
      })
    }) ?? []
  const projects = new Set(charts?.map((it) => it.projects).flat(Number.POSITIVE_INFINITY) as string[])
  return [...projectsWithMetrics, ...projects]
}

const warnings: Ref<Map<string, Accident[]> | undefined> = ref()

computedAsync(async () => {
  warnings.value = await getAccidentsFromMetaDb(getProjectAndProjectWithMetrics(props.charts), timeRangeConfigurator.value)
})
provide(accidentsKeys, warnings)
</script>
