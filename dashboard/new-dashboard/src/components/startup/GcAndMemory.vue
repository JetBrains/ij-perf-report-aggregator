<template>
  <StartupPage>
    <template #default="{ configurators }">
      <Divider label="Classloading Files Size" />
      <LineChart
        :skip-zero-values="false"
        :measures="[
          'metrics.classLoadingMetrics/totalSizeKb',
          'metrics.classLoadingMetrics/companionSizeKb',
          'metrics.classLoadingMetrics/lambdaSizeKb',
          'metrics.classLoadingMetrics/actionSizeKb',
          'metrics.classLoadingMetrics/inspectionSizeKb',
          'metrics.classLoadingMetrics/kotlinReflectSizeKb',
          'metrics.classLoadingMetrics/androidSizeKb',
        ]"
        :configurators="configurators"
        title=""
        tooltip-trigger="axis"
      />
      <Divider label="Memory" />
      <LineChart
        :skip-zero-values="false"
        :measures="['metrics.memory/usedMb', 'metrics.memory/metaspaceMb', 'metrics.memory/maxMb']"
        :configurators="configurators"
        title=""
        tooltip-trigger="axis"
      />
      <Divider label="GC" />
      <div class="grid grid-cols-2 gap-6">
        <LineChart
          class="col-span-2"
          :skip-zero-values="false"
          :measures="['metrics.gc/totalHeapUsedMax', 'metrics.gc/freedMemoryByGC']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
        <LineChart
          :skip-zero-values="false"
          title="Number of pauses"
          :measures="['metrics.gc/gcPauseCount']"
          :configurators="configurators"
          :with-measure-name="true"
          tooltip-trigger="axis"
        />
        <LineChart
          :skip-zero-values="false"
          :measures="['metrics.gc/fullGCPause', 'metrics.gc/gcPause']"
          :configurators="configurators"
          title=""
          tooltip-trigger="axis"
        />
      </div>
    </template>
  </StartupPage>
</template>
<script setup lang="ts">
import LineChart from "../charts/LineChart.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
</script>
