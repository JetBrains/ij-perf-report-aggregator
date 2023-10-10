<template>
  <DashboardPage
    v-slot="{ averagesConfigurators }"
    db-name="perfint"
    table="webstorm"
    persistent-id="webstorm_dashboard"
    initial-machine="linux-blade-hetzner"
  >
    <section class="flex gap-6">
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#JavaScript'"
          :title="'Indexing JavaScript (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
      <div class="flex-1 min-w-0">
        <AggregationChart
          :configurators="averagesConfigurators"
          :aggregated-measure="'processingSpeed#TypeScript'"
          :title="'Indexing TypeScript (kB/s)'"
          :chart-color="'#219653'"
          :value-unit="'counter'"
        />
      </div>
    </section>
    <section>
      <GroupProjectsChart
        label="Local Inspections"
        measure="localInspections"
        :projects="[
          'aws_cdk/localInspection/logging',
          'WEB_5976/localInspection/react_mui',
          'toh-pt6/localInspection/hero.service.ts',
          'toh-pt6/localInspection/heroes.component.html',
          'vue3-admin-vite/localInspection/index.vue',
          'eslint-plugin-jest/localInspection/misc.ts',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Indexing"
        measure="indexingTimeWithoutPauses"
        :projects="[
          'aws_cdk/indexing',
          'vkui/indexing',
          'ring-ui/indexing',
          'angular-toh-pt6/indexing',
          'angular/indexing',
          'vue3-admin-vite/indexing',
          'eslint-plugin-jest/indexing',
          'dxos/indexing',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Scanning"
        measure="scanningTimeWithoutPauses"
        :projects="[
          'aws_cdk/indexing',
          'vkui/indexing',
          'ring-ui/indexing',
          'angular-toh-pt6/indexing',
          'angular/indexing',
          'vue3-admin-vite/indexing',
          'eslint-plugin-jest/indexing',
          'dxos/indexing',
        ]"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Completion"
        measure="completion"
        :projects="['vkui/completion', 'ring-ui/completion', 'vue3-admin-vite/completion', 'toh-pt6/completion']"
      />
    </section>
    <section>
      <GroupProjectsChart
        label="Code Vision (JSReferencesCodeVisionProvider)"
        measure="JSReferencesCodeVisionProvider"
        :projects="[
          'aws_cdk/localInspection/logging',
          'WEB_5976/localInspection/react_mui',
          'toh-pt6/localInspection/hero.service.ts',
          'vue3-admin-vite/localInspection/index.vue',
          'eslint-plugin-jest/localInspection/misc.ts',
        ]"
      />
    </section>
  </DashboardPage>
</template>

<script setup lang="ts">
import AggregationChart from "../charts/AggregationChart.vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import DashboardPage from "../common/DashboardPage.vue"
</script>
