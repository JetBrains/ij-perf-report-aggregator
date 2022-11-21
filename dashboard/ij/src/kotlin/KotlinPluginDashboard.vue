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
      label="Completion empty kotlin project"
      measure="completion#mean_value"
      :projects="[
        'kotlin_empty/completion/empty_place_with_library_cache_k1',
        'kotlin_empty/completion/empty_place_without_library_cache_k1',
        'kotlin_empty/completion/empty_place_with_library_cache_k2',
        'kotlin_empty/completion/empty_place_without_library_cache_k2',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion mean value with Library cache K1"
      measure="completion#mean_value"
      :projects="[
        'intellij_commit/completion/empty_place_with_library_cache_k1',
        'intellij_commit/completion/after_parameter_with_library_cache_k1',
        'kotlin_lang/completion/after_parameter_with_library_cache_k1',
        'kotlin_lang/completion/empty_place_with_library_cache_k1'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion mean value with Library cache K2"
      measure="completion#mean_value"
      :projects="[
        'intellij_commit/completion/empty_place_with_library_cache_k2',
        'intellij_commit/completion/after_parameter_with_library_cache_k2',
        'kotlin_lang/completion/after_parameter_with_library_cache_k2',
        'kotlin_lang/completion/empty_place_with_library_cache_k2',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion mean value without Library cache K1"
      measure="completion#mean_value"
      :projects="[
        'intellij_commit/completion/empty_place_without_library_cache_k1',
        'intellij_commit/completion/after_parameter_without_library_cache_k1',
        'kotlin_lang/completion/after_parameter_without_library_cache_k1',
        'kotlin_lang/completion/empty_place_without_library_cache_k1',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion mean value without Library cache K2"
      measure="completion#mean_value"
      :projects="[
        'intellij_commit/completion/empty_place_without_library_cache_k2',
        'intellij_commit/completion/after_parameter_without_library_cache_k2',
        'kotlin_lang/completion/after_parameter_without_library_cache_k2',
        'kotlin_lang/completion/empty_place_without_library_cache_k2',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Highlight on random files - time to line mean value"
      measure="highlighting#timeToLines#mean_value"
      :projects="[
       'intellij_commit/highlightOnRandomFiles',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Highlight mean value with Library cache K1"
      measure="highlighting#mean_value"
      :projects="[
       'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k1',
        'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k1',
        'intellij_commit/highlight/OraIntrospector_with_library_cache_k1',
        'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k1',
        'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k1',
        'kotlin_lang/highlight/CommonParser_with_library_cache_k1',
        'kotlin_lang/highlight/FirErrors_with_library_cache_k1',
        'kotlin_lang/highlight/Flag_with_library_cache_k1',
        'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k1',
        'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k1',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Highlight mean value with Library cache K2"
      measure="highlighting#mean_value"
      :projects="[
        'intellij_commit/highlight/DexInlineCallStackComparisonTest_with_library_cache_k2',
        'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_with_library_cache_k2',
        'intellij_commit/highlight/OraIntrospector_with_library_cache_k2',
        'intellij_commit/highlight/SolutionModel.Generated_with_library_cache_k2',
        'intellij_commit/highlight/UIAutomationInteractionModel.Generated_with_library_cache_k2',
        'kotlin_lang/highlight/CommonParser_with_library_cache_k2',
        'kotlin_lang/highlight/FirErrors_with_library_cache_k2',
        'kotlin_lang/highlight/Flag_with_library_cache_k2',
        'kotlin_lang/highlight/KtFirDataClassConverters_with_library_cache_k2',
        'kotlin_lang/highlight/DefaultArgumentStubGenerator_with_library_cache_k2',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Highlight mean value without Library cache K1"
      measure="highlighting#mean_value"
      :projects="[
        'intellij_commit/highlight/DexInlineCallStackComparisonTest_without_library_cache_k1',
        'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_without_library_cache_k1',
        'intellij_commit/highlight/OraIntrospector_without_library_cache_k1',
        'intellij_commit/highlight/SolutionModel.Generated_without_library_cache_k1',
        'intellij_commit/highlight/UIAutomationInteractionModel.Generated_without_library_cache_k1',
        'kotlin_lang/highlight/CommonParser_without_library_cache_k1',
        'kotlin_lang/highlight/FirErrors_without_library_cache_k1',
        'kotlin_lang/highlight/Flag_without_library_cache_k1',
        'kotlin_lang/highlight/KtFirDataClassConverters_without_library_cache_k1',
        'kotlin_lang/highlight/DefaultArgumentStubGenerator_without_library_cache_k1',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Highlight mean value without Library cache K2"
      measure="highlighting#mean_value"
      :projects="[
        'intellij_commit/highlight/DexInlineCallStackComparisonTest_without_library_cache_k2',
        'intellij_commit/highlight/DexLocalVariableTableBreakpointTest_without_library_cache_k2',
        'intellij_commit/highlight/OraIntrospector_without_library_cache_k2',
        'intellij_commit/highlight/SolutionModel.Generated_without_library_cache_k2',
        'intellij_commit/highlight/UIAutomationInteractionModel.Generated_without_library_cache_k2',
        'kotlin_lang/highlight/CommonParser_without_library_cache_k2',
        'kotlin_lang/highlight/FirErrors_without_library_cache_k2',
        'kotlin_lang/highlight/Flag_without_library_cache_k2',
        'kotlin_lang/highlight/KtFirDataClassConverters_without_library_cache_k2',
        'kotlin_lang/highlight/DefaultArgumentStubGenerator_without_library_cache_k2',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages mean value with Library cache K1"
      measure="findUsages#mean_value"
      :projects="[
        'intellij_commit/findUsages/loadModuleEntity_with_library_cache_k1',
        'intellij_commit/findUsages/SqlBlock_with_library_cache_k1',
        'intellij_commit/findUsages/OraIntrospector_with_library_cache_k1',
        'intellij_commit/findUsages/UIAutomationInteractionModel.Generated_with_library_cache_k1',
        'intellij_commit/findUsages/SolutionModel.Generated_with_library_cache_k1',
        'kotlin_lang/findUsages/FirErrors_with_library_cache_k1',
        'kotlin_lang/findUsages/Flag_with_library_cache_k1',
        'kotlin_lang/findUsages/CommonParser_with_library_cache_k1',
        'kotlin_lang/findUsages/SymbolTable_with_library_cache_k1',
        'kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache_k1',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages mean value with Library cache K2"
      measure="findUsages#mean_value"
      :projects="[
        'intellij_commit/findUsages/loadModuleEntity_with_library_cache_k2',
        'intellij_commit/findUsages/SqlBlock_with_library_cache_k2',
        'intellij_commit/findUsages/OraIntrospector_with_library_cache_k2',
        'intellij_commit/findUsages/UIAutomationInteractionModel.Generated_with_library_cache_k2',
        'intellij_commit/findUsages/SolutionModel.Generated_with_library_cache_k2',
        'kotlin_lang/findUsages/FirErrors_with_library_cache_k2',
        'kotlin_lang/findUsages/Flag_with_library_cache_k2',
        'kotlin_lang/findUsages/CommonParser_with_library_cache_k2',
        'kotlin_lang/findUsages/SymbolTable_with_library_cache_k2',
        'kotlin_lang/findUsages/DefaultArgumentStubGenerator_with_library_cache_k2',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages mean value without Library cache K1"
      measure="findUsages#mean_value"
      :projects="[
        'intellij_commit/findUsages/loadModuleEntity_without_library_cache_k1',
        'intellij_commit/findUsages/SqlBlock_without_library_cache_k1',
        'intellij_commit/findUsages/OraIntrospector_without_library_cache_k1',
        'intellij_commit/findUsages/UIAutomationInteractionModel.Generated_without_library_cache_k1',
        'intellij_commit/findUsages/SolutionModel.Generated_without_library_cache_k1',
        'kotlin_lang/findUsages/FirErrors_without_library_cache_k1',
        'kotlin_lang/findUsages/Flag_without_library_cache_k1',
        'kotlin_lang/findUsages/CommonParser_without_library_cache_k1',
        'kotlin_lang/findUsages/SymbolTable_without_library_cache_k1',
        'kotlin_lang/findUsages/DefaultArgumentStubGenerator_without_library_cache_k1',
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="FindUsages mean value without Library cache K2"
      measure="findUsages#mean_value"
      :projects="[
        'intellij_commit/findUsages/loadModuleEntity_without_library_cache_k2',
        'intellij_commit/findUsages/SqlBlock_without_library_cache_k2',
        'intellij_commit/findUsages/OraIntrospector_without_library_cache_k2',
        'intellij_commit/findUsages/UIAutomationInteractionModel.Generated_without_library_cache_k2',
        'intellij_commit/findUsages/SolutionModel.Generated_without_library_cache_k2',
        'kotlin_lang/findUsages/FirErrors_without_library_cache_k2',
        'kotlin_lang/findUsages/Flag_without_library_cache_k2',
        'kotlin_lang/findUsages/CommonParser_without_library_cache_k2',
        'kotlin_lang/findUsages/SymbolTable_without_library_cache_k2',
        'kotlin_lang/findUsages/DefaultArgumentStubGenerator_without_library_cache_k2',
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
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider(false)

const persistentStateManager = new PersistentStateManager("kotlinplugin_dashboard", {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "kotlin")
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
  triggeredByConfigurator,
  releaseConfigurator,
]
initDataComponent(configurators)
</script>