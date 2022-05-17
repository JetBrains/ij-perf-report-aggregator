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
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>
    <GroupLineChart
      label="Indexing"
      measure="indexing"
      :projects="['community/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning"
      measure="scanning"
      :projects="['community/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Rebuild"
      measure="build_compilation_duration"
      :projects="['community/rebuild','intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Kotlin Builder"
      measure="kotlin_builder_time"
      :projects="['community/rebuild','intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Inspection"
      measure="inspection_execution_time"
      :projects="['java/inspection', 'grails/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Java"
      measure="find_usages_execution_time"
      :projects="['community/findUsages/PsiManager_getInstance_Before', 'community/findUsages/PsiManager_getInstance_After']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Kotlin"
      measure="find_usages_execution_time"
      :projects="['community/findUsages/LanguageInjectionCondition_getId_Before', 'community/findUsages/LanguageInjectionCondition_getId_After']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Local Inspection"
      measure="local_inspection_execution_time"
      :projects="['intellij_sources/localInspection/java_file','intellij_sources/localInspection/kotlin_file']"
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

const persistentStateManager = new PersistentStateManager("idea_dashboard", {
  machine: "windows-blade",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "idea")
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true)
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [])
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistentStateManager)

const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator
]
initDataComponent(configurators)
</script>