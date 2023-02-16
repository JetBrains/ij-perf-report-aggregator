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
        <MeasureSelect
          :configurator="measureConfigurator"
          title="Metrics"
          :selected-label="metricsSelectLabelFormat"
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
        <LineChartCard />
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
import MeasureSelect from "shared/src/components/MeasureSelect.vue"
import { AggregationOperatorConfigurator } from "shared/src/configurators/AggregationOperatorConfigurator"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { MeasureConfigurator } from "shared/src/configurators/MeasureConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { aggregationOperatorConfiguratorKey, chartStyleKey, chartToolTipKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { metricsSelectLabelFormat } from "../../shared/labels"
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
const container = ref<HTMLElement>()

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager("ij-explore")

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistentStateManager,
  [timeRangeConfigurator],
)

const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager)

const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager)
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager)
const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  measureConfigurator,
  productConfigurator,
  projectConfigurator
]

provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))
initDataComponent(configurators)
function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}
</script>