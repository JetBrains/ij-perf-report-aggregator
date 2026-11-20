<template>
  <div>
    <div class="relative flex py-5 items-center">
      <div class="grow border-t" />
      <span class="shrink mx-4 text-lg">{{ measure }}</span>
      <div class="grow border-t" />
    </div>
    <section class="flex gap-x-6">
      <div
        v-for="platform in platforms"
        :key="platform.label"
        class="flex-1 min-w-0"
      >
        <GroupProjectsChart
          :label="platform.label"
          :measure="measure"
          :projects="platform.projects"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { jbrLinuxConfigurations, jbrMacConfigurations, jbrWindowsConfigurations } from "./configurations"

// One measure of a JBR benchmark, charted side by side for each platform's configurations (projects are
// named "<benchmark>_<configuration>").
const { measure, benchmark } = defineProps<{
  measure: string
  benchmark: string
}>()

const platforms = computed(() =>
  [
    { label: "macOS", configurations: jbrMacConfigurations },
    { label: "Ubuntu", configurations: jbrLinuxConfigurations },
    { label: "Windows", configurations: jbrWindowsConfigurations },
  ].map(({ label, configurations }) => ({ label, projects: configurations.map((configuration) => benchmark + "_" + configuration) }))
)
</script>
