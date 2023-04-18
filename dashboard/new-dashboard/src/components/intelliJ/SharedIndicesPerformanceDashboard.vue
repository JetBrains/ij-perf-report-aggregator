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
            label="Indexing (IntelliJ project)"
            measure="indexing"
            :projects="['ijx-intellij-speed/downloading-shared-indexes',
                        'ijx-intellij-speed/shared-indexes',
                        'ijx-intellij-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning (IntelliJ project)"
            measure="scanning"
            :projects="['ijx-intellij-speed/downloading-shared-indexes',
                        'ijx-intellij-speed/shared-indexes',
                        'ijx-intellij-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing runs (IntelliJ project)"
            measure="numberOfIndexingRuns"
            :projects="['ijx-intellij-speed/downloading-shared-indexes',
                        'ijx-intellij-speed/shared-indexes',
                        'ijx-intellij-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by usual indexes (IntelliJ project)"
            measure="numberOfIndexedFiles"
            :projects="['ijx-intellij-speed/downloading-shared-indexes',
                        'ijx-intellij-speed/shared-indexes',
                        'ijx-intellij-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by shared indexes (IntelliJ project)"
            measure="numberOfFilesIndexedByExtensions"
            :projects="['ijx-intellij-speed/downloading-shared-indexes',
                        'ijx-intellij-speed/shared-indexes',
                        'ijx-intellij-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Updating time (IntelliJ project)"
            measure="updatingTime"
            :projects="['ijx-intellij-speed/downloading-shared-indexes',
                        'ijx-intellij-speed/shared-indexes',
                        'ijx-intellij-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing (Community project)"
            measure="indexing"
            :projects="['ijx-intellij-community-speed/downloading-shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-community-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-community-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning (Community project)"
            measure="scanning"
            :projects="['ijx-intellij-community-speed/downloading-shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-community-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-community-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing runs (Community project)"
            measure="numberOfIndexingRuns"
            :projects="['ijx-intellij-community-speed/downloading-shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-community-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-community-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by usual indexes (Community project)"
            measure="numberOfIndexedFiles"
            :projects="['ijx-intellij-community-speed/downloading-shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-community-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-community-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by shared indexes (Community project)"
            measure="numberOfFilesIndexedByExtensions"
            :projects="['ijx-intellij-community-speed/downloading-shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-community-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-community-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Updating time (Community project)"
            measure="updatingTime"
            :projects="['ijx-intellij-community-speed/downloading-shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes',
                        'ijx-intellij-community-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-intellij-community-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-intellij-community-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing (JDK Maven project)"
            measure="indexing"
            :projects="['ijx-jdk-maven-speed/downloading-shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-maven-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-maven-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning (JDK Maven project)"
            measure="scanning"
            :projects="['ijx-jdk-maven-speed/downloading-shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-maven-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-maven-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing runs (JDK Maven project)"
            measure="numberOfIndexingRuns"
            :projects="['ijx-jdk-maven-speed/downloading-shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-maven-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-maven-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by usual indexes (JDK Maven project)"
            measure="numberOfIndexedFiles"
            :projects="['ijx-jdk-maven-speed/downloading-shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-maven-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-maven-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by shared indexes (JDK Maven project)"
            measure="numberOfFilesIndexedByExtensions"
            :projects="['ijx-jdk-maven-speed/downloading-shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-maven-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-maven-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Updating time (JDK Maven project)"
            measure="updatingTime"
            :projects="['ijx-jdk-maven-speed/downloading-shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes',
                        'ijx-jdk-maven-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-maven-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-maven-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing (JDK Gradle project)"
            measure="indexing"
            :projects="['ijx-jdk-gradle-speed/downloading-shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-gradle-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning (JDK Gradle project)"
            measure="scanning"
            :projects="['ijx-jdk-gradle-speed/downloading-shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-gradle-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing runs (JDK Gradle project)"
            measure="numberOfIndexingRuns"
            :projects="['ijx-jdk-gradle-speed/downloading-shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-gradle-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by usual indexes (JDK Gradle project)"
            measure="numberOfIndexedFiles"
            :projects="['ijx-jdk-gradle-speed/downloading-shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-gradle-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing files by shared indexes (JDK Gradle project)"
            measure="numberOfFilesIndexedByExtensions"
            :projects="['ijx-jdk-gradle-speed/downloading-shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-gradle-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Updating time (JDK Gradle project)"
            measure="updatingTime"
            :projects="['ijx-jdk-gradle-speed/downloading-shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-archive-and-git-hashes',
                        'ijx-jdk-gradle-speed/shared-indexes-with-hashes-without-javaScript-plugin',
                        'ijx-jdk-gradle-speed/usual-indexes']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
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

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "ideaSharedIndices"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("shared_indexes_dashboard", {
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
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>