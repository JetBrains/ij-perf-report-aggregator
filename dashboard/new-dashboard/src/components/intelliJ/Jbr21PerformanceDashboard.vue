<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="jbr_21_idea_dashboard"
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
    projects: ["21jbr-community/indexing", "community/indexing"],
  },
  {
    labels: ["Compilation"],
    measures: ["build_compilation_duration"],
    projects: ["21jbr-community/rebuild", "community/rebuild"],
  },
  {
    labels: ["Find Usages Library.getName() before compilation", "First Code Analysis"],
    measures: ["findUsages", "firstCodeAnalysis"],
    projects: ["21jbr-community/findUsages/Library_getName_Before", "community/findUsages/Library_getName_Before"],
  },
  {
    labels: ["Find Usages Library.getName() after compilation"],
    measures: ["findUsages"],
    projects: ["21jbr-community/findUsages/Library_getName_After", "community/findUsages/Library_getName_After"],
  },
  {
    labels: ["Find Usages PsiManager.getInstance() before compilation", "First Code Analysis"],
    measures: ["findUsages", "firstCodeAnalysis"],
    projects: ["21jbr-community/findUsages/PsiManager_getInstance_Before", "community/findUsages/PsiManager_getInstance_Before"],
  },
  {
    labels: ["Find Usages PsiManager.getInstance() after compilation"],
    measures: ["findUsages"],
    projects: ["21jbr-community/findUsages/PsiManager_getInstance_After", "community/findUsages/PsiManager_getInstance_After"],
  },
  {
    labels: ["SE: go-to-action Editor"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-action/Editor/typingLetterByLetter", "community/go-to-action/Editor/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-action Kotlin"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-action/Kotlin/typingLetterByLetter", "community/go-to-action/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-action Runtime"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-action/Runtime/typingLetterByLetter", "community/go-to-action/Runtime/typingLetterByLetter"],
  },

  {
    labels: ["SE: go-to-all Editor"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-all/Editor/typingLetterByLetter", "community/go-to-all/Editor/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-all Kotlin"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-all/Kotlin/typingLetterByLetter", "community/go-to-all/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-all Runtime"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-all/Runtime/typingLetterByLetter", "community/go-to-all/Runtime/typingLetterByLetter"],
  },

  {
    labels: ["SE: go-to-class Editor"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-class/Editor/typingLetterByLetter", "community/go-to-class/Editor/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-class Kotlin"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-class/Kotlin/typingLetterByLetter", "community/go-to-class/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-class Runtime"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-class/Runtime/typingLetterByLetter", "community/go-to-class/Runtime/typingLetterByLetter"],
  },

  {
    labels: ["SE: go-to-file Editor"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-file/Editor/typingLetterByLetter", "community/go-to-file/Editor/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-file Kotlin"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-file/Kotlin/typingLetterByLetter", "community/go-to-file/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-file Runtime"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-file/Runtime/typingLetterByLetter", "community/go-to-file/Runtime/typingLetterByLetter"],
  },

  {
    labels: ["SE: go-to-symbol Editor"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-symbol/Editor/typingLetterByLetter", "community/go-to-symbol/Editor/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-symbol Kotlin"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-symbol/Kotlin/typingLetterByLetter", "community/go-to-symbol/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-symbol Runtime"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-symbol/Runtime/typingLetterByLetter", "community/go-to-symbol/Runtime/typingLetterByLetter"],
  },

  {
    labels: ["SE: go-to-text Editor"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-text/Editor/typingLetterByLetter", "community/go-to-text/Editor/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-text Kotlin"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-text/Kotlin/typingLetterByLetter", "community/go-to-text/Kotlin/typingLetterByLetter"],
  },
  {
    labels: ["SE: go-to-text Runtime"],
    measures: ["searchEverywhere"],
    projects: ["21jbr-community/go-to-text/Runtime/typingLetterByLetter", "community/go-to-text/Runtime/typingLetterByLetter"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
