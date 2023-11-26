<template>
  <section
    v-for="definition in definitions"
    v-show="definition.label.indexOf('firstElementShown') === -1"
    :key="definition.label"
    class="flex gap-x-6"
  >
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`${definition.label} ${metrics.value} K1`"
        :measure="metrics.value"
        :projects="definition.projects.map((s) => `${s}_k1`)"
        :legend-formatter="replaceKotlinName"
      />
    </div>
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`${definition.label} ${metrics.value} K2`"
        :measure="metrics.value"
        :projects="definition.projects.map((s) => `${s}_k2`)"
        :legend-formatter="replaceKotlinName"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import { Ref } from "vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { replaceKotlinName } from "./label-formatter"
import { ProjectsChartDefinition } from "./projects"

interface Props {
  definitions: ProjectsChartDefinition[]
  metrics: Ref<string | string[]>
}

defineProps<Props>()
</script>
