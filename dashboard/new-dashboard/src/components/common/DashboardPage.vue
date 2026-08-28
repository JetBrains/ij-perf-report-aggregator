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

    <SelectButton
      v-if="canCompare"
      v-model="renderMode"
      :options="renderModeOptions"
      option-label="label"
      option-value="value"
      :allow-empty="false"
      class="self-start"
    />

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <CompareTable v-if="renderMode === 'compare'" />
        <!-- v-show keeps the slot mounted in compare mode (chart components must stay mounted to keep -->
        <!-- their CompareSectionsRegistry entries alive) but hides it so non-compare-aware dashboard -->
        <!-- content doesn't sit alongside the table. -->
        <div
          v-show="renderMode === 'charts'"
          class="flex flex-col gap-6"
        >
          <slot
            :averages-configurators="averagesConfigurators"
            :machine-configurator="machineConfigurator"
            :branch-configurator="branchConfigurator"
          />
        </div>
      </div>
      <InfoSidebar :timerange-configurator="timeRangeConfigurator" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, provide, ref, useTemplateRef, watch } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dashboardProjectsFilter } from "../../configurators/DashboardProjectsFilter"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { nightly, ReleaseNightlyConfigurator, ReleaseType } from "../../configurators/ReleaseNightlyConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import {
  accidentsConfiguratorKey,
  branchConfiguratorKey,
  compareSectionsRegistryKey,
  containerKey,
  dashboardConfiguratorsKey,
  renderModeKey,
  serverConfiguratorKey,
  sidebarVmKey,
} from "../../shared/keys"
import { Chart, extractUniqueProjects } from "../charts/DashboardCharts"
import CompareTable from "../charts/CompareTable.vue"
import { CompareSectionsRegistry, RenderMode } from "../charts/compareMode"
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
  withMode?: boolean
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
  withMode = true,
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

const compareRegistry = new CompareSectionsRegistry()
// Most dashboards declare their projects inline on the charts in the slot rather than through the
// `charts` prop, so take the union of both: every chart registers itself here as it mounts.
const chartsProjects = charts == null ? [] : extractUniqueProjects(charts)
const dashboardProjects = computed(() => [...new Set([...chartsProjects, ...compareRegistry.sections.value.flatMap((section) => section.projects)])].toSorted())

const branchConfigurator = branch == null ? null : createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const filters = []
filters.push(timeRangeConfigurator, dashboardProjectsFilter(dashboardProjects))
if (branchConfigurator != null) {
  filters.push(branchConfigurator)
}
const testModeConfigurator =
  withMode && dbTypeStore().isModeSupported() ? createTestModeConfigurator(serverConfigurator, persistenceForDashboard, filters, "mode", true, initialMode) : null

// The machine list is narrowed by the selected mode as well: a mode can run on a single hardware
// class (goland's `wsl` only on windows-azure), so the groups differ per mode.
const machineFilters = testModeConfigurator == null ? filters : [...filters, testModeConfigurator]
const machineConfigurator = initialMachine == null ? undefined : new MachineConfigurator(serverConfigurator, persistenceForDashboard, machineFilters)
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

const dashboardConfigurators = [timeRangeConfigurator, triggeredByConfigurator] as FilterConfigurator[]
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

if (testModeConfigurator != null) {
  dashboardConfigurators.push(testModeConfigurator)
}

provide(dashboardConfiguratorsKey, dashboardConfigurators)
provide(branchConfiguratorKey, branchConfigurator)

const renderMode = ref<RenderMode>("charts")
const renderModeOptions = [
  { label: "Charts", value: "charts" },
  { label: "Compare with base", value: "compare" },
]
provide(renderModeKey, renderMode)
provide(compareSectionsRegistryKey, compareRegistry)

const canCompare = computed(() => {
  if (branchConfigurator == null) return false
  const sel = branchConfigurator.selected.value
  const count = sel == null ? 0 : Array.isArray(sel) ? sel.length : 1
  return count >= 2
})

// Drop back to charts when the user narrows to a single branch: comparing master to itself is degenerate.
watch(canCompare, (allowed) => {
  if (!allowed && renderMode.value === "compare") {
    renderMode.value = "charts"
  }
})

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const updateConfigurators = (configurator: FilterConfigurator) => {
  dashboardConfigurators.push(configurator)
}
</script>
