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
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>
    <GroupLineChart
      label="Gradle: gradle.sync.duration"
      measure="gradle.sync.duration"
      :projects="[
        'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
        'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
        'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
        'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
        'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
        'project-import-gradle-openliberty/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Gradle: GRADLE_CALL"
      measure="GRADLE_CALL"
      :projects="[
        'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
        'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
        'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
        'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
        'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
        'project-import-gradle-openliberty/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Gradle: PROJECT_RESOLVERS"
      measure="PROJECT_RESOLVERS"
      :projects="[
        'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
        'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
        'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
        'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
        'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
        'project-import-gradle-openliberty/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Gradle: DATA_SERVICES"
      measure="DATA_SERVICES"
      :projects="[
        'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
        'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
        'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
        'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
        'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
        'project-import-gradle-openliberty/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Gradle: WORKSPACE_MODEL_APPLY"
      measure="WORKSPACE_MODEL_APPLY"
      :projects="[
        'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
        'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
        'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
        'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
        'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
        'project-import-gradle-openliberty/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Gradle: CPU | Load | 75th pctl"
      measure="CPU | Load | 75th pctl"
      :projects="[
        'project-import-gradle-monolith-51-modules-4000-dependencies-2000000-files/measureStartup',
        'project-import-gradle-micronaut/measureStartup', 'project-import-gradle-hibernate-orm/measureStartup',
        'project-import-gradle-cas/measureStartup', 'project-import-gradle-500-modules/measureStartup',
        'project-import-gradle-1000-modules/measureStartup', 'project-reimport-space/measureStartup',
        'project-import-space/measureStartup', 'project-import-open-telemetry/measureStartup',
        'project-import-gradle-openliberty/measureStartup'
      ]"
      :server-configurator="serverConfigurator"
    />

    <GroupLineChart
      label="Maven: maven.sync.duration"
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
      label="Maven: maven.projects.processor.reading.task"
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
      label="Maven: maven.projects.processor.resolving.task"
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
      label="Maven: maven.projects.processor.plugin.resolving.task"
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
      label="Maven: maven.import.stats.importing.task"
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
      label="Maven: maven.import.stats.importing.task.old"
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
      label="Maven: maven.import.stats.applying.model.task"
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
      label="Maven: maven.import.after.import.configuration"
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
      label="Maven: maven.project.importer.base.refreshing.files.task"
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
      label="Maven: maven.projects.processor.wait.for.completion.task"
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
      label="Maven: jpa.facet.importer.reimport.facet"
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
      label="Maven: quarkus.maven.importer.base.task"
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
      label="Maven: quarkus.maven.post.processor.task"
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
      label="Maven: jpa.facet.importer.lambdas"
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