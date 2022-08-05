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
      label="Slow Inspections"
      measure="globalInspections"
      :projects="['drupal8-master-with-plugin/inspection', 'magento/inspection', 'wordpress/inspection',
                  'laravel-io/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Fast Inspections"
      measure="globalInspections"
      :projects="['mediawiki/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Slow Local Inspections"
      measure="localInspections"
      :projects="['mpdf/localInspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Slow Indexing"
      measure="indexing"
      :projects="['drupal8-master-with-plugin/indexing', 'laravel-io/indexing','wordpress/indexing','mediawiki/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Medium Indexing"
      measure="indexing"
      :projects="['magento/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="PHP Typing Time"
      measure="typing"
      :projects="['mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="PHP Typing Average Responsiveness"
      measure="test#average_awt_delay"
      :projects="[ 'mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="PHP Typing Responsiveness"
      measure="test#max_awt_delay"
      :projects="['mpdf/typing', 'mpdf_powersave/typing'
      ]"
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

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("phpstorm_dashboard", {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "phpstormWithPlugins")
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