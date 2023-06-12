<template>
  <LineChart
    :title="props.label"
    :value-unit="props.valueUnit"
    :measures="Array.isArray(measure) ? measure : [measure]"
    :configurators="configurators as DataQueryConfigurator[]"
    :skip-zero-values="false"
  />
</template>

<script setup lang="ts">
import { inject, onMounted } from "vue"
import { dimensionConfigurator } from "../../configurators/DimensionConfigurator"
import { ServerConfigurator } from "../../configurators/ServerConfigurator"
import { FilterConfigurator } from "../../configurators/filter"
import { dashboardConfiguratorsKey, serverConfiguratorKey } from "../../shared/keys"
import { ValueUnit } from "../common/chart"
import { DataQueryConfigurator } from "../common/dataQuery"
import LineChart from "./LineChart.vue"

interface Props {
  label: string
  measure: string | string[]
  projects: string[]
  valueUnit?: ValueUnit
}

const props = withDefaults(defineProps<Props>(), {
  valueUnit: "ms",
})

const serverConfigurator = inject(serverConfiguratorKey) as ServerConfigurator
const dashboardConfigurators = inject(dashboardConfiguratorsKey) as DataQueryConfigurator[] | FilterConfigurator[]
const scenarioConfigurator = dimensionConfigurator("project", serverConfigurator, null, true, [...(dashboardConfigurators as FilterConfigurator[])])
const configurators = [...dashboardConfigurators, scenarioConfigurator, serverConfigurator]

onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
