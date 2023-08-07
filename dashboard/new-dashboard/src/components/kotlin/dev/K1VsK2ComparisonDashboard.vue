<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="TimeRangeConfigurator.timeRanges"
        :value="timeRangeConfigurator.value.value"
        :on-change="onTimeRangeChange"
      >
        <template #icon>
          <CalendarIcon class="w-4 h-4 text-gray-500" />
        </template>
      </TimeRangeSelect>
      <BranchSelect
        :branch-configurator="branchConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
        :selection-limit="1"
      />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      >
        <template #icon>
          <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
        </template>
      </DimensionHierarchicalSelect>
      <span class="p-buttonset ml-4">
        <Button
          v-for="(table, _, index) in tables"
          :key="table.name"
          :label="table.name"
          :outlined="activeTab != index"
          @click="setActiveTab(index)"
        />
      </span>
    </template>
  </Toolbar>

  <TabView
    v-model:activeIndex="activeTab"
    class="k1-vs-k2-comparison-tab-view"
  >
    <TabPanel
      v-for="table in tables"
      :key="table.name"
      :header="table.name"
    >
      <section class="flex flex-col w-full mt-8">
        <div class="flex flex-row gap-6">
          <div class="basis-1/2">
            <h3 class="text-2xl mb-3">{{ table.name }}</h3>
            <p class="text-sm text-gray-600">Measure: {{ table.measure }}</p>
          </div>
        </div>
        <TestComparisonTable
          :headline="table.name"
          :measure="table.measure"
          :comparisons="table.projects.map(transformToTestComparison)"
          :configurators="configurators"
          baseline-column-label="K1"
          current-column-label="K2"
          difference-column-label="Improvement (%)"
          class="mt-8"
        />
        <p class="text-sm text-gray-500 text-right mt-4">The table only displays the results of the last build from the selected branch.</p>
      </section>
    </TabPanel>
  </TabView>
</template>

<script setup lang="ts">
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../../configurators/BranchConfigurator"
import { MachineConfigurator } from "../../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../../configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "../../../configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "../../../configurators/TimeRangeConfigurator"
import { serverConfiguratorKey } from "../../../shared/keys"
import DimensionHierarchicalSelect from "../../charts/DimensionHierarchicalSelect.vue"
import BranchSelect from "../../common/BranchSelect.vue"
import { PersistentStateManager } from "../../common/PersistentStateManager"
import TestComparisonTable, { TestComparison } from "../../common/TestComparisonTable.vue"
import TimeRangeSelect from "../../common/TimeRangeSelect.vue"

const dbName = "perfintDev"
const dbTable = "kotlin"
const initialMachine = "linux-blade-hetzner"
const persistentId = "kotlinDev_k1VsK2Comparison"

const router = useRouter()

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
provide(serverConfiguratorKey, serverConfigurator)

const persistentStateManager = new PersistentStateManager(
  persistentId,
  {
    machine: initialMachine,
    branch: "kt-master",
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])

const configurators = [timeRangeConfigurator, branchConfigurator, triggeredByConfigurator, machineConfigurator]

const activeTab = ref(0)

function onTimeRangeChange(value: TimeRange) {
  timeRangeConfigurator.value.value = value
}

function setActiveTab(index: number) {
  activeTab.value = index
}

function transformToTestComparison(projectName: string): TestComparison {
  // We want to compare K1 and K2 tests against each other, and they are respectively suffixed with "_k1" and "_k2".
  return {
    label: projectName,
    baselineTestName: `${projectName}_k1`,
    currentTestName: `${projectName}_k2`,
  }
}

/**
 * This property contains the data with which separate comparison tables are built. The data for each category includes the performance test names without their variant (k1, k2)
 * suffixes.
 */
const tables = {
  completion: {
    name: "Completion",
    measure: "completion#mean_value",
    projects: ["kotlin_empty/completion/empty_place_with_library_cache", "kotlin_empty/completion/empty_place_typing_with_library_cache"],
  },
}
</script>

<style>
/* This together with the button set is basically a workaround for missing tabview styles, but I also like the button set design more
 * than the tab nav.
 */
.k1-vs-k2-comparison-tab-view .p-tabview-nav {
  display: none;
}
</style>
