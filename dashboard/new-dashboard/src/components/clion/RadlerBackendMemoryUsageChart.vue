<template>
  <Tabs v-model:value="selectedIndex">
    <TabList>
      <Tab
        v-for="(item, index) in tabDescription"
        :key="item"
        :value="index"
      >
        {{ item }}
      </Tab>
    </TabList>

    <TabPanels>
      <TabPanel
        v-for="(item, index) in tabDescription"
        :key="item"
        :value="index"
      >
        <section>
          <GroupProjectsChart
            :key="`radler-${selectedIndex}`"
            :label="`[Radler] ${label}, Mb`"
            :measure="[backendMeasures[selectedIndex]]"
            :projects="[radlerProject]"
            :legend-formatter="legendFormatter"
          />
        </section>
      </TabPanel>
    </TabPanels>
  </Tabs>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ref } from "vue"

const { label, measure, project, metric } = defineProps<{
  label: string
  measure: string
  project: string
  metric: "rd.memory.allocatedManagedMemoryMb" | "rd.memory.privateMemorySizeMb" | "rd.memory.workingSetMb"
}>()

const selectedIndex = ref(1)

const getAllMeasures = (prefix: string) => [`${prefix}/beforeGC`, prefix, `${prefix}/idle`]
const tabDescription = ["Before GC", "After GC", "After GC (idle)"]
const radlerProject = `radler/${project}`
const backendMeasurePrefix = `${metric}/${measure}`

const backendMeasures = getAllMeasures(backendMeasurePrefix)

const lengthDescComparer = (a: string, b: string) => b.length - a.length

const legendFormatter = (name: string) => {
  for (const backendMeasure of backendMeasures.toSorted(lengthDescComparer)) {
    name = name.replace(backendMeasure, ".NET (Backend)")
  }
  return name
}
</script>
