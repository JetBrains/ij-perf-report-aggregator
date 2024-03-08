<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="idea"
    persistent-id="ideaDev_performance_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :is-like="true"
          :title="'Completion'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'searchEverywhere\_%'"
          :is-like="true"
          :title="'Search Everywhere'"
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
        :aliases="chart.aliases"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["VFS Refresh"],
    measures: ["vfs_initial_refresh"],
    projects: ["intellij_commit/vfsRefresh/default", "intellij_commit/vfsRefresh/with-1-thread(s)", "intellij_commit/vfsRefresh/git-status"],
    aliases: ["default", "1 thread", "git status"],
  },
  {
    labels: ["Rebuild (Big projects)"],
    measures: ["build_compilation_duration"],
    projects: ["community/rebuild"],
  },
  {
    labels: ["Rebuild"],
    measures: ["build_compilation_duration"],
    projects: ["grails/rebuild", "java/rebuild", "spring_boot/rebuild"],
  },
  {
    labels: ["Inspection"],
    measures: ["globalInspections"],
    projects: ["java/inspection", "grails/inspection", "spring_boot_maven/inspection", "spring_boot/inspection", "kotlin/inspection", "kotlin_coroutines/inspection"],
  },
  {
    labels: ["Local Inspection", "First Code Analysis"],
    measures: ["localInspections", "firstCodeAnalysis"],
    projects: [
      "intellij_commit/localInspection/java_file",
      "intellij_commit/localInspection/kotlin_file",
      "kotlin/localInspection",
      "kotlin_coroutines/localInspection",
      "intellij_commit/localInspection/java_file-ContentManagerImpl",
      "intellij_commit/localInspection/kotlin_file-DexInlineTest",
    ],
  },
  {
    labels: ["Completion"],
    measures: ["completion"],
    projects: [
      "community/completion/kotlin_file",
      "grails/completion/groovy_file",
      "grails/completion/java_file",
      "intellij_commit/completion/kotlin_file",
      "intellij_commit/completion/java_file",
    ],
  },
  {
    labels: ["Debug run configuration", "Debug step into"],
    measures: ["debugRunConfiguration", "debugStep_into"],
    projects: ["kotlin_petclinic/debug"],
  },
  {
    labels: ["Show Intentions (average awt delay)", "Show Intentions (showQuickFixes)", "Show Intentions (awt dispatch time)"],
    measures: ["test#average_awt_delay", "showQuickFixes", "AWTEventQueue.dispatchTimeTotal"],
    projects: ["grails/showIntentions/Find cause", "kotlin/showIntention/Import", "spring_boot/showIntentions"],
  },
  {
    labels: ["Show File History"],
    measures: ["showFileHistory"],
    projects: ["intellij_commit/showFileHistory/EditorImpl"],
  },
  {
    labels: ["Expand Project Menu"],
    measures: ["%expandProjectMenu"],
    projects: ["intellij_commit/expandProjectMenu"],
  },
  {
    labels: ["Expand Main Menu"],
    measures: ["%expandMainMenu"],
    projects: ["intellij_commit/expandMainMenu"],
  },
  {
    labels: ["Expand Editor Menu"],
    measures: ["%expandEditorMenu"],
    projects: ["intellij_commit/expandEditorMenu"],
  },
  {
    labels: ["Highlight"],
    measures: ["highlighting"],
    projects: ["kotlin/highlight", "kotlin_coroutines/highlight"],
  },
  {
    labels: ["FileStructure"],
    measures: ["FileStructurePopup"],
    projects: ["intellij_commit/FileStructureDialog/java_file", "intellij_commit/FileStructureDialog/kotlin_file"],
  },
  {
    labels: ["Creating a new file"],
    measures: [["createJavaFile", "createKotlinFile"]],
    projects: ["intellij_commit/createJavaClass", "intellij_commit/createKotlinClass"],
  },
  {
    labels: ["Search Everywhere Class (slow typing)", "SE Dialog Shown Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-class/Kotlin/typingLetterByLetter", "community/go-to-class/Editor/typingLetterByLetter", "community/go-to-class/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere File (slow typing)", "SE Dialog Shown File (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-file/Editor/typingLetterByLetter", "community/go-to-file/Kotlin/typingLetterByLetter", "community/go-to-file/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere All (slow typing)", "SE Dialog Shown All (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-all/Editor/typingLetterByLetter", "community/go-to-all/Kotlin/typingLetterByLetter", "community/go-to-all/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere Symbol (slow typing)", "SE Dialog Shown Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-symbol/Editor/typingLetterByLetter", "community/go-to-symbol/Kotlin/typingLetterByLetter", "community/go-to-symbol/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["Search Everywhere Text (slow typing)", "SE Dialog Shown Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_dialog_shown"],
    projects: ["community/go-to-text/Editor/typingLetterByLetter", "community/go-to-text/Kotlin/typingLetterByLetter", "community/go-to-text/Runtime/typingLetterByLetter"],
  },
  {
    labels: ["IsUpToDateCheck duration"],
    measures: ["isUpToDateCheck"],
    projects: ["community/findUsages/PsiManager_getInstance_Before", "community/findUsages/PsiManager_getInstance_After"],
  },
  {
    labels: ["FindUsages PsiManager#getInstance Before and After Compilation"],
    measures: ["findUsages"],
    projects: ["community/findUsages/PsiManager_getInstance_Before", "community/findUsages/PsiManager_getInstance_After"],
  },
  {
    labels: ["FindUsages Library#getName (all usages)", "FindUsages Library#getName (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Library_getName_Before", "community/findUsages/Library_getName_After", "intellij_commit/findUsages/Library_getName"],
  },
  {
    labels: ["FindUsages LocalInspectionTool#getID Before and After Compilation"],
    measures: ["findUsages"],
    projects: ["community/findUsages/LocalInspectionTool_getID_Before", "community/findUsages/LocalInspectionTool_getID_After"],
  },
  {
    labels: ["FindUsages ActionsKt#runReadAction and Application#runReadAction Before and After Compilation"],
    measures: ["findUsages"],
    projects: [
      "community/findUsages/ActionsKt_runReadAction_Before",
      "community/findUsages/ActionsKt_runReadAction_After",
      "community/findUsages/Application_runReadAction_Before",
      "community/findUsages/Application_runReadAction_After",
      "intellij_commit/findUsages/ActionsKt_runReadAction",
      "intellij_commit/findUsages/Application_runReadAction",
    ],
  },
  {
    labels: ["FindUsages Persistent#absolutePath and PropertyMapping#value Before and After Compilation"],
    measures: ["findUsages"],
    projects: [
      "community/findUsages/Persistent_absolutePath_Before",
      "community/findUsages/Persistent_absolutePath_After",
      "community/findUsages/PropertyMapping_value_Before",
      "community/findUsages/PropertyMapping_value_After",
      "intellij_commit/findUsages/Persistent_absolutePath",
    ],
  },
  {
    labels: ["FindUsages Object#hashCode and Path#toString Before and After Compilation"],
    measures: ["findUsages"],
    projects: [
      "community/findUsages/Object_hashCode_Before",
      "community/findUsages/Object_hashCode_After",
      "community/findUsages/Path_toString_Before",
      "community/findUsages/Path_toString_After",
    ],
  },
  {
    labels: ["FindUsages Objects#hashCode Before and After Compilation", "FindUsages Objects#hashCode Before and After Compilation (first usage)"],
    measures: [
      ["findUsages", "fus_find_usages_all"],
      ["findUsages_firstUsage", "fus_find_usages_first"],
    ],
    projects: ["community/findUsages/Objects_hashCode_Before", "community/findUsages/Objects_hashCode_After"],
  },
  {
    labels: ["FindUsages Path#div Before and After Compilation"],
    measures: ["findUsages"],
    projects: ["community/findUsages/Path_div_Before", "community/findUsages/Path_div_After", "intellij_commit/findUsages/Path_div"],
  },
  {
    labels: ["Find Usages with idea.is.internal=true Before Compilation"],
    measures: ["findUsages"],
    projects: [
      "intellij_commit/findUsages/PsiManager_getInstance_firstCall",
      "intellij_commit/findUsages/PsiManager_getInstance_secondCall",
      "intellij_commit/findUsages/PsiManager_getInstance_thirdCallInternalMode",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
