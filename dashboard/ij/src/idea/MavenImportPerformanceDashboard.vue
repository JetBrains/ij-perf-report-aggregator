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
      label="maven.sync.duration"
      measure="maven.sync.duration"
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
      label="maven.import.after.import.configuration"
      measure="maven.import.after.import.configuration"
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
      label="maven.import.stats.applying.model.task"
      measure="maven.import.stats.applying.model.task"
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
      label="maven.import.stats.importing.task"
      measure="maven.import.stats.importing.task"
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
      label="maven.import.stats.importing.task.old"
      measure="maven.import.stats.importing.task.old"
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
      label="maven.project.importer.base.refreshing.files.task"
      measure="maven.project.importer.base.refreshing.files.task"
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
      label="maven.projects.processor.plugin.resolving.task"
      measure="maven.projects.processor.plugin.resolving.task"
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
      label="maven.projects.processor.reading.task"
      measure="maven.projects.processor.reading.task"
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
      label="maven.projects.processor.resolving.task"
      measure="maven.projects.processor.resolving.task"
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
      label="maven.projects.processor.wait.for.completion.task"
      measure="maven.projects.processor.wait.for.completion.task"
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
      label="quarkus.maven.importer.base.task"
      measure="quarkus.maven.importer.base.task"
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
      label="quarkus.maven.post.processor.task"
      measure="quarkus.maven.post.processor.task"
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
      label="jpa.facet.importer.lambdas"
      measure="jpa.facet.importer.lambdas"
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
      label="jpa.facet.importer.reimport.facet"
      measure="jpa.facet.importer.reimport.facet"
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
      label="importer_run.com.intellij.jpa.importer.maven.JpaFacetImporter.total_duration_ms"
      measure="importer_run.com.intellij.jpa.importer.maven.JpaFacetImporter.total_duration_ms"
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
      label="importer_run.com.intellij.quarkus.run.maven.QsMavenImporter.total_duration_ms"
      measure="importer_run.com.intellij.quarkus.run.maven.QsMavenImporter.total_duration_ms"
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
      label="importer_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms"
      measure="importer_run.com.intellij.spring.facet.importer.maven.SpringFacetImporter.total_duration_ms"
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
      label="importer_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms"
      measure="importer_run.com.intellij.spring.mvc.importer.boot.SpringBootWebFacetImporter.total_duration_ms"
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
      label="importer_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms"
      measure="importer_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms"
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
      label="importer_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.total_duration_ms"
      measure="importer_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.total_duration_ms"
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
      label="importer_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms"
      measure="importer_run.org.jetbrains.idea.maven.importing.MavenAnnotationProcessorConfigurator.total_duration_ms"
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
      label="importer_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms"
      measure="importer_run.org.jetbrains.idea.maven.importing.MavenCompilerConfigurator.total_duration_ms"
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
      label="importer_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms"
      measure="importer_run.org.jetbrains.idea.maven.importing.MavenEncodingConfigurator.total_duration_ms"
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
      label="importer_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms"
      measure="importer_run.org.jetbrains.idea.maven.importing.MavenExternalAnnotationsConfigurator.total_duration_ms"
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
      label="importer_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms"
      measure="importer_run.org.jetbrains.idea.maven.importing.MavenRemoteRepositoriesConfigurator.total_duration_ms"
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
      label="importer_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporter.total_duration_ms"
      measure="importer_run.org.jetbrains.kotlin.idea.maven.KotlinMavenImporter.total_duration_ms"
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
      label="legacy_import.create_modules.duration_ms"
      measure="legacy_import.create_modules.duration_ms"
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
      label="legacy_import.delete_obsolete.duration_ms"
      measure="legacy_import.delete_obsolete.duration_ms"
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
      label="legacy_import.duration_ms"
      measure="legacy_import.duration_ms"
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
      label="legacy_import.importers.duration_ms"
      measure="legacy_import.importers.duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.total_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.collect_folders_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.config_modules_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.before_apply_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.web.WebFacetImporter.after_apply_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.total_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.collect_folders_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.collect_folders_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.config_modules_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.config_modules_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.before_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.before_apply_duration_ms"
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
      label="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.after_apply_duration_ms"
      measure="workspace_import.configurator_run.org.jetbrains.idea.maven.ext.javaee.ear.EarFacetImporter.after_apply_duration_ms"
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
      label="maven.project.importer.post.importing.task.marker"
      measure="maven.project.importer.post.importing.task.marker"
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
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
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
const branchConfigurator = createBranchConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator])
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