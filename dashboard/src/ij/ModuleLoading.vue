<template>
  <div class="grid grid-cols-2 gap-4">
    <LineChartCard
      :measures='["cache loading", "modules loading without cache", "modules loading with cache"]'
    />
    <LineChartCard
      :measures='[
        "project files loading", "serializers creation",
        "applying loaded changes (in queue)", "applying loaded changes"
      ]'
    />
    <LineChartCard
      :measures='["modules instantiation", "libraries instantiation", "loading entities from files", "unloaded modules loading"]'
    />
    <LineChartCard
      :measures='["postponed events sending", "events sending (in queue)", "events sending", "firing modules_added event"]'
    />
  </div>
</template>
<script lang="ts">
import LineChartCard from "../shared/components/LineChartCard.vue"
import { defineComponent } from "vue"

export default defineComponent({
  name: "ModuleLoading",
  components: {LineChartCard},
})
</script>