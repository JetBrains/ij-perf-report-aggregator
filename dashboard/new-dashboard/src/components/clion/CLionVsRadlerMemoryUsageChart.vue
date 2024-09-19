<template>
  <TabGroup
    v-slot="{ selectedIndex }"
    :default-index="1"
  >
    <div class="border-b border-gray-200">
      <TabList class="-mb-px flex space-x-8">
        <Tab
          v-for="item in tabDescription"
          :key="item"
          v-slot="{ selected }"
          as="template"
        >
          <button
            :class="[
              selected ? 'border-indigo-500 text-indigo-600' : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
              'whitespace-nowrap py-1 px-1 border-b-2 font-medium',
            ]"
          >
            {{ item }}
          </button>
        </Tab>
      </TabList>
    </div>

    <TabPanels>
      <TabPanel
        v-for="item in tabDescription"
        :key="item"
        class="p-3"
        :unmount="false"
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
  </TabGroup>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from "@headlessui/vue"

const { label, measure, project } = defineProps<{
  label: string
  measure: string
  project: string
}>()

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
