<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_gradle_dashboard"
    initial-machine="linux-blade-hetzner"
    :charts="charts"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="testConfigurator"
        title="Test"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4" />
        </template>
      </MeasureSelect>
    </template>
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
import { computed } from "vue"
import { SimpleMeasureConfigurator } from "../../../../configurators/SimpleMeasureConfigurator"
import { ChartDefinition, combineCharts } from "../../../charts/DashboardCharts"
import GroupProjectsChart from "../../../charts/GroupProjectsChart.vue"
import MeasureSelect from "../../../charts/MeasureSelect.vue"
import DashboardPage from "../../../common/DashboardPage.vue"
import { GRADLE_METRICS } from "./gradle-metrics"
import { GRADLE_PROJECTS } from "./gradle-projects"

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(GRADLE_PROJECTS)

const charts = computed(() => {
  const chartsDeclaration: ChartDefinition[] = GRADLE_METRICS.map((metric) => {
    return {
      labels: [metric],
      measures: [metric],
      projects: testConfigurator.selected.value ?? [],
    }
  })
  return combineCharts(chartsDeclaration)
})
</script>
