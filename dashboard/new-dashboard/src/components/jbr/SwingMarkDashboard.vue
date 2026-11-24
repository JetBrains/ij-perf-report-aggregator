<template>
  <JbrBenchmarkDashboard
    benchmark="SwingMark"
    :groups="[{ title: null, scenarios }]"
    description="SwingMark's score: the total milliseconds of its run, lower is better."
    :default-configurations="['macOS15aarch64Metal', 'Ubuntu2404x86_64MM_Wayland', 'Windows11x86_64']"
  />
</template>

<script setup lang="ts">
import JbrBenchmarkDashboard from "./JbrBenchmarkDashboard.vue"

const scenarios = ["Default"]
</script>
