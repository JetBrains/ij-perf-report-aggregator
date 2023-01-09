<template>
  <Dashboard>
    <template #toolbar>
      <DimensionSelect
        label="Branch"
        :dimension="branchConfigurator"
      />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <DimensionSelect
        label="Nightly/Release"
        :dimension="releaseConfigurator"
      />
      <DimensionSelect
        label="Triggered by"
        :dimension="triggeredByConfigurator"
      />
      <TimeRangeSelect :configurator="timeRangeConfigurator"/>
    </template>
    <GroupLineChart
      label="workspace_commit.attempts"
      measure="workspace_commit.attempts"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_commit.duration_in_background_ms"
      measure="workspace_commit.duration_in_background_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_commit.duration_in_write_action_ms"
      measure="workspace_commit.duration_in_write_action_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_commit.duration_of_workspace_update_call_ms"
      measure="workspace_commit.duration_of_workspace_update_call_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.commit.duration_ms"
      measure="workspace_import.commit.duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.after_apply_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.before_apply_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.config_modules_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.after_apply_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.before_apply_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.config_modules_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms"
      measure="workspace_import.configurator_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.importing.MavenWslTargetConfigurator.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.after_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.before_apply_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.collect_folders_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.config_modules_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.plugins.groovy.GroovyPluginConfigurator.total_duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.duration_ms"
      measure="workspace_import.duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.legacy_importers.duration_ms"
      measure="workspace_import.legacy_importers.duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="workspace_import.populate.duration_ms"
      measure="workspace_import.populate.duration_ms"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />

    <GroupLineChart
      label="Maven: CPU | Load | 75th pctl"
      measure="CPU | Load | 75th pctl"
      :projects="[
        'project-import-maven-quarkus/measureStartup',
        'project-import-maven-500-modules/measureStartup', 'project-import-maven-1000-modules/measureStartup',
        'project-import-maven-keycloak/measureStartup', 'project-import-maven-javaee7/measureStartup',
        'project-import-maven-javaee8/measureStartup', 'project-import-maven-jersey/measureStartup',
        'project-import-maven-flink/measureStartup', 'project-import-maven-drill/measureStartup',
        'project-import-maven-azure-sdk-java/measureStartup', 'project-import-maven-hive/measureStartup',
        'project-import-maven-quarkus-to-legacy-model/measureStartup', 'project-import-maven-500-modules-to-legacy-model/measureStartup',
        'project-import-maven-1000-modules-to-legacy-model/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
  </Dashboard>
</template>

<script lang="ts" setup>
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import Dashboard from "shared/src/components/Dashboard.vue"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "shared/src/components/DimensionSelect.vue"
import GroupLineChart from "shared/src/components/GroupLineChart.vue"
import TimeRangeSelect from "shared/src/components/TimeRangeSelect.vue"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("import_dashboard", {
  machine: "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "idea")
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true, [timeRangeConfigurator], (a, _) => {
  return a.includes("/") ? 1 : -1
})
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const releaseConfigurator = new ReleaseNightlyConfigurator(persistentStateManager)
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])

const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator,
  triggeredByConfigurator,
]
initDataComponent(configurators)
</script>