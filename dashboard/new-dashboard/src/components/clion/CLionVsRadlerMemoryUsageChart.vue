<template>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <!-- TODO: add aliases -->
        <GroupProjectsChart
          :label="radlerLabel"
          :measure="[backendMeasure, frontendMeasure]"
          :projects="[radlerProject]"
        />
      </section>
    </div>

    <div class="flex-1 min-w-0">
      <section>
        <!-- TODO: add aliases -->
        <GroupProjectsChart
          :label="clionLabel"
          :measure="[frontendMeasure]"
          :projects="[clionProject]"
        />
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"

const props = defineProps<{
  label: string
  measure: string
  project: string
}>()

const clionProject = `clion/${props.project}`
const radlerProject = `radler/${props.project}`
const clionLabel = `[CLion] ${props.label}, Mb`
const radlerLabel = `[Radler] ${props.label}, Mb`
const frontendMeasure = `JVM.heapUsageMb/${props.measure}`
const backendMeasure = `rd.memory.allocatedManagedMemoryMb/${props.measure}`
</script>
