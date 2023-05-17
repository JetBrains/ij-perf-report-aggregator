<template>
  <div class="flex flex-col gap-5">
    <DashboardToolbar
      :branch-configurator="branchConfigurator"
      :on-change-range="onChangeRange"
      :time-range-configurator="timeRangeConfigurator"
      :triggered-by-configurator="triggeredByConfigurator"
    />

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
            <div class="flex-1 min-w-0">
              <GroupProjectsChart
                label="Ubuntu"
                :measure="metric"
                :projects="ubuntuConfigurations"
                :server-configurator="serverConfigurator"
                :configurators="dashboardConfigurators"
              />
            </div>
            <div class="flex-1 min-w-0">
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
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardToolbar from "../common/DashboardToolbar.vue"

provideReportUrlProvider(false, true)

const dbName = "jbr"
const dbTable = "report"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("jbr_javadraw_dashboard", {
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

const metricsNames = ["Plus_200_Random_Small_Circles", "Plus_2_SweepGradient_Circles", "Plus_320_Long_Lines", "Plus_4000_Random_Small_Circles"]
const ubuntuConfigurations = ["Ubuntu2004x64", "Ubuntu2004x64OGL", "Ubuntu2204x64", "Ubuntu2204x64OGL"].map(config => "JavaDraw_" + config)
const macOSConfigurations = ["macOS13x64OGL", "macOS13x64Metal", "macOS13aarch64OGL", "macOS13aarch64Metal", "macOS12x64OGL", "macOS12x64Metal", "macOS12aarch64OGL",
  "macOS12aarch64Metal"].map(config => "JavaDraw_" + config)
const windowsConfigurations = ["Windows10x64"].map(config => "JavaDraw_" + config)

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

</script>