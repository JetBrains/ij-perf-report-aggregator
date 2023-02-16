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
        <div class="relative flex py-5 items-center">
          <div class="flex-grow border-t border-gray-400" />
          <span class="flex-shrink mx-4 text-gray-400 text-lg">Remote Mode</span>
          <div class="flex-grow border-t border-gray-400" />
        </div>
        <!-- :skip-zero-values="false"
:configurators="dashboardConfigurators" because computed measures cannot be filtered -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
          <LineChart
            title="editor appeared"
            :measures='["editor appeared.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="time to edit"
            :measures='["time to edit.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="terminal ready"
            :measures='["terminal ready.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="file tree rendered"
            :measures='["file tree rendered.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="highlighting done"
            :measures='["highlighting done.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="window appeared"
            :measures='["window appeared.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
        </div>

        <div class="relative flex py-5 items-center">
          <div class="flex-grow border-t border-gray-400" />
          <span class="flex-shrink mx-4 text-gray-400 text-lg">ShortCircuit</span>
          <div class="flex-grow border-t border-gray-400" />
        </div>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
          <LineChart
            title="editor appeared"
            :measures='["shortCircuit.editor appeared.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="time to edit"
            :measures='["shortCircuit.time to edit.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="terminal ready"
            :measures='["shortCircuit.terminal ready.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="file tree rendered"
            :measures='["shortCircuit.file tree rendered.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="highlighting done"
            :measures='["shortCircuit.highlighting done.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
          <LineChart
            title="window appeared"
            :measures='["shortCircuit.window appeared.end"]'
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
          />
        </div>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>


<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import LineChart from "../charts/LineChart.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider()

const dbName = "fleet"
const dbTable = "report"
const initialMachine = "Linux Munich i7-3770, 32 Gb"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("fleetStartup_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator],
)

const dashboardConfigurators = [
  serverConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
]

function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}
</script>
