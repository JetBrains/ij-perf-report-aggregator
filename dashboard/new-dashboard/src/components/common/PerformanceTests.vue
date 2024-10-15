<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <CopyLink :timerange-configurator="timeRangeConfigurator" />
        <TimeRangeSelect :timerange-configurator="timeRangeConfigurator" />
        <BranchSelect
          v-if="releaseNightlyConfigurator != null && branchConfigurator != null"
          :branch-configurator="branchConfigurator"
          :release-configurator="releaseNightlyConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <BranchSelect
          v-else-if="branchConfigurator != null"
          :branch-configurator="branchConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <span
          v-if="testMetricSwitcher == TestMetricSwitcher.Tests"
          class="flex flex-row justify-between items-center"
        >
          <DimensionSelect
            label="Tests"
            :selected-label="testsSelectLabelFormat"
            :dimension="scenarioConfigurator"
            class="mr-2"
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
        </span>
        <span
          v-else-if="testMetricSwitcher == TestMetricSwitcher.Metrics"
          class="flex flex-row justify-between items-center"
        >
          <MeasureSelect
            title="Metrics"
            :selected-label="metricsSelectLabelFormat"
            :configurator="measureConfigurator"
            class="mr-2"
          >
            <template #icon>
              <BeakerIcon class="w-4 h-4" />
            </template>
          </MeasureSelect>
          <DimensionSelect
            label="Tests"
            :selected-label="testsSelectLabelFormat"
            :dimension="scenarioConfigurator"
          >
            <template #icon>
              <ChartBarIcon class="w-4 h-4" />
            </template>
          </DimensionSelect>
        </span>
        <DimensionSelect
          v-if="testModeConfigurator != null && testModeConfigurator.values.value.length > 1"
          label="Mode"
          :dimension="testModeConfigurator"
          :selected-label="modeSelectLabelFormat"
        >
          <template #icon>
            <AdjustmentsVerticalIcon class="w-4 h-4" />
          </template>
        </DimensionSelect>
        <MachineSelect
          v-if="machineConfigurator != null"
          :machine-configurator="machineConfigurator"
        />
        <SelectButton
          v-model="testMetricSwitcher"
          :allow-empty="false"
          :options="testMetricSwitcherOptions"
          class="flex flex-1"
        />
      </template>
      <template #end>
        <PlotSettings @update:configurators="updateConfigurators" />
      </template>
    </StickyToolbar>
    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <span v-if="testMetricSwitcher == TestMetricSwitcher.Tests && configuratorsUpdated">
          <template
            v-for="measure in measureConfigurator.selected.value"
            :key="measure"
          >
            <LineChart
              :title="measure"
              :measures="[measure]"
              :configurators="[...configurators, scenarioConfigurator]"
              :skip-zero-values="false"
              :value-unit="unit"
              :chart-type="unit == 'ns' ? 'scatter' : 'line'"
              :legend-formatter="(name: string) => name"
              :can-be-closed="true"
              @chart-closed="onTestChartClosed"
            />
          </template>
        </span>
        <span v-else-if="testMetricSwitcher == TestMetricSwitcher.Metrics && measureConfigurator.selected.value != null && measureConfigurator.selected.value?.length > 0">
          <template
            v-for="scenario in scenarios"
            :key="scenario"
          >
            <GroupProjectsChart
              :measure="measureConfigurator.selected.value"
              :projects="[scenario]"
              :label="scenario"
              :can-be-closed="true"
              @chart-closed="onMeasureChartClosed"
            />
          </template>
        </span>
      </div>
      <InfoSidebar :timerange-configurator="timeRangeConfigurator" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { provide, Ref, ref, watch, WatchStopHandle, useTemplateRef } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { nightly, ReleaseNightlyConfigurator, ReleaseType } from "../../configurators/ReleaseNightlyConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { accidentsConfiguratorKey, containerKey, dashboardConfiguratorsKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import { testsSelectLabelFormat, metricsSelectLabelFormat, modeSelectLabelFormat } from "../../shared/labels"
import DimensionSelect from "../charts/DimensionSelect.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import LineChart from "../charts/LineChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import CopyLink from "../settings/CopyLink.vue"
import PlotSettings from "../settings/PlotSettings.vue"
import MachineSelect from "./MachineSelect.vue"
import { PersistentStateManager } from "./PersistentStateManager"
import StickyToolbar from "./StickyToolbar.vue"
import { DataQueryConfigurator } from "./dataQuery"
import { provideReportUrlProvider } from "./lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "./sideBar/InfoSidebar"
import InfoSidebar from "./sideBar/InfoSidebar.vue"
import { AccidentsConfiguratorForTests } from "../../configurators/accidents/AccidentsConfiguratorForTests"
import { createTestModeConfigurator, defaultModeName } from "../../configurators/TestModeConfigurator"
import { dbTypeStore } from "../../shared/dbTypes"

