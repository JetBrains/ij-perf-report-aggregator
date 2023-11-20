<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <CopyLink :timerange-configurator="timeRangeConfigurator" />
        <TimeRangeSelect
          :ranges="timeRangeConfigurator.timeRanges"
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
        <div v-if="testMetricSwitcher == TestMetricSwitcher.Tests">
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
        </div>
        <div v-else-if="testMetricSwitcher == TestMetricSwitcher.Metrics">
          <MeasureSelect
            title="Metrics"
            :selected-label="metricsSelectLabelFormat"
            :configurator="measureConfigurator"
          >
            <template #icon>
              <BeakerIcon class="w-4 h-4 text-gray-500" />
            </template>
          </MeasureSelect>
          <DimensionSelect
            label="Tests"
            :selected-label="testsSelectLabelFormat"
            :dimension="scenarioConfigurator"
          >
            <template #icon>
              <ChartBarIcon class="w-4 h-4 text-gray-500" />
            </template>
          </DimensionSelect>
        </div>
        <MachineSelect :machine-configurator="machineConfigurator" />
        <SelectButton
          v-model="testMetricSwitcher"
          :options="testMetricSwitcherOptions"
        />
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
        <div v-if="testMetricSwitcher == TestMetricSwitcher.Tests">
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
        <div v-else-if="testMetricSwitcher == TestMetricSwitcher.Metrics">
          <div
            v-if="measureConfigurator.selected.value != null"
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
              />
            </template>
          </div>
        </div>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { provide, Ref, ref, watch, WatchStopHandle } from "vue"
import { useRouter } from "vue-router"
import { AccidentsConfiguratorForTests } from "../../configurators/AccidentsConfigurator"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { DimensionConfigurator, dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { NoOpConfigurator } from "../../configurators/NoOpConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "../../configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getDBType } from "../../shared/dbTypes"
import { accidentsConfiguratorKey, containerKey, dashboardConfiguratorsKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import { testsSelectLabelFormat, metricsSelectLabelFormat } from "../../shared/labels"
import DimensionSelect from "../charts/DimensionSelect.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import LineChart from "../charts/PerformanceLineChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import CopyLink from "../settings/CopyLink.vue"
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

enum TestMetricSwitcher {
  Tests = "Tests",
  Metrics = "Metrics",
}

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
    type: TestMetricSwitcher.Tests,
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
if (machineConfigurator.selected.value.length === 0) {
  machineConfigurator.selected.value = [props.initialMachine]
}
let scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, [branchConfigurator, timeRangeConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
let measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator, branchConfigurator, timeRangeConfigurator], true, "line")

const accidentsConfigurator = new AccidentsConfiguratorForTests(scenarioConfigurator.selected, measureConfigurator.selected, timeRangeConfigurator)
provide(accidentsConfiguratorKey, accidentsConfigurator)

const configurators: DataQueryConfigurator[] = [serverConfigurator, branchConfigurator, machineConfigurator, timeRangeConfigurator, triggeredByConfigurator, accidentsConfigurator]

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
provide(dashboardConfiguratorsKey, configurators)

const testMetricSwitcher: Ref<TestMetricSwitcher | null> = ref(TestMetricSwitcher.Tests)
const testMetricSwitcherOptions = [TestMetricSwitcher.Tests, TestMetricSwitcher.Metrics]
persistentStateManager.add("type", testMetricSwitcher)
let previousValue: TestMetricSwitcher | null = null
let watchStopHandle: WatchStopHandle | null = null
watch(
  testMetricSwitcher,
  (value) => {
    if (value == null) {
      testMetricSwitcher.value = previousValue
      return
    }
    if (value == previousValue) return
    previousValue = value
    switch (value) {
      case TestMetricSwitcher.Tests: {
        scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, [branchConfigurator, timeRangeConfigurator])
        measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator, branchConfigurator, timeRangeConfigurator], true, "line")
        const index = configurators.findIndex((configurator) => configurator instanceof NoOpConfigurator)
        if (index == -1) {
          configurators.push(scenarioConfigurator)
        } else {
          configurators[index] = scenarioConfigurator
        }
        break
      }
      case TestMetricSwitcher.Metrics: {
        measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator], true, "line")
        scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, [branchConfigurator, timeRangeConfigurator, measureConfigurator])
        if (watchStopHandle != null) watchStopHandle()
        watchStopHandle = watch(scenarioConfigurator.selected, (value) => {
          if (value?.length != 0) {
            scenarios = toArray(value)
          }
        })
        configurators[configurators.findIndex((configurator) => configurator instanceof DimensionConfigurator && configurator.name == "project")] = new NoOpConfigurator()
        break
      }
    }
  },
  { immediate: true }
)

function toArray(value: string | string[] | null): string[] {
  if (value == null) return []
  if (Array.isArray(value)) {
    return value
  }
  return [value]
}

let scenarios = toArray(scenarioConfigurator.selected.value)
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>
