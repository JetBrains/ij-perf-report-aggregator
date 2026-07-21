<template>
  <GroupProjectsChart
    :label="label"
    :measure="measure"
    :projects="resolvedProjects"
    :aliases="aliases"
  />
</template>

<script setup lang="ts">
import { computed } from "vue"
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { toNewSeProjects, useNewSearchEverywhere } from "./searchEverywhereMetrics"

const {
  label,
  measure,
  projects,
  aliases = null,
} = defineProps<{
  label: string
  measure: string | string[]
  projects: string[]
  aliases?: string[] | null
}>()

const useNewSe = useNewSearchEverywhere()
const resolvedProjects = computed(() => (useNewSe.value ? toNewSeProjects(projects) : projects))
</script>
