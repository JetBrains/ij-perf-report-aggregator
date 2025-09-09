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
      <GroupProjectsWithClientChart
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
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsWithClientChart from "../charts/GroupProjectsWithClientChart.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["VFS Refresh"],
    measures: ["vfs_initial_refresh"],
    projects: ["intellij_commit/vfsRefresh/default", "intellij_commit/vfsRefresh/with-1-thread(s)", "intellij_commit/vfsRefresh/git-status"],
    aliases: ["default", "1 thread", "git status"],
  },
  {
    labels: ["VFS Refresh after Mass Changes"],
    measures: [["vfsRefreshAfterMassCreate", "vfsRefreshAfterMassModify", "vfsRefreshAfterMassDelete"]],
    projects: ["empty_project/vfs-mass-update-txt", "empty_project/vfs-mass-update-java", "empty_project/vfs-mass-update-kt"],
    aliases: ["txt", "java", "kotlin"],
  },
  {
    labels: ["Inspection"],
    measures: ["globalInspections"],
    projects: ["kotlin/inspection", "kotlin_coroutines/inspection", "intellij_commit/jvm-inspection"],
  },
  {
    labels: ["JVM Total Time to safepoints", "Full GC Pause", "JVM GC collection times"],
    measures: ["JVM.totalTimeToSafepointsMs", "fullGCPause", "JVM.GC.collectionTimesMs"],
    projects: ["kotlin/inspection", "kotlin_coroutines/inspection"],
  },
  {
    labels: ["Local Inspection", "First Code Analysis"],
    measures: ["localInspections", "firstCodeAnalysis"],
    projects: [
      "intellij_commit/localInspection/java_file",
      "kotlin/localInspection",
      "kotlin_coroutines/localInspection",
      "intellij_commit/localInspection/java_file-ContentManagerImpl",
    ],
  },
  {
    labels: ["Highlighting - remove symbol", "Highlighting - remove symbol warmup"],
    measures: ["typing_EditorBackSpace_duration", "typing_EditorBackSpace_warmup_duration"],
    projects: ["intellij_commit/editor-highlighting", "intellij_commit/editor-kotlin-highlighting"],
  },
  {
    labels: ["Highlighting - type symbol", "Highlighting - type symbol warmup"],
    measures: ["typing_}_duration", "typing_}_warmup_duration"],
    projects: ["intellij_commit/editor-highlighting", "intellij_commit/editor-kotlin-highlighting"],
  },
  {
    labels: ["Highlighting - remove method"],
    measures: ["replaceTextCodeAnalysis"],
    projects: ["intellij_commit/red-code-kotlin"],
  },
  {
    labels: ["Debug run configuration", "Debug step into"],
    measures: ["debugRunConfiguration", "debugStep_into"],
    projects: ["kotlin_petclinic/debug"],
  },
  {
    labels: ["Show Intentions (average awt delay)", "Show Intentions (awt dispatch time)"],
    measures: ["test#average_awt_delay", "AWTEventQueue.dispatchTimeTotal"],
    projects: ["kotlin/showIntention/Import"],
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
    labels: ["Creating a new file"],
    measures: ["createKotlinFile"],
    projects: ["intellij_commit/createKotlinClass"],
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
    labels: ["Find in Files"],
    measures: [["findInFiles#openDialog", "findInFiles#search: newInstance", "findInFiles#search: intellij-ide-starter"]],
    projects: ["intellij_commit/find-in-files", "intellij_commit/find-in-files-old"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
