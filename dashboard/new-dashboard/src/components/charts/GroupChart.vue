<template>
  <div class="flex flex-col gap-y-2.5 py-3 px-5 border border-solid rounded-md border-zinc-200">
    <h3 class="uppercase m-0 text-sm">
      {{ label }}
    </h3>
    <LineChart
      :value-unit="props.valueUnit"
      :measures="[measure]"
      :configurators="configurators"
      :skip-zero-values="false"
    />
  </div>
</template>

<script setup lang="ts">
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { DataQueryConfigurator } from "shared/src/dataQuery"
import { onMounted } from "vue"
import LineChart from "./LineChart.vue"

interface GroupChartProps {
  label: string
  measure: string
  projects: Array<string>
  serverConfigurator: ServerConfigurator
  configurators: Array<DataQueryConfigurator>
  valueUnit?: "ns" | "ms"
}

const props = withDefaults(defineProps<GroupChartProps>(), {
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
