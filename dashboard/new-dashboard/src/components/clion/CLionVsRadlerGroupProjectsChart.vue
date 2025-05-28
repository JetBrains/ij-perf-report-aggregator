<template>
  <GroupProjectsChart
    :label="`[CLion vs Radler] ${label}`"
    :measure="measureArray"
    :projects="projects"
  />
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const { label, measure, project } = defineProps<{
  label: string
  measure: string | string[]
  project: string
}>()

const measureArray = Array.isArray(measure) ? measure : [measure]
const projects = [`clion/${project}`, `radler/${project}`]
</script>
