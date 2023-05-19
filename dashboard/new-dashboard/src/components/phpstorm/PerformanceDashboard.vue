<template>
  <DashboardPage
    v-slot="{averagesConfigurators}"
    db-name="perfint"
    table="phpstorm"
    persistent-id="phpstorm_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#PHP'"
          :title="'Indexing PHP (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'completion\_%'"
          :is-like="true"
          :title="'Completion'"
        />
      </div>
      <div class="flex-1 min-w-0">
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
      />
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Batch Inspections"
          measure="globalInspections"
          :projects="['mediawiki/inspection','php-cs-fixer/inspection', 'proxyManager/inspection']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Batch Inspections"
          measure="globalInspections"
          :projects="['akaunting/inspection','aggregateStitcher/inspection', 'prestaShop/inspection', 'kunstmaanBundlesCMS/inspection']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local Inspections"
          measure="localInspections"
          :projects="['mpdf/localInspection', 'WI_65655/localInspection']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local Inspections"
          measure="localInspections"
          :projects="['WI_59961/localInspection', 'bitrix/localInspection', 'WI_65893/localInspection']"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        measure="updatingTime"
        :projects="['b2c-demo-shop/indexing', 'bitrix/indexing', 'oro/indexing', 'ilias/indexing', 'magento2/indexing', 'drupal8-master-with-plugin/indexing',
                    'laravel-io/indexing','wordpress/indexing','mediawiki/indexing', 'WI_66681/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        measure="updatingTime"
        :projects="['akaunting/indexing', 'aggregateStitcher/indexing', 'prestaShop/indexing', 'kunstmaanBundlesCMS/indexing', 'shopware/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        measure="updatingTime"
        :projects="['WI_39333-5x/indexing', 'php-cs-fixer/indexing','many_classes/indexing', 'magento/indexing', 'proxyManager/indexing',
                    'dql/indexing', 'tcpdf/indexing', 'WI_51645/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        measure="updatingTime"
        :projects="['empty_project/indexing','complex_meta/indexing', 'WI_53502-10x/indexing', 'many_array_access/indexing-10x', 'WI_66279-10x/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion"
        :projects="['many_classes/completion/classes','magento2/completion/function_var', 'magento2/completion/function_stlr', 'magento2/completion/classes','dql/completion',
                    'WI_64694/completion','WI_58919/completion', 'WI_58807/completion', 'WI_58306/completion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="PHP Typing Time"
        measure="typing"
        :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                    'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                    'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing']"
      />
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="PHP SearchEverywhere Class"
          measure="searchEverywhere_class"
          :projects="['bitrix/go-to-class/BCCo', 'magento2/go-to-class/MaAdMUser']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Code Vision (PhpReferencesCodeVisionProvider)"
          measure="PhpReferencesCodeVisionProvider"
          :projects="['mpdf/localInspection', 'WI_65655/localInspection', 'laravel-io/localInspection/HasAuthor', 'laravel-io/localInspection/Tag']"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="PHP Typing Average Responsiveness"
          measure="test#average_awt_delay"
          :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                      'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                      'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="PHP Typing Responsiveness"
          measure="test#max_awt_delay"
          :projects="['WI_29056/typing', 'WI_41934/typing', 'WI_44525/typing', 'WI_60709/typing', 'bitrix/typing', 'heredoc/typing', 'html_in_fragment/typing',
                      'html_in_fragment_powersave/typing', 'html_in_literal/typing', 'html_in_literal_powersave/typing', 'large_method_phpdoc/typing', 'large_phpdoc/typing',
                      'large_phpdoc_comment/typing', 'lots_phpdoc_methods/typing', 'mpdf/typing', 'mpdf_powersave/typing']"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Blade Typing Time"
          measure="typing"
          :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                      'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Blade Average Responsiveness"
          measure="test#average_awt_delay"
          :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                      'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Blade Responsiveness"
          measure="test#max_awt_delay"
          :projects="['blade_in_php_fragment_large_file/typing', 'blade_in_blade_fragment_large_file/typing', 'blade_new_line_large_file/typing',
                      'blade_in_blade_fragment_laravel/typing', 'blade_in_php_fragment_laravel/typing']"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Index size"
        measure="indexSize"
        :projects="['akaunting/indexing', 'aggregateStitcher/indexing', 'prestaShop/indexing', 'kunstmaanBundlesCMS/indexing']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Inline Rename"
        measure="startInlineRename"
        :projects="['mpdf/inlineRename']"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import { DataQuery, DataQueryExecutorConfiguration } from "shared/src/dataQuery"
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"

const typingOnlyConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({f: "project", v: "%typing", o: "like"})
    return true
  },
  createObservable() {
    return null
  },
}
</script>