<template>
  <ElDivider>
    Bootstrap
  </ElDivider>
  <div class="grid grid-cols-2 gap-4">
    <LineChartCard
      :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
    />
    <LineChartCard
      :measures='["pluginDescriptorInitV18_d", "appComponentCreation_d", "projectComponentCreation_d"]'
    />
  </div>

  <ElDivider>
    Class and Resource Loading
  </ElDivider>
  <div class="grid grid-cols-2 gap-4">
    <LineChartCard
      :measures='["classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime"]'
    />
    <LineChartCard
      :measures='["classLoadingCount", "resourceLoadingCount"]'
    />
  </div>

  <el-divider>
    Services
  </el-divider>
  <div class="grid grid-cols-2 gap-4">
    <LineChartCard
      :measures='["appComponentCreation_d", "serviceSyncPreloading_d", "serviceAsyncPreloading_d"]'
    />
    <LineChartCard
      :skip-zero-values="false"
      :measures='["projectComponentCreation_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d"]'
    />
  </div>

  <el-divider>
    Post-opening
  </el-divider>
  <div class="grid grid-cols-2 gap-4">
    <LineChartCard
      :measures='["projectDumbAware", "editorRestoring", "editorRestoringTillPaint"]'
    />
    <LineChartCard
      :measures='["splash_i", "startUpCompleted"]'
    />
  </div>
</template>
<script lang="ts">
import LineChartCard from "shared/src/components/LineChartCard.vue"
import { defineComponent } from "vue"

export default defineComponent({
  name: "Pulse",
  components: {LineChartCard},
})
</script>