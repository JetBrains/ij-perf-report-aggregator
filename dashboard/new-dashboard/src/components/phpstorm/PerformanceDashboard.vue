<template>
  <div class="flex flex-col gap-5">
    <Toolbar class="customToolbar">
      <template #start>
        <TimeRangeSelect
          :ranges="TimeRangeConfigurator.timeRanges"
          :value="timeRangeConfigurator.value.value"
          :on-change="onChangeRange"
        >
          <template #icon>
            <CalendarIcon class="w-4 h-4 text-gray-500" />
          </template>
        </TimeRangeSelect>
        <BranchSelect
          :branch-configurator="branchConfigurator"
          :release-configurator="releaseConfigurator"
          :triggered-by-configurator="triggeredByConfigurator"
        />
        <DimensionHierarchicalSelect
          label="Machine"
          :dimension="machineConfigurator"
        >
          <template #icon>
            <ComputerDesktopIcon class="w-4 h-4 text-gray-500" />
          </template>
        </DimensionHierarchicalSelect>
      </template>
    </Toolbar>

    <main class="flex">
      <div
        ref="container"
        class="flex flex-1 flex-col gap-6 overflow-hidden"
      >
        <section class="flex gap-6">
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'processingSpeed#PHP'"
              :title="'Indexing PHP (kB/s)'"
              :chart-color="'#219653'"
              :value-unit="'counter'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="averagesConfigurators"
              :aggregated-measure="'completion\_%'"
              :is-like="true"
              :title="'Completion'"
            />
          </div>
          <div class="flex-1">
            <AggregationChart
              :configurators="[...averagesConfigurators, typingOnlyConfigurator]"
              :aggregated-measure="'test#average_awt_delay'"
              :title="'UI responsiveness during typing'"
              :chart-color="'#F2994A'"
            />
          </div>
        </section>
        <section>
          <GroupProjectsChart
            label="Batch Inspections"
            measure="globalInspections"
            :projects="['drupal8-master-with-plugin/inspection', 'shopware/inspection', 'b2c-demo-shop/inspection', 'magento/inspection', 'wordpress/inspection',
                        'laravel-io/inspection']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1">
            <GroupProjectsChart
              label="Batch Inspections"
              measure="globalInspections"
              :projects="['mediawiki/inspection','php-cs-fixer/inspection', 'proxyManager/inspection']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1">
            <GroupProjectsChart
              label="Batch Inspections"
              measure="globalInspections"
              :projects="['akaunting/inspection','aggregateStitcher/inspection', 'prestaShop/inspection', 'kunstmaanBundlesCMS/inspection']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>

        <section class="flex gap-x-6">
          <div class="flex-1">
            <GroupProjectsChart
              label="Local Inspections"
              measure="localInspections"
              :projects="['mpdf/localInspection', 'WI_65655/localInspection']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1">
            <GroupProjectsChart
              label="Local Inspections"
              measure="localInspections"
              :projects="['WI_59961/localInspection', 'bitrix/localInspection', 'WI_65893/localInspection']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="indexing"
            :projects="['b2c-demo-shop/indexing', 'bitrix/indexing', 'oro/indexing', 'ilias/indexing', 'magento2/indexing', 'drupal8-master-with-plugin/indexing', 
                        'laravel-io/indexing','wordpress/indexing','mediawiki/indexing', 'WI_66681/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="indexing"
            :projects="['akaunting/indexing', 'aggregateStitcher/indexing', 'prestaShop/indexing', 'kunstmaanBundlesCMS/indexing', 'shopware/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="indexing"
            :projects="['WI_39333/indexing', 'php-cs-fixer/indexing','many_classes/indexing', 'magento/indexing', 'proxyManager/indexing',
                        'dql/indexing', 'tcpdf/indexing', 'WI_51645/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Indexing"
            measure="indexing"
            :projects="['phpcs/indexing','empty_project/indexing','complex_meta/indexing', 'WI_53502/indexing','heredoc/indexing', 'many_array_access/indexing',
                        'WI_66279/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="Completion"
            measure="completion"
            :projects="['many_classes/completion/classes','magento2/completion/function_var', 'magento2/completion/function_stlr', 'magento2/completion/classes','dql/completion',
                        'WI_64694/completion','WI_58919/completion', 'WI_58807/completion', 'WI_58306/completion']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section>
          <GroupProjectsChart
            label="PHP Typing Time"
            measure="typing"
            :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing', 
                        'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                        'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1">
            <GroupProjectsChart
              label="PHP SearchEverywhere Class"
              measure="searchEverywhere_class"
              :projects="['bitrix/go-to-class/BCCo', 'magento2/go-to-class/MaAdMUser']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1">
            <GroupProjectsChart
              label="Code Vision (PhpReferencesCodeVisionProvider)"
              measure="PhpReferencesCodeVisionProvider"
              :projects="['mpdf/localInspection', 'WI_65655/localInspection', 'laravel-io/localInspection/HasAuthor', 'laravel-io/localInspection/Tag']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1">
            <GroupProjectsChart
              label="PHP Typing Average Responsiveness"
              measure="test#average_awt_delay"
              :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                          'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                          'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1">
            <GroupProjectsChart
              label="PHP Typing Responsiveness"
              measure="test#max_awt_delay"
              :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                          'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                          'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>
        <section class="flex gap-x-6">
          <div class="flex-1">
            <GroupProjectsChart
              label="Blade Typing Time"
              measure="typing"
              :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                          'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1">
            <GroupProjectsChart
              label="Blade Average Responsiveness"
              measure="test#average_awt_delay"
              :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                          'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
          <div class="flex-1">
            <GroupProjectsChart
              label="Blade Responsiveness"
              measure="test#max_awt_delay"
              :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                          'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
              :server-configurator="serverConfigurator"
              :configurators="dashboardConfigurators"
            />
          </div>
        </section>
        <section>
          <GroupProjectsChart
            label="Index size"
            measure="indexSize"
            :projects="['akaunting/indexing', 'aggregateStitcher/indexing', 'prestaShop/indexing', 'kunstmaanBundlesCMS/indexing']"
            :server-configurator="serverConfigurator"
            :configurators="dashboardConfigurators"
          />
        </section>
      </div>
      <InfoSidebar />
    </main>
  </div>
