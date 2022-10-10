<template>
  <div class="flex flex-col gap-y-2.5 py-3 px-5 border border-solid rounded-md border-zinc-200">
    <h3 class="uppercase m-0 text-sm font-semibold">
      {{ label }}
    </h3>
    <LineChart
      :compound-tooltip="true"
      :chart-type="'line'"
      :value-unit="props.valueUnit"
      :measures="[measure]"
      :configurators="configurators"
      :skip-zero-values="false"
      trigger="item"
    />
  </div>
</template>

<script setup lang="ts">
import { dimensionConfigurator } from "shared/src/configurators/DimensionConfigurator"
import { ServerConfigurator } from "shared/src/configurators/ServerConfigurator"
import { configuratorListKey } from "shared/src/injectionKeys"
import { inject, onMounted } from "vue"
import LineChart from "./LineChart.vue"

interface GroupChartProps {
  label: string
  measure: string
  projects: Array<string>
  serverConfigurator: ServerConfigurator
  valueUnit?: "ns" | "ms"
}

const props = withDefaults(defineProps<GroupChartProps>(), {
  valueUnit: "ms"
})
const providedConfigurators = inject(configuratorListKey, null)

if (providedConfigurators == null) {
  throw new Error("`dataQueryExecutor` is not provided")
}

const scenarioConfigurator = dimensionConfigurator("project", props.serverConfigurator, null, true)
const configurators = [...providedConfigurators, scenarioConfigurator]

onMounted(() => {
  scenarioConfigurator.selected.value = props.projects
})
</script>
