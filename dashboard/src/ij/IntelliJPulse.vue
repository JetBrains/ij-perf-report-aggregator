<template>
  <ChartRow title="Bootstrap">
    <LineChartCard
      :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
    />
    <LineChartCard
      :measures='["pluginDescriptorInitV18_d", "appComponentCreation_d", "projectComponentCreation_d"]'
    />
  </ChartRow>

  <ChartRow title="Class and Resource Loading">
    <LineChartCard
      :measures='["classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime"]'
    />
    <LineChartCard
      :measures='["classLoadingCount", "resourceLoadingCount"]'
    />
  </ChartRow>

  <ChartRow title="Services">
    <LineChartCard
      :measures='["appComponentCreation_d", "serviceSyncPreloading_d", "serviceAsyncPreloading_d"]'
    />
    <LineChartCard
      :skip-zero-values="false"
      :measures='["projectComponentCreation_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d"]'
    />
  </ChartRow>

  <ChartRow title="Post-opening">
    <LineChartCard
      :measures='["projectDumbAware", "editorRestoring", "editorRestoringTillPaint"]'
    />
    <LineChartCard
      :measures='["splash_i", "startUpCompleted"]'
    />
  </ChartRow>
</template>
<script setup lang="ts">
import LineChartCard from "../shared/components/LineChartCard.vue"
import ChartRow from "./ChartRow.vue"
</script>