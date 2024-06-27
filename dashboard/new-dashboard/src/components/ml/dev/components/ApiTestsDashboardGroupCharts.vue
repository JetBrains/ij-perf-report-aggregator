<template>
  <section
    v-for="definition in definitions"
    :key="definition.labels[0]"
    class="flex gap-x-6"
  >
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`stage: ${definition.labels[0]}`"
        :measure="definition.measures as string[] | string"
        :projects="definition.projects.map((s) => `stage: ${s}`)"
        :machines="definition.machines"
      />
    </div>
    <div class="flex-1 min-w-0">
      <GroupProjectsChart
        :label="`prod: ${definition.labels[0]}`"
        :measure="definition.measures as string[] | string"
        :projects="definition.projects.map((s) => `prod: ${s}`)"
        :machines="definition.machines"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import { ChartDefinition } from "../../../charts/DashboardCharts"
import GroupProjectsChart from "../../../charts/GroupProjectsChart.vue"

interface Props {
  definitions: ChartDefinition[]
}

defineProps<Props>()
</script>
