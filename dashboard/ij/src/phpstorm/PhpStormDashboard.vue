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
      measure="inspection_execution_time"
      :projects="['drupal8-master-with-plugin/inspection', 'shopware/inspection', 'b2c-demo-shop/inspection', 'magento/inspection', 'wordpress/inspection',
                  'laravel-io/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Fast Inspections"
      measure="inspection_execution_time"
      :projects="['mediawiki/inspection','php-cs-fixer/inspection', 'proxyManager/inspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Slow Local Inspections"
      measure="local_inspection_execution_time"
      :projects="['mpdf/localInspection', 'WI_65655/localInspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Fast Local Inspections"
      measure="local_inspection_execution_time"
      :projects="['WI_59961/localInspection', 'bitrix/localInspection', 'WI_65893/localInspection']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Slow Indexing"
      measure="indexing"
      :projects="['b2c-demo-shop/indexing', 'bitrix/indexing', 'oro/indexing', 'ilias/indexing', 'magento2/indexing', 'drupal8-master-with-plugin/indexing', 
                  'laravel-io/indexing','wordpress/indexing','mediawiki/indexing', 'WI_66681/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Medium Indexing"
      measure="indexing"
      :projects="['WI_39333/indexing', 'php-cs-fixer/indexing','many_classes/indexing', 'magento/indexing', 'proxyManager/indexing',
                  'shopware/indexing', 'dql/indexing', 'tcpdf/indexing', 'WI_51645/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Fast Indexing"
      measure="indexing"
      :projects="['phpcs/indexing','empty_project/indexing','complex_meta/indexing', 'WI_53502/indexing','heredoc/indexing', 'many_array_access/indexing', 'WI_66279/indexing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion"
      measure="completion_execution_time"
      :projects="['many_classes/completion/classes','magento2/completion/function_var', 'magento2/completion/function_stlr', 'magento2/completion/classes','dql/completion',
                  'WI_64694/completion','WI_58919/completion', 'WI_58807/completion', 'WI_58306/completion']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="PHP Typing Time"
      measure="typing_total_time"
      :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing', 
                  'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                  'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="PHP Typing Average Responsiveness"
      measure="average_responsiveness_time"
      :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                  'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                  'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="PHP Typing Responsiveness"
      measure="responsiveness_time"
      :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                  'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                  'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Blade Typing Time"
      measure="typing_total_time"
      :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                  'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Blade Average Responsiveness"
      measure="average_responsiveness_time"
      :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                  'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Blade Responsiveness"
      measure="responsiveness_time"
      :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                  'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
      :server-configurator="serverConfigurator"
    />
  </Dashboard>
</template>

<script lang="ts" setup>
import { initDataComponent } from "shared/src/DataQueryExecutor"
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import { chartDefaultStyle } from "shared/src/chart"
import Dashboard from "shared/src/components/Dashboard.vue"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { chartStyleKey } from "shared/src/injectionKeys"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide } from "vue"
import { useRouter } from "vue-router"
import DimensionHierarchicalSelect from "../../../shared/src/components/DimensionHierarchicalSelect.vue"
import DimensionSelect from "../../../shared/src/components/DimensionSelect.vue"
import GroupLineChart from "../../../shared/src/components/GroupLineChart.vue"
import TimeRangeSelect from "../../../shared/src/components/TimeRangeSelect.vue"

provide(chartStyleKey, {
  ...chartDefaultStyle,
})

provideReportUrlProvider()

const persistentStateManager = new PersistentStateManager("phpstorm_dashboard", {
  machine: "linux-blade-hetzner",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "phpstorm")
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true, [timeRangeConfigurator], (a, _) => {
  return a.includes("/") ? 1 : -1
})
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [])
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