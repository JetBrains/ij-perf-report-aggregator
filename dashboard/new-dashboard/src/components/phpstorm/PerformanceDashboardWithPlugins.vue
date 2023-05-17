<template>
  <div class="flex flex-col gap-5">
    <DashboardToolbar
      :branch-configurator="branchConfigurator"
      :machine-configurator="machineConfigurator"
      :release-configurator="releaseConfigurator"
      :on-change-range="onChangeRange"
      :time-range-configurator="timeRangeConfigurator"
      :triggered-by-configurator="triggeredByConfigurator"
    />

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <section>
          <GroupProjectsChart
            label="Batch Inspections"
            measure="globalInspections"
            :projects="['drupal8-master-with-plugin/inspection', 'magento/inspection', 'wordpress/inspection', 'laravel-io/inspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Batch Inspections"
              measure="globalInspections"
              :projects="['mediawiki/inspection']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Local Inspections"
              measure="localInspections"
              :projects="['mpdf/localInspection']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Indexing"
              measure="updatingTime"
              :projects="['drupal8-master-with-plugin/indexing', 'laravel-io/indexing','wordpress/indexing','mediawiki/indexing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Indexing"
              measure="updatingTime"
              :projects="['magento/indexing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>

        <section class="flex gap-x-6">
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Typing Time"
              measure="typing"
              :projects="['mpdf/typing', 'mpdf_powersave/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Typing Average Responsiveness"
              measure="test#average_awt_delay"
              :projects="[ 'mpdf/typing', 'mpdf_powersave/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1 min-w-0">
            <GroupProjectsChart
              label="Typing Responsiveness"
              measure="test#max_awt_delay"
              :projects="['mpdf/typing', 'mpdf_powersave/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
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

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "phpstormWithPlugins"
const initialMachine = "linux-blade-hetzner"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("phpstorm_plugins_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator, branchConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistenceForDashboard)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [branchConfigurator, timeRangeConfigurator],
)


const dashboardConfigurators = [
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator,
  triggeredByConfigurator,
]

function onChangeRange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}
</script>