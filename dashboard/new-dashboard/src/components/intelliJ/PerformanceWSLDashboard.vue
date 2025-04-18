<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="ideaDev_wsl_performance_dashboard"
    initial-machine="Windows Munich i7-13700, 64 Gb"
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
    labels: ["Indexing"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["community/indexing", "spring-pet-clinic-gradle/indexing", "spring-pet-clinic-maven/indexing"],
  },
  {
    labels: ["Scanning"],
    measures: ["scanningTimeWithoutPauses"],
    projects: ["community/indexing", "spring-pet-clinic-gradle/indexing", "spring-pet-clinic-maven/indexing"],
  },
  {
    labels: ["Number of indexed files"],
    measures: ["numberOfIndexedFiles"],
    projects: ["community/indexing", "spring-pet-clinic-gradle/indexing", "spring-pet-clinic-maven/indexing"],
  },
  {
    labels: ["Rebuild of Community"],
    measures: ["build_compilation_duration"],
    projects: ["community/rebuild"],
  },
  {
    labels: ["Eslint indexing WSL/Windows"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["indexing-eslint/wsl", "indexing-eslint/windows"],
  },
  {
    labels: ["Eslint scanning WSL/Windows"],
    measures: ["scanningTimeWithoutPauses"],
    projects: ["indexing-eslint/wsl", "indexing-eslint/windows"],
  },
  {
    labels: ["Eslint number of indexed files WSL/Windows"],
    measures: ["numberOfIndexedFiles"],
    projects: ["indexing-eslint/wsl", "indexing-eslint/windows"],
  },
  {
    labels: ["Design Patterns indexing WSL/Windows"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["indexing-java/wsl", "indexing-java/windows"],
  },
  {
    labels: ["Design Patterns scanning WSL/Windows"],
    measures: ["scanningTimeWithoutPauses"],
    projects: ["indexing-java/wsl", "indexing-java/windows"],
  },
  {
    labels: ["Design Patterns number of indexed files WSL/Windows"],
    measures: ["numberOfIndexedFiles"],
    projects: ["indexing-java/wsl", "indexing-java/windows"],
  },
  {
    labels: ["Community indexing WSL/Windows"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["indexing-sources-intellij-community-master/wsl", "indexing-sources-intellij-community-master/windows"],
  },
  {
    labels: ["Community scanning WSL/Windows"],
    measures: ["scanningTimeWithoutPauses"],
    projects: ["indexing-sources-intellij-community-master/wsl", "indexing-sources-intellij-community-master/windows"],
  },
  {
    labels: ["Community number of indexed files WSL/Windows"],
    measures: ["numberOfIndexedFiles"],
    projects: ["indexing-sources-intellij-community-master/wsl", "indexing-sources-intellij-community-master/windows"],
  },
  {
    labels: ["Indexing on WSL"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["wsl-java/indexing", "wsl-spring-pet-clinic-gradle/indexing", "wsl-spring-pet-clinic-maven/indexing", "wsl-empty_project/indexing", "wsl-grails/indexing"],
  },
  {
    labels: ["Scanning on WSL"],
    measures: ["scanningTimeWithoutPauses"],
    projects: ["wsl-java/indexing", "wsl-spring-pet-clinic-gradle/indexing", "wsl-spring-pet-clinic-maven/indexing", "wsl-empty_project/indexing", "wsl-grails/indexing"],
  },
  {
    labels: ["Number of indexed files on WSL"],
    measures: ["numberOfIndexedFiles"],
    projects: ["wsl-java/indexing", "wsl-spring-pet-clinic-gradle/indexing", "wsl-spring-pet-clinic-maven/indexing", "wsl-empty_project/indexing", "wsl-grails/indexing"],
  },
  {
    labels: ["Rebuild on WSL"],
    measures: ["build_compilation_duration"],
    projects: ["wsl-empty_project/rebuild", "wsl-java/rebuild", "wsl-grails/rebuild"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
