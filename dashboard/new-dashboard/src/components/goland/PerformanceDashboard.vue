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
            label="Indexing: Lightweight projects"
            measure="indexing"
            :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number Of Indexed Files: Lightweight projects"
            measure="numberOfIndexedFiles"
            :projects="['flux/indexing', 'delve/indexing', 'istio/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing: Heavyweight projects"
            measure="indexing"
            :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number Of Indexed Files: Heavyweight projects"
            measure="numberOfIndexedFiles"
            :projects="['moby/indexing', 'mattermost-server/indexing', 'cockroach/indexing', 'kubernetes/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Inspection execution time: Lightweight projects"
            measure="globalInspections"
            :projects="['istio/inspection', 'moby/inspection', 'flux/inspection', 'delve/inspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Inspection execution time: Heavyweight projects"
            measure="globalInspections"
            :projects="['cockroach/inspection', 'kubernetes/inspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Local inspection execution time"
            measure="localInspections"
            :projects="['kubernetes/localInspection', 'mattermost-server/localInspection', 'GO-5422/localInspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Typing: average responsiveness time"
            measure="test#average_awt_delay"
            :projects="['mattermost-server/typing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Typing: total time"
            measure="typing"
            :projects="['mattermost-server/typing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Find Usages execution time"
            measure="findUsages"
            :projects="['vault/findUsages/Backend', 'vault/findUsages/List', 'vault/findUsages/Path', 'vault/findUsages/String']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Find Usages number of found usages"
            measure="findUsages#number"
            :projects="['vault/findUsages/Backend', 'vault/findUsages/List', 'vault/findUsages/Path', 'vault/findUsages/String']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
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
import { refToObservable } from "shared/src/configurators/rxjs"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { Accident, getAccidentsFromMetaDb } from "shared/src/meta"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardToolbar from "../common/DashboardToolbar.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "goland"
const initialMachine = "linux-blade-hetzner"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("goland_dashboard", {
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
const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getAccidentsFromMetaDb(warnings, null, data)
})
</script>