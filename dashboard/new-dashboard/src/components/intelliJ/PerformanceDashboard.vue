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
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'processingSpeed#JAVA'"
              :title="'Indexing Java (kB/s)'"
              :chart-color="'#219653'"
              :value-unit="'counter'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'processingSpeed#Kotlin'"
              :title="'Indexing Kotlin (kB/s)'"
              :chart-color="'#9B51E0'"
              :value-unit="'counter'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'completion\_%'"
              :is-like="true"
              :title="'Completion'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="[...averagesConfigurators, typingOnlyConfigurator]"
              :aggregated-measure="'test#average_awt_delay'"
              :title="'UI responsiveness during typing'"
              :chart-color="'#F2994A'"
            />
          </div>
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing (Big projects)"
            measure="indexing"
            :projects="['community/indexing', 'intellij_sources/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning (Big projects)"
            measure="scanning"
            :projects="['community/indexing', 'intellij_sources/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing runs long (Big projects)"
            measure="numberOfIndexingRuns"
            :projects="['community/indexing', 'intellij_sources/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Processing speed of JAVA files"
            measure="processingSpeed#JAVA"
            :projects="['community/indexing', 'intellij_sources/indexing', 'empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 
                        'kotlin_coroutines/indexing', 'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Processing speed of KOTLIN files"
            measure="processingSpeed#Kotlin"
            :projects="['community/indexing', 'intellij_sources/indexing', 'empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 
                        'kotlin_coroutines/indexing', 'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing with the new record storage (IntelliJ project)"
            measure="indexing"
            :projects="['vfs-record-storage/in-memory-intellij_sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij_sources/indexing',
                        'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij_sources/indexing', 
                        'vfs-record-storage/in-memory-with-streamlined-attributes-intellij_sources/indexing', 'vfs-record-storage/lock-free-intellij_sources/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning with the new record storage (IntelliJ project)"
            measure="scanning"
            :projects="['vfs-record-storage/in-memory-intellij_sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij_sources/indexing',
                        'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij_sources/indexing', 
                        'vfs-record-storage/in-memory-with-streamlined-attributes-intellij_sources/indexing','vfs-record-storage/lock-free-intellij_sources/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing runs with the new record storage (IntelliJ project)"
            measure="numberOfIndexingRuns"
            :projects="['vfs-record-storage/in-memory-intellij_sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij_sources/indexing',
                        'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij_sources/indexing', 
                        'vfs-record-storage/in-memory-with-streamlined-attributes-intellij_sources/indexing', 'vfs-record-storage/lock-free-intellij_sources/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="indexing"
            :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                        'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Scanning"
            measure="scanning"
            :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                        'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Number of indexing runs"
            measure="numberOfIndexingRuns"
            :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing',
                        'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="VFS Refresh"
            measure="vfs_initial_refresh"
            :projects="['intellij_sources/vfsRefresh/default', 'intellij_sources/vfsRefresh/with-1-thread(s)', 'intellij_sources/vfsRefresh/git-status']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Rebuild (Big projects)"
            measure="build_compilation_duration"
            :projects="['community/rebuild','intellij_sources/rebuild']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Rebuild"
            measure="build_compilation_duration"
            :projects="['grails/rebuild','java/rebuild','spring_boot/rebuild']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Inspection"
            measure="globalInspections"
            :projects="['java/inspection', 'grails/inspection', 'spring_boot_maven/inspection', 'spring_boot/inspection', 'kotlin/inspection', 'kotlin_coroutines/inspection']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="FindUsages PsiManager#getInstance Before and After Compilation"
            measure="findUsages"
            :projects="['community/findUsages/PsiManager_getInstance_Before', 'community/findUsages/PsiManager_getInstance_After'
            ]"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="FindUsages Library#getName Before and After Compilation"
            measure="findUsages"
            :projects="['community/findUsages/Library_getName_Before', 'community/findUsages/Library_getName_After'
            ]"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="FindUsages LocalInspectionTool#getID Before and After Compilation"
            measure="findUsages"
            :projects="['community/findUsages/LocalInspectionTool_getID_Before', 'community/findUsages/LocalInspectionTool_getID_After'
            ]"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="FindUsages ActionsKt#runReadAction and Application#runReadAction Before and After Compilation"
            measure="findUsages"
            :projects="['community/findUsages/ActionsKt_runReadAction_Before', 'community/findUsages/Application_runReadAction_Before',
                        'community/findUsages/ActionsKt_runReadAction_After', 'community/findUsages/Application_runReadAction_After'
            ]"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="FindUsages Persistent#absolutePath and PropertyMapping#value Before and After Compilation"
            measure="findUsages"
            :projects="['community/findUsages/Persistent_absolutePath_After', 'community/findUsages/PropertyMapping_value_After',
                        'community/findUsages/Persistent_absolutePath_Before', 'community/findUsages/PropertyMapping_value_Before'
            ]"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Find Usages with idea.is.internal=true Before Compilation"
            measure="findUsages"
            :projects="['intellij_sources/findUsages/PsiManager_getInstance_firstCall', 'intellij_sources/findUsages/PsiManager_getInstance_secondCall',
                        'intellij_sources/findUsages/PsiManager_getInstance_thirdCallInternalMode']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Local Inspection"
            measure="localInspections"
            :projects="['intellij_sources/localInspection/java_file','intellij_sources/localInspection/kotlin_file', 'kotlin/localInspection',
                        'kotlin_coroutines/localInspection', 'gradle_kts_vulnerable_dep/localInspection']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="First Code Analysis"
            measure="firstCodeAnalysis"
            :projects="['intellij_sources/localInspection/java_file','intellij_sources/localInspection/kotlin_file', 'kotlin/localInspection',
                        'kotlin_coroutines/localInspection', 'gradle_kts_vulnerable_dep/localInspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Completion"
            measure="completion"
            :projects="['community/completion/kotlin_file','grails/completion/groovy_file', 'grails/completion/java_file']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Search Everywhere"
            measure="searchEverywhere"
            :projects="['community/go-to-action/SharedIndex', 'community/go-to-class/EditorImpl','community/go-to-class/SharedIndex',
                        'community/go-to-file/EditorImpl','community/go-to-file/SharedIndex']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Search Everywhere (Dialog Shown)"
            measure="searchEverywhere_dialog_shown"
            :projects="['community/go-to-action/SharedIndex', 'community/go-to-class/EditorImpl','community/go-to-class/SharedIndex',
                        'community/go-to-file/EditorImpl','community/go-to-file/SharedIndex']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Search Everywhere (Items Loaded)"
            measure="searchEverywhere_items_loaded"
            :projects="['community/go-to-action/SharedIndex', 'community/go-to-class/EditorImpl','community/go-to-class/SharedIndex',
                        'community/go-to-file/EditorImpl','community/go-to-file/SharedIndex']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Show Intentions (average awt delay)"
            measure="test#average_awt_delay"
            :projects="['grails/showIntentions/Find cause', 'kotlin/showIntention/Import', 'spring_boot/showIntentions']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Show File History"
            measure="showFileHistory"
            :projects="['intellij_sources/showFileHistory/EditorImpl']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Expand Menu"
            measure="expandActionGroup"
            :projects="['intellij_sources/expandProjectMenu', 'intellij_sources/expandMainMenu', 'intellij_sources/expandEditorMenu']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Highlight"
            measure="highlighting"
            :projects="['kotlin/highlight', 'kotlin_coroutines/highlight']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="FileStructure"
            measure="FileStructurePopup"
            :projects="['intellij_sources/FileStructureDialog/java_file', 'intellij_sources/FileStructureDialog/kotlin_file']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Typing during indexing (with changed count of indexing threads)"
            measure="typing"
            :projects="['typingInJavaFile_16Threads/typing', 'typingInJavaFile_4Threads/typing', 'typingInKotlinFile_16Threads/typing', 'typingInKotlinFile_4Threads/typing']"
            :server-configurator="serverConfigurator" 
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Typing during indexing (average awt delay)"
            measure="test#average_awt_delay"
            :projects="['typingInJavaFile_16Threads/typing', 'typingInJavaFile_4Threads/typing', 'typingInKotlinFile_16Threads/typing', 'typingInKotlinFile_4Threads/typing']"
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
import { DataQuery, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
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
const dbTable = "idea"
const initialMachine = "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("idea_dashboard", {
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
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator,
  triggeredByConfigurator,
]

const typingOnlyConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({f: "project", v: "%typing", o: "like"})
    return true
  },
  createObservable() {
    return null
  },
}

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