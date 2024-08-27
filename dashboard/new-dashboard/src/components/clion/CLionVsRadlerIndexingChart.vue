<template>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <!-- TODO: find better way of presenting two projects with different metrics on the same chart -->
        <GroupProjectsChart
          :label="label"
          :measure="['fus_dumb_indexing_time', 'backendIndexingTimeMs']"
          :projects="[clionProject, radlerProject]"
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
</script>
