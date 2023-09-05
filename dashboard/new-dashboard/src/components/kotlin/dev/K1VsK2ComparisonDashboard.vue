<template>
  <Toolbar class="customToolbar">
    <template #start>
      <TimeRangeSelect
        :ranges="TimeRangeConfigurator.timeRanges"
        :value="timeRangeConfigurator.value.value"
        :on-change="onTimeRangeChange"
      />
      <BranchSelect
        :branch-configurator="branchConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
        :selection-limit="1"
      />
      <MachineSelect :machine-configurator="machineConfigurator" />
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
      <K1VsK2ComparisonTable
        :name="table.name"
        :measure="table.measure"
        :projects="table.projects"
        :configurators="configurators"
      />
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
import BranchSelect from "../../common/BranchSelect.vue"
import MachineSelect from "../../common/MachineSelect.vue"
import { PersistentStateManager } from "../../common/PersistentStateManager"
import TimeRangeSelect from "../../common/TimeRangeSelect.vue"
import { completionProjects, findUsagesProjects, highlightingProjects } from "../projects"
import K1VsK2ComparisonTable from "./K1VsK2ComparisonTable.vue"

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

/**
 * This property contains the data with which separate comparison tables are built. The data for each category includes the performance test names without their variant (k1, k2)
 * suffixes.
 */
const tables = {
  completion: {
    name: "Completion (all elements)",
    measure: "completion#mean_value",
    projects: flattenProjectCategories(completionProjects),
  },
  completionFirstElement: {
    name: "Completion (first element)",
    measure: "completion#firstElementShown#mean_value",
    projects: flattenProjectCategories(completionProjects),
  },
  highlighting: {
    name: "Semantic highlighting",
    measure: "semanticHighlighting#mean_value",
    projects: flattenProjectCategories(highlightingProjects),
  },
  localInspections: {
    name: "Local Inspections",
    measure: "localInspections#mean_value",
    projects: flattenProjectCategories(highlightingProjects),
  },
  findUsages: {
    name: "Find Usages",
    measure: "findUsages#mean_value",
    projects: flattenProjectCategories(findUsagesProjects),
  },
}

function flattenProjectCategories(projectsByCategory: Record<string, string[]>) {
  const result: string[] = []
  for (const [_, projects] of Object.entries(projectsByCategory)) {
    result.push(...projects)
  }
  return result
}
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}

/* This together with the button set is basically a workaround for missing tabview styles, but I also like the button set design more
 * than the tab nav.
 */
.k1-vs-k2-comparison-tab-view .p-tabview-nav {
  display: none;
}
</style>
