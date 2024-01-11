<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_new_gradle_dashboard"
    initial-machine="linux-blade-hetzner"
    :charts="charts"
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
import { GRADLE_PROJECTS } from "./gradle-projects"

const metricsDeclaration = [
  // total sync time
  "ExternalSystemSyncProjectTask",
  // full time of the sink operation, with all our overhead for preparation
  "GradleExecution",
  // work inside Gradle connection, operations that are performed inside connection
  "GradleConnection",
  // resolving models from daemon
  "GradleCall",
  // processing the data we received from Gradle
  "ExternalSystemSyncResultProcessing",
  // work of data services
  "ProjectDataServices",
  // project resolve
  "GradleProjectResolvers",
  // apply ws model
  "WorkspaceModelApply",
]

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(GRADLE_PROJECTS)

const charts = computed(() => {
  const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
    return {
      labels: [metric],
      measures: [metric],
      projects: testConfigurator.selected.value ?? [],
    }
  })
  return combineCharts(chartsDeclaration)
})
</script>
