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
      label="Completion with Library cache"
      measure="completion"
      :projects="[
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_after_dot_with_library_cache_k1',
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_on_empty_place_with_library_cache_k1',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_after_parameter_with_library_cache_k1',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_on_empty_placer_with_library_cache_k1',
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_after_dot_with_library_cache_k2',
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_on_empty_place_with_library_cache_k2',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_after_parameter_with_library_cache_k2',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_on_empty_placer_with_library_cache_k2'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion without Library cache"
      measure="completion"
      :projects="[
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_after_dot_without_library_cache_k1',
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_on_empty_place_without_library_cache_k1',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_after_parameter_without_library_cache_k1',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_on_empty_placer_without_library_cache_k1',
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_after_dot_without_library_cache_k2',
        'intellij_sources_specific_commit/completion/intellij_completion_in_method_on_empty_place_without_library_cache_k2',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_after_parameter_without_library_cache_k2',
        'kotlin_lang/completion/kotlin_lang_completion_in_method_on_empty_placer_without_library_cache_k2',
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

const persistentStateManager = new PersistentStateManager("kotlinplugin_dashboard", {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("kotlin", "kotlinPlugin")
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
  triggeredByConfigurator
]
initDataComponent(configurators)
</script>