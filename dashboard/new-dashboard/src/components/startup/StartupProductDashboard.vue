<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="timeRangeConfigurator.timeRanges"
        :value="timeRangeConfigurator.value.value"
        :on-change="onChangeRange"
      />
      <BranchSelect
        :branch-configurator="branchConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
      />
      <DimensionSelect
        label="Project"
        :value-to-label="getProjectName"
        :dimension="projectConfigurator"
      />
      <MachineSelect :machine-configurator="machineConfigurator" />
      <slot name="toolbar" />
    </template>
    <template #end>
      <PlotSettings @update:configurators="updateConfigurators" />
    </template>
  </Toolbar>
  <main class="flex">
    <div
      ref="container"
      class="flex flex-1 flex-col gap-5 overflow-hidden pt-5"
    >
      <Divider label="Bootstrap" />
      <section class="grid grid-cols-2 gap-x-6">
        <StartupLineChart
          :measures="['appInit_d', 'app initialization.end']"
          title="App Initialization"
        />
        <StartupLineChart
          :measures="['bootstrap_d']"
          title="Bootstrap"
        />
      </section>

      <section class="grid grid-cols-2 gap-x-6">
        <StartupLineChart
          :measures="['classLoadingPreparedCount', 'classLoadingLoadedCount']"
          title="Class Loading (Count)"
        />
        <StartupLineChart
          :measures="['editorRestoring']"
          title="Editor restoring"
        />
      </section>

      <span v-if="highlightingPasses">
        <Divider label="Highlighting Passes" />
        <StartupLineChart :measures="highlightingPasses" />
        <StartupLineChart :measures="['metrics.codeAnalysisDaemon/fusExecutionTime', 'metrics.runDaemon/executionTime']" />
      </span>

      <Divider label="Notifications" />
      <StartupLineChart
        :measures="['metrics.notifications/number']"
        :skip-zero-values="false"
      />

      <Divider label="Exit" />
      <StartupLineChart :measures="['metrics.exitMetrics/application.exit', 'metrics.exitMetrics/saveSettingsOnExit', 'metrics.exitMetrics/disposeProjects']" />
    </div>
    <InfoSidebarStartup />
  </main>
</template>
<script setup lang="ts">
import { provide, Ref, ref } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { getDBType } from "../../shared/dbTypes"
import { chartStyleKey, chartToolTipKey, configuratorListKey } from "../../shared/injectionKeys"
import { containerKey, sidebarStartupKey } from "../../shared/keys"
import ChartTooltip from "../charts/ChartTooltip.vue"
import DimensionSelect from "../charts/DimensionSelect.vue"
import StartupLineChart from "../charts/StartupLineChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import Divider from "../common/Divider.vue"
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
import { fetchHighlightingPasses } from "./utils"

interface StartupProductDashboard {
  product: string
  defaultProject: string
}

const props = defineProps<StartupProductDashboard>()

const container = ref<HTMLElement>()

const tooltip = ref<typeof ChartTooltip>()
provide(chartToolTipKey, tooltip as Ref<typeof ChartTooltip>)
provide(containerKey, container)

provideReportUrlProvider()
provide(chartStyleKey, {
  ...chartDefaultStyle,
})

const dbName = "ij"
const dbTable = "report"

const sidebarVm = new InfoSidebarImpl<InfoDataFromStartup>(getDBType(dbName, dbTable))
provide(sidebarStartupKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager(
  `${props.product}-startup-dashboard`,
  {
    project: props.defaultProject,
    machine: "Windows Munich i7-3770, 32Gb",
    branch: "master",
  },
  useRouter()
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const productConfigurator = dimensionConfigurator("product", serverConfigurator, persistentStateManager, false, [timeRangeConfigurator, branchConfigurator])
productConfigurator.selected.value = props.product
const projectConfigurator = createProjectConfigurator(productConfigurator, serverConfigurator, persistentStateManager, [
  productConfigurator,
  timeRangeConfigurator,
  branchConfigurator,
])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  productConfigurator,
  projectConfigurator,
  branchConfigurator,
  triggeredByConfigurator,
] as DataQueryConfigurator[]

provide(configuratorListKey, configurators)

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

const highlightingPasses = fetchHighlightingPasses()
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>
