<template>
  <StartupPage>
    <template #default="{ configurators }">
      <div class="grid grid-cols-2 gap-6 pt-5">
        <LineChart
          :measures="['cache loading', 'modules loading without cache', 'modules loading with cache', 'workspace file index initialization']"
          :configurators="configurators"
          title=""
        />
        <LineChart
          :measures="['project files loading', 'serializers creation', 'applying loaded changes (in queue)', 'applying loaded changes']"
          :configurators="configurators"
          title=""
        />
        <LineChart
          :measures="['modules instantiation', 'libraries instantiation', 'loading entities from files', 'unloaded modules loading']"
          :configurators="configurators"
          title=""
        />
        <LineChart
          :measures="['postponed events sending', 'events sending (in queue)', 'events sending', 'firing modules_added event']"
          :configurators="configurators"
          title=""
        />
      </div>
    </template>
  </StartupPage>
</template>
<script setup lang="ts">
import LineChart from "../charts/LineChart.vue"
import StartupPage from "./StartupPage.vue"
</script>
