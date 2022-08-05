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
        label="Nightly/Release"
        :dimension="releaseConfigurator"
      />
      <DimensionSelect
        label="Triggered by"
        :dimension="triggeredByConfigurator"
      />
    </template>
    <GroupLineChart
      label="Indexing Long"
      measure="indexing"
      :projects="['community/indexing', 'lock-free-vfs-record-storage-intellij_sources/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning Long"
      measure="scanning"
      :projects="['community/indexing', 'lock-free-vfs-record-storage-intellij_sources/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Indexing Short"
      measure="indexing"
      :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                  'spring_boot/indexing', 'spring_boot_maven/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning Short"
      measure="scanning"
      :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                  'spring_boot/indexing', 'spring_boot_maven/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Rebuild Long"
      measure="build_compilation_duration"
      :projects="['community/rebuild','intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Kotlin Builder Long"
      measure="kotlin_builder_time"
      :projects="['community/rebuild','intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Java Builder Long"
      measure="java_time"
      :projects="['community/rebuild','intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Rebuild Short"
      measure="build_compilation_duration"
      :projects="['grails/rebuild','java/rebuild','spring_boot/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Kotlin Builder Short"
      measure="kotlin_builder_time"
      :projects="['grails/rebuild','java/rebuild','spring_boot/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Java Builder Short"
      measure="java_time"
      :projects="['grails/rebuild','java/rebuild','spring_boot/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Inspection"
      measure="globalInspections"
      :projects="['java/inspection', 'grails/inspection', 'spring_boot_maven/inspection', 'spring_boot/inspection', 'kotlin/inspection', 'kotlin_coroutines/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Java"
      measure="findUsages"
      :projects="['community/findUsages/PsiManager_getInstance_Before', 'community/findUsages/PsiManager_getInstance_After']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Kotlin"
      measure="findUsages"
      :projects="['community/findUsages/LocalInspectionTool_getID_Before', 'community/findUsages/LocalInspectionTool_getID_After']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Local Inspection"
      measure="localInspections"
      :projects="['intellij_sources/localInspection/java_file','intellij_sources/localInspection/kotlin_file']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion"
      measure="completion"
      :projects="['community/completion/kotlin_file','grails/completion/groovy_file', 'grails/completion/java_file']"
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
import { ReleaseNightlyConfigurator } from "../../shared/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "../../shared/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "../../shared/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "../../shared/injectionKeys"
import { provideReportUrlProvider } from "../../shared/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("idea_dashboard", {
  machine: "macMini Intel 3.2, 16GB",
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