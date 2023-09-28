<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect
          :ranges="TimeRangeConfigurator.timeRanges"
          :value="timeRangeConfigurator.value.value"
          :on-change="onChangeRange"
        />
        <BranchSelect
          v-if="releaseConfigurator != null"
          :branch-configurator="branchConfigurator"
          :release-configurator="releaseConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <BranchSelect
          v-else
          :branch-configurator="branchConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <DimensionSelect
          label="Tests"
          :selected-label="testsSelectLabelFormat"
          :dimension="scenarioConfigurator"
        >
          <template #icon>
            <ChartBarIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionSelect>
        <MeasureSelect
          title="Metrics"
          :selected-label="metricsSelectLabelFormat"
          :configurator="measureConfigurator"
        >
          <template #icon>
            <BeakerIcon class="w-4 h-4 text-gray-500" />
          </template>
        </MeasureSelect>
        <MachineSelect :machine-configurator="machineConfigurator" />
      </template>
      <template #end>
        <PlotSettings @update:configurators="updateConfigurators" />
      </template>
    </Toolbar>
    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <template
          v-for="measure in measureConfigurator.selected.value"
          :key="measure"
        >
          <LineChart
            :title="measure"
            :measures="[measure]"
            :configurators="configurators"
            :skip-zero-values="false"
            :value-unit="props.unit"
            :chart-type="props.unit == 'ns' ? 'scatter' : 'line'"
            :legend-formatter="(name) => name"
          />
        </template>
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
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "../../configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getDBType } from "../../shared/dbTypes"
import { accidentsKeys, containerKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import { testsSelectLabelFormat, metricsSelectLabelFormat } from "../../shared/labels"
import { Accident, getAccidentsFromMetaDb } from "../../util/meta"
import DimensionSelect from "../charts/DimensionSelect.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import LineChart from "../charts/PerformanceLineChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import PlotSettings from "../settings/PlotSettings.vue"
import MachineSelect from "./MachineSelect.vue"
import { PersistentStateManager } from "./PersistentStateManager"
import { DataQueryConfigurator } from "./dataQuery"
import { provideReportUrlProvider } from "./lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "./sideBar/InfoSidebar"
import { InfoDataPerformance } from "./sideBar/InfoSidebarPerformance"
import InfoSidebar from "./sideBar/InfoSidebarPerformance.vue"

interface PerformanceTestsProps {
  dbName: string
  table: string
  initialMachine: string
  withInstaller?: boolean
  unit?: "ns" | "ms"
}

const props = withDefaults(defineProps<PerformanceTestsProps>(), {
  withInstaller: true,
  unit: "ms",
})

provideReportUrlProvider(props.withInstaller)

const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarImpl<InfoDataPerformance>(getDBType(props.dbName, props.table))

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(props.dbName, props.table)
provide(serverConfiguratorKey, serverConfigurator)
const persistentStateManager = new PersistentStateManager(
  `${props.dbName}-${props.table}-dashboard`,
  {
    machine: props.initialMachine,
    branch: "master",
    project: [],
    measure: [],
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, [branchConfigurator, timeRangeConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator, branchConfigurator, timeRangeConfigurator], true, "line")

const configurators: DataQueryConfigurator[] = [serverConfigurator, scenarioConfigurator, branchConfigurator, machineConfigurator, timeRangeConfigurator, triggeredByConfigurator]

const releaseConfigurator = props.withInstaller ? new ReleaseNightlyConfigurator(persistentStateManager) : null
if (releaseConfigurator != null) {
  configurators.push(releaseConfigurator)
}

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}

const projectAndMetrics: string[] = []
const projects = scenarioConfigurator.selected.value
const measures = measureConfigurator.selected.value
if (projects != null && measures != null) {
  if (Array.isArray(projects)) {
    projectAndMetrics.push(...projects)
  } else {
    projectAndMetrics.push(projects)
  }

  if (Array.isArray(projects)) {
    projectAndMetrics.push(...projects.map((project) => measures.map((metric) => `${project}/${metric}`)).flat(100))
  } else {
    projectAndMetrics.push(...measures.map((metric) => `${projects}/${metric}`))
  }
}

const warnings: Ref<Map<string, Accident[]> | undefined> = ref()
computedAsync(async () => {
  warnings.value = await getAccidentsFromMetaDb(projectAndMetrics, timeRangeConfigurator.value)
})
provide(accidentsKeys, warnings)
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>
