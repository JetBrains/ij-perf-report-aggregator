<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_workspace_model_benchmarks_dashboard"
    initial-machine="linux-blade-hetzner"
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
import { ChartDefinition, combineCharts } from "../charts/DashboardCharts"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const metricsDeclaration = [
  "Deserialization",
  "DeserializationFromFile",
  "Serialization",
  "SerializationFromFile",
  "Serialization Size",
  "Named entities adding (2000000)",
  "Soft linked entities adding (2000000)",
  "duration",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "serialize-community-project",
      "requesting-same-entity",
      "renaming-named-entities",
      "refers-named-entities",
      "rbs-new-on-many-content-roots",
      "adding-storage-recreating",
      "adding-soft-linked-entities",
      "10-000-orphan-content-roots-to-modules",
      "10-000-orphan-source-roots-to-many-content-roots-to-modules",
      "10-000-orphan-source-roots-to-modules"
    ],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
