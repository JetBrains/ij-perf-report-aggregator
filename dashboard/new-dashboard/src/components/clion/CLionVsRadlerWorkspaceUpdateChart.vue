<template>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="`[Radler] ${label}, Total Time`"
          measure="cidr.workspace.metrics#duration_total_ms"
          :projects="radlerProjects"
          :legend-formatter="legendFormatter"
        />
      </section>
    </div>
    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="`[CLion] ${label}, Total Time`"
          measure="cidr.workspace.metrics#duration_total_ms"
          :projects="clionProjects"
          :legend-formatter="legendFormatter"
        />
      </section>
    </div>
  </section>
  <section class="flex gap-x-6 flex-col md:flex-row">
    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="`[Radler] ${label}, In Write Action Time`"
          measure="cidr.workspace.metrics#duration_in_write_action_ms"
          :projects="radlerProjects"
          :legend-formatter="legendFormatter"
        />
      </section>
    </div>
    <div class="flex-1 min-w-0">
      <section>
        <GroupProjectsChart
          :label="`[CLion] ${label}, In Write Action Time`"
          measure="cidr.workspace.metrics#duration_in_write_action_ms"
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

const { label, projects, names } = defineProps<{
  label: string
  projects: string[]
  names: string[]
}>()

console.assert(projects.length == names.length)
const getAllProjects = (prefix: string) => projects.map((project) => `${prefix}/${project}`)
const clionProjects = getAllProjects("clion")
const radlerProjects = getAllProjects("radler")

const projectToNameMap = new Map<string, string>()
for (const [i, name] of names.entries()) {
  const clionProject = removeCommonSegments(clionProjects)[i]
  const radlerProject = removeCommonSegments(radlerProjects)[i]
  projectToNameMap.set(clionProject, name)
  projectToNameMap.set(radlerProject, name)
}

const legendFormatter = (name: string) => {
  for (const [project, prettyName] of projectToNameMap) {
    name = name.replace(project, prettyName)
  }

  return name
}
</script>
