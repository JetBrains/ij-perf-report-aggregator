<template>
  <StartupPage>
    <template #default="{ configurators }">
      <Divider label="Classloading Files Size" />
      <PerformanceLineChart
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
      />
      <Divider label="Memory" />
      <PerformanceLineChart
        :skip-zero-values="false"
        :measures="['metrics.memory/usedMb', 'metrics.memory/metaspaceMb', 'metrics.memory/maxMb']"
        :configurators="configurators"
        title=""
      />
      <Divider label="GC" />
      <div class="grid grid-cols-2 gap-6">
        <PerformanceLineChart
          class="col-span-2"
          :skip-zero-values="false"
          :measures="['metrics.gc/totalHeapUsedMax', 'metrics.gc/freedMemoryByGC']"
          :configurators="configurators"
          title=""
        />
        <PerformanceLineChart
          :skip-zero-values="false"
          title="Number of pauses"
          :measures="['metrics.gc/gcPauseCount']"
          :configurators="configurators"
          :with-measure-name="true"
        />
        <PerformanceLineChart
          :skip-zero-values="false"
          :measures="['metrics.gc/fullGCPause', 'metrics.gc/gcPause']"
          :configurators="configurators"
          title=""
        />
      </div>
    </template>
  </StartupPage>
</template>
<script setup lang="ts">
import PerformanceLineChart from "../charts/PerformanceLineChart.vue"
import Divider from "../common/Divider.vue"
import StartupPage from "./StartupPage.vue"
</script>
