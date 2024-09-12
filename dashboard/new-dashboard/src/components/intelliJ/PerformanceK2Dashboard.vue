<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_performance-k2_dashboard_dev"
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
    labels: ["FindUsages PsiManager_getInstance", "Number FindUsages PsiManager_getInstance"],
    measures: ["findUsages", "findUsages#number"],
    projects: [
      "intellij_commit-changedDefault/findUsages/PsiManager_getInstance_firstCall",
      "intellij_commit/findUsages/PsiManager_getInstance_firstCall",
      "intellij_commit-k1/findUsages/PsiManager_getInstance_firstCall",
      "intellij_commit-k2/findUsages/PsiManager_getInstance_firstCall",
    ],
    aliases: ["findUsages-getInstance-k2", "findUsages-getInstance-default", "findUsages-getInstance-k1", "findUsages-getInstance-k2"],
  },
  {
    labels: ["FindUsages String_toString", "Number FindUsages String_toString"],
    measures: ["findUsages", "findUsages#number"],
    projects: [
      "intellij_commit-changedDefault/findUsages/String_toString",
      "intellij_commit/findUsages/String_toString",
      "intellij_commit-k1/findUsages/String_toString",
      "intellij_commit-k2/findUsages/String_toString",
    ],
    aliases: ["findUsages-toString-k2", "findUsages-toString-default", "findUsages-toString-k1", "findUsages-toString-k2"],
  },
  {
    labels: ["Local inspections .kt Kotlin Serialization"],
    measures: ["localInspections"],
    projects: ["kotlin-changedDefault/localInspection", "kotlin/localInspection", "kotlin-k1/localInspection", "kotlin-k2/localInspection"],
    aliases: ["localInspections-k2", "localInspections-default", "localInspections-k1", "localInspections-k2"],
  },
  {
    labels: ["Completion .java IntelliJ", "Completion .java 90p IntelliJ"],
    measures: ["completion", "fus_completion_duration_90p"],
    projects: [
      "intellij_commit-changedDefault/completion/java_file",
      "intellij_commit/completion/java_file",
      "intellij_commit-k1/completion/java_file",
      "intellij_commit-k2/completion/java_file",
    ],
    aliases: ["intellij-completion-k2", "intellij-completion-default", "intellij-completion-k1", "intellij-completion-k2"],
  },
  {
    labels: ["Completion .kt Toolbox", "Completion .kt 90p Toolbox"],
    measures: ["completion", "fus_completion_duration_90p"],
    projects: [
      "toolbox_enterprise-changedDefault/ultimateCase/UserRepository",
      "toolbox_enterprise/ultimateCase/UserRepository",
      "toolbox_enterprise-k1/ultimateCase/UserRepository",
      "toolbox_enterprise-k2/ultimateCase/UserRepository",
    ],
    aliases: ["toolbox-completion-k2", "toolbox-completion-default", "toolbox-completion-k1", "toolbox-completion-k2"],
  },
  {
    labels: ["Search Everywhere Go to All"],
    measures: ["searchEverywhere"],
    projects: [
      "community-changedDefault/go-to-all/Editor/typingLetterByLetter",
      "community/go-to-all/Editor/typingLetterByLetter",
      "community-k1/go-to-all/Editor/typingLetterByLetter",
      "community-k2/go-to-all/Editor/typingLetterByLetter",
    ],
    aliases: ["SE-go-to-all-k2", "SE-go-to-all-default", "SE-go-to-all-k1", "SE-go-to-all-k2"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
