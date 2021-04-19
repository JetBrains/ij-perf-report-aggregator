<template>
  <AggregationOperatorSelect />
  <el-row :gutter="5">
    <el-col :span="12">
      <BarChartCard
        :measures='[
          "bootstrap_d", "appInitPreparation_d", "appInit_d", "pluginDescriptorLoading_d",
          "appComponentCreation_d", "projectComponentCreation_d",
        ]'
      />
    </el-col>
    <el-col :span="12">
      <BarChartCard
        :measures='["splash_i", "startUpCompleted"]'
      />
    </el-col>
  </el-row>
  <el-row
    :gutter="5"
    style="margin-top: 5px;"
  >
    <!-- todo "moduleLoading_d" when it will be fixed -->
    <el-col :span="12">
      <BarChartCard
        :measures='[
          "appStarter_d",
          "serviceSyncPreloading_d", "serviceAsyncPreloading_d",
          "projectServiceSyncPreloading_d", "projectServiceAsyncPreloading_d",
        ]'
      />
    </el-col>
    <el-col :span="12">
      <BarChartCard
        :measures='[
          "projectDumbAware", "editorRestoring", "editorRestoringTillPaint"
        ]'
      />
    </el-col>
  </el-row>
</template>
<script lang="ts">
import AggregationOperatorSelect from "shared/src/components/AggregationOperatorSelect.vue"
import BarChartCard from "shared/src/components/BarChartCard.vue"
import { defineComponent } from "vue"

export default defineComponent({
  name: "ProgressOverTime",
  components: {BarChartCard, AggregationOperatorSelect},
})
</script>