interface PerformanceTestsProps {
  dbName: string
  table: string
  initialMachine: string | null
  withInstaller?: boolean
  unit?: "ns" | "ms"
  releaseConfigurator?: ReleaseType
  branch?: string | null
}

const { dbName, table, initialMachine, withInstaller = true, unit = "ms", releaseConfigurator = nightly, branch = "master" } = defineProps<PerformanceTestsProps>()

enum TestMetricSwitcher {
  Tests = "Tests",
  Metrics = "Metrics",
}

const container = useTemplateRef<HTMLElement>("container")
const router = useRouter()
const sidebarVm = new InfoSidebarImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, table)

provideReportUrlProvider(withInstaller)
provide(serverConfiguratorKey, serverConfigurator)
const persistentStateManager = new PersistentStateManager(
  `${dbName}-${table}-dashboard`,
  {
    machine: initialMachine ?? "",
    branch: branch ?? "",
    project: [],
    measure: [],
    type: TestMetricSwitcher.Tests,
    releaseConfigurator,
    mode: defaultModeName,
  },
  router
)

const filters = []
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
filters.push(timeRangeConfigurator)

const branchConfigurator = branch == null ? null : createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
if (branchConfigurator != null) filters.push(branchConfigurator)

const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, filters)

const measureScenarioFilters = [triggeredByConfigurator, timeRangeConfigurator]
if (branchConfigurator != null) {
  measureScenarioFilters.push(branchConfigurator)
}
let scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, measureScenarioFilters)
filters.push(scenarioConfigurator)

const testModeConfigurator = dbTypeStore().isModeSupported() ? createTestModeConfigurator(serverConfigurator, persistentStateManager, [...filters]) : null
if (testModeConfigurator != null) {
  filters.push(testModeConfigurator)
}

let measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, measureScenarioFilters, true, "line")
const machineConfigurator = initialMachine == null ? null : new MachineConfigurator(serverConfigurator, persistentStateManager, filters)
if (initialMachine != null && machineConfigurator != null && machineConfigurator.selected.value.length === 0) {
  machineConfigurator.selected.value = [initialMachine]
}

const accidentsConfigurator = new AccidentsConfiguratorForTests(serverConfigurator.serverUrl, scenarioConfigurator.selected, measureConfigurator.selected, timeRangeConfigurator)
provide(accidentsConfiguratorKey, accidentsConfigurator)

const configurators: DataQueryConfigurator[] = [serverConfigurator, timeRangeConfigurator, triggeredByConfigurator, accidentsConfigurator]
if (branchConfigurator != null) {
  configurators.push(branchConfigurator)
}
if (machineConfigurator != null) {
  configurators.push(machineConfigurator)
}
if (testModeConfigurator != null) {
  configurators.push(testModeConfigurator)
}

const releaseNightlyConfigurator = withInstaller ? new ReleaseNightlyConfigurator(persistentStateManager) : null
if (releaseNightlyConfigurator != null) {
  configurators.push(releaseNightlyConfigurator)
}

const configuratorsUpdated = ref(false)
const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configuratorsUpdated.value = true
  configurators.push(configurator)
}
provide(dashboardConfiguratorsKey, configurators)

function onTestChartClosed(metric: Ref<string[]>) {
  measureConfigurator.setSelected(measureConfigurator.selected.value?.filter((item) => !metric.value.includes(item)) as string[])
}

function onMeasureChartClosed(projects: string[]) {
  if (Array.isArray(scenarioConfigurator.selected.value)) {
    scenarioConfigurator.selected.value = scenarioConfigurator.selected.value.filter((item) => !projects.includes(item))
  } else if (scenarioConfigurator.selected.value != null && projects.includes(scenarioConfigurator.selected.value)) {
    scenarioConfigurator.selected.value = null
  }
}

const testMetricSwitcher: Ref<TestMetricSwitcher | null> = ref(TestMetricSwitcher.Tests)
const testMetricSwitcherOptions = [TestMetricSwitcher.Tests, TestMetricSwitcher.Metrics]
persistentStateManager.add("type", testMetricSwitcher)
let watchStopHandle: WatchStopHandle | null = null
watch(
  testMetricSwitcher,
  (value) => {
    switch (value) {
      case TestMetricSwitcher.Tests: {
        scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, measureScenarioFilters)
        measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator, ...measureScenarioFilters], true, "line")
        break
      }
      case TestMetricSwitcher.Metrics: {
        measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, measureScenarioFilters, true, "line")
        scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, [...measureScenarioFilters, measureConfigurator])
        if (watchStopHandle != null) watchStopHandle()
        watchStopHandle = watch(scenarioConfigurator.selected, (value) => {
          scenarios = toArray(value)
        })
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
