<template>
  <StartupPage>
    <Divider label="Classloading Files Size" />
    <LineChartCard
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
    />
    <Divider label="Memory" />
    <LineChartCard
      :skip-zero-values="false"
      :measures="['metrics.memory/usedMb', 'metrics.memory/metaspaceMb', 'metrics.memory/maxMb']"
    />
    <Divider label="GC" />
    <div class="grid grid-cols-2 gap-6">
      <LineChartCard
        class="col-span-2"
        :skip-zero-values="false"
        :measures="['metrics.gc/totalHeapUsedMax', 'metrics.gc/freedMemoryByGC']"
      />
      <LineChartCard
        :skip-zero-values="false"
        label="Number of pauses"
        :measures="['metrics.gc/gcPauseCount']"
      />
      <LineChartCard
        :skip-zero-values="false"
        :measures="['metrics.gc/fullGCPause', 'metrics.gc/gcPause']"
      />
    </div>
  </StartupPage>
</template>
<script setup lang="ts">
import LineChartCard from "../charts/StartupLineChart.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
</script>
