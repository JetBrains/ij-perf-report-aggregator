<template>
  <GroupProjectsChart
    :can-be-closed="false"
    :value-unit="valueUnit"
    :legend-formatter="legendFormatter"
    :machines="machines"
    :label="label"
    :projects="projectWithClient"
    :measure="measure"
    :aliases="aliasesWithClient"
  />
</template>

<script setup lang="ts">
import { ValueUnit } from "../common/chart"
import GroupProjectsChart from "./GroupProjectsChart.vue"

interface Props {
  label: string
  measure: string | string[]
  projects: string[]
  machines?: string[] | null
  valueUnit?: ValueUnit
  legendFormatter?: (name: string) => string
  aliases?: string[] | null
}

const { label, measure, projects, machines = null, valueUnit = "ms", legendFormatter = (name: string) => name, aliases = null } = defineProps<Props>()

const projectWithClient = [...projects, ...projects.map((project) => `${project}/embeddedClient`)]
const aliasesWithClient = aliases != null ? [...aliases, ...aliases] : null
</script>
