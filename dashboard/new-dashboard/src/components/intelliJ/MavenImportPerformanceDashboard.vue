<template>
  <DashboardPage
    db-name="perfint"
    table="idea"
    persistent-id="idea_maven_dashboard"
    initial-machine="linux-blade-hetzner"
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

const metricsName = ["after_apply_duration_ms", "before_apply_duration_ms", "collect_folders_duration_ms", "config_modules_duration_ms", "total_duration_ms"]
const metricsDeclaration = [
  "maven.sync.duration",
  "maven.import.after.import.configuration",
  "maven.import.stats.applying.model.task",
  "maven.import.stats.importing.task",
  "maven.import.stats.importing.task.old",
  "maven.project.importer.base.refreshing.files.task",
  "maven.projects.processor.plugin.resolving.task",
  "maven.projects.processor.reading.task",
  "maven.projects.processor.resolving.task",
  "maven.projects.processor.wait.for.completion.task",
  "quarkus.maven.importer.base.task",
  "quarkus.maven.post.processor.task",
  "jpa.facet.importer.lambdas",
  "jpa.facet.importer.reimport.facet",
  "importer_run.com.intellij.jpa.importer.maven.JpaFacetImporter.total_duration_ms",
  "importer_run.com.intellij.quarkus.run.maven.QsMavenImporter.total_duration_ms",
  "importer_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms",
  "importer_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms",
  "importer_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporter.total_duration_ms",
  "legacy_import.create_modules.duration_ms",
  "legacy_import.delete_obsolete.duration_ms",
  "legacy_import.duration_ms",
  "legacy_import.importers.duration_ms",
  "workspace_commit.attempts",
  "workspace_commit.duration_in_background_ms",
  "workspace_commit.duration_in_write_action_ms",
  "workspace_commit.duration_of_workspace_update_call_ms",
  "workspace_import.commit.duration_ms",
  ...metricsName.map((metric) => "workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter" + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx." + metric),
  ...metricsName.map((metric) => "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx." + metric),
  "workspace_import.duration_ms",
  "workspace_import.legacy_importers.duration_ms",
  "workspace_import.legacy_importers.stats.duration_of_bridges_creation_ms",
  "workspace_import.legacy_importers.stats.duration_of_bridges_commit_ms",
  "workspace_import.populate.duration_ms",
  "maven.project.importer.post.importing.task.marker",
  "post_import_tasks_run.total_duration_ms",

  "AWTEventQueue.dispatchTimeTotal",
  "CPU | Load |Total % 95th pctl",
  "Memory | IDE | RESIDENT SIZE (MB) 95th pctl",
  "Memory | IDE | VIRTUAL SIZE (MB) 95th pctl",
  "gcPause",
  "gcPauseCount",
  "fullGCPause",
  "freedMemoryByGC",
  "totalHeapUsedMax",
]

const chartsDeclaration: ChartDefinition[] = metricsDeclaration.map((metric) => {
  return {
    labels: [metric],
    measures: [metric],
    projects: [
      "project-import-maven-quarkus/measureStartup",
      "project-reimport-maven-quarkus/measureStartup",
      "project-import-from-cache-maven-quarkus/measureStartup",
      "project-import-maven-1000-modules/measureStartup",
      "project-import-maven-5000-modules/measureStartup",
      "project-import-maven-keycloak/measureStartup",
      "project-import-maven-javaee7/measureStartup",
      "project-import-maven-javaee8/measureStartup",
      "project-import-maven-jersey/measureStartup",
      "project-import-maven-flink/measureStartup",
      "project-import-maven-drill/measureStartup",
      "project-import-maven-azure-sdk-java/measureStartup",
      "project-import-maven-hive/measureStartup",
      "project-import-maven-quarkus-to-legacy-model/measureStartup",
      "project-import-maven-1000-modules-to-legacy-model/measureStartup",
    ],
  }
})
const charts = combineCharts(chartsDeclaration)
</script>
