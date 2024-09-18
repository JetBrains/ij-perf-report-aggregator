<template>
  <StickyToolbar>
    <template #start>
      <TimeRangeSelect :timerange-configurator="timeRangeConfigurator" />
      <BranchSelect
        :branch-configurator="branchConfigurator"
        :triggered-by-configurator="triggeredByConfigurator"
        :selection-limit="1"
      />
      <K1VsK2ComparisonProjectCategoryFilter
        :initial-project-categories="initialProjectCategories"
        @update:selected-project-categories="(newValue: string[]) => (selectedProjectCategories = newValue)"
      />
      <span class="p-buttonset ml-4">
        <Button
          v-for="table in tables"
          :key="table.name"
          :label="table.name"
          :outlined="activeTab != table.name"
          @click="setActiveTab(table.name)"
        />
      </span>
    </template>
  </StickyToolbar>

  <Tabs
    class="k1-vs-k2-comparison-tab-view"
    v-model:value="activeTab"
  >
    <TabPanels>
      <TabList>
        <Tab
          v-for="table in tables"
          :value="table.name"
          >{{ table.name }}</Tab
        >
      </TabList>
      <TabPanel
        v-for="table in tables"
        :key="table.name"
        :value="table.name"
      >
        <K1VsK2ComparisonTable
          :name="table.name"
          :measure="table.measure"
          :projects="table.projects"
          :allowed-project-categories="selectedProjectCategories"
          :configurators="configurators"
        />
      </TabPanel>
    </TabPanels>
  </Tabs>
</template>

<script setup lang="ts">
import { provide, Ref, ref } from "vue"
import { useRouter } from "vue-router"
import { createBranchConfigurator } from "../../../configurators/BranchConfigurator"
import { MachineConfigurator } from "../../../configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../../configurators/PrivateBuildConfigurator"
import { ServerWithCompressConfigurator } from "../../../configurators/ServerWithCompressConfigurator"
import { TimeRangeConfigurator } from "../../../configurators/TimeRangeConfigurator"
import { serverConfiguratorKey } from "../../../shared/keys"
import BranchSelect from "../../common/BranchSelect.vue"
import { PersistentStateManager } from "../../common/PersistentStateManager"
import StickyToolbar from "../../common/StickyToolbar.vue"
import TimeRangeSelect from "../../common/TimeRangeSelect.vue"
import { completionProjects, findUsagesProjects, highlightingProjects, MACHINES } from "../projects"
import K1VsK2ComparisonProjectCategoryFilter from "./K1VsK2ComparisonProjectCategoryFilter.vue"
import K1VsK2ComparisonTable from "./K1VsK2ComparisonTable.vue"

interface Props {
  dbName: string
  persistentId: string
  initialBranch: string
}

const { dbName, persistentId, initialBranch } = defineProps<Props>()

const dbTable = "kotlin"

const router = useRouter()

const serverConfigurator = new ServerWithCompressConfigurator(dbName, dbTable)
provide(serverConfiguratorKey, serverConfigurator)

const persistentStateManager = new PersistentStateManager(
  persistentId,
  {
    branch: initialBranch,
    projectCategories: [],
  },
  router
)

const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator], true, [MACHINES.linux, MACHINES.mac])

const configurators = [timeRangeConfigurator, branchConfigurator, triggeredByConfigurator, machineConfigurator]

const selectedProjectCategories: Ref<string[]> = ref([])

persistentStateManager.add("projectCategories", selectedProjectCategories, (existingValue) => {
  return typeof existingValue === "string" ? [existingValue] : existingValue
})

// The initial selected project categories are taken from the initial state of the persistent state manager.
const initialProjectCategories = selectedProjectCategories.value

const activeTab = ref()

function setActiveTab(index: string) {
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
  localInspections: {
    name: "Code analysis",
    measure: "localInspections#mean_value",
    projects: flattenProjectCategories(highlightingProjects),
  },
  highlighting: {
    name: "Semantic highlighting",
    measure: "semanticHighlighting#mean_value",
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
/* This together with the button set is basically a workaround for missing tabview styles, but I also like the button set design more
 * than the tab nav.
 */
.k1-vs-k2-comparison-tab-view .p-tabview-nav {
  display: none;
}
</style>
