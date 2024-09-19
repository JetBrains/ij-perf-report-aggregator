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
        <section class="flex gap-x-6 flex-col md:flex-row">
          <div class="flex-1 min-w-0">
            <section>
              <GroupProjectsChart
                :key="`radler-${selectedIndex}`"
                :label="`[Radler] ${label}, Mb`"
                :measure="[backendMeasures[selectedIndex], frontendMeasures[selectedIndex]]"
                :projects="[radlerProject]"
                :legend-formatter="legendFormatter"
              />
            </section>
          </div>

          <div class="flex-1 min-w-0">
            <section>
              <GroupProjectsChart
                :key="`clion-${selectedIndex}`"
                :label="`[CLion] ${label}, Mb`"
                :measure="[frontendMeasures[selectedIndex]]"
                :projects="[clionProject]"
                :legend-formatter="legendFormatter"
              />
            </section>
          </div>
        </section>
      </TabPanel>
    </TabPanels>
  </Tabs>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { ref } from "vue"

const { label, measure, project } = defineProps<{
  label: string
  measure: string
  project: string
}>()

const selectedIndex = ref(0)

const getAllMeasures = (prefix: string) => [`${prefix}/beforeGC`, prefix, `${prefix}/idle`]
const tabDescription = ["Before GC", "After GC", "After GC (idle)"]
const clionProject = `clion/${project}`
const radlerProject = `radler/${project}`
const backendMeasurePrefix = `rd.memory.allocatedManagedMemoryMb/${measure}`
const frontendMeasurePrefix = `JVM.heapUsageMb/${measure}`

const backendMeasures = getAllMeasures(backendMeasurePrefix)
const frontendMeasures = getAllMeasures(frontendMeasurePrefix)

const legendFormatter = (name: string) => {
  const lengthDescComparer = (a: string, b: string) => b.length - a.length

  name = name.replace(clionProject, "JVM (Frontend)") // HACK: remove this line when CLion Classic gets more than one memory metric
  for (const frontendMeasure of frontendMeasures.toSorted(lengthDescComparer)) {
    name = name.replace(frontendMeasure, "JVM (Frontend)")
  }
  for (const backendMeasure of backendMeasures.toSorted(lengthDescComparer)) {
    name = name.replace(backendMeasure, ".NET (Backend)")
  }
  return name
}
</script>
