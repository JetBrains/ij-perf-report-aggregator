<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect
          :ranges="timeRangeConfigurator.timeRanges"
          :value="timeRangeConfigurator.value.value"
          :on-change="onChangeRange"
        />
        <BranchSelect :branch-configurator="branchConfigurator" />
        <DimensionSelect
          label="Product"
          :value-to-label="(it: string) => productCodeToName.get(it) ?? it"
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
        <LineChartCard />
      </div>
      <InfoSidebarStartup />
    </main>
  </div>
  <ChartTooltip ref="tooltip" />
</template>
<script setup lang="ts">
import { provide, Ref, ref } from "vue"
import { AggregationOperatorConfigurator } from "../../configurators/AggregationOperatorConfigurator"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getDBType } from "../../shared/dbTypes"
import { aggregationOperatorConfiguratorKey, chartStyleKey, chartToolTipKey, configuratorListKey } from "../../shared/injectionKeys"
import { containerKey, sidebarStartupKey } from "../../shared/keys"
import { metricsSelectLabelFormat } from "../../shared/labels"
import ChartTooltip from "../charts/ChartTooltip.vue"
import DimensionSelect from "../charts/DimensionSelect.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import LineChartCard from "../charts/StartupLineChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import MachineSelect from "../common/MachineSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { chartDefaultStyle } from "../common/chart"
import { DataQueryConfigurator } from "../common/dataQuery"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "../common/sideBar/InfoSidebar"
import { InfoDataFromStartup } from "../common/sideBar/InfoSidebarStartup"
import InfoSidebarStartup from "../common/sideBar/InfoSidebarStartup.vue"
import PlotSettings from "../settings/PlotSettings.vue"
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
provide(chartToolTipKey, tooltip as Ref<typeof ChartTooltip>)

const dbName = "ij"
const dbTable = "report"
const container = ref<HTMLElement>()
provide(containerKey, container)

const sidebarVm = new InfoSidebarImpl<InfoDataFromStartup>(getDBType(dbName, dbTable))
provide(sidebarStartupKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager("ij-explore")

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])

const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])

const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager, false, [timeRangeConfigurator, branchConfigurator])
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  measureConfigurator,
  productConfigurator,
  projectConfigurator,
  branchConfigurator,
] as DataQueryConfigurator[]

provide(aggregationOperatorConfiguratorKey, new AggregationOperatorConfigurator(persistentStateManager))
provide(configuratorListKey, configurators)

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}
</script>
