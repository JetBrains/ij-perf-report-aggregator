<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="custom_jbr_idea_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
  >
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
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing", "Scanning"],
    measures: ["indexingTimeWithoutPauses", "scanningTimeWithoutPauses"],
    projects: [
      "21jbr-community/indexing",
      "community/indexing",
      "21jbr-intellij_sources/indexing",
      "intellij_sources/indexing",
      "17jbr-community/indexing",
      "17jbr-intellij_sources/indexing",
    ],
  },
  {
    labels: ["VFS refresh"],
    measures: ["vfs_initial_refresh"],
    projects: ["21jbr-intellij_sources/vfsRefresh/default", "intellij_sources/vfsRefresh/default", "17jbr-intellij_sources/vfsRefresh/default"],
  },
  {
    labels: ["Compilation"],
    measures: ["build_compilation_duration"],
    projects: [
      "21jbr-community/rebuild",
      "community/rebuild",
      "21jbr-intellij_sources/rebuild",
      "intellij_sources/rebuild",
      "17jbr-community/rebuild",
      "17jbr-intellij_sources/rebuild",
    ],
  },
  {
    labels: ["Find Usages Library.getName() before compilation", "First Code Analysis"],
    measures: ["findUsages", "firstCodeAnalysis"],
    projects: ["21jbr-community/findUsages/Library_getName_Before", "community/findUsages/Library_getName_Before", "17jbr-community/findUsages/Library_getName_Before"],
  },
  {
    labels: ["Find Usages Library.getName() after compilation"],
    measures: ["findUsages"],
    projects: ["21jbr-community/findUsages/Library_getName_After", "community/findUsages/Library_getName_After", "17jbr-community/findUsages/Library_getName_After"],
  },
  {
    labels: ["Find Usages PsiManager.getInstance() before compilation", "First Code Analysis"],
    measures: ["findUsages", "firstCodeAnalysis"],
    projects: [
      "21jbr-community/findUsages/PsiManager_getInstance_Before",
      "community/findUsages/PsiManager_getInstance_Before",
      "17jbr-community/findUsages/PsiManager_getInstance_Before",
    ],
  },
  {
    labels: ["Find Usages PsiManager.getInstance() after compilation"],
    measures: ["findUsages"],
    projects: [
      "21jbr-community/findUsages/PsiManager_getInstance_After",
      "community/findUsages/PsiManager_getInstance_After",
      "17jbr-community/findUsages/PsiManager_getInstance_After",
    ],
  },
  {
    labels: ["SE: go-to-action Editor"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-action/Editor/typingLetterByLetter",
      "community/go-to-action/Editor/typingLetterByLetter",
      "17jbr-community/go-to-action/Editor/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-action Kotlin"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-action/Kotlin/typingLetterByLetter",
      "community/go-to-action/Kotlin/typingLetterByLetter",
      "17jbr-community/go-to-action/Kotlin/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-action Runtime"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-action/Runtime/typingLetterByLetter",
      "community/go-to-action/Runtime/typingLetterByLetter",
      "17jbr-community/go-to-action/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-all Editor"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-all/Editor/typingLetterByLetter", "community/go-to-all/Editor/typingLetterByLetter", "17jbr-community/go-to-all/Editor/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-all Kotlin"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-all/Kotlin/typingLetterByLetter", "community/go-to-all/Kotlin/typingLetterByLetter", "17jbr-community/go-to-all/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-all Runtime"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-all/Runtime/typingLetterByLetter",
      "community/go-to-all/Runtime/typingLetterByLetter",
      "17jbr-community/go-to-all/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-class Editor"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-class/Editor/typingLetterByLetter",
      "community/go-to-class/Editor/typingLetterByLetter",
      "17jbr-community/go-to-class/Editor/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-class Kotlin"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-class/Kotlin/typingLetterByLetter",
      "community/go-to-class/Kotlin/typingLetterByLetter",
      "17jbr-community/go-to-class/Kotlin/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-class Runtime"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-class/Runtime/typingLetterByLetter",
      "community/go-to-class/Runtime/typingLetterByLetter",
      "17jbr-community/go-to-class/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-file Editor"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-file/Editor/typingLetterByLetter",
      "community/go-to-file/Editor/typingLetterByLetter",
      "17jbr-community/go-to-file/Editor/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-file Kotlin"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-file/Kotlin/typingLetterByLetter",
      "community/go-to-file/Kotlin/typingLetterByLetter",
      "17jbr-community/go-to-file/Kotlin/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-file Runtime"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-file/Runtime/typingLetterByLetter",
      "community/go-to-file/Runtime/typingLetterByLetter",
      "17jbr-community/go-to-file/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-symbol Editor"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-symbol/Editor/typingLetterByLetter",
      "community/go-to-symbol/Editor/typingLetterByLetter",
      "17jbr-community/go-to-symbol/Editor/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-symbol Kotlin"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-symbol/Kotlin/typingLetterByLetter",
      "community/go-to-symbol/Kotlin/typingLetterByLetter",
      "17jbr-community/go-to-symbol/Kotlin/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-symbol Runtime"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-symbol/Runtime/typingLetterByLetter",
      "community/go-to-symbol/Runtime/typingLetterByLetter",
      "17jbr-community/go-to-symbol/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-text Editor"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-text/Editor/typingLetterByLetter",
      "community/go-to-text/Editor/typingLetterByLetter",
      "17jbr-community/go-to-text/Editor/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-text Kotlin"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-text/Kotlin/typingLetterByLetter",
      "community/go-to-text/Kotlin/typingLetterByLetter",
      "17jbr-community/go-to-text/Kotlin/typingLetterByLetter",
    ],
  },
  {
    labels: ["SE: go-to-text Runtime"],
    measures: ["searchEverywhere"],
    projects: [
      "21jbr-community/go-to-text/Runtime/typingLetterByLetter",
      "community/go-to-text/Runtime/typingLetterByLetter",
      "17jbr-community/go-to-text/Runtime/typingLetterByLetter",
    ],
  },
  {
    labels: ["Local inspections java file"],
    measures: ["localInspections"],
    projects: ["21jbr-intellij_sources/localInspection/java_file", "intellij_sources/localInspection/java_file", "17jbr-intellij_sources/localInspection/java_file"],
  },
  {
    labels: ["Local inspections kotlin file"],
    measures: ["localInspections"],
    projects: ["21jbr-intellij_sources/localInspection/kotlin_file", "intellij_sources/localInspection/kotlin_file", "17jbr-intellij_sources/localInspection/kotlin_file"],
  },
  {
    labels: ["File History"],
    measures: ["showFileHistory"],
    projects: ["21jbr-intellij_sources/showFileHistory/EditorImpl", "intellij_sources/showFileHistory/EditorImpl", "17jbr-intellij_sources/showFileHistory/EditorImpl"],
  },
  {
    labels: ["File Structure dialogue java file"],
    measures: ["FileStructurePopup"],
    projects: ["21jbr-intellij_sources/FileStructureDialog/java_file", "intellij_sources/FileStructureDialog/java_file", "17jbr-intellij_sources/FileStructureDialog/java_file"],
  },
  {
    labels: ["File Structure dialogue kotlin file"],
    measures: ["FileStructurePopup"],
    projects: [
      "21jbr-intellij_sources/FileStructureDialog/kotlin_file",
      "intellij_sources/FileStructureDialog/kotlin_file",
      "17jbr-intellij_sources/FileStructureDialog/kotlin_file",
    ],
  },
  {
    labels: ["Expand Editor menu"],
    measures: ["%expandEditorMenu"],
    projects: ["21jbr-intellij_sources/expandEditorMenu", "intellij_sources/expandEditorMenu", "17jbr-intellij_sources/expandEditorMenu"],
  },
  {
    labels: ["Expand Main menu"],
    measures: ["%expandMainMenu"],
    projects: ["21jbr-intellij_sources/expandMainMenu", "intellij_sources/expandMainMenu", "17jbr-intellij_sources/expandMainMenu"],
  },
  {
    labels: ["Expand Project menu"],
    measures: ["%expandProjectMenu"],
    projects: ["21jbr-intellij_sources/expandProjectMenu", "intellij_sources/expandProjectMenu", "17jbr-intellij_sources/expandProjectMenu"],
  },
  {
    labels: ["Create new java class"],
    measures: ["createJavaFile"],
    projects: ["21jbr-intellij_sources/createJavaClass", "intellij_sources/createJavaClass", "17jbr-intellij_sources/createJavaClass"],
  },
  {
    labels: ["Create new kotlin class"],
    measures: ["createKotlinFile"],
    projects: ["21jbr-intellij_sources/createKotlinClass", "intellij_sources/createKotlinClass", "17jbr-intellij_sources/createKotlinClass"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
