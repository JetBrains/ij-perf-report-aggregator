<template>
  <JbrBenchmarkDashboard
    benchmark="JavaDraw"
    :groups="[{ title: null, scenarios }]"
    value-unit="fps"
    better-direction="higher"
    description="Frames per second, higher is better."
    :default-configurations="['macOS15aarch64Metal', 'Ubuntu2404x86_64MM_Wayland', 'Windows11x86_64']"
  />
</template>

<script setup lang="ts">
import JbrBenchmarkDashboard from "./JbrBenchmarkDashboard.vue"

const scenarios = [
  "Display_PNG_Bitmap_Twice_Pass_1",
  "Display_PNG_Bitmap_Twice_Pass_2",
  "Plus_2_SweepGradient_Circles",
  "Plus_200_Random_Small_Circles",
  "Plus_320_Long_Lines",
  "Plus_4000_Random_Small_Circles",
]
</script>
