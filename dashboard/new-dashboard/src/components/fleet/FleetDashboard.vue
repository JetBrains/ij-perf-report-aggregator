<template>
  <div class="flex flex-col gap-5">
    <StickyToolbar>
      <template #start>
        <TimeRangeSelect :timerange-configurator="timeRangeConfigurator" />
        <MachineSelect :machine-configurator="machineConfigurator" />
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
        <Divider title="Remote Mode" />
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
          <LineChart
            v-for="metric in metricNames"
            :key="metric"
            :title="metric"
            :measures="[metric + '.end']"
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
            :with-measure-name="true"
          />
        </div>

        <Divider title="ShortCircuit" />
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
          <LineChart
            v-for="metric in metricNames"
            :key="'shortCircuit.' + metric"
            :title="'shortCircuit.' + metric"
            :measures="['shortCircuit.' + metric + '.end']"
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
            :with-measure-name="true"
          />
        </div>

        <Divider title="ShortCircuit CDS" />
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 pt-4">
          <LineChart
            v-for="metric in metricNames"
            :key="'cds.shortCircuit.' + metric"
            :title="'cds.shortCircuit.' + metric"
            :measures="['cds.shortCircuit.' + metric + '.end']"
            :skip-zero-values="false"
            :configurators="dashboardConfigurators"
            :with-measure-name="true"
          />
        </div>
      </div>
      <InfoSidebar :timerange-configurator="timeRangeConfigurator" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { AccidentsConfiguratorForTests } from "../../configurators/AccidentsConfigurator"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"
import { ServerWithCompressConfigurator } from "../../configurators/ServerWithCompressConfigurator"
import { TimeRangeConfigurator } from "../../configurators/TimeRangeConfigurator"
import { accidentsConfiguratorKey, containerKey, serverConfiguratorKey, sidebarVmKey } from "../../shared/keys"
import LineChart from "../charts/LineChart.vue"
import Divider from "../common/Divider.vue"
import MachineSelect from "../common/MachineSelect.vue"
import { PersistentStateManager } from "../common/PersistentStateManager"
import StickyToolbar from "../common/StickyToolbar.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"
import { DataQueryConfigurator } from "../common/dataQuery"
import { provideReportUrlProvider } from "../common/lineChartTooltipLinkProvider"
import { InfoSidebarImpl } from "../common/sideBar/InfoSidebar"
import InfoSidebar from "../common/sideBar/InfoSidebar.vue"
import PlotSettings from "../settings/PlotSettings.vue"

provideReportUrlProvider()

const dbName = "fleet"
const dbTable = "report"
const initialMachine = "Linux Munich i7-3770, 32 Gb"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const metricNames = ["editor appeared", "time to edit", "terminal ready", "file tree rendered", "highlighting done", "window appeared"]

const serverConfigurator = new ServerWithCompressConfigurator(dbName, dbTable)
provide(serverConfiguratorKey, serverConfigurator)
const persistenceForDashboard = new PersistentStateManager(
  "fleetStartup_dashboard",
  {
    machine: initialMachine,
    project: [],
    branch: "master",
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const machineConfigurator = new MachineConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const accidentsConfigurator = new AccidentsConfiguratorForTests(serverConfigurator.serverUrl, ref("fleet"), ref(null), timeRangeConfigurator)
provide(accidentsConfiguratorKey, accidentsConfigurator)

const dashboardConfigurators = [serverConfigurator, machineConfigurator, timeRangeConfigurator, accidentsConfigurator] as DataQueryConfigurator[]

const updateConfigurators = (configurator: DataQueryConfigurator) => {
  dashboardConfigurators.push(configurator)
}
</script>
