<template>
  <LineChart
    :title="props.label"
    :value-unit="props.valueUnit"
    :measures="[measure]"
    :configurators="configurators"
    :skip-zero-values="false"
  />
</template>

<script setup lang="ts">
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { DataQueryConfigurator } from "shared/src/dataQuery"
import { onMounted } from "vue"
import LineChart from "./LineChart.vue"

interface Props {
  label: string
  measure: string
  projects: Array<string>
  serverConfigurator: ServerConfigurator
  configurators: Array<DataQueryConfigurator>
  valueUnit?: "ns" | "ms"
}

const props = withDefaults(defineProps<Props>(), {
  valueUnit: "ms"
})

const scenarioConfigurator = dimensionConfigurator(
  "project", 
  props.serverConfigurator, 
  null, 
  true
)
const configurators = [...props.configurators, scenarioConfigurator]

onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
