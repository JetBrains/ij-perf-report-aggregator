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
        <section class="flex gap-6">
          <div class="w-1/2">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'processingSpeed#Rust'"
              :title="'Indexing Rust (kB/s)'"
              :chart-color="'#219653'"
              :value-unit="'counter'"
            />
          </div>
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="indexing"
            :projects="['intelli-j-with-rust-test/test-rustling-cargo-sync', 'intelli-j-with-rust-test/run-ide-with-rust-plugin',
                        'intelli-j-with-rust-test/test-drogue-cloud-c-p-u-usage']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Inspection execution time"
            measure="globalInspections"
            :projects="['intelli-j-with-rust-test/test-cargo-inspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Inspection execution time"
            measure="globalInspections"
            :projects="['intelli-j-with-rust-test/test-cargo-inspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Completion"
            measure="completion"
            :projects="['intelli-j-with-rust-test/test-arrow-rs-completion', 'intelli-j-with-rust-test/test-vec-completion']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Find Usages"
            measure="findUsages"
            :projects="['intelli-j-with-rust-test/run-ide-with-rust-plugin-find-usages', 'intelli-j-with-rust-test/run-ide-with-rust-plugin-wasm-find-usages']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Local Inspections (on file open)"
            measure="firstCodeAnalysis"
            :projects="[ 'intelli-j-with-rust-test/test-arrow-rs-highlighting', 'intelli-j-with-rust-test/test-cargo-highlighting', 
                         'intelli-j-with-rust-test/test-my-sql-async-highlighting']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
            :accidents="warnings"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Cargo Sync"
            measure="cargo_sync_execution_time"
            :projects="[ 'intelli-j-with-rust-test/test-rustling-cargo-sync']"
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
import { Accident, getAccidentsFromMetaDb } from "shared/src/meta"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "rust"
const initialMachine = "Linux EC2 m5d.xlarge or 5d.xlarge or m5ad.xlarge"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("rust_dashboard", {
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

const averagesConfigurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
]

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
  getAccidentsFromMetaDb(warnings, null, data as TimeRange)
})
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>