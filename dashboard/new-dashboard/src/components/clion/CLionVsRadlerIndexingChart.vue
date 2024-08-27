<template>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="label"
          :measure="['ocSymbolBuildingTimeMs', 'backendIndexingTimeMs']"
          :projects="[clionProject, radlerProject]"
          :legend-formatter="legendFormatter"
        />
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const props = defineProps<{
  label: string
  project: string
}>()

const clionProject = `clion/${props.project}`
const radlerProject = `radler/${props.project}`
const label = `[CLion vs Radler] ${props.label}`
const frontendMetric = `${clionProject.replace("/indexing", "")} – ocSymbolBuildingTimeMs`
const backendMetric = `${radlerProject.replace("/indexing", "")} – backendIndexingTimeMs`

const legendFormatter = (name: string) => name.replace(frontendMetric, "CLion").replace(backendMetric, "Radler")
</script>
