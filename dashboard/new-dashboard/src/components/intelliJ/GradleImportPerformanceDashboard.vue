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
          :release-configurator="releaseConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
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
        <section>
          <GroupProjectsChart
            label="CPU | Load |Total % 95th pctl"
            measure="CPU | Load |Total % 95th pctl"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="CPU | IDE | Threads 95th pctl"
            measure="CPU | IDE | Threads 95th pctl"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Memory | IDE | RESIDENT SIZE (MB) 95th pctl"
            measure="Memory | IDE | RESIDENT SIZE (MB) 95th pctl"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Memory | IDE | VIRTUAL SIZE (MB) 95th pctl"
            measure="Memory | IDE | VIRTUAL SIZE (MB) 95th pctl"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="gcPause"
            measure="gcPause"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="gcPauseCount"
            measure="gcPauseCount"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="fullGCPause"
            measure="fullGCPause"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="freedMemoryByGC"
            measure="freedMemoryByGC"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="totalHeapUsedMax"
            measure="totalHeapUsedMax"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="gradle.sync.duration"
            measure="gradle.sync.duration"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="GRADLE_CALL"
            measure="GRADLE_CALL"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="PROJECT_RESOLVERS"
            measure="PROJECT_RESOLVERS"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="DATA_SERVICES"
            measure="DATA_SERVICES"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="WORKSPACE_MODEL_APPLY"
            measure="WORKSPACE_MODEL_APPLY"
            :projects="[
              'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
              'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
              'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
              'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
              'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
              'project-import-gradle-openliberty/measureStartup'
            ]"
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
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { Accident, getWarningFromMetaDb } from "shared/src/meta"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "idea"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("idea_gradle_dashboard", {
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

function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}

const warnings = ref<Array<Accident>>()
refToObservable(timeRangeConfigurator.value).subscribe(data => {
  getWarningFromMetaDb(warnings, null, data as TimeRange)
})
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>