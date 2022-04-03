<template>
  <Divider align="center">
    Bootstrap
  </Divider>
  <div :class="$style.twoCard">
    <LineChartCard
      :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
    />
    <LineChartCard
      :measures='["pluginDescriptorInitV18_d", "appComponentCreation_d", "projectComponentCreation_d"]'
    />
  </div>

  <Divider align="center">
    Class and Resource Loading
  </Divider>
  <div :class="$style.twoCard">
    <LineChartCard
      :measures='["classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime"]'
    />
    <LineChartCard
      :measures='["classLoadingCount", "resourceLoadingCount"]'
    />
  </div>

  <Divider align="center">
    Services
  </Divider>
  <div :class="$style.twoCard">
    <LineChartCard
      :measures='["appComponentCreation_d", "serviceSyncPreloading_d", "serviceAsyncPreloading_d"]'
    />
    <LineChartCard
      :skip-zero-values="false"
      :measures='["projectComponentCreation_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d"]'
    />
  </div>

  <Divider align="center">
    Post-opening
  </Divider>
  <div :class="$style.twoCard">
    <LineChartCard
      :measures='["projectDumbAware", "editorRestoring", "editorRestoringTillPaint"]'
    />
    <LineChartCard
      :measures='["splash_i", "startUpCompleted"]'
    />
  </div>
</template>
<script setup lang="ts">
import LineChartCard from "shared/src/components/LineChartCard.vue"
</script>
<style module>
.twoCard {
  @apply grid grid-cols-1 lg:grid-cols-2 gap-4
}
</style>