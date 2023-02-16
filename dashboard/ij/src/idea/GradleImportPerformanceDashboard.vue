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
      label="gradle.sync.duration"
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
      label="GRADLE_CALL"
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
      label="PROJECT_RESOLVERS"
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
      label="DATA_SERVICES"
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
      label="WORKSPACE_MODEL_APPLY"
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
      label="CPU | Load | 75th pctl"
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