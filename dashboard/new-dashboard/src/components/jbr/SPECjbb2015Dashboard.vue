<template>
  <JbrBenchmarkDashboard
    benchmark="SPECjbb2015"
    :groups="[{ title: null, scenarios }]"
    value-unit="jOPS"
    better-direction="higher"
    description="SPECjbb2015's max-jOPS (throughput) and critical-jOPS (throughput under response time limits), higher is better."
    :default-configurations="['macOS15aarch64', 'Ubuntu2404x86_64', 'Windows11x86_64']"
  />
</template>

<script setup lang="ts">
import JbrBenchmarkDashboard from "./JbrBenchmarkDashboard.vue"

const scenarios = ["jbb2015.result.metric.max-jOPS", "jbb2015.result.metric.critical-jOPS"]
</script>
