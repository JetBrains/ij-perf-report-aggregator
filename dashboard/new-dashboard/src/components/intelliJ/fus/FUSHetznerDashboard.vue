<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="fus_import_dashboard"
    initial-machine="linux-blade-hetzner"
    :charts="chartsImport"
    :with-installer="false"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="chart in chartsImport"
          :key="chart.definition.label"
          :label="chart.definition.label"
          :measure="chart.definition.measure"
          :projects="chart.projects"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { ChartDefinition, combineCharts } from "../../charts/DashboardCharts"
import GroupProjectsChart from "../../charts/GroupProjectsChart.vue"
import DashboardPage from "../../common/DashboardPage.vue"

const chartsImportDeclaration: ChartDefinition[] = [
  {
    labels: ["Maven Sync"],
    measures: ["maven.import.stats.sync.project.task"],
    projects: ["project-import-maven-flink/fastInstaller"],
  },
  {
    labels: ["Gradle Sync"],
    measures: ["fus_gradle.sync"],
    projects: ["project-import-gradle-android-extra-large/fastInstaller"],
  },
]

const chartsImport = combineCharts(chartsImportDeclaration)
</script>
