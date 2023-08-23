<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="fus_dev_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="chartsDev"
    :with-installer="false"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="chart in chartsDev"
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

const chartsDevDeclaration: ChartDefinition[] = [
  {
    labels: ["Completion time to show 90p", "Completion duration 90p"],
    measures: ["fus_time_to_show_90p", "fus_completion_duration_90p"],
    projects: ["intellij_commit/completion/java_file", "intellij_commit/completion/kotlin_file"],
  },
  {
    labels: ["FindUsages first usage", "FindUsages all usages"],
    measures: ["fus_find_usages_first", "fus_find_usages_all"],
    projects: ["intellij_commit/findUsages/PsiManager_getInstance"],
  },
]

const chartsDev = combineCharts(chartsDevDeclaration)
</script>
