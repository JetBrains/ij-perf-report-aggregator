<template>
  <DashboardPage
    db-name="perfintDev"
    table="idea"
    persistent-id="idea_maven_importers_and_configurators_dashboard_dev_server"
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
import { MAVEN_METRICS_CONFIGURATORS } from "./maven-metrics"
import { MAVEN_PROJECTS_FAST_INSTALLERS } from "./maven-projects"

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(MAVEN_PROJECTS_FAST_INSTALLERS)

const charts = computed(() => {
  const chartsDeclaration: ChartDefinition[] = MAVEN_METRICS_CONFIGURATORS.map((metric) => {
    return {
      labels: [metric],
      measures: [metric],
      projects: testConfigurator.selected.value ?? [],
    }
  })
  return combineCharts(chartsDeclaration)
})
</script>
