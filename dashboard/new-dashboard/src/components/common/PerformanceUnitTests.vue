<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <CopyLink :timerange-configurator="timeRangeConfigurator" />
        <TimeRangeSelect :timerange-configurator="timeRangeConfigurator" />
        <BranchSelect
          :branch-configurator="branchConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <DimensionSelect
          label="Tests"
          :selected-label="testsSelectLabelFormat"
          :dimension="scenarioConfigurator"
        >
          <template #icon>
            <ChartBarIcon class="w-4 h-4" />
          </template>
        </DimensionSelect>
        <MeasureSelect
          title="Metrics"
          :selected-label="metricsSelectLabelFormat"
          :configurator="measureConfigurator"
        >
          <template #icon>
            <BeakerIcon class="w-4 h-4" />
          </template>
        </MeasureSelect>
        <MachineSelect :machine-configurator="machineConfigurator" />
      </template>
      <template #end>
        <PlotSettings @update:configurators="updateConfigurators" />
      </template>
    </StickyToolbar>
    <main class="flex">
      <div
        v-if="scenarios.length > 0 && measureConfigurator.selected.value != null && measureConfigurator.selected.value.length > 0"
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <template
          v-for="scenario in scenarios"
          :key="scenario"
        >
          <GroupProjectsChart
            :measure="measureConfigurator.selected.value"
            :projects="[scenario]"
            :label="scenario"
            :can-be-closed="true"
            @chart-closed="onChartClosed"
          />
        </template>
      </div>
      <div
        v-else-if="scenarios.length === 0"
        class="flex flex-1 items-center justify-center text-gray-400 dark:text-gray-500 py-8"
      >
        Select a test and a metric to see charts
      </div>
      <div
        v-else
        class="flex flex-1 items-center justify-center text-gray-400 dark:text-gray-500 py-8"
      >
        Select a metric to see charts
      </div>
      <InfoSidebar :timerange-configurator="timeRangeConfigurator" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { Observable } from "rxjs"
import { computed, provide, useTemplateRef, watch } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { accidentsConfiguratorKey, containerKey, dashboardConfiguratorsKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import { testsSelectLabelFormat, metricsSelectLabelFormat } from "../../shared/labels"
import { MAIN_METRICS } from "../../util/mainMetrics"
import DimensionSelect from "../charts/DimensionSelect.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import CopyLink from "../settings/CopyLink.vue"
import PlotSettings from "../settings/PlotSettings.vue"
import MachineSelect from "./MachineSelect.vue"
import { PersistentStateManager } from "./PersistentStateManager"
import StickyToolbar from "./StickyToolbar.vue"
import { DataQuery, DataQueryConfigurator } from "./dataQuery"
import { provideReportUrlProvider } from "./lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "./sideBar/InfoSidebar"
import InfoSidebar from "./sideBar/InfoSidebar.vue"
import { AccidentsConfiguratorForTests } from "../../configurators/accidents/AccidentsConfiguratorForTests"

export interface PerformanceUnitTestsProps {
  dbName: string
  table: string
  initialMachine: string
  withInstaller?: boolean
  projectFilter?: string
  persistentId?: string
  preselectAll?: boolean
}

const { dbName, table, initialMachine, withInstaller = true, projectFilter, persistentId, preselectAll = false } = defineProps<PerformanceUnitTestsProps>()

class ProjectLikeFilter implements FilterConfigurator {
  constructor(private readonly pattern: string) {}

  configureFilter(query: DataQuery): boolean {
    query.addFilter({ f: "", q: `project like '${this.pattern}'` })
    return true
  }

  createObservable(): Observable<unknown> {
    return new Observable((subscriber) => {
      subscriber.next(null)
      subscriber.complete()
    })
  }
}

provideReportUrlProvider(withInstaller)

const container = useTemplateRef<HTMLElement>("container")
const router = useRouter()
const sidebarVm = new InfoSidebarImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, table)
provide(serverConfiguratorKey, serverConfigurator)
const persistentStateManager = new PersistentStateManager(
  persistentId ?? `${dbName}-${table}-dashboard`,
  {
    machine: initialMachine,
    branch: "master",
    project: [],
    measure: [],
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
if (machineConfigurator.selected.value.length === 0) {
  machineConfigurator.selected.value = [initialMachine]
}
const projectFilters: FilterConfigurator[] = projectFilter == null ? [] : [new ProjectLikeFilter(projectFilter)]
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, [...projectFilters, branchConfigurator, timeRangeConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator, branchConfigurator, timeRangeConfigurator], true, "line")

if (preselectAll) {
  let projectsSelected = false
  watch(
    scenarioConfigurator.values,
    (values) => {
      if (projectsSelected || values.length === 0) return
      const current = scenarioConfigurator.selected.value
      if (current == null || (Array.isArray(current) && current.length === 0)) {
        scenarioConfigurator.selected.value = values.filter((it): it is string => typeof it === "string")
      }
      projectsSelected = true
    },
    { immediate: true }
  )

  let measuresSelected = false
  watch(
    measureConfigurator.data,
    (available) => {
      if (measuresSelected || available.length === 0) return
      const current = measureConfigurator.selected.value
      if (current == null || current.length === 0) {
        const mainMetrics = available.filter((it) => MAIN_METRICS.includes(it))
        if (mainMetrics.length > 0) {
          measureConfigurator.setSelected(mainMetrics)
        }
      }
      measuresSelected = true
    },
    { immediate: true }
  )
}

const accidentsConfigurator = new AccidentsConfiguratorForTests(serverConfigurator.serverUrl, scenarioConfigurator.selected, measureConfigurator.selected, timeRangeConfigurator)
provide(accidentsConfiguratorKey, accidentsConfigurator)

const configurators: DataQueryConfigurator[] = [branchConfigurator, machineConfigurator, timeRangeConfigurator, triggeredByConfigurator]

provide(dashboardConfiguratorsKey, configurators)

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}

function onChartClosed(projects: string[]) {
  if (Array.isArray(scenarioConfigurator.selected.value)) {
    scenarioConfigurator.selected.value = scenarioConfigurator.selected.value.filter((item) => !projects.includes(item))
  } else if (scenarioConfigurator.selected.value != null && projects.includes(scenarioConfigurator.selected.value)) {
    scenarioConfigurator.selected.value = null
  }
}

const scenarios = computed(() => {
  if (scenarioConfigurator.selected.value == null) return []
  if (Array.isArray(scenarioConfigurator.selected.value)) {
    return scenarioConfigurator.selected.value
  }
  return [scenarioConfigurator.selected.value]
})
</script>
