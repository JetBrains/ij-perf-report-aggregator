<template>
  <section
    v-for="definition in definitions"
    :key="definition.labels[0]"
    class="flex-1 min-w-0"
  >
    <GroupProjectsChart
      :label="definition.labels[0]"
      :measure="definition.measures as string[] | string"
      :projects="definition.projects.map((s) => `${s}_k2`)"
      :legend-formatter="replaceKotlinName"
      :machines="definition.machines"
    />
  </section>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ChartDefinition } from "../charts/DashboardCharts"
import { replaceKotlinName } from "./label-formatter"

interface Props {
  definitions: ChartDefinition[]
}

const { definitions } = defineProps<Props>()
</script>
