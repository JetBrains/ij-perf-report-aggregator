<template>
  <StartupPage>
    <div class="grid grid-cols-2 gap-6 pt-5">
      <LineChartCard
        :measures='["cache loading", "modules loading without cache", "modules loading with cache", "workspace file index initialization"]'
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
  </StartupPage>
</template>
<script setup lang="ts">
import LineChartCard from "../charts/LineChartCard.vue"
import StartupPage from "./StartupPage.vue"
</script>