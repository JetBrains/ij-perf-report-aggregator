<template>
  <IntelliJPulse>
    <Divider label="Memory" />
    <LineChartCard
      :skip-zero-values="false"
      :measures='["metrics.memory/usedMb", "metrics.memory/metaspaceMb","metrics.memory/maxMb"]'
    />
    <Divider label="GC" />
    <div class="grid grid-cols-2 gap-6">
      <LineChartCard
        class="col-span-2"
        :skip-zero-values="false"
        :measures='["metrics.gc/totalHeapUsedMax", "metrics.gc/freedMemoryByGC"]'
      />
      <LineChartCard
        :skip-zero-values="false"
        label="Number of pauses"
        :measures='["metrics.gc/gcPauseCount"]'
      />
      <LineChartCard
        :skip-zero-values="false"
        :measures='["metrics.gc/fullGCPause","metrics.gc/gcPause"]'
      />
    </div>
  </IntelliJPulse>
</template>
<script setup lang="ts">
import LineChartCard from "../charts/LineChartCard.vue"
import Divider from "../common/Divider.vue"
import IntelliJPulse from "./IntelliJPulse.vue"

</script>