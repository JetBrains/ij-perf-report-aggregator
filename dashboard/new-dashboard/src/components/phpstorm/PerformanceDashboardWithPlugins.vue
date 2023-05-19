<template>
  <DashboardPage
    db-name="perfint"
    table="phpstormWithPlugins"
    persistent-id="phpstorm_plugins_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section>
      <GroupProjectsChart
        label="Batch Inspections"
        measure="globalInspections"
        :projects="['drupal8-master-with-plugin/inspection', 'magento/inspection', 'wordpress/inspection', 'laravel-io/inspection']"
      />
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Batch Inspections"
          measure="globalInspections"
          :projects="['mediawiki/inspection']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Local Inspections"
          measure="localInspections"
          :projects="['mpdf/localInspection']"
        />
      </div>
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="updatingTime"
          :projects="['drupal8-master-with-plugin/indexing', 'laravel-io/indexing','wordpress/indexing','mediawiki/indexing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="updatingTime"
          :projects="['magento/indexing']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing Time"
          measure="typing"
          :projects="['mpdf/typing', 'mpdf_powersave/typing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing Average Responsiveness"
          measure="test#average_awt_delay"
          :projects="[ 'mpdf/typing', 'mpdf_powersave/typing']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Typing Responsiveness"
          measure="test#max_awt_delay"
          :projects="['mpdf/typing', 'mpdf_powersave/typing']"
        />
      </div>
    </section>
  </DashboardPage>>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>