</template>

<script setup lang="ts">
import { PersistentStateManager } from "shared/src/PersistentStateManager"
import DimensionHierarchicalSelect from "shared/src/components/DimensionHierarchicalSelect.vue"
import { createBranchConfigurator } from "shared/src/configurators/BranchConfigurator"
import { MachineConfigurator } from "shared/src/configurators/MachineConfigurator"
import { privateBuildConfigurator } from "shared/src/configurators/PrivateBuildConfigurator"
import { ReleaseNightlyConfigurator } from "shared/src/configurators/ReleaseNightlyConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { TimeRangeConfigurator } from "shared/src/configurators/TimeRangeConfigurator"
import { DataQuery, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import { provideReportUrlProvider } from "shared/src/lineChartTooltipLinkProvider"
import { provide, ref } from "vue"
import { useRouter } from "vue-router"
import { containerKey, sidebarVmKey } from "../../shared/keys"
import InfoSidebar from "../InfoSidebar.vue"
import { InfoSidebarVmImpl } from "../InfoSidebarVm"
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import BranchSelect from "../common/BranchSelect.vue"
import TimeRangeSelect from "../common/TimeRangeSelect.vue"

provideReportUrlProvider()

const dbName = "perfint"
const dbTable = "phpstorm"
const initialMachine = "linux-blade-hetzner"
const container = ref<HTMLElement>()
const router = useRouter()
const sidebarVm = new InfoSidebarVmImpl()

provide(containerKey, container)
provide(sidebarVmKey, sidebarVm)

const serverConfigurator = new ServerConfigurator(dbName, dbTable)
const persistenceForDashboard = new PersistentStateManager("phpstorm_dashboard", {
  machine: initialMachine,
  project: [],
  branch: "master",
}, router)

const timeRangeConfigurator = new TimeRangeConfigurator(persistenceForDashboard)

const branchConfigurator = createBranchConfigurator(serverConfigurator, persistenceForDashboard, [timeRangeConfigurator])
const machineConfigurator = new MachineConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [timeRangeConfigurator, branchConfigurator],
)
const releaseConfigurator = new ReleaseNightlyConfigurator(persistenceForDashboard)
const triggeredByConfigurator = privateBuildConfigurator(
  serverConfigurator,
  persistenceForDashboard,
  [branchConfigurator, timeRangeConfigurator],
)

const averagesConfigurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
]

const dashboardConfigurators = [
  serverConfigurator,
  branchConfigurator,
  machineConfigurator,
  timeRangeConfigurator,
  releaseConfigurator,
  triggeredByConfigurator,
]

const typingOnlyConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({f: "project", v: "%typing", o: "like"})
    return true
  },
  createObservable() {
    return null
  },
}

function onChangeRange(value: string) {
  timeRangeConfigurator.value.value = value
}
</script>

<style>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}
</style>