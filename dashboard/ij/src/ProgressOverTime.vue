<template>
  <AggregationOperatorSelect class="py-2" />
  <div class="grid grid-cols-2 gap-4">
    <BarChartCard
      :measures='[
        "bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d",
      ]'
    />
    <BarChartCard
      :measures='["splash_i", "startUpCompleted"]'
    />

    <!-- todo "moduleLoading_d" when it will be fixed -->
    <BarChartCard
      :measures='[
        "appStarter_d",
        "serviceSyncPreloading_d", "serviceAsyncPreloading_d",
        "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d",
      ]'
    />
    <BarChartCard
      :measures='[
        "projectDumbAware", "editorRestoring", "editorRestoringTillPaint"
      ]'
    />
  </div>
</template>
<script setup lang="ts">
import AggregationOperatorSelect from "shared/src/components/AggregationOperatorSelect.vue"
import BarChartCard from "shared/src/components/BarChartCard.vue"</script>