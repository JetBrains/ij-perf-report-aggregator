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
            v-for="chart in charts"
            :key="chart.definition.label"
            :label="chart.definition.label"
            :measure="chart.definition.measure"
            :projects="chart.projects"
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
import { combineLatest } from "rxjs"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRange, TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { refToObservable } from "shared/src/configurators/rxjs"
import { DataQuery, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { Accident, getWarningFromMetaDb } from "shared/src/meta"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import AggregationChart from "../charts/AggregationChart.vue"
import { ChartDefinition, combineCharts, extractUniqueProjects } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import AccidentWarning from "../common/AccidentWarning.vue"
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

const chartsDeclaration: Array<ChartDefinition> = [{
  labels: ["Indexing (Big projects)", "Scanning (Big projects)", "Number of indexing runs long (Big projects)"],
  measures: ["indexing", "scanning", "numberOfIndexingRuns"],
  projects: ["community/indexing", "intellij_sources/indexing"],
}, {
  labels: ["Processing speed of JAVA files", "Processing speed of KOTLIN files"],
  measures: ["processingSpeed#JAVA", "processingSpeed#Kotlin"],
  projects: ["community/indexing", "intellij_sources/indexing", "empty_project/indexing", "grails/indexing", "java/indexing", "kotlin/indexing",
    "kotlin_coroutines/indexing", "spring_boot/indexing", "spring_boot_maven/indexing", "kotlin_petclinic/indexing"],
}, {
  labels: ["Indexing with the new record storage (IntelliJ project)", "Scanning with the new record storage (IntelliJ project)",
    "Number of indexing runs with the new record storage (IntelliJ project)"],
  measures: ["indexing", "scanning", "numberOfIndexingRuns"],
  projects: ["vfs-record-storage/in-memory-intellij_sources/indexing", "vfs-record-storage/in-memory-with-non-strict-names-intellij_sources/indexing",
    "vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij_sources/indexing",
    "vfs-record-storage/in-memory-with-streamlined-attributes-intellij_sources/indexing", "vfs-record-storage/lock-free-intellij_sources/indexing"],
}, {
  labels: ["Indexing", "Scanning", "Number of indexing runs"],
  measures: ["indexing", "scanning", "numberOfIndexingRuns"],
  projects: ["empty_project/indexing", "grails/indexing", "java/indexing", "kotlin/indexing", "kotlin_coroutines/indexing",
    "spring_boot/indexing", "spring_boot_maven/indexing", "kotlin_petclinic/indexing"],
}, {
  labels: ["VFS Refresh"],
  measures: ["vfs_initial_refresh"],
  projects: ["intellij_sources/vfsRefresh/default", "intellij_sources/vfsRefresh/with-1-thread(s)", "intellij_sources/vfsRefresh/git-status"],
}, {
  labels: ["Rebuild (Big projects)"],
  measures: ["build_compilation_duration"],
  projects: ["community/rebuild", "intellij_sources/rebuild"],
}, {
  labels: ["Rebuild"],
  measures: ["build_compilation_duration"],
  projects: ["grails/rebuild", "java/rebuild", "spring_boot/rebuild"],
}, {
  labels: ["Inspection"],
  measures: ["globalInspections"],
  projects: ["java/inspection", "grails/inspection", "spring_boot_maven/inspection", "spring_boot/inspection", "kotlin/inspection", "kotlin_coroutines/inspection"],
}, {
  labels: ["FindUsages PsiManager#getInstance Before and After Compilation"],
  measures: ["findUsages"],
  projects: ["community/findUsages/PsiManager_getInstance_Before", "community/findUsages/PsiManager_getInstance_After"],
}, {
  labels: ["FindUsages Library#getName Before and After Compilation"],
  measures: ["findUsages"],
  projects: ["community/findUsages/Library_getName_Before", "community/findUsages/Library_getName_After"],
}, {
  labels: ["FindUsages LocalInspectionTool#getID Before and After Compilation"],
  measures: ["findUsages"],
  projects: ["community/findUsages/LocalInspectionTool_getID_Before", "community/findUsages/LocalInspectionTool_getID_After"],
}, {
  labels: ["FindUsages ActionsKt#runReadAction and Application#runReadAction Before and After Compilation"],
  measures: ["findUsages"],
  projects: ["community/findUsages/ActionsKt_runReadAction_Before", "community/findUsages/ActionsKt_runReadAction_After",
    "community/findUsages/Application_runReadAction_Before", "community/findUsages/Application_runReadAction_After"],
}, {
  labels: ["FindUsages Persistent#absolutePath and PropertyMapping#value Before and After Compilation"],
  measures: ["findUsages"],
  projects: ["community/findUsages/Persistent_absolutePath_Before", "community/findUsages/Persistent_absolutePath_After",
    "community/findUsages/PropertyMapping_value_Before", "community/findUsages/PropertyMapping_value_After"
  ],
}, {
    labels: ["FindUsages Object#hashCode and Path#toString Before and After Compilation"],
    measures: ["findUsages"],
    projects: ["community/findUsages/Object_hashCode_Before", "community/findUsages/Object_hashCode_After",
      "community/findUsages/Path_toString_Before", "community/findUsages/Path_toString_After"
    ],
  }, {
  labels: ["FindUsages Objects#hashCode Before and After Compilation"],
  measures: ["findUsages"],
  projects: ["community/findUsages/Objects_hashCode_Before", "community/findUsages/Objects_hashCode_After"],
}, {
  labels: ["FindUsages Path#div Before and After Compilation"],
  measures: ["findUsages"],
  projects: ["community/findUsages/Path_div_Before", "community/findUsages/Path_div_After"],
}, {
  labels: ["Find Usages with idea.is.internal=true Before Compilation"],
  measures: ["findUsages"],
  projects: ["intellij_sources/findUsages/PsiManager_getInstance_firstCall", "intellij_sources/findUsages/PsiManager_getInstance_secondCall",
    "intellij_sources/findUsages/PsiManager_getInstance_thirdCallInternalMode"],
}, {
  labels: ["Local Inspection", "First Code Analysis"],
  measures: ["localInspections", "firstCodeAnalysis"],
  projects: ["intellij_sources/localInspection/java_file", "intellij_sources/localInspection/kotlin_file", "kotlin/localInspection",
    "kotlin_coroutines/localInspection"],
}, {
  labels: ["Completion"],
  measures: ["completion"],
  projects: ["community/completion/kotlin_file", "grails/completion/groovy_file", "grails/completion/java_file"],
}, {
  labels: ["Search Everywhere", "Search Everywhere (Dialog Shown)", "Search Everywhere (Items Loaded)"],
  measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
  projects: ["community/go-to-action/SharedIndex", "community/go-to-class/EditorImpl", "community/go-to-class/SharedIndex",
    "community/go-to-file/EditorImpl", "community/go-to-file/SharedIndex"],
}, {
    labels: ["Search Everywhere Slow Typing", "Search Everywhere Slow Typing (Dialog Shown)", "Search Everywhere Slow Typing (Items Loaded)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown", "searchEverywhere_items_loaded"],
    projects: ["community/go-to-action/Runtime", "community/go-to-class/Kotlin", "community/go-to-file/properties"],
}, {
  labels: ["Show Intentions (average awt delay)"],
  measures: ["test#average_awt_delay"],
  projects: ["grails/showIntentions/Find cause", "kotlin/showIntention/Import", "spring_boot/showIntentions"],
}, {
  labels: ["Show File History"],
  measures: ["showFileHistory"],
  projects: ["intellij_sources/showFileHistory/EditorImpl"],
}, {
  labels: ["Expand Menu"],
  measures: ["expandActionGroup"],
  projects: ["intellij_sources/expandProjectMenu", "intellij_sources/expandMainMenu", "intellij_sources/expandEditorMenu"],
}, {
  labels: ["Highlight"],
  measures: ["highlighting"],
  projects: ["kotlin/highlight", "kotlin_coroutines/highlight"],
}, {
  labels: ["FileStructure"],
  measures: ["FileStructurePopup"],
  projects: ["intellij_sources/FileStructureDialog/java_file", "intellij_sources/FileStructureDialog/kotlin_file"],
}, {
  labels: ["Typing during indexing (with changed count of indexing threads)", "Typing during indexing (average awt delay)", "Typing during indexing (max awt delay)"],
  measures: ["typing", "test#average_awt_delay", "test#max_awt_delay"],
  projects: ["typingInJavaFile_16Threads/typing", "typingInJavaFile_4Threads/typing", "typingInKotlinFile_16Threads/typing", "typingInKotlinFile_4Threads/typing"],
}]

const charts = combineCharts(chartsDeclaration)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("idea_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)
const scenarioConfigurator = dimensionConfigurator(
  "project",
  serverConfigurator,
  null,
  true,
)
scenarioConfigurator.selected.value = extractUniqueProjects(chartsDeclaration)

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator, branchConfigurator, scenarioConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistenceForDashboard)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [branchConfigurator, timeRangeConfigurator, scenarioConfigurator],
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

const projects = chartsDeclaration.map(it => it.projects).flat(Number.POSITIVE_INFINITY) as Array<string>
const warnings = ref<Array<Accident>>()
combineLatest([refToObservable(branchConfigurator.selected), refToObservable(timeRangeConfigurator.value), ]).subscribe(data => {
  getWarningFromMetaDb(warnings, data[0], projects, dbName+"_"+dbTable, data[1] as TimeRange)
})

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