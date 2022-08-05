<template>
  <Dashboard>
    <template #toolbar>
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
      <DimensionSelect
        label="Branch"
        :dimension="branchConfigurator"
      />
      <DimensionHierarchicalSelect
        label="Machine"
        :dimension="machineConfigurator"
      />
      <DimensionSelect
        label="Triggered by"
        :dimension="triggeredByConfigurator"
      />
    </template>
    <GroupLineChart
      label="Indexing"
      measure="indexing"
      :projects="['intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning"
      measure="scanning"
      :projects="['intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Rebuild"
      measure="build_compilation_duration"
      :projects="['intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Kotlin Builder"
      measure="kotlin_builder_time"
      :projects="['intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Java Builder"
      measure="java_time"
      :projects="['intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Java"
      measure="findUsages"
      :projects="['intellij_sources/findUsages/Application_runReadAction', 'intellij_sources/findUsages/LocalInspectionTool_getID',
                  'intellij_sources/findUsages/PsiManager_getInstance', 'intellij_sources/findUsages/PropertyMapping_value']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Kotlin"
      measure="findUsages"
      :projects="['intellij_sources/findUsages/ActionsKt_runReadAction', 'intellij_sources/findUsages/DynamicPluginListener_TOPIC', 'intellij_sources/findUsages/Path_div',
                  'intellij_sources/findUsages/Persistent_absolutePath', 'intellij_sources/findUsages/RelativeTextEdit_rangeTo']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Local Inspection"
      measure="localInspections"
      :projects="['intellij_sources/localInspection/java_file','intellij_sources/localInspection/kotlin_file']"
      :server-configurator="serverConfigurator"
    />
  </Dashboard>
</template>

<script lang="ts" setup>
import { initDataComponent } from "../../shared/DataQueryExecutor"
import { PersistentStateManager } from "../../shared/PersistentStateManager"
import { chartDefaultStyle } from "../../shared/chart"
import Dashboard from "../../shared/components/Dashboard.vue"
import DimensionHierarchicalSelect from "../../shared/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../../shared/components/DimensionSelect.vue"
import GroupLineChart from "../../shared/components/GroupLineChart.vue"
import TimeRangeSelect from "../../shared/components/TimeRangeSelect.vue"
import { dimensionConfigurator } from "../../shared/configurators/DimensionConfigurator"
import { MachineConfigurator } from "../../shared/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "../../shared/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "../../shared/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "../../shared/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "../../shared/injectionKeys"
import { provideReportUrlProvider } from "../../shared/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider(false)

const persistentStateManager = new PersistentStateManager("idea_dashboard", {
  machine: "Linux EC2 C6i.8xlarge (32 vCPU Xeon, 64 GB)",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfintDev", "idea")
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true, [timeRangeConfigurator], (a, _) => {
  return a.includes("/") ? 1 : -1
})
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [timeRangeConfigurator, branchConfigurator])
const triggeredByConfigurator = privateBuildConfigurator(serverConfigurator, persistentStateManager, [branchConfigurator, timeRangeConfigurator])

const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  triggeredByConfigurator,
]
initDataComponent(configurators)
</script>