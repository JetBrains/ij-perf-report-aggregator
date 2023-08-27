<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="fus_dashboard"
    initial-machine="Linux EC2 C6id.8xlarge (32 vCPU Xeon, 64 GB)"
    :charts="charts"
  >
    <section>
      <div>
        <GroupProjectsChart
          v-for="chart in charts"
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

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Indexing time", "Scanning time"],
    measures: ["fus_dumb_indexing_time", "fus_scanning_time"],
    projects: ["intellij_sources/indexing"],
  },
  {
    labels: ["Completion time to show 90p", "Completion duration 90p"],
    measures: ["fus_time_to_show_90p", "fus_completion_duration_90p"],
    projects: ["grails/completion/groovy_file", "grails/completion/java_file", "keycloak_release_20/completion/QuarkusRuntimePomXml"],
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
