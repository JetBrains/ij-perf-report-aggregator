<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_maven_importers_and_configurators_dashboard"
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
import { MAVEN_PROJECTS } from "./maven-projects"

const importerRunMetricTypes = ["total_duration_ms", "number_of_modules"]

const workspaceConfiguratorsMetricTypes = [
  ...importerRunMetricTypes,
  "after_apply_duration_ms",
  "before_apply_duration_ms",
  "collect_folders_duration_ms",
  "config_modules_duration_ms",
]

const importerRunMetrics = [
  "importer_run.com.intellij.quarkus.run.maven.QsMavenImporter",
  "importer_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator",
  "importer_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator",
  "importer_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator",
  "importer_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator",
  "importer_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator",
  "importer_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporter",
  "importer_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporterEx",
]

const workspaceConfiguratorsMetrics = [
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenShadePluginConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator",
  "workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter",
  "workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx",
  "workspace_import.configurator_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporterEx",
]

const metricsDeclaration = [
  // Importer run metrics
  ...importerRunMetrics.flatMap((metric) => importerRunMetricTypes.map((type) => `${metric}.${type}`)),
  // Workspace configuration metrics
  ...workspaceConfiguratorsMetrics.flatMap((metric) => workspaceConfiguratorsMetricTypes.map((type) => `${metric}.${type}`)),
]

const testConfigurator = new SimpleMeasureConfigurator("project", null)
testConfigurator.initData(MAVEN_PROJECTS)

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
