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
        <DimensionSelect
          label="Product"
          :value-to-label="it => productCodeToName.get(it) ?? it"
          :dimension="productConfigurator"
        />
        <DimensionSelect
          label="Project"
          :value-to-label="getProjectName"
          :dimension="projectConfigurator"
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
        <Divider title="Bootstrap" />
        <section class="flex gap-x-6">
          <div class="flex-1">
            <LineChartCard
              :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
            />
          </div>
          <div class="flex-1">
            <LineChartCard
              :measures='["pluginDescriptorInitV18_d", "RunManager initialization"]'
            />
          </div>
        </section>

        <Divider title="Class and Resource Loading" />
        <LineChartCard
          :measures='["classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime"]'
        />
        <section class="flex gap-x-6">
          <div class="flex-1">
            <LineChartCard
              :measures='["classLoadingCount", "resourceLoadingCount", "classLoadingPreparedCount", "classLoadingLoadedCount"]'
            />
          </div>
          <div class="flex-1">
            <LineChartCard
              :measures='["metrics.classLoadingMetrics/inlineCount", "metrics.classLoadingMetrics/companionCount",
                          "metrics.classLoadingMetrics/lambdaCount", "metrics.classLoadingMetrics/methodHandleCount"]'
            />
          </div>
        </section>
        <Divider title="Services" />
        <section class="flex gap-x-6">
          <div class="flex-1">
            <LineChartCard
              :measures='["appComponentCreation_d", "projectComponentCreation_d"]'
            />
          </div>
          <div class="flex-1">
            <LineChartCard
              :skip-zero-values="false"
              :measures='["serviceSyncPreloading_d", "serviceAsyncPreloading_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d"]'
            />
          </div>
        </section>

        <Divider title="Post-opening" />
        <section class="flex gap-x-6">
          <div class="flex-1">
            <LineChartCard
              :measures='["projectDumbAware", "editorRestoring", "editorRestoringTillPaint"]'
            />
          </div>
          <div class="flex-1">
            <LineChartCard
              :measures='["splash_i", "startUpCompleted"]'
            />
          </div>
        </section>
        <Divider title="Memory" />
        <LineChartCard
          :skip-zero-values="false"
          :measures='["metrics.memory/usedMb", "metrics.memory/metaspaceMb","metrics.memory/maxMb"]'
        />
        <div class="relative flex py-5 items-center">
          <div class="flex-grow border-t border-gray-400" />
          <span class="flex-shrink mx-4 text-gray-400 text-lg">GC</span>
          <div class="flex-grow border-t border-gray-400" />
        </div>
        <LineChartCard
          :skip-zero-values="false"
          :measures='["metrics.gc/totalHeapUsedMax", "metrics.gc/freedMemoryByGC"]'
        />
        <section class="flex gap-x-6">
          <div class="flex-1">
            <LineChartCard
              :skip-zero-values="false"
              label="Number of pauses"
              :measures='["metrics.gc/gcPauseCount"]'
            />
          </div>
          <div class="flex-1">
            <LineChartCard
              :skip-zero-values="false"
              :measures='["metrics.gc/fullGCPause","metrics.gc/gcPause"]'
            />
          </div>
        </section>
      </div>
    </main>
  </div>
  <ChartTooltip
    ref="tooltip"
  />
</template>
<script setup lang="ts">
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import ChartTooltip from "shared/src/components/ChartTooltip.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import LineChartCard from "shared/src/components/LineChartCard.vue"
import { AggregationOperatorConfigurator } from "shared/src/configurators/AggregationOperatorConfigurator"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey, chartToolTipKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import Divider from "../common/Divider.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { createProjectConfigurator, getProjectName } from "./projectNameMapping"

const productCodeToName = new Map([
  ["DB", "DataGrip"],
  ["IU", "IntelliJ IDEA"],
  ["PS", "PhpStorm"],
  ["WS", "WebStorm"],
  ["GO", "GoLand"],
  ["PY", "PyCharm Professional"],
  ["RM", "RubyMine"],
])

provideReportUrlProvider()
provide(chartStyleKey, {
  ...chartDefaultStyle,
  // a lot of bars, as result, height of bar is not enough to make label readable
  barSeriesLabelPosition: "right",
})
const tooltip = ref<typeof ChartTooltip>()
provide(chartToolTipKey, tooltip)

const dbName = "ij"
const dbTable = "report"
const initialMachine = "linux-blade-hetzner"
const container = ref<HTMLElement>()
const router = useRouter()


const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager("ij-dashboard", {
  product: "IU",
  project: "simple for IJ",
  machine: "macMini M1, 16GB",
}, useRouter())

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator],
)
const measureConfigurator = new MeasureConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator],
  true,
  "line",
)

const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager)
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager)
const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  productConfigurator,
  projectConfigurator
]

provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))
initDataComponent(configurators)
function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}
</script>