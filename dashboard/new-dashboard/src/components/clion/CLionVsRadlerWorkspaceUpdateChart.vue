<template>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="`[Radler] ${props.label}, Total Time`"
          :measure="measure"
          :projects="radlerProjects"
          :legend-formatter="legendFormatter"
        />
      </section>
    </div>

    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="`[CLion] ${props.label}, Total Time`"
          :measure="measure"
          :projects="clionProjects"
          :legend-formatter="legendFormatter"
        />
      </section>
    </div>
  </section>
</template>

<script setup lang="ts">
import GroupProjectsChart from "../charts/GroupProjectsChart.vue"
import { removeCommonSegments } from "../../util/removeCommonPrefixes"

const props = defineProps<{
  label: string
  projects: string[]
  names: string[]
}>()

console.assert(props.projects.length == props.names.length)
const getAllProjects = (prefix: string) => props.projects.map((project) => `${prefix}/${project}`)
const clionProjects = getAllProjects("clion")
const radlerProjects = getAllProjects("radler")
const measure = "workspaceModel.updates.ms"

const projectToNameMap = new Map<string, string>()
for (let i = 0; i < props.names.length; i++) {
  const clionProject = removeCommonSegments(clionProjects)[i]
  const radlerProject = removeCommonSegments(radlerProjects)[i]
  projectToNameMap.set(clionProject, props.names[i])
  projectToNameMap.set(radlerProject, props.names[i])
}

const legendFormatter = (name: string) => {
  for (const [project, prettyName] of projectToNameMap) {
    name = name.replace(project, prettyName)
  }

  return name
}
</script>
