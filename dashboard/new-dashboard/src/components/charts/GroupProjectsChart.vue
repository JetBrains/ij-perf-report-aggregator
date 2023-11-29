<template>
  <LineChart
    :title="props.label"
    :value-unit="props.valueUnit"
    :measures="measureArray"
    :configurators="configurators"
    :skip-zero-values="false"
    :legend-formatter="props.legendFormatter"
  />
</template>

<script setup lang="ts">
import { computed, Ref, watch } from "vue"
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
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, null, true, [...(dashboardConfigurators as FilterConfigurator[])], null)
const configurators = [...dashboardConfigurators, scenarioConfigurator, serverConfigurator]
const measureArray: Ref<string[]> = computed(() => {
  return Array.isArray(props.measure) ? props.measure : [props.measure]
})

watch(
  () => [props.projects, props.aliases],
  ([projects, aliases]) => {
    if (projects != null) {
      scenarioConfigurator.aliases = aliases ?? removeCommonSegments(projects)
    }
    scenarioConfigurator.selected.value = projects
  },
  { immediate: true }
)
</script>
