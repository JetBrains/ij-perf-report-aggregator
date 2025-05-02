<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="phpstorm"
    persistent-id="phpstorm_code_editing_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-6">
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
        label="Completion"
        measure="completion"
        :projects="[
          'many_classes/completion/classes',
          'magento2/completion/function_var',
          'magento2/completion/function_stlr',
          'magento2/completion/classes',
          'dql/completion',
          'WI_64694/completion',
          'WI_58919/completion',
          'WI_58807/completion',
          'WI_58306/completion',
        ]"
      />
    </section>

    <section>
      <GroupProjectsChart
        label="PHP Typing Time"
        measure="typing"
        :projects="[
          'WI_29056/typing',
          'WI_41934/typing',
          'WI_44525/typing',
          'WI_60709/typing',
          'bitrix/typing',
          'heredoc/typing',
          'html_in_fragment/typing',
          'html_in_fragment_powersave/typing',
          'html_in_literal/typing',
          'html_in_literal_powersave/typing',
          'large_method_phpdoc/typing',
          'large_phpdoc/typing',
          'large_phpdoc_comment/typing',
          'lots_phpdoc_methods/typing',
          'mpdf/typing',
          'mpdf_powersave/typing',
        ]"
      />
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="PHP Typing Average Responsiveness"
          measure="test#average_awt_delay"
          :projects="[
            'WI_29056/typing',
            'WI_41934/typing',
            'WI_44525/typing',
            'WI_60709/typing',
            'bitrix/typing',
            'heredoc/typing',
            'html_in_fragment/typing',
            'html_in_fragment_powersave/typing',
            'html_in_literal/typing',
            'html_in_literal_powersave/typing',
            'large_method_phpdoc/typing',
            'large_phpdoc/typing',
            'large_phpdoc_comment/typing',
            'lots_phpdoc_methods/typing',
            'mpdf/typing',
            'mpdf_powersave/typing',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="PHP Typing Responsiveness"
          measure="test#max_awt_delay"
          :projects="[
            'WI_29056/typing',
            'WI_41934/typing',
            'WI_44525/typing',
            'WI_60709/typing',
            'bitrix/typing',
            'heredoc/typing',
            'html_in_fragment/typing',
            'html_in_fragment_powersave/typing',
            'html_in_literal/typing',
            'html_in_literal_powersave/typing',
            'large_method_phpdoc/typing',
            'large_phpdoc/typing',
            'large_phpdoc_comment/typing',
            'lots_phpdoc_methods/typing',
            'mpdf/typing',
            'mpdf_powersave/typing',
          ]"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Blade Typing Time"
          measure="typing"
          :projects="[
            'blade_in_php_fragment_large_file/typing',
            'blade_in_blade_fragment_large_file/typing',
            'blade_new_line_large_file/typing',
            'blade_in_blade_fragment_laravel/typing',
            'blade_in_php_fragment_laravel/typing',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Blade Average Responsiveness"
          measure="test#average_awt_delay"
          :projects="[
            'blade_in_php_fragment_large_file/typing',
            'blade_in_blade_fragment_large_file/typing',
            'blade_new_line_large_file/typing',
            'blade_in_blade_fragment_laravel/typing',
            'blade_in_php_fragment_laravel/typing',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Blade Responsiveness"
          measure="test#max_awt_delay"
          :projects="[
            'blade_in_php_fragment_large_file/typing',
            'blade_in_blade_fragment_large_file/typing',
            'blade_new_line_large_file/typing',
            'blade_in_blade_fragment_laravel/typing',
            'blade_in_php_fragment_laravel/typing',
          ]"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
import { DataQuery, DataQueryExecutorConfiguration } from "../common/dataQuery"

const typingOnlyConfigurator = {
  configureQuery(query: DataQuery, _configuration: DataQueryExecutorConfiguration): boolean {
    query.addFilter({ f: "project", v: "%typing", o: "like" })
    return true
  },
  createObservable() {
    return null
  },
}
</script>
