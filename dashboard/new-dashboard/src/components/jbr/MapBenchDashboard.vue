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
        <BranchSelect
          :branch-configurator="branchConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
      </template>
    </Toolbar>

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <div
          v-for="metric in metricsNames"
          :key="metric"
        >
          <div class="relative flex py-5 items-center">
            <div class="flex-grow border-t border-gray-400" />
            <span class="flex-shrink mx-4 text-gray-400 text-lg">{{ metric }}</span>
            <div class="flex-grow border-t border-gray-400" />
          </div>
          <section>
            <GroupProjectsChart
              label="macOS"
              :measure="metric"
              :projects="macOSConfigurations"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </section>
          <section class="flex gap-x-6">
            <div class="flex-1">
              <GroupProjectsChart
                label="Ubuntu"
                :measure="metric"
                :projects="ubuntuConfigurations"
                :server-configurator="serverConfigurator"
                :configurators="dashboardConfigurators"
              />
            </div>
            <div class="flex-1">
              <GroupProjectsChart
                label="Windows"
                :measure="metric"
                :projects="windowsConfigurations"
                :server-configurator="serverConfigurator"
                :configurators="dashboardConfigurators"
              />
            </div>
          </section>
        </div>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider(false, true)

const dbName = "jbr"
const dbTable = "report"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("jbr_mapbench_dashboard", {
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])

const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [branchConfigurator, timeRangeConfigurator],
)

const dashboardConfigurators = [
  branchConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator,
]

const metricsNames = ["CircleTests", "EllipseTests-fill-false", "EllipseTests-fill-true", "spiralTest-dash-false", "spiralTest-fill", "dc_boulder_2013-13-30-06-13-17",
  "dc_boulder_2013-13-30-06-13-20", "dc_shp_alllayers_2013-00-30-07-00-43", "dc_shp_alllayers_2013-00-30-07-00-47", "dc_spearfish_2013-11-30-06-11-15",
  "dc_spearfish_2013-11-30-06-11-19", "test_z_625k"].flatMap(test => {
  return ["ser.Pct95"].map(stat => test + "." + stat)
})
const ubuntuConfigurations = ["Ubuntu2004x64", "Ubuntu2004x64OGL", "Ubuntu2204x64", "Ubuntu2204x64OGL"].map(config => "Mapbench_" + config)
const macOSConfigurations = ["macOS13x64OGL", "macOS13x64Metal", "macOS13aarch64OGL", "macOS13aarch64Metal", "macOS12x64OGL", "macOS12x64Metal", "macOS12aarch64OGL",
  "macOS12aarch64Metal"].map(config => "Mapbench_" + config)
const windowsConfigurations = ["Windows10x64"].map(config => "Mapbench_" + config)

function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}

</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>