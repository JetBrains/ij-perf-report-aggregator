<template>
  <GroupProjectsChart
    :label="label"
    :measure="props.measure"
    :projects="projects"
  />
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const props = defineProps<{
  label: string
  measure: string
  project: string
}>()

const projects = [`clion/${props.project}`, `radler/${props.project}`]
const label = `[CLion vs Radler] ${props.label}`
</script>
