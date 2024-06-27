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
    labels: ["Inspection", "JVM Total Time to safepoints", "Full GC Pause", "JVM GC collection times"],
    measures: ["globalInspections", "JVM.totalTimeToSafepointsMs", "fullGCPause", "JVM.GC.collectionTimesMs"],
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
    labels: ["Editor Scrolling AWT Delay"],
    measures: [["scrollEditor#max_awt_delay", "scrollEditor#average_awt_delay"]],
    projects: ["intellij_commit/scrollEditor/java_file_ContentManagerImpl"],
  },
  {
    labels: ["Editor Scrolling CPU Load"],
    measures: [["scrollEditor#max_cpu_load", "scrollEditor#average_cpu_load"]],
    projects: ["intellij_commit/scrollEditor/java_file_ContentManagerImpl"],
  },
  {
    labels: ["Project View"],
    measures: [["projectViewInit", "projectViewInit#cachedNodesLoaded"]],
    projects: ["intellij_commit/projectView"],
  },
  {
    labels: ["Popups"],
    measures: [
      [
        "popupShown#EditorContextMenu",
        "popupShown#ProjectViewContextMenu",
        "popupShown#ProjectWidget",
        "popupShown#RunConfigurations",
        "popupShown#VcsLogBranchFilter",
        "popupShown#VcsLogDateFilter",
        "popupShown#VcsLogPathFilter",
        "popupShown#VcsLogUserFilter",
        "popupShown#VcsWidget",
      ],
    ],
    projects: ["popups-performance-test/test-popups"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
