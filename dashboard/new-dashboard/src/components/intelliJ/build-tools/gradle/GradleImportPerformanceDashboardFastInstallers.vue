<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_new_gradle_dashboard"
    initial-machine="linux-blade-hetzner"
    :charts="charts"
    :with-installer="false"
  >
    <template #configurator>
      <MeasureSelect
        :configurator="testConfigurator"
        title="Test"
      >
        <template #icon>
          <ChartBarIcon class="w-4 h-4 text-gray-500" />
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
import { GRADLE_METRICS_NEW_DASHBOARD } from "./gradle-metrics"
import { GRADLE_PROJECTS_FAST_INSTALLERS } from "./gradle-projects"

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(GRADLE_PROJECTS_FAST_INSTALLERS)

const charts = computed(() => {
  const chartsDeclaration: ChartDefinition[] = GRADLE_METRICS_NEW_DASHBOARD.map((metric) => {
    return {
      labels: [metric],
      measures: [metric],
      projects: testConfigurator.selected.value ?? [],
    }
  })
  return combineCharts(chartsDeclaration)
})
</script>
