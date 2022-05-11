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
      <TimeRangeSelect :configurator="timeRangeConfigurator" />
    </template>
    <GroupLineChart
      label="Slow Inspections"
      :measures="['inspection_execution_time']"
      :projects="['drupal8-master-with-plugin/inspection', 'shopware/inspection', 'b2c-demo-shop/inspection', 'magento/inspection', 'wordpress/inspection',
                  'laravel-io/inspection']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Fast Inspections"
      :measures="['inspection_execution_time']"
      :projects="['mediawiki/inspection','php-cs-fixer/inspection', 'proxyManager/inspection']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Slow Local Inspections"
      :measures="['local_inspection_execution_time']"
      :projects="['mpdf/localInspection', 'WI_65655/localInspection']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Fast Local Inspections"
      :measures="['local_inspection_execution_time']"
      :projects="['WI_59961/localInspection', 'bitrix/localInspection']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Slow Indexing"
      :measures="['indexing']"
      :projects="['b2c-demo-shop/indexing', 'bitrix/indexing', 'oro/indexing', 'ilias/indexing', 'magento2/indexing', 'drupal8-master-with-plugin/indexing', 
                  'laravel-io/indexing','wordpress/indexing','mediawiki/indexing']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Medium Indexing"
      :measures="['indexing']"
      :projects="['WI_39333/indexing','many_array_access/indexing', 'php-cs-fixer/indexing','many_classes/indexing', 'magento/indexing', 'proxyManager/indexing', 
                  'shopware/indexing', 'dql/indexing', 'tcpdf/indexing', 'WI_51645/indexing']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Fast Indexing"
      :measures="['indexing']"
      :projects="['phpcs/indexing','empty_project/indexing','complex_meta/indexing','broken_phpdoc/indexing','WI_53502/indexing','heredoc/indexing']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Completion"
      :measures="['completion_execution_time']"
      :projects="['many_classes/completion/classes','magento2/completion/function_var','magento2/completion/classes','dql/completion','WI_64694/completion','WI_58919/completion',
                  'WI_58807/completion', 'WI_58306/completion']"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing Time"
      :measures="['typing_total_time']"
      :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing', 
                  'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                  'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing Average Responsiveness"
      :measures="['average_responsiveness_time']"
      :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                  'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                  'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :configurators="configurators"
      :server-configurator="serverConfigurator"
    />
    <GroupLineChart
      label="Typing Responsiveness"
      :measures="['responsiveness_time']"
      :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                  'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                  'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing'
      ]"
      :configurators="configurators"
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

const persistentStateManager = new PersistentStateManager("phpstorm_dashboard", {
  machine: "linux-blade",
  project: [],
  branch: "master",
}, useRouter())

const serverConfigurator = new ServerConfigurator("perfint", "phpstorm")
const branchConfigurator = dimensionConfigurator("branch", serverConfigurator, persistentStateManager, true)
const machineConfigurator = new MachineConfigurator(serverConfigurator, persistentStateManager, [])
const timeRangeConfigurator = new TimeRangeConfigurator(persistentStateManager)
const configurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
]
initDataComponent(configurators)
</script>