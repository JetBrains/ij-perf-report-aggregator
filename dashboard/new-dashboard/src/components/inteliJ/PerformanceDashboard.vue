<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect :configurator="timeRangeConfigurator" />
        <DimensionSelect
          label="Branch"
          :dimension="branchConfigurator"
        />
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        />
      </template>
    </Toolbar>

    <div class="flex flex-col gap-6">
      <section>
        <GroupChart
          label="Indexing Long"
          measure="indexing"
          :projects="['community/indexing', 'lock-free-vfs-record-storage-intellij_sources/indexing', 'intellij_sources/indexing']"
          :server-configurator="serverConfigurator"
        />
      </section>

      <section class="flex gap-x-6">
        <div class="flex-1">
          <GroupChart
            label="Kotlin Builder Long"
            measure="kotlin_builder_time"
            :projects="['community/rebuild','intellij_sources/rebuild']"
            :server-configurator="serverConfigurator"
          />
        </div>
        <div class="flex-1">
          <GroupChart
            label="Rebuild Long"
            measure="build_compilation_duration"
            :projects="['community/rebuild','intellij_sources/rebuild']"
            :server-configurator="serverConfigurator"
          />
        </div>
      </section>
    </div>
  </div>
</template>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>

<script setup lang="ts">
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"
import GroupChart from "../common/GroupChart.vue"

provide(chartStyleKey, {
    ...chartDefaultStyle,
})

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("idea_performance_dashboard", {
    machine: "macMini Intel 3.2, 16GB",
    project: [],
    branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "idea")
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true, [timeRangeConfigurator], (a, _) => {
    return a.includes("/") ? 1 : -1
})
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const releaseConfigurator = new ReleaseNightlyConfigurator(persistentStateManager)
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])

const configurators = [
    serverConfigurator,
    branchConfigurator,
    machineConfigurator,
    timeRangeConfigurator,
    releaseConfigurator,
    triggeredByConfigurator,
]
initDataComponent(configurators)
</script>