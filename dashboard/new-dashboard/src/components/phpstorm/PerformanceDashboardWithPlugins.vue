<template>
  <DashboardPage
    db-name="perfint"
    table="phpstormWithPlugins"
    persistent-id="phpstorm_plugins_dashboard"
    initial-machine="linux-blade-hetzner"
    branch="master"
  >
    <section>
      <GroupProjectsChart
        label="Batch Inspections"
        measure="globalInspections"
        :projects="[
          'drupal8-master-with-plugin/inspection',
          'drupal8-master-with-plugin/inspectionWithOldPlugins',
          'magento/inspection',
          'magento/inspectionWithOldPlugins',
          'wordpress/inspection',
          'wordpress/inspectionWithOldPlugins',
          'laravel-io/inspection',
          'laravel-io/inspectionWithOldPlugins',
        ]"
      />
    </section>
    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Batch Inspections"
          measure="globalInspections"
          :projects="['mediawiki/inspection', 'mediawiki/inspectionWithOldPlugins']"
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
          :projects="[
            'drupal8-master-with-plugin/indexing',
            'drupal8-master-with-plugin/indexingWithOldPlugins',
            'laravel-io/indexing',
            'laravel-io/indexingWithOldPlugins',
            'wordpress/indexing',
            'wordpress/indexingWithOldPlugins',
            'mediawiki/indexing',
            'mediawiki/indexingWithOldPlugins',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="updatingTime"
          :projects="['magento/indexing', 'magento/indexingWithOldPlugins']"
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
          :projects="['mpdf/typing', 'mpdf_powersave/typing']"
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
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
