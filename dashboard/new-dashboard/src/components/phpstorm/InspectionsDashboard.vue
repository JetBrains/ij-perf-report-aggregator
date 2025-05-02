<template>
  <DashboardPage
    db-name="perfintDev"
    table="phpstorm"
    persistent-id="phpstorm_inspections_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section>
      <GroupProjectsChart
        label="Batch Inspections"
        measure="globalInspections"
        :projects="[
          'drupal8-master-with-plugin/inspection',
          'shopware/inspection',
          'b2c-demo-shop/inspection',
          'magento/inspection',
          'wordpress/inspection',
          'laravel-io/inspection',
        ]"
      />
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Batch Inspections"
          measure="globalInspections"
          :projects="['mediawiki/inspection', 'php-cs-fixer/inspection', 'proxyManager/inspection']"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Batch Inspections"
          measure="globalInspections"
          :projects="['akaunting/inspection', 'aggregateStitcher/inspection', 'prestaShop/inspection', 'kunstmaanBundlesCMS/inspection']"
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

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Code Vision (PhpReferencesCodeVisionProvider)"
          measure="PhpReferencesCodeVisionProvider"
          :projects="['mpdf/localInspection', 'WI_65655/localInspection', 'laravel-io/localInspection/HasAuthor', 'laravel-io/localInspection/Tag']"
        />
      </div>
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
