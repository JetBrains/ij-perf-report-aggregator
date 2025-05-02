<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfintDev"
    table="phpstorm"
    persistent-id="phpstorm_indexing_dashboard"
    initial-machine="linux-blade-hetzner"
    :with-installer="false"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeedAvg#PHP'"
          :title="'Indexing PHP (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="[
            'b2c-demo-shop/indexing',
            'bitrix/indexing',
            'oro/indexing',
            'ilias/indexing',
            'magento2/indexing',
            'drupal8-master-with-plugin/indexing',
            'laravel-io/indexing',
            'wordpress/indexing',
            'mediawiki/indexing',
            'WI_66681/indexing',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="['akaunting/indexing', 'aggregateStitcher/indexing', 'prestaShop/indexing', 'kunstmaanBundlesCMS/indexing', 'shopware/indexing']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="[
            'WI_39333-5x/indexing',
            'php-cs-fixer/indexing',
            'many_classes/indexing',
            'magento/indexing',
            'proxyManager/indexing',
            'dql/indexing',
            'tcpdf/indexing',
            'WI_51645/indexing',
            'WI_66279-10x/indexing',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Indexing"
          measure="indexingTimeWithoutPauses"
          :projects="['empty_project/indexing', 'complex_meta/indexing', 'WI_53502-10x/indexing', 'many_array_access/indexing-10x']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="[
            'b2c-demo-shop/indexing',
            'bitrix/indexing',
            'oro/indexing',
            'ilias/indexing',
            'magento2/indexing',
            'drupal8-master-with-plugin/indexing',
            'laravel-io/indexing',
            'wordpress/indexing',
            'mediawiki/indexing',
            'WI_66681/indexing',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="['akaunting/indexing', 'aggregateStitcher/indexing', 'prestaShop/indexing', 'kunstmaanBundlesCMS/indexing', 'shopware/indexing']"
        />
      </div>
    </section>

    <section class="flex gap-x-6">
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="[
            'WI_39333-5x/indexing',
            'php-cs-fixer/indexing',
            'many_classes/indexing',
            'magento/indexing',
            'proxyManager/indexing',
            'dql/indexing',
            'tcpdf/indexing',
            'WI_51645/indexing',
          ]"
        />
      </div>
      <div class="flex-1 min-w-0">
        <GroupProjectsChart
          label="Scanning"
          measure="scanningTimeWithoutPauses"
          :projects="['empty_project/indexing', 'complex_meta/indexing', 'WI_53502-10x/indexing', 'many_array_access/indexing-10x', 'WI_66279-10x/indexing']"
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
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
