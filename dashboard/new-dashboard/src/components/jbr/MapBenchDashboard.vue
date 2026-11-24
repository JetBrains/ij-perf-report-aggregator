<template>
  <JbrBenchmarkDashboard
    benchmark="Mapbench"
    :groups="[{ title: null, scenarios }]"
    value-unit="fps"
    better-direction="higher"
    description="Frames per second of drawing each recorded MapBench file, higher is better."
    :default-configurations="['macOS15aarch64Metal', 'Ubuntu2404x86_64MM_Wayland', 'Windows11x86_64']"
  />
</template>

<script setup lang="ts">
import JbrBenchmarkDashboard from "./JbrBenchmarkDashboard.vue"

const scenarios = [
  "CircleTests",
  "EllipseTests-fill-false",
  "EllipseTests-fill-true",
  "dc_boulder_2013-13-30-06-13-17",
  "dc_boulder_2013-13-30-06-13-20",
  "dc_shp_alllayers_2013-00-30-07-00-43",
  "dc_shp_alllayers_2013-00-30-07-00-47",
  "dc_spearfish_2013-11-30-06-11-15",
  "dc_spearfish_2013-11-30-06-11-19",
  "dc_topp:states_2013-11-30-06-11-06",
  "dc_topp:states_2013-11-30-06-11-07",
  "spiralTest-dash-false",
  "spiralTest-fill",
  "test_z_625k",
].map((scenario) => scenario + ".ser")
</script>
