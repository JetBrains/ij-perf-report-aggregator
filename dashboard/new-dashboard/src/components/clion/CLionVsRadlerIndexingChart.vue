<template>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="`[CLion vs Radler]` + label"
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

const { label, project } = defineProps<{
  label: string
  project: string
}>()

const clionProject = `clion/${project}`
const radlerProject = `radler/${project}`
const frontendMetric = `${clionProject.replace("/indexing", "")} – ocSymbolBuildingTimeMs`
const backendMetric = `${radlerProject.replace("/indexing", "")} – backendIndexingTimeMs`

const legendFormatter = (name: string) => name.replace(frontendMetric, "CLion").replace(backendMetric, "Radler")
</script>
