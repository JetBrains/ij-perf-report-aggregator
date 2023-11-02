<template>
  <LineChart
    :title="props.label"
    :value-unit="props.valueUnit"
    :measures="Array.isArray(measure) ? measure : [measure]"
    :configurators="configurators"
    :skip-zero-values="false"
    :legend-formatter="props.legendFormatter"
  />
</template>

<script setup lang="ts">
import { watch } from "vue"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { injectOrError } from "../../shared/injectionKeys"
import { dashboardConfiguratorsKey, serverConfiguratorKey } from "../../shared/keys"
import { ValueUnit } from "../common/chart"
import LineChart from "./PerformanceLineChart.vue"

interface Props {
  label: string
  measure: string | string[]
  projects: string[]
  valueUnit?: ValueUnit
  legendFormatter?: (name: string) => string
}

const props = withDefaults(defineProps<Props>(), {
  valueUnit: "ms",
  legendFormatter: (name: string) => name,
})

const serverConfigurator = injectOrError(serverConfiguratorKey)
const dashboardConfigurators = injectOrError(dashboardConfiguratorsKey)
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, null, true, [...(dashboardConfigurators as FilterConfigurator[])])
const configurators = [...dashboardConfigurators, scenarioConfigurator, serverConfigurator]

watch(
  () => props.projects,
  (projects) => {
    scenarioConfigurator.selected.value = projects
  },
  { immediate: true }
)
</script>
