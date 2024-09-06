<template>
  <GroupProjectsChart
    :label="`[CLion vs Radler] ${label}`"
    :measure="measure"
    :projects="projects"
  />
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const { label, measure, project } = defineProps<{
  label: string
  measure: string
  project: string
}>()

const projects = [`clion/${project}`, `radler/${project}`]
</script>
