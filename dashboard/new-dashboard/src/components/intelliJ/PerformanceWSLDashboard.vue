<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="ideaDev_wsl_performance_dashboard"
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
    labels: ["Indexing on WSL"],
    measures: ["indexingTimeWithoutPauses"],
    projects: ["wsl-java/indexing", "wsl-spring-pet-clinic-gradle/indexing", "wsl-spring-pet-clinic-maven/indexing", "wsl-empty_project/indexing", "wsl-grails/indexing"],
  },
  {
    labels: ["Indexing on WSL"],
    measures: ["scanningTimeWithoutPauses"],
    projects: ["wsl-java/indexing", "wsl-spring-pet-clinic-gradle/indexing", "wsl-spring-pet-clinic-maven/indexing", "wsl-empty_project/indexing", "wsl-grails/indexing"],
  },
  {
    labels: ["Rebuild on WSL"],
    measures: ["build_compilation_duration"],
    projects: ["wsl-empty_project/rebuild"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
