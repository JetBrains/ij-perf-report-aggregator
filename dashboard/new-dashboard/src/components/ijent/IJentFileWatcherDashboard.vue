<template>
  <DashboardPage
    db-name="perfintDev"
    table="ijent"
    persistent-id="ijent_filewatcher_dashboard"
    initial-machine="Linux Munich i7-13700, 64 Gb"
    :with-installer="false"
    :with-mode="false"
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

const ENVS = ["Local", "Docker", "WSL"] as const
const projects = ENVS.map((env) => `eel-filewatcher-ai-agent/${env}`)
const aliases = [...ENVS]

const chartsDeclaration: ChartDefinition[] = [
  {
    labels: ["Apply Patch (Create) — IDE Pickup", "Apply Patch (Update) — IDE Pickup", "Shell Batch — IDE Pickup", "IJent Events Count"],
    measures: ["agent.applyPatch.create.idePickup", "agent.applyPatch.update.idePickup", "agent.shellBatch.idePickup", "ijent.events.count"],
    projects,
    aliases,
  },
]

const charts = combineCharts(chartsDeclaration)
</script>
