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
      label="Indexing Long"
      measure="indexing"
      :projects="['community/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning Long"
      measure="scanning"
      :projects="['community/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="VFS Record Storage Indexing"
      measure="indexing"
      :projects="['vfs-record-storage/in-memory-intellij-sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij-sources/indexing',
      'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij-sources/indexing', 'vfs-record-storage/in-memory-with-streamlined-attributes-intellij-sources/indexing',
      'vfs-record-storage/lock-free-intellij-sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="VFS Record Storage Scanning"
      measure="scanning"
      :projects="['vfs-record-storage/in-memory-intellij-sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij-sources/indexing',
      'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij-sources/indexing', 'vfs-record-storage/in-memory-with-streamlined-attributes-intellij-sources/indexing',
      'vfs-record-storage/lock-free-intellij-sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Indexing Short"
      measure="indexing"
      :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                  'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning Short"
      measure="scanning"
      :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                  'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="VFS Refresh"
      measure="vfs_initial_refresh"
      :projects="['intellij_sources/vfsRefresh/default', 'intellij_sources/vfsRefresh/with-1-thread(s)', 'intellij_sources/vfsRefresh/with-6-thread(s)']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Rebuild Long"
      measure="build_compilation_duration"
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
      label="Inspection"
      measure="globalInspections"
      :projects="['java/inspection', 'grails/inspection', 'spring_boot_maven/inspection', 'spring_boot/inspection', 'kotlin/inspection', 'kotlin_coroutines/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Long"
      measure="findUsages"
      :projects="['community/findUsages/PsiManager_getInstance_Before', 'community/findUsages/PsiManager_getInstance_After', 'community/findUsages/LocalInspectionTool_getID_Before', 'community/findUsages/LocalInspectionTool_getID_After']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages Short"
      measure="findUsages"
      :projects="['community/findUsages/ActionsKt_runReadAction_Before', 'community/findUsages/ActionsKt_runReadAction_After',
      'community/findUsages/Application_runReadAction_Before', 'community/findUsages/Application_runReadAction_After',
      'community/findUsages/Persistent_absolutePath_Before', 'community/findUsages/Persistent_absolutePath_After',
      'community/findUsages/PropertyMapping_value_Before', 'community/findUsages/PropertyMapping_value_After',
      'community/findUsages/RelativeTextEdit_rangeTo_Before', 'community/findUsages/RelativeTextEdit_rangeTo_After',
      'community/findUsages/DynamicPluginListener_TOPIC_Before', 'community/findUsages/DynamicPluginListener_TOPIC_After',
      'community/findUsages/Path_div_Before', 'community/findUsages/Path_div_After']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Local Inspection"
      measure="localInspections"
      :projects="['intellij_sources/localInspection/java_file','intellij_sources/localInspection/kotlin_file', 'kotlin/localInspection', 'kotlin_coroutines/localInspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion"
      measure="completion"
      :projects="['community/completion/kotlin_file','grails/completion/groovy_file', 'grails/completion/java_file']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Go to Action"
      measure="searchEverywhere_action"
      :projects="['community/go-to-action/SharedIndex']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Go to Class"
      measure="searchEverywhere_class"
      :projects="['community/go-to-class/EditorImpl','community/go-to-class/SharedIndex']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Go to File"
      measure="searchEverywhere_file"
      :projects="['community/go-to-file/EditorImpl','community/go-to-file/SharedIndex']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Show Intentions"
      measure="show_intentions"
      :projects="['grails/showIntentions/Find cause', 'kotlin/showIntention/Import', 'spring_boot/showIntentions/']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Git Status"
      measure="git#status"
      :projects="['intellij_sources']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Expand Menu"
      measure="ExpandActionGroup"
      :projects="['intellij_sources/expandProjectMenu', 'intellij_sources/expandMainMenu', 'intellij_sources/expandEditorMenu']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Highlight"
      measure="highlighting"
      :projects="['kotlin/highlight', 'kotlin_coroutines/highlight']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing"
      measure="typing"
      :projects="['typingInJavaFile_16Threads/typing', 'typingInJavaFile_4Threads/typing', 'typingInKotlinFile_16Threads/typing', 'typingInKotlinFile_4Threads/typing']"
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