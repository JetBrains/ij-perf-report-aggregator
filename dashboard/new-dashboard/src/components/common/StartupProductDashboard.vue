<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <TimeRangeSelect :timerange-configurator="timeRangeConfigurator" />
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
        <CopyLink :timerange-configurator="timeRangeConfigurator" />
        <slot name="toolbar" />
      </template>
      <template #end>
        <PlotSettings @update:configurators="updateConfigurators" />
      </template>
    </StickyToolbar>
    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-5 overflow-hidden pt-5"
      >
        <Divider label="Main Metrics" />
        <section>
          <LineChart
            title="FUS Total startup"
            :measures="['metrics.startup/fusTotalDuration', 'metrics.reopenProjectPerformance/fusCodeVisibleInEditorDurationMs']"
            :configurators="configurators"
          />
        </section>
        <Accordion :lazy="true">
          <AccordionTab header="Additional metrics">
            <Divider label="Bootstrap" />
            <section class="grid grid-cols-2 gap-x-6">
              <LineChart
                :measures="['appInit_d', 'app initialization.end']"
                title="App Initialization"
                :configurators="configurators"
                tooltip-trigger="axis"
              />
              <LineChart
                :measures="['bootstrap_d']"
                title="Bootstrap"
                :configurators="configurators"
                :with-measure-name="true"
              />
            </section>

            <section class="grid grid-cols-2 gap-x-6">
              <LineChart
                :measures="['classLoadingPreparedCount', 'classLoadingLoadedCount']"
                title="Class Loading (Count)"
                :configurators="configurators"
                tooltip-trigger="axis"
              />
              <LineChart
                :configurators="configurators"
                :measures="['editorRestoring']"
                title="Editor restoring"
                :with-measure-name="true"
              />
            </section>

            <span v-if="highlightingPasses">
              <Divider label="Highlighting Passes" />
              <span v-if="showAllPasses">
                <LineChart
                  title="Highlighting Passes"
                  :measures="highlightingPasses"
                  :configurators="configurators"
                />
              </span>
              <LineChart
                title="Code Analysis"
                :measures="['metrics.codeAnalysisDaemon/fusExecutionTime', 'metrics.runDaemon/executionTime']"
                :configurators="configurators"
                tooltip-trigger="axis"
              />
            </span>
            <slot :configurators="configurators"></slot>
            <Divider label="Notifications" />
            <LineChart
              title="Notifications"
              :measures="['metrics.notifications/number']"
              :skip-zero-values="false"
              :configurators="configurators"
              :with-measure-name="true"
            />

            <Divider label="Exit" />
            <LineChart
              title="Exit Metrics"
              :measures="['metrics.exitMetrics/application.exit', 'metrics.exitMetrics/saveSettingsOnExit', 'metrics.exitMetrics/disposeProjects']"
              :configurators="configurators"
              tooltip-trigger="axis"
            />
          </AccordionTab>
        </Accordion>
      </div>

      <InfoSidebar />
    </main>
  </div>
</template>
<script setup lang="ts">
import { computed, provide, ref } from "vue"
import { useRouter } from "vue-router"
import { AccidentsConfiguratorForStartup } from "../../configurators/AccidentsConfigurator"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { configuratorListKey } from "../../shared/injectionKeys"
import { accidentsConfiguratorKey, containerKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import DimensionSelect from "../charts/DimensionSelect.vue"
import LineChart from "../charts/LineChart.vue"
import CopyLink from "../settings/CopyLink.vue"
import PlotSettings from "../settings/PlotSettings.vue"
import { createProjectConfigurator, getProjectName } from "../startup/projectNameMapping"
import { fetchHighlightingPasses } from "../startup/utils"
import BranchSelect from "./BranchSelect.vue"
import Divider from "./Divider.vue"
import MachineSelect from "./MachineSelect.vue"
import { PersistentStateManager } from "./PersistentStateManager"
import StickyToolbar from "./StickyToolbar.vue"
import TimeRangeSelect from "./TimeRangeSelect.vue"
import { DataQueryConfigurator } from "./dataQuery"
import { provideReportUrlProvider } from "./lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "./sideBar/InfoSidebar"

import InfoSidebar from "./sideBar/InfoSidebar.vue"

interface StartupProductDashboard {
  product: string
  defaultProject: string
  persistentId?: string | null
  withInstaller?: boolean
}

const props = withDefaults(defineProps<StartupProductDashboard>(), {
  persistentId: null,
  withInstaller: false,
})
provideReportUrlProvider(props.withInstaller)

const container = ref<HTMLElement>()

const dbName = props.withInstaller ? "ij" : "ijDev"
const dbTable = "report"

const sidebarVm = new InfoSidebarImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, dbTable)
provide(serverConfiguratorKey, serverConfigurator)
const persistentStateManager = new PersistentStateManager(
  props.persistentId ?? `${props.product}-startup-dashboard`,
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

const metrics = ref([
  "appInit_d",
  "app initialization.end",
  "classLoadingPreparedCount",
  "classLoadingLoadedCount",
  "editorRestoring",
  "startup/fusTotalDuration",
  "codeAnalysisDaemon/fusExecutionTime",
  "runDaemon/executionTime",
  "notifications/number",
  "exitMetrics/application.exit",
  "exitMetrics/saveSettingsOnExit",
  "exitMetrics/disposeProjects",
])

const accidentsConfigurator = new AccidentsConfiguratorForStartup(serverConfigurator.serverUrl, ref(props.product), projectConfigurator.selected, metrics, timeRangeConfigurator)
provide(accidentsConfiguratorKey, accidentsConfigurator)

const configurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  productConfigurator,
  projectConfigurator,
  branchConfigurator,
  triggeredByConfigurator,
  accidentsConfigurator,
] as DataQueryConfigurator[]

provide(configuratorListKey, configurators)

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}

const highlightingPasses = fetchHighlightingPasses()

const showAllPasses = computed(() => {
  return projectConfigurator.selected.value == null || projectConfigurator.selected.value.length == 1 || typeof projectConfigurator.selected.value == "string"
})
</script>
