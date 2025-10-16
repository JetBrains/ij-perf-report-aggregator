<template>
  <LineChart
    :title="label"
    :value-unit="valueUnit"
    :measures="measureArray"
    :configurators="configurators"
    :skip-zero-values="false"
    :legend-formatter="legendFormatter"
    :can-be-closed="canBeClosed"
    @chart-closed="onChartClosed"
  />
</template>

<script setup lang="ts">
import { computed, Ref, watch } from "vue"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { injectOrError } from "../../shared/injectionKeys"
import { dashboardConfiguratorsKey, serverConfiguratorKey } from "../../shared/keys"
import { removeCommonSegments } from "../../util/removeCommonPrefixes"
import { ValueUnit } from "../common/chart"
import LineChart from "./LineChart.vue"
import { MachineConfigurator } from "../../configurators/MachineConfigurator"

interface Props {
  label: string
  measure: string | string[]
  projects: string[]
  machines?: string[] | null
  valueUnit?: ValueUnit
  legendFormatter?: (name: string) => string
  aliases?: string[] | null
  canBeClosed?: boolean
}

const { label, measure, projects, machines = null, valueUnit = "ms", legendFormatter = (name: string) => name, aliases = null, canBeClosed = false } = defineProps<Props>()

const serverConfigurator = injectOrError(serverConfiguratorKey)
const dashboardConfigurators = injectOrError(dashboardConfiguratorsKey)
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, null, true)
const configurators = [...dashboardConfigurators, scenarioConfigurator, serverConfigurator]

if (machines != null) {
  const machineConfigurator = new MachineConfigurator(serverConfigurator, undefined, [], true, machines)
  configurators.push(machineConfigurator)
}

const measureArray: Ref<string[]> = computed(() => {
  return Array.isArray(measure) ? measure : [measure]
})

const emit = defineEmits<{
  chartClosed: [projects: string[]]
}>()

function onChartClosed(): void {
  emit("chartClosed", projects)
}

watch(
  () => [projects, aliases],
  ([projects, aliases]) => {
    if (projects != null) {
      const aliasesToUse = aliases ?? removeCommonSegments(projects)
      scenarioConfigurator.aliases = new Map(projects.map((key, index) => [key, aliasesToUse[index]]))
    }
    scenarioConfigurator.selected.value = projects
  },
  { immediate: true }
)
</script>
