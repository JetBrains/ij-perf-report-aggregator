<template>
  <JbrBenchmarkDashboard
    benchmark="J2DBench"
    :groups="[{ title: null, scenarios }]"
    value-unit="pixelsPerSecond"
    better-direction="higher"
    description="Pixels drawn per second, as J2DBench averages them over its runs, higher is better."
    :default-configurations="['macOS15aarch64Metal', 'Ubuntu2404x86_64MM_Wayland', 'Windows11x86_64']"
  />
</template>

<script setup lang="ts">
import JbrBenchmarkDashboard from "./JbrBenchmarkDashboard.vue"

const scenarios = ["poly250", "poly250-rand_col", "poly250-AA-rand_col"]
</script>
