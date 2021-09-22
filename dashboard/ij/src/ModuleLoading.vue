<template>
  <div class="grid grid-cols-2 gap-4">
    <LineChartCard
      :measures='["module cache loading", "modules loading without cache", "modules loading with cache"]'
    />
    <LineChartCard
      :measures='[
        "project files loading", "serializers creation",
        "project model changes saving (in queue)", "project model changes saving"
      ]'
    />
    <LineChartCard
      :measures='["modules instantiation", "project libraries loading", "entities loading", "unloaded modules loading"]'
    />
    <LineChartCard
      :measures='["events sending after modules are loaded", "events sending (in queue)", "events sending", "events modules added"]'
    />
  </div>
</template>
<script lang="ts">
import LineChartCard from "shared/src/components/LineChartCard.vue"
import { defineComponent } from "vue"

export default defineComponent({
  name: "ModuleLoading",
  components: {LineChartCard},
})
</script>