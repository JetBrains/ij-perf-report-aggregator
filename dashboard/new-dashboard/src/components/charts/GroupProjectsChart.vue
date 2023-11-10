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
import { removeCommonSegments } from "../../util/removeCommonPrefixes"
import { ValueUnit } from "../common/chart"
import LineChart from "./PerformanceLineChart.vue"

interface Props {
  label: string
  measure: string | string[]
  projects: string[]
  valueUnit?: ValueUnit
  legendFormatter?: (name: string) => string
  aliases?: string[] | null
}

const props = withDefaults(defineProps<Props>(), {
  valueUnit: "ms",
  legendFormatter: (name: string) => name,
  aliases: null,
})

const serverConfigurator = injectOrError(serverConfiguratorKey)
const dashboardConfigurators = injectOrError(dashboardConfiguratorsKey)
const aliases = props.aliases ?? removeCommonSegments(props.projects)
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, null, true, [...(dashboardConfigurators as FilterConfigurator[])], null, aliases)
const configurators = [...dashboardConfigurators, scenarioConfigurator, serverConfigurator]

watch(
  () => props.projects,
  (projects) => {
    scenarioConfigurator.selected.value = projects
  },
  { immediate: true }
)
</script>
