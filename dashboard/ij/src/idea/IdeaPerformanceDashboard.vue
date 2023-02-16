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
      label="Indexing (Big projects)"
      measure="indexing"
      :projects="['community/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning (Big projects)"
      measure="scanning"
      :projects="['community/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Number of indexing runs long (Big projects)"
      measure="numberOfIndexingRuns"
      :projects="['community/indexing', 'intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Processing speed of JAVA files"
      measure="processingSpeed#JAVA"
      :projects="['community/indexing', 'intellij_sources/indexing', 'empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing',
                  'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Processing speed of KOTLIN files"
      measure="processingSpeed#Kotlin"
      :projects="['community/indexing', 'intellij_sources/indexing', 'empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing',
                  'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Indexing with the new record storage (IntelliJ project)"
      measure="indexing"
      :projects="['vfs-record-storage/in-memory-intellij_sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij_sources/indexing',
                  'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij_sources/indexing', 
                  'vfs-record-storage/in-memory-with-streamlined-attributes-intellij_sources/indexing', 'vfs-record-storage/lock-free-intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning with the new record storage (IntelliJ project)"
      measure="scanning"
      :projects="['vfs-record-storage/in-memory-intellij_sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij_sources/indexing',
                  'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij_sources/indexing', 
                  'vfs-record-storage/in-memory-with-streamlined-attributes-intellij_sources/indexing','vfs-record-storage/lock-free-intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Number of indexing runs with the new record storage (IntelliJ project)"
      measure="numberOfIndexingRuns"
      :projects="['vfs-record-storage/in-memory-intellij_sources/indexing', 'vfs-record-storage/in-memory-with-non-strict-names-intellij_sources/indexing',
                  'vfs-record-storage/in-memory-with-non-strict-names-streamlined-attributes-intellij_sources/indexing', 
                  'vfs-record-storage/in-memory-with-streamlined-attributes-intellij_sources/indexing', 'vfs-record-storage/lock-free-intellij_sources/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Indexing"
      measure="indexing"
      :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                  'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Scanning"
      measure="scanning"
      :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing', 
                  'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Number of indexing runs"
      measure="numberOfIndexingRuns"
      :projects="['empty_project/indexing', 'grails/indexing', 'java/indexing', 'kotlin/indexing', 'kotlin_coroutines/indexing',
                  'spring_boot/indexing', 'spring_boot_maven/indexing', 'kotlin_petclinic/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="VFS Refresh"
      measure="vfs_initial_refresh"
      :projects="['intellij_sources/vfsRefresh/default', 'intellij_sources/vfsRefresh/with-1-thread(s)', 'intellij_sources/vfsRefresh/git-status']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Rebuild (Big projects)"
      measure="build_compilation_duration"
      :projects="['community/rebuild','intellij_sources/rebuild']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Rebuild"
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
      label="FindUsages PsiManager#getInstance Before and After Compilation"
      measure="findUsages"
      :projects="['community/findUsages/PsiManager_getInstance_Before', 'community/findUsages/PsiManager_getInstance_After'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages Library#getName Before and After Compilation"
      measure="findUsages"
      :projects="['community/findUsages/Library_getName_Before', 'community/findUsages/Library_getName_After'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages LocalInspectionTool#getID Before and After Compilation"
      measure="findUsages"
      :projects="['community/findUsages/LocalInspectionTool_getID_Before', 'community/findUsages/LocalInspectionTool_getID_After'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages ActionsKt#runReadAction and Application#runReadAction Before and After Compilation"
      measure="findUsages"
      :projects="['community/findUsages/ActionsKt_runReadAction_Before', 'community/findUsages/Application_runReadAction_Before',
                  'community/findUsages/ActionsKt_runReadAction_After', 'community/findUsages/Application_runReadAction_After'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages Persistent#absolutePath and PropertyMapping#value Before and After Compilation"
      measure="findUsages"
      :projects="['community/findUsages/Persistent_absolutePath_After', 'community/findUsages/PropertyMapping_value_After',
                  'community/findUsages/Persistent_absolutePath_Before', 'community/findUsages/PropertyMapping_value_Before'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Find Usages with idea.is.internal=true Before Compilation"
      measure="findUsages"
      :projects="['intellij_sources/findUsages/PsiManager_getInstance_firstCall', 'intellij_sources/findUsages/PsiManager_getInstance_secondCall',
                  'intellij_sources/findUsages/PsiManager_getInstance_thirdCallInternalMode']"
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
      label="Search Everywhere (Go to Action)"
      measure="searchEverywhere_action"
      :projects="['community/go-to-action/SharedIndex']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Search Everywhere (Go to Class)"
      measure="searchEverywhere_class"
      :projects="['community/go-to-class/EditorImpl','community/go-to-class/SharedIndex']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Search Everywhere (Go to File)"
      measure="searchEverywhere_file"
      :projects="['community/go-to-file/EditorImpl','community/go-to-file/SharedIndex']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Show Intentions (average awt delay)"
      measure="test#average_awt_delay"
      :projects="['grails/showIntentions/Find cause', 'kotlin/showIntention/Import', 'spring_boot/showIntentions/']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Expand Menu"
      measure="expandActionGroup"
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
      label="Typing during indexing (with changed count of indexing threads)"
      measure="typing"
      :projects="['typingInJavaFile_16Threads/typing', 'typingInJavaFile_4Threads/typing', 'typingInKotlinFile_16Threads/typing', 'typingInKotlinFile_4Threads/typing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing during indexing (average awt delay)"
      measure="test#average_awt_delay"
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

const persistentStateManager = new PersistentStateManager("idea_dashboard", {
  machine: "macMini Intel 3.2, 16GB",
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