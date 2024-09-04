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
        <template
          v-for="measure in measureConfigurator.selected.value"
          :key="measure"
        >
          <LineChart
            :title="measure"
            :measures="[measure]"
            :configurators="configurators"
            :skip-zero-values="false"
          />
        </template>
      </div>
      <InfoSidebar :timerange-configurator="timeRangeConfigurator" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { provide, useTemplateRef } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../configurators/BranchConfigurator"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { MeasureConfigurator } from "../../configurators/MeasureConfigurator"
import { privateBuildConfigurator } from "../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import { testsSelectLabelFormat, metricsSelectLabelFormat } from "../../shared/labels"
import DimensionSelect from "../charts/DimensionSelect.vue"
import LineChart from "../charts/LineChart.vue"
import MeasureSelect from "../charts/MeasureSelect.vue"
import BranchSelect from "../common/BranchSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import StickyToolbar from "../common/StickyToolbar.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { DataQueryConfigurator } from "../common/dataQuery"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "../common/sideBar/InfoSidebar"
import InfoSidebar from "../common/sideBar/InfoSidebar.vue"
import PlotSettings from "../settings/PlotSettings.vue"

provideReportUrlProvider(false, true)

const dbName = "jbr"
const dbTable = "report"
const initialMachine = "Linux Munich i7-3770, 32 Gb"
const container = useTemplateRef<HTMLElement>("container")
const router = useRouter()
const sidebarVm = new InfoSidebarImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerWithCompressConfigurator(dbName, dbTable)
const persistentStateManager = new PersistentStateManager(
  `${dbName}-${dbTable}-dashboard`,
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
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, persistentStateManager, true, [branchConfigurator, timeRangeConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])
const measureConfigurator = new MeasureConfigurator(serverConfigurator, persistentStateManager, [scenarioConfigurator, branchConfigurator, timeRangeConfigurator], true, "line")

const configurators = [serverConfigurator, scenarioConfigurator, branchConfigurator, timeRangeConfigurator, triggeredByConfigurator] as DataQueryConfigurator[]

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  configurators.push(configurator)
}
</script>
