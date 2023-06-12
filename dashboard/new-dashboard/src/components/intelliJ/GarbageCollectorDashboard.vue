<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_gc_dashboard"
    initial-machine="Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
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
    labels: ["Finding Usages", "First Code Analysis", "Completion", "Dispatch total time"],
    measures: ["findUsages", "firstCodeAnalysis", "completion", "AWTEventQueue.dispatchTimeTotal"],
    projects: ["userScenario_defaultGC/userScenario", "userScenario_ZGC/userScenario", "userScenario_ZGC_ZUncommit/userScenario", "userScenario_ZGC_transparentHugePages/userScenario", "userScenario_ZGC_largePages/userScenario"]
  }
]

const charts = combineCharts(chartsDeclaration)
</script>
