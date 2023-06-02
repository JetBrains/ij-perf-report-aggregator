<template>
  <LineChart
    :title="props.label"
    :value-unit="props.valueUnit"
    :measures="[measure]"
    :configurators="configurators as DataQueryConfigurator[]"
    :skip-zero-values="false"
    :accidents="accidents"
  />
</template>

<script setup lang="ts">
import { ValueUnit } from "shared/src/chart"
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { FilterConfigurator } from "shared/src/configurators/filter"
import { DataQueryConfigurator } from "shared/src/dataQuery"
import { inject, onMounted } from "vue"
import { accidentsKeys, dashboardConfiguratorsKey, serverConfiguratorKey } from "../../shared/keys"
import LineChart from "./LineChart.vue"

interface Props {
  label: string
  measure: string
  projects: string[]
  valueUnit?: ValueUnit
}

const props = withDefaults(defineProps<Props>(), {
  valueUnit: "ms",
})

const serverConfigurator = inject(serverConfiguratorKey) as ServerConfigurator
const accidents = inject(accidentsKeys)
const dashboardConfigurators = inject(dashboardConfiguratorsKey) as DataQueryConfigurator[]|FilterConfigurator[]
const scenarioConfigurator = dimensionConfigurator(
  "project", 
  serverConfigurator,
  null, 
  true,
  [...dashboardConfigurators as FilterConfigurator[]]
)
const configurators = [...dashboardConfigurators, scenarioConfigurator, serverConfigurator]

onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
