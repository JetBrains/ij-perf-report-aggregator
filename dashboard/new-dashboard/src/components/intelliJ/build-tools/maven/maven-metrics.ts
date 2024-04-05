export const MAVEN_METRICS_DASHBOARD = [
  // Main flow
  "maven.sync.duration",
  "maven.projects.processor.resolving.task",
  "maven.projects.processor.reading.task",
  "maven.import.stats.plugins.resolving.task",
  "maven.import.stats.applying.model.task",
  "maven.import.stats.sync.project.task",
  "maven.import.after.import.configuration",
  "maven.import.configure.post.process",
  "maven.project.importer.base.refreshing.files.task",
  "maven.project.importer.post.importing.task.marker",
  "post_import_tasks_run.total_duration_ms",

  // Workspace commit
  "workspace_commit.attempts",
  "workspace_commit.duration_in_background_ms",
  "workspace_commit.duration_in_write_action_ms",
  "workspace_commit.duration_of_workspace_update_call_ms",

  // Workspace import
  "workspace_import.commit.duration_ms",
  "workspace_import.duration_ms",
  "workspace_import.legacy_importers.duration_ms",
  "workspace_import.legacy_importers.stats.duration_of_bridges_creation_ms",
  "workspace_import.legacy_importers.stats.duration_of_bridges_commit_ms",
  "workspace_import.populate.duration_ms",

  // Legacy import
  "legacy_import.create_modules.duration_ms",
  "legacy_import.delete_obsolete.duration_ms",
  "legacy_import.duration_ms",
  "legacy_import.importers.duration_ms",
]

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
  "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporterEx",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporterEx",
  "workspace_import.configurator_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporterEx",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.buildHelper.MavenBuildHelperPluginConfigurator",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.workspaceModel.LegacyToWorkspaceConfiguratorBridgeDynamic",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.importing.workspaceModel.LegacyToWorkspaceConfiguratorBridgeStatic",
  "workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.buildHelper.MavenBuildHelperPluginConfigurator",
]

export const MAVEN_METRICS_CONFIGURATORS = [
  // Importer run metrics
  ...importerRunMetrics.flatMap((metric) => importerRunMetricTypes.map((type) => `${metric}.${type}`)),
  // Workspace configuration metrics
  ...workspaceConfiguratorsMetrics.flatMap((metric) => workspaceConfiguratorsMetricTypes.map((type) => `${metric}.${type}`)),
]
