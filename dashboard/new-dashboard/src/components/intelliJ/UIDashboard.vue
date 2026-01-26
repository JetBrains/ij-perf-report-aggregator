<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_java_dashboard_devserver"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        v-for="chart in charts"
        :key="chart.definition.label"
        :label="chart.definition.label"
        :measure="chart.definition.measure"
        :projects="chart.projects"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
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
  {
    labels: ["EditorContextMenu"],
    measures: ["popupShown#EditorContextMenu"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["ProjectViewContextMenu"],
    measures: ["popupShown#ProjectViewContextMenu"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["ProjectWidget"],
    measures: ["popupShown#ProjectWidget"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["RunConfigurations"],
    measures: ["popupShown#RunConfigurations"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogBranchFilter"],
    measures: ["popupShown#VcsLogBranchFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogDateFilter"],
    measures: ["popupShown#VcsLogDateFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogPathFilter"],
    measures: ["popupShown#VcsLogPathFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsLogUserFilter"],
    measures: ["popupShown#VcsLogUserFilter"],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["VcsWidget"],
    measures: [["popupShown#VcsWidget", "afterShow#GitBranchesTreePopup", "afterShow#GitDefaultBranchesPopup"]],
    projects: ["popups-performance-test/test-popups"],
  },
  {
    labels: ["FileStructurePopup"],
    measures: ["FileStructurePopup"],
    projects: ["intellij_commit/FileStructureDialog/java_file", "intellij_commit/FileStructureDialog/kotlin_file"],
  },
  {
    labels: ["ShowQuickFixes", "Show Intentions (average awt delay)", "Show Intentions (awt dispatch time)"],
    measures: ["showQuickFixes", "test#average_awt_delay", "AWTEventQueue.dispatchTimeTotal"],
    projects: ["grails/showIntentions/Find cause", "kotlin/showIntention/OldDarkScheme/Import", "kotlin/showIntention/Import", "spring_boot/showIntentions"],
  },
  {
    labels: ["FindUsagesPopup"],
    measures: ["findUsage_popup"],
    projects: ["intellij_commit/findUsages/PsiManager_getInstance_firstCall", "intellij_commit/findUsages/Library_getName"],
  },
  {
    labels: ["SearchEverywherePopup"],
    measures: ["searchEverywhere_dialog_shown"],
    projects: [
      "community/go-to-all/Runtime/typingLetterByLetter",
      "community/go-to-all-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-all/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Cold Search Everywhere Action (slow typing)", "Cold SE First Element Added Action (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-action/Kotlin/typingLetterByLetter",
      "community/go-to-action/Editor/typingLetterByLetter",
      "community/go-to-action/Runtime/typingLetterByLetter",
      "community/go-to-action-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-action/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-action/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-action/Editor/typingLetterByLetter",
      "intellij_commit/go-to-action/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-action-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-action/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-action/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-action/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Warm Search Everywhere Action (slow typing)", "Warm SE First Element Added Action (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-action-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-action-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-action-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-action-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-action-with-warmup/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-action-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-action-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-action-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-action-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-action-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-action-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-action-with-warmup/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Cold Search Everywhere Class (slow typing)", "Cold SE First Element Added Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-class/Kotlin/typingLetterByLetter",
      "community/go-to-class/Editor/typingLetterByLetter",
      "community/go-to-class/Runtime/typingLetterByLetter",
      "community/go-to-class-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-class/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-class/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-class/Editor/typingLetterByLetter",
      "intellij_commit/go-to-class/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-class-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-class/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-class/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-class/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Warm Search Everywhere Class (slow typing)", "Warm SE First Element Added Class (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-class-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-class-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-class-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-class-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-class-with-warmup/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-class-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-class-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-class-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-class-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-class-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-class-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-class-with-warmup/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Cold Search Everywhere File (slow typing)", "Cold SE First Element Added File(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-file/Editor/typingLetterByLetter",
      "community/go-to-file/Kotlin/typingLetterByLetter",
      "community/go-to-file/Runtime/typingLetterByLetter",
      "community/go-to-file-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-file/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-file/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-file/Editor/typingLetterByLetter",
      "intellij_commit/go-to-file/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-file-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-file/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-file/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-file/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Warm Search Everywhere File (slow typing)", "Warm SE First Element Added File(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-file-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-file-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-file-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-file-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-file-with-warmup/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-file-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-file-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-file-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-file-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-file-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-file-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-file-with-warmup/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Cold Search Everywhere All (slow typing)", "Cold SE First Element Added All(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-all/Editor/typingLetterByLetter",
      "community/go-to-all/Kotlin/typingLetterByLetter",
      "community/go-to-all/Runtime/typingLetterByLetter",
      "community/go-to-all-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-all/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-all/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-all/Editor/typingLetterByLetter",
      "intellij_commit/go-to-all/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-all-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-all/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-all/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-all/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Warm Search Everywhere All (slow typing)", "Warm SE First Element Added All(slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-all-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-all-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-all-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-all-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-all-with-warmup/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-all-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-all-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-all-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-all-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-all-with-warmup/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Cold Search Everywhere Symbol (slow typing)", "Cold SE First Element Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-symbol/Editor/typingLetterByLetter",
      "community/go-to-symbol/Kotlin/typingLetterByLetter",
      "community/go-to-symbol/Runtime/typingLetterByLetter",
      "community/go-to-symbol-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-symbol/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-symbol/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-symbol/Editor/typingLetterByLetter",
      "intellij_commit/go-to-symbol/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-symbol-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-symbol/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-symbol/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-symbol/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Warm Search Everywhere Symbol (slow typing)", "Warm SE First Element Symbol (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-symbol-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-symbol-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-symbol-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-symbol-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-symbol-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-symbol-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-symbol-with-warmup/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Cold Search Everywhere Text (slow typing)", "Cold SE First Element Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-text/Editor/typingLetterByLetter",
      "community/go-to-text/Kotlin/typingLetterByLetter",
      "community/go-to-text/Runtime/typingLetterByLetter",
      "community/go-to-text-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-text/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-text/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-text/Editor/typingLetterByLetter",
      "intellij_commit/go-to-text/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-text-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-text/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-text/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-text/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Warm Search Everywhere Text (slow typing)", "Warm SE First Element Text (slow typing)"],
    measures: ["searchEverywhere", "searchEverywhere_first_elements_added"],
    projects: [
      "community/go-to-text-with-warmup/Editor/typingLetterByLetter",
      "community/go-to-text-with-warmup/Kotlin/typingLetterByLetter",
      "community/go-to-text-with-warmup/Runtime/typingLetterByLetter",
      "community/go-to-text-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",
      "java/go-to-text-with-warmup/Runtime/typingLetterByLetter",

      "intellij_commit/go-to-text-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/go-to-text-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/go-to-text-with-warmup/Runtime/typingLetterByLetter",
      "intellij_commit/go-to-text-with-warmup-finished-embeddings/Runtime/typingLetterByLetter",

      "intellij_commit/new-se-go-to-text-with-warmup/Kotlin/typingLetterByLetter",
      "intellij_commit/new-se-go-to-text-with-warmup/Editor/typingLetterByLetter",
      "intellij_commit/new-se-go-to-text-with-warmup/Runtime/typingLetterByLetter",
    ],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
