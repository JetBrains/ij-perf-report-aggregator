<template>
  <section
    v-for="definition in definitions"
    :key="definition.labels[0]"
    class="flex gap-x-6"
  >
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`${definition.labels[0]} K1`"
        :measure="definition.measures"
        :projects="definition.projects.map((s) => `${s}_k1`)"
        :legend-formatter="replaceKotlinName"
        :machines="definition.machines"
      />
    </div>
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`${definition.labels[0]} K2`"
        :measure="definition.measures"
        :projects="definition.projects.map((s) => `${s}_k2`)"
        :legend-formatter="replaceKotlinName"
        :machines="definition.machines"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition } from "../charts/DashboardCharts"
import { replaceKotlinName } from "./label-formatter"

interface Props {
  definitions: ChartDefinition[]
}

defineProps<Props>()
</script>
