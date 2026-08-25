<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_lsp_client_lua"
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
import DashboardPage from "../common/DashboardPage.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"

const projects = ["lua_lsp_completion/completion/imported_module", "lua_lsp_completion/completion/local_variables", "lua_lsp_completion/completion/table_fields"]

const declaration: ChartDefinition[] = [
  {
    labels: ["Completion", "Completion First Element Shown", "Completion Number"],
    measures: ["completion", "completion#firstElementShown", "completion#number"],
    projects,
  },
]

const charts = combineCharts(declaration)
</script>
