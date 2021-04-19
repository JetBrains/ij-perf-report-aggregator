<template>
  <el-divider>
    Bootstrap
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d", "euaShowing_d", "appStarter_d"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='["pluginDescriptorInitV18_d", "appComponentCreation_d", "projectComponentCreation_d"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Class and Resource Loading
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["classLoadingTime", "classLoadingSearchTime", "classLoadingDefineTime"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='["classLoadingCount", "resourceLoadingCount"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Services
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["appComponentCreation_d", "serviceSyncPreloading_d", "serviceAsyncPreloading_d"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :skip-zero-values="false"
        :measures='["projectComponentCreation_d", "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d"]'
      />
    </el-col>
  </el-row>

  <el-divider>
    Post-opening
  </el-divider>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["projectDumbAware", "editorRestoring", "editorRestoringTillPaint"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='["splash_i", "startUpCompleted"]'
      />
    </el-col>
  </el-row>
</template>
<script lang="ts">
import LineChartCard from "shared/src/components/LineChartCard.vue"
import { defineComponent } from "vue"

export default defineComponent({
  name: "Pulse",
  components: {LineChartCard},
})
</script>