<template>
  <section
    v-for="definition in definitions"
    :key="definition.label"
    class="flex gap-x-6"
  >
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`${definition.label} K1`"
        :measure="definition.measure"
        :projects="definition.projects.map((s) => `${s}_k1`)"
        :legend-formatter="replaceKotlinName"
      />
    </div>
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`${definition.label} K2`"
        :measure="definition.measure"
        :projects="definition.projects.map((s) => `${s}_k2`)"
        :legend-formatter="replaceKotlinName"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { replaceKotlinName } from "./label-formatter"
import { ProjectsChartDefinition } from "./projects"

interface Props {
  definitions: ProjectsChartDefinition[]
}

defineProps<Props>()
</script>
