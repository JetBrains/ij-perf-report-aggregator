<template>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["module cache loading", "modules loading without cache", "modules loading with cache"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='[
          "project files loading", "serializers creation",
          "project model changes saving (in queue)", "project model changes saving"
        ]'
      />
    </el-col>
  </el-row>
  <el-row :gutter="5">
    <el-col :span="12">
      <LineChartCard
        :measures='["modules instantiation", "project libraries loading", "entities loading", "unloaded modules loading"]'
      />
    </el-col>
    <el-col :span="12">
      <LineChartCard
        :measures='["events sending after modules are loaded", "events sending (in queue)", "events sending", "events modules added"]'
      />
    </el-col>
  </el-row>
</template>
<script lang="ts">
import LineChartCard from "shared/src/components/LineChartCard.vue"
import { defineComponent } from "vue"

export default defineComponent({
  name: "ModuleLoading",
  components: {LineChartCard},
})
</